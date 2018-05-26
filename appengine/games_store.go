package appengine

import (
	"github.com/bigblind/makker/games"
	"google.golang.org/appengine/datastore"
	"context"
)

const gameInstanceKind = "gameInstance"

type gameInstanceEntity struct {
	MetaState int8
	GameName string
	GameVersion int
	Moves []games.Move `datastore:",noindex"`
	State interface{} `datastore:",noindex"`
	AdminUserId string
}

func entityFromInstance(i *games.GameInstance) *gameInstanceEntity {
	ent := gameInstanceEntity{
		MetaState: int8(i.MetaState),
		GameName: i.GameName,
		GameVersion: i.GameVersion,
		Moves: i.Moves,
		State: i.State,
		AdminUserId: i.AdminUserId,
	}

	return &ent
}

func (ent *gameInstanceEntity) toInstance(key *datastore.Key) *games.GameInstance {
	i := games.GameInstance{
		Id: key.Encode(),
		GameName: ent.GameName,
		GameVersion: ent.GameVersion,
		Moves: ent.Moves,
		State: ent.State.(games.GameState),
		AdminUserId: ent.AdminUserId,
		MetaState: games.MetaState(ent.MetaState),
	}

	return &i
}

type AppEngineGameStore struct {
	ctx context.Context
}

func NewGameStore(ctx context.Context) games.GameStore {
	return AppEngineGameStore{ctx}
}

func (gs AppEngineGameStore) SaveInstance(instance *games.GameInstance) error {
	ent := entityFromInstance(instance)

	var key *datastore.Key
	var err error
	if instance.Id != "" {
		key, err = datastore.DecodeKey(instance.Id)
		if err != nil {
			return err
		}
	} else {
		key = datastore.NewIncompleteKey(gs.ctx, gameInstanceKind, nil)
	}

	key, err = datastore.Put(gs.ctx, key, ent)
	if err != nil {
		return err
	}

	instance.Id = key.Encode()
	return nil
}

func (gs AppEngineGameStore) GetInstanceById(id string) (*games.GameInstance, error) {
	key, err := datastore.DecodeKey(id)
	if err != nil {
		return nil, err
	}

	var ent gameInstanceEntity
	err = datastore.Get(gs.ctx, key, &ent)
	if err != nil {
		return nil, err
	}

	return ent.toInstance(key), nil
}

func (gs AppEngineGameStore) GetInstancesByGame(gameName string) (*[]games.GameInstance, error) {
	q := datastore.NewQuery(gameInstanceKind)
	q = q.Filter("GameName =", gameName)
	q = q.Filter("MetaState =", games.WaitingForPlayers)

	var res []gameInstanceEntity
	keys, err := q.GetAll(gs.ctx, &res)
	if err != nil {
		return nil, err
	}

	insts := make([]games.GameInstance, len(res))
	for i := range insts {
		insts[i] = *res[i].toInstance(keys[i])
	}

	return &insts, nil
}

func (gs AppEngineGameStore) GetInstancesByGameVersion(game games.Game) (*[]games.GameInstance, error) {
	q := datastore.NewQuery(gameInstanceKind)
	inf := game.Info()
	q = q.Filter("GameName =", inf.Name)
	q = q.Filter("GameVersion =", inf.Version)
	q = q.Filter("MetaState =", games.WaitingForPlayers)

	var res []gameInstanceEntity
	keys, err := q.GetAll(gs.ctx, &res)
	if err != nil {
		return nil, err
	}

	insts := make([]games.GameInstance, len(res))
	for i := range insts {
		insts[i] = *res[i].toInstance(keys[i])
	}

	return &insts, nil
}

func (gs AppEngineGameStore) DeleteGameInstance(instance *games.GameInstance) error {
	key, err := datastore.DecodeKey(instance.Id)
	if err != nil {
		return err
	}

	return datastore.Delete(gs.ctx, key)
}

