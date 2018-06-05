package games

import "github.com/stretchr/testify/mock"

type MockGamesStore struct {
	mock.Mock
}

func (mgs MockGamesStore) SaveInstance(instance *GameInstance) error {
	args := mgs.Called(instance)
	return args.Error(0)
}

func (mgs MockGamesStore) GetInstanceById(id string) (*GameInstance, error) {
	args := mgs.Called(id)
	return args.Get(0).(*GameInstance), args.Error(1)
}

func (mgs MockGamesStore) GetInstancesByGame(gameName string, state... MetaState) (*[]GameInstance, error) {
	args := mgs.Called(gameName)
	return args.Get(0).(*[]GameInstance), args.Error(1)
}

func (mgs MockGamesStore) GetInstancesByGameVersion(game Game, state... MetaState) (*[]GameInstance, error) {
	args := mgs.Called(game)
	return args.Get(0).(*[]GameInstance), args.Error(1)
}

func (mgs MockGamesStore) DeleteGameInstance(instance *GameInstance) error {
	args := mgs.Called(instance)
	return args.Error(0)
}
