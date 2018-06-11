package interactors

import (
	"context"
	"fmt"
	"github.com/bigblind/makker/channels"
	"github.com/bigblind/makker/di"
	"github.com/bigblind/makker/games"
	"strings"
	"time"
)

var Interactor GamesInteractor

func init() {
	Interactor = NewInteractor()

	initChannels()
}

func initChannels() {
	err := di.Graph.Invoke(func(gs games.GameStore, cp channels.ChannelProvider) {

		cp.SetUserChecker("games", func(ctx context.Context, channel channels.Channel, userId string) error {

			parts := strings.Split(channel.Id(), ";")

			if userId != parts[1] {
				return fmt.Errorf("this is not your private channel")
			}

			instId := parts[0]
			inst, err := gs.GetInstanceById(ctx, instId)
			if err != nil {
				return err
			}

			if inst.MetaState == games.WaitingForPlayers && !inst.HasPlayer(userId) {
				return Interactor.joinInstance(ctx, inst, userId)
			}

			if !inst.HasPlayer(userId) {
				return fmt.Errorf("you're not in this game.")
			}

			return nil
		})

		cp.OnLeave("games", func(ctx context.Context, channel channels.Channel, userId, socketId string) {
			inter := NewInteractor()

			parts := strings.Split(channel.Id(), ";")

			instId := parts[0]
			inter.LeaveGame(ctx, instId, userId)
		})
	})

	if err != nil {
		panic(err)
	}
}

type GamesInteractor struct {
	store games.GameStore
	cp    channels.ChannelProvider
}

func NewInteractor() GamesInteractor {
	var inter GamesInteractor
	err := di.Graph.Invoke(func(gs games.GameStore, cp channels.ChannelProvider) {
		inter = GamesInteractor{
			gs,
			cp,
		}
	})

	if err != nil {
		panic(err)
	}

	return inter
}

func (inter GamesInteractor) CreateInstance(ctx context.Context, gameName, userId string) (instanceResponse, error) {
	g, err := games.Registry.GetGameLatestVersion(gameName)

	if err != nil {
		return instanceResponse{}, err
	}

	inst := games.NewInstance(g, userId)
	inst.AddPlayer(userId)

	err = inter.store.SaveInstance(ctx, inst)
	if err != nil {
		return instanceResponse{}, err
	}

	return instanceToResponse(inst, userId, inter.cp), nil
}

func (inter GamesInteractor) JoinGame(ctx context.Context, instanceId, userId string) error {
	inst, err := inter.store.GetInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	return inter.joinInstance(ctx, inst, userId)
}

func (inter GamesInteractor) joinInstance(ctx context.Context, inst *games.GameInstance, userId string) error {
	if inst.HasPlayer(userId) {
		return fmt.Errorf("%v is already in the game.", userId)
	}

	inst.AddPlayer(userId)

	inter.EmitPublic(ctx, inst, "player_join", map[string]string{"user_id": userId})

	return inter.store.SaveInstance(ctx, inst)
}

func (inter GamesInteractor) LeaveGame(ctx context.Context, instanceId, userId string) error {
	inst, err := inter.store.GetInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	inst.RemovePlayer(userId)

	inter.EmitPublic(ctx, inst, "player_leave", map[string]string{"user_id": userId})

	return inter.store.SaveInstance(ctx, inst)
}

func (inter GamesInteractor) StartGame(ctx context.Context, instanceId, userId string) error {
	inst, err := inter.store.GetInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if inst.AdminUserId != userId {
		return fmt.Errorf("you're not the admin of this game")
	}

	inst.ShufflePlayers()
	inst.MetaState = games.InProgress
	inst.Game().InitializeState(&inst.State)

	err = inter.store.SaveInstance(ctx, inst)
	if err != nil {
		return err
	}

	inter.EmitMetaState(ctx, inst)

	return nil
}

func (inter GamesInteractor) GetInstance(ctx context.Context, instanceId string, userId ...string) (instanceResponse, error) {
	inst, err := inter.store.GetInstanceById(ctx, instanceId)
	if err != nil {
		return instanceResponse{}, err
	}

	uid := ""
	if len(userId) != 0 {
		uid = userId[0]
	}

	return instanceToResponse(inst, uid, inter.cp), err
}

