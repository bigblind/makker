package deps

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/bigblind/makker/games"
	"google.golang.org/appengine/datastore"
	"github.com/bigblind/makker/di"
)

func init()  {
	di.Graph.Provide(NewGameStore)
}

const gameInstanceKind = "gameInstance"

type PlayerEntity struct {
	UserId       string
	Score        int32  `datastore:",noindex"`
	PrivateState []byte `datastore:",noindex"`
	PublicState  []byte `datastore:",noindex"`
}

type GameStateEntity struct {
	Players     []PlayerEntity
	SharedState []byte
}

type gameInstanceEntity struct {
	MetaState   int8
	GameName    string
	GameVersion int
	Moves       []byte `datastore:",noindex"`
	State       GameStateEntity
	AdminUserId string
}

func entityFromInstance(i *games.GameInstance) *gameInstanceEntity {
	players := make([]PlayerEntity, len(i.State.Players))
	for i, p := range i.State.Players {
		players[i] = PlayerEntity{
			UserId:       p.UserId,
			Score:        p.Score,
			PrivateState: gobEncode(p.PrivateState),
			PublicState:  gobEncode(p.PublicState),
		}
	}
	ent := gameInstanceEntity{
		MetaState:   int8(i.MetaState),
		GameName:    i.GameName,
		GameVersion: i.GameVersion,
		Moves:       gobEncode(i.Moves),
		State: GameStateEntity{
			Players:     players,
			SharedState: gobEncode(i.State.SharedState),
		},
		AdminUserId: i.AdminUserId,
	}

	return &ent
}

func (ent *gameInstanceEntity) toInstance(key *datastore.Key) *games.GameInstance {
	players := make([]games.PlayerState, len(ent.State.Players))
	for i, p := range ent.State.Players {
		players[i] = games.PlayerState{
			UserId: p.UserId,
			Score:  p.Score,
		}
		gobDecode(&players[i].PrivateState, p.PrivateState)
		gobDecode(&players[i].PublicState, p.PublicState)
	}
	i := games.GameInstance{
		Id:          key.Encode(),
		GameName:    ent.GameName,
		GameVersion: ent.GameVersion,
		State: games.GameState{
			Players: players,
		},
		AdminUserId: ent.AdminUserId,
		MetaState:   games.MetaState(ent.MetaState),
	}
	gobDecode(&i.Moves, ent.Moves)
	gobDecode(&i.State.SharedState, ent.State.SharedState)
	return &i
}

type appEngineGameStore struct {
}

func NewGameStore() games.GameStore {
	return appEngineGameStore{}
}

func (gs appEngineGameStore) SaveInstance(ctx context.Context, instance *games.GameInstance) error {
	ent := entityFromInstance(instance)

	var key *datastore.Key
	var err error
	if instance.Id != "" {
		key, err = datastore.DecodeKey(instance.Id)
		if err != nil {
			return err
		}
	} else {
		key = datastore.NewIncompleteKey(ctx, gameInstanceKind, nil)
	}

	key, err = datastore.Put(ctx, key, ent)
	if err != nil {
		return err
	}

	instance.Id = key.Encode()
	return nil
}

func (gs appEngineGameStore) GetInstanceById(ctx context.Context, id string) (*games.GameInstance, error) {
	key, err := datastore.DecodeKey(id)
	if err != nil {
		return nil, err
	}

	var ent gameInstanceEntity
	err = datastore.Get(ctx, key, &ent)
	if err != nil {
		return nil, err
	}

	return ent.toInstance(key), nil
}

func (gs appEngineGameStore) GetInstancesByGame(ctx context.Context, gameName string, state ...games.MetaState) (*[]games.GameInstance, error) {
	q := datastore.NewQuery(gameInstanceKind)
	q = q.Filter("GameName =", gameName)
	if len(state) != 0 {
		q = q.Filter("MetaState =", state[0])
	}

	var res []gameInstanceEntity
	keys, err := q.GetAll(ctx, &res)
	if err != nil {
		return nil, err
	}

	insts := make([]games.GameInstance, len(res))
	for i := range insts {
		insts[i] = *res[i].toInstance(keys[i])
	}

	return &insts, nil
}

func (gs appEngineGameStore) GetInstancesByGameVersion(ctx context.Context, game games.Game, state ...games.MetaState) (*[]games.GameInstance, error) {
	q := datastore.NewQuery(gameInstanceKind)
	inf := game.Info()
	q = q.Filter("GameName =", inf.Name)
	q = q.Filter("GameVersion =", inf.Version)
	if len(state) != 0 {
		q = q.Filter("MetaState =", state[0])
	}

	var res []gameInstanceEntity
	keys, err := q.GetAll(ctx, &res)
	if err != nil {
		return nil, err
	}

	insts := make([]games.GameInstance, len(res))
	for i := range insts {
		insts[i] = *res[i].toInstance(keys[i])
	}

	return &insts, nil
}

func (gs appEngineGameStore) DeleteGameInstance(ctx context.Context, instance *games.GameInstance) error {
	key, err := datastore.DecodeKey(instance.Id)
	if err != nil {
		return err
	}

	return datastore.Delete(ctx, key)
}

func gobEncode(v interface{}) []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	enc := gob.NewEncoder(buf)
	enc.Encode(v)
	return buf.Bytes()
}

func gobDecode(ptr interface{}, data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(ptr)
}
