package deps

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/bigblind/makker/di"
	"github.com/bigblind/makker/games"
	"google.golang.org/appengine/datastore"
	"time"
	"github.com/bigblind/makker/logging"
	"reflect"
	"io"
)

func init() {
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
	CreatedAt	time.Time
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
		CreatedAt:	 i.CreatedAt,
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
	game, _ := games.Registry.GetGame(ent.GameName, ent.GameVersion)
	info := game.Info()
	players := make([]games.PlayerState, len(ent.State.Players))
	for i, p := range ent.State.Players {
		players[i] = games.PlayerState{
			UserId: p.UserId,
			Score:  p.Score,
		}
		gobDecode(&players[i].PrivateState, p.PrivateState, info.PrivateStateType)
		gobDecode(&players[i].PublicState, p.PublicState, info.PublicStateType)
	}

	i := games.GameInstance{
		Id:          key.Encode(),
		CreatedAt:	 ent.CreatedAt,
		GameName:    ent.GameName,
		GameVersion: ent.GameVersion,
		State: games.GameState{
			Players: players,
		},
		AdminUserId: ent.AdminUserId,
		MetaState:   games.MetaState(ent.MetaState),
	}

	//TODO: implement move type decoding
	// gobDecode(&i.Moves, ent.Moves)
	gobDecode(&i.State.SharedState, ent.State.SharedState, info.SharedStateType)
	return &i
}

func (ent *gameInstanceEntity) registerTypes()  {
	game, err := games.Registry.GetGame(ent.GameName, ent.GameVersion)
	if err != nil {
		panic(err)
	}

	info := game.Info()
	if info.PublicStateType != nil {
		gob.Register(info.PublicStateType)
	}

	if info.PrivateStateType != nil {
		gob.Register(info.PrivateStateType)
	}

	if info.SharedStateType != nil {
		gob.Register(info.SharedStateType)
	}
}

type appEngineGameStore struct {
	logger *logging.StructuredLogger
}

func NewGameStore(logger *logging.StructuredLogger) games.GameStore {
	return appEngineGameStore{logger}
}

func (gs appEngineGameStore) SaveInstance(ctx context.Context, instance *games.GameInstance) error {
	ent := entityFromInstance(instance)
	gs.logger.Debugf(ctx, "Created entity: %#v\ninstance:%#v", ent, instance)

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
	gs.logger.Debugf(ctx, "Loaded entity: %#v", ent)

	inst := ent.toInstance(key)
	gs.logger.Debugf(ctx, "decoded entity: %#v", inst)
	return inst, nil
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

func gobDecode(ptr interface{}, data []byte, typ interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	// If the game has no type for a particular state, they don't use it,
	// so e don't need to decode anything
	if typ == nil {
		return nil
	}
	val := reflect.New(reflect.TypeOf(typ))

	err := dec.DecodeValue(val)
	// DecodeValue will return io.EOF when it has no more objects to decode.
	// And since we're only encoding one object per []byte, it'll always return io.EOF, unless
	// something went wrong.
	if err != io.EOF {
		return err
	}

	// get a reflect.Value from our pointer
	// e take its Elem(), because you can only Set() the
	// pointer's value, not the pointer itself. It's similar to how you'd
	// dereference a pointer when changing the underlying value.
	ptrval := reflect.ValueOf(ptr).Elem()
	ptrval.Set(val)
	return nil
}