func (inter GamesInteractor) MakeMove(ctx context.Context, instanceId, userId string, moveData interface{}) (instanceResponse, error) {
	inst, err := inter.store.GetInstanceById(ctx, instanceId)
	if err != nil {
		return instanceResponse{}, err
	}

	idx := inst.GetPlayerIndex(userId)
	if idx < 0 {
		return instanceToResponse(inst, userId, inter.cp), fmt.Errorf("you're not in this game.")
	}

	game := inst.Game()

	if !game.CanPlayerMove(int(idx), &inst.State) {
		return instanceToResponse(inst, userId, inter.cp), fmt.Errorf("you can't make a move right now")
	}

	move := games.Move{
		Data:   moveData,
		Player: int8(idx),
		Time:   time.Now(),
	}

	err = game.HandleUpdate(&inst.State, move)
	if err != nil {
		return instanceToResponse(inst, userId, inter.cp), err
	}

	inst.Moves = append(inst.Moves, move)

	if game.IsGameOver(&inst.State) {
		inst.MetaState = games.GameOver
	}

	err = inter.store.SaveInstance(ctx, inst)
	if err != nil {
		return instanceToResponse(inst, userId, inter.cp), err
	}

	events := make([]channels.Event, 0, len(inst.State.Players)*2+1)

	pubc := inter.cp.NewChannel(ctx, "games", inst.Channels("").Public, true)
	events = append(events, channels.Event{
		Channel: pubc,
		Name:    "state",
		Data:    inst.State.SharedState,
	})

	for _, p := range inst.State.Players {
		privc := inter.cp.NewChannel(ctx, "games", inst.Channels(p.UserId).Private, false)

		events = append(events, channels.Event{
			Channel: privc,
			Name:    "private_state",
			Data:    p.PrivateState,
		},
			channels.Event{
				Channel: pubc,
				Name:    "public_state",
				Data: map[string]interface{}{
					"user_id": p.UserId,
					"data":    p.PublicState,
					"score":   p.Score,
				},
			})
	}

	inter.cp.EmitBatch(ctx, events)

	if inst.MetaState == games.GameOver {
		inter.EmitMetaState(ctx, inst)
	}

	return instanceToResponse(inst, userId, inter.cp), nil
}

func (inter GamesInteractor) ListInstances(ctx context.Context, gname string, state ...games.MetaState) (*[]instanceResponse, error) {
	insts, err := inter.store.GetInstancesByGame(ctx, gname, state...)
	if err != nil {
		return nil, err
	}

	ris := make([]instanceResponse, len(*insts))
	for i, inst := range *insts {
		ris[i] = instanceToResponse(&inst, "", inter.cp)
	}

	return &ris, nil
}

func (inter GamesInteractor) EmitMetaState(ctx context.Context, inst *games.GameInstance) {
	inter.EmitPublic(ctx, inst, "meta_state", map[string]games.MetaState{"state": inst.MetaState})
}

func (inter GamesInteractor) EmitPublic(ctx context.Context, inst *games.GameInstance, event string, data interface{}) {
	cs := inst.Channels("")

	c := inter.cp.NewChannel(ctx, "games", cs.Public, true)
	c.Emit(event, data)
}

func (inter GamesInteractor) EmitPrivate(ctx context.Context, inst *games.GameInstance, userId, event string, data interface{}) {
	cs := inst.Channels(userId)

	c := inter.cp.NewChannel(ctx, "games", cs.Private, true)
	c.Emit(event, data)
}

type instanceResponsePlayer struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Score  int32  `json:"score"`
}

type instanceResponse struct {
	Id       string                   `json:"id"`
	CreatedAt time.Time				  `json:"created_at"`
	GameInfo games.GameInfo           `json:"game_info"`
	State    games.MetaState          `json:"state"`
	Players  []instanceResponsePlayer `json:"players"`

	PublicChannel  string `json:"public_channel"`
	PrivateChannel string `json:"private_channel"`
}

func instanceToResponse(i *games.GameInstance, uid string, cp channels.ChannelProvider) instanceResponse {
	chanIds := i.Channels(uid)
	ps := make([]instanceResponsePlayer, len(i.State.Players))
	for j, p := range i.State.Players {
		ps[j] = instanceResponsePlayer{
			UserId: p.UserId,
			Name:   p.UserId,
			Score:  p.Score,
		}
	}

	return instanceResponse{
		Id:       i.Id,
		CreatedAt:i.CreatedAt,
		GameInfo: i.Game().Info(),
		State:    i.MetaState,
		Players:  ps,

		PublicChannel:  cp.NewChannel(nil, "games", chanIds.Public, true).ClientId(),
		PrivateChannel: cp.NewChannel(nil, "games", chanIds.Private, false).ClientId(),
	}
}