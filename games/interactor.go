package games

import (
	"fmt"
	"time"
	"context"
	"github.com/bigblind/makker/di"
	"github.com/bigblind/makker/channels"
	"strings"
	"go/ast"
)

func init()  {
	di.Graph.Invoke(func(cp channels.ChannelProvider) {

		cp.SetUserChecker("games", func(ctx context.Context, channel channels.Channel, userId string) error {
			parts := strings.Split(channel.Id(), ";")

			if userId != parts[1] {
				return fmt.Errorf("this is not your private channel")
			}

			instId := parts[0]
			inter := NewInteractor(ctx)
			inst, err := inter.GetInstance(instId)
			if err != nil {
				return err
			}

			if inst.MetaState == WaitingForPlayers && !inst.HasPlayer(userId) {
				return inter.JoinGame(inst.Id, userId)
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
}

func NewInteractor(ctx context.Context) GamesInteractor {
	var inter GamesInteractor
	di.Graph.Invoke(func(sc StoreConstructor) {
		inter = GamesInteractor{
			sc(ctx),
		}
	})

	return inter
}

func (inter GamesInteractor) CreateInstance(gameName, userId string) (GameInstance, error) {
	g, err := Registry.GetGameLatestVersion(gameName)

	if err != nil {
		return GameInstance{}, err
	}

	inst := NewInstance(g, userId)
	inst.AddPlayer(userId)
	err = inter.store.SaveInstance(inst)
	return *inst, err
}

func (inter GamesInteractor) JoinGame(instanceId, userId string) error {
	inst, err := inter.store.GetInstanceById(instanceId)
	if err != nil {
		return err
	}

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

func (inter GamesInteractor) GetInstance(instanceId string) (GameInstance, error) {
	inst, err := inter.store.GetInstanceById(instanceId)
	if err != nil {
		return GameInstance{}, err
	}

	// If the game isn't in progress, we can return it as-is
	if inst.MetaState != InProgress {
		return *inst, nil
	}

	return *inst, err
}

func (inter GamesInteractor) MakeMove(instanceId, userId string, moveData interface{}) (GameInstance, error) {
	inst, err := inter.store.GetInstanceById(instanceId)
	if err != nil {
		return *inst, err
	}

	idx := inst.GetPlayerIndex(userId)
	if idx < 0 {
		return *inst, fmt.Errorf("you're not in this game.")
	}

	game := inst.Game()

	if !game.CanPlayerMove(int(idx), &inst.State) {
		return *inst, fmt.Errorf("you can't make a move right now")
	}

	move := Move{
		Data: moveData,
		Player: int8(idx),
		Time: time.Now(),
	}

	err = game.HandleUpdate(&inst.State, move)
	if err != nil {
		return *inst, err
	}

	inst.Moves = append(inst.Moves, move)

	if game.IsGameOver(&inst.State) {
		inst.MetaState = GameOver
	}

	err = inter.store.SaveInstance(inst)
	if err != nil {
		return *inst, err
	}

	return *inst, nil
}

func (inter GamesInteractor) ListInstances(gname string, state... MetaState) (*[]GameInstance, error) {
	return inter.store.GetInstancesByGame(gname, state...)
}

func transformInstanceForPlayer(inst GameInstance, userId string) (GameInstance, error) {
	idx := int(inst.GetPlayerIndex(userId))
	// If the user isn't a player in the game,
	// act as if they're the user at position 0, keeping the list as it is.
	if idx < 0 {
		idx = 0
	}

	// Rotate the list of players so that the current user is first.
	// Also remove private state from other players
	ps := inst.State.Players
	n := len(ps)
	inst.State.Players = make([]PlayerState, n)
	for i, p := range ps {
		if p.UserId != userId {
			p.PrivateState = nil
		}
		inst.State.Players[(n + i - idx) % n] = p
	}
	return inst, nil
}