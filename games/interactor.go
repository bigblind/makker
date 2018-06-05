package games

import (
	"context"
	"fmt"
	"github.com/bigblind/makker/channels"
	"github.com/bigblind/makker/di"
	"strings"
	"time"
)

func init() {
	di.Graph.Invoke(func(gs GameStore, cp channels.ChannelProvider) {

		cp.SetUserChecker("games", func(ctx context.Context, channel channels.Channel, userId string) error {

			parts := strings.Split(channel.Id(), ";")

			if userId != parts[1] {
				return fmt.Errorf("this is not your private channel")
			}

			instId := parts[0]
			inst, err := gs.GetInstanceById(instId)
			if err != nil {
				return err
			}

			inter := NewInteractor(ctx)
			if inst.MetaState == WaitingForPlayers && !inst.HasPlayer(userId) {
				return inter.joinInstance(inst, userId)
			}

			if !inst.HasPlayer(userId) {
				return fmt.Errorf("you're not in this game.")
			}

			return nil
		})

		cp.OnLeave("games", func(ctx context.Context, channel channels.Channel, userId, socketId string) {
			inter := NewInteractor(ctx)

			parts := strings.Split(channel.Id(), ";")

			instId := parts[0]
			inter.LeaveGame(instId, userId)
		})
	})
}

type GamesInteractor struct {
	store GameStore
	cp    channels.ChannelProvider
}

func NewInteractor(ctx context.Context) GamesInteractor {
	var inter GamesInteractor
	di.Graph.Invoke(func(sc StoreConstructor, cp channels.ChannelProvider) {
		inter = GamesInteractor{
			sc(ctx),
			cp,
		}
	})

	return inter
}

func (inter GamesInteractor) CreateInstance(gameName, userId string) (instanceResponse, error) {
	g, err := Registry.GetGameLatestVersion(gameName)

	if err != nil {
		return instanceResponse{}, err
	}

	inst := NewInstance(g, userId)
	inst.AddPlayer(userId)
	err = inter.store.SaveInstance(inst)
	return instanceToResponse(inst, userId, inter.cp), err
}

func (inter GamesInteractor) JoinGame(instanceId, userId string) error {
	inst, err := inter.store.GetInstanceById(instanceId)
	if err != nil {
		return err
	}

	return inter.joinInstance(inst, userId)
}

func (inter GamesInteractor) joinInstance(inst *GameInstance, userId string) error {
	if inst.HasPlayer(userId) {
		return fmt.Errorf("%v is already in the game.", userId)
	}

	inst.AddPlayer(userId)

	return inter.store.SaveInstance(inst)
}

func (inter GamesInteractor) LeaveGame(instanceId, userId string) error {
	inst, err := inter.store.GetInstanceById(instanceId)
	if err != nil {
		return err
	}

	inst.RemovePlayer(userId)

	return inter.store.SaveInstance(inst)
}

func (inter GamesInteractor) StartGame(instanceId, userId string) error {
	inst, err := inter.store.GetInstanceById(instanceId)
	if err != nil {
		return err
	}

	if inst.AdminUserId != userId {
		return fmt.Errorf("you're not the admin of this game")
	}

	inst.ShufflePlayers()
	inst.MetaState = InProgress
	inst.Game().InitializeState(&inst.State)

	return inter.store.SaveInstance(inst)
}

func (inter GamesInteractor) GetInstance(instanceId string, userId ...string) (instanceResponse, error) {
	inst, err := inter.store.GetInstanceById(instanceId)
	if err != nil {
		return instanceResponse{}, err
	}

	uid := ""
	if len(userId) == 0 {
		uid = userId[0]
	}

	return instanceToResponse(inst, uid, inter.cp), err
}

func (inter GamesInteractor) MakeMove(instanceId, userId string, moveData interface{}) (instanceResponse, error) {
	inst, err := inter.store.GetInstanceById(instanceId)
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

	move := Move{
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
		inst.MetaState = GameOver
	}

	err = inter.store.SaveInstance(inst)
	if err != nil {
		return instanceToResponse(inst, userId, inter.cp), err
	}

	return instanceToResponse(inst, userId, inter.cp), nil
}

func (inter GamesInteractor) ListInstances(gname string, state ...MetaState) (*[]instanceResponse, error) {
	insts, err := inter.store.GetInstancesByGame(gname, state...)
	if err != nil {
		return nil, err
	}

	ris := make([]instanceResponse, len(*insts))
	for i, inst := range *insts {
		ris[i] = instanceToResponse(&inst, "", inter.cp)
	}

	return &ris, nil
}

type instanceResponsePlayer struct {
	UserId string `json:"user_id"`
	Score  int32  `json:"score"`
}

type instanceResponse struct {
	Id       string                   `json:"id"`
	GameInfo GameInfo                 `json:"game_info"`
	State    MetaState                `json:"state"`
	Players  []instanceResponsePlayer `json:"player_ids"`

	PublicChannel  string `json:"public_channel"`
	PrivateChannel string `json:"private_channel"`
}

func instanceToResponse(i *GameInstance, uid string, cp channels.ChannelProvider) instanceResponse {
	chanIds := i.Channels(uid)
	ps := make([]instanceResponsePlayer, len(i.State.Players))
	for j, p := range i.State.Players {
		ps[j] = instanceResponsePlayer{
			UserId: p.UserId,
			Score:  p.Score,
		}
	}

	return instanceResponse{
		Id:       i.Id,
		GameInfo: i.Game().Info(),
		State:    i.MetaState,
		Players:  ps,

		PublicChannel:  cp.NewChannel(nil, "games", chanIds.Public, true).ClientId(),
		PrivateChannel: cp.NewChannel(nil, "games", chanIds.Private, false).ClientId(),
	}
}
