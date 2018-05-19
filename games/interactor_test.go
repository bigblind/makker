package games

import (
	"github.com/stretchr/testify/mock"
	"testing"
	"github.com/stretchr/testify/require"
	"fmt"
)

type MockGamesStore struct {
		mock.Mock
}

func (mgs MockGamesStore) SaveInstance(instance *GameInstance) error {
	args :=mgs.Called(instance)
	return args.Error(0)
}

func (mgs MockGamesStore) GetInstanceById(id string) (*GameInstance, error) {
	args :=mgs.Called(id)
	return args.Get(0).(*GameInstance), args.Error(1)
}

func (mgs MockGamesStore) GetInstancesByGame(gameName string) (*[]GameInstance, error) {
	args :=mgs.Called(gameName)
	return args.Get(0).(*[]GameInstance), args.Error(1)
}

func (mgs MockGamesStore) GetInstancesByGameVersion(game Game) (*[]GameInstance, error) {
	args :=mgs.Called(game)
	return args.Get(0).(*[]GameInstance), args.Error(1)
}

func (mgs MockGamesStore) DeleteGameInstance(instance *GameInstance) error {
	args :=mgs.Called(instance)
	return args.Error(0)
}

func createInteractor() (GamesInteractor, *MockGamesStore) {
	mgs := new(MockGamesStore)
	return GamesInteractor{mgs}, mgs
}

func TestGamesInteractor_CreateInstance(t *testing.T) {
	req := require.New(t)
	int, mgs := createInteractor()
	game := makeGame("myGame", 1)
	game2 := makeGame("myGame", 2)
	Registry.Register(game)
	Registry.Register(game2)

	mgs.On("SaveInstance", mock.AnythingOfType("*games.GameInstance")).Return(nil).Once()

	inst, err := int.CreateInstance("myGame", "userId")

	req.NoError(err)
	req.Equal("myGame", inst.GameName)
	req.Equal(2, inst.GameVersion)
	req.Equal("userId", inst.AdminUserId)
	req.Equal(1, len(inst.State.Players))
	req.Equal("userId", inst.State.Players[0].userId)

	// error cases
	// The game does not exist
	_, err = int.CreateInstance("nonExistentGame", "foo")
	req.Error(err, "Should throw an error when there's no game with the given name")

	// The GameStore returned an error
	mgs.On("SaveInstance", mock.AnythingOfType("*games.GameInstance")).Return(fmt.Errorf("foo"))
	_, err = int.CreateInstance("myGame", "oo")
	mgs.AssertExpectations(t)
	req.Error(err, "Should return an error when the GameStore returns an error")
}

func TestGamesInteractor_JoinGame(t *testing.T) {
	req := require.New(t)
	int, mgs := createInteractor()
	g := makeGame("", 1)
	inst := NewInstance(g, "adminId")
	mgs.On("GetInstanceById", "instanceId").Return(inst, nil).Once()
	mgs.On("SaveInstance", inst).Return(nil).Once()

	err := int.JoinGame("instanceId", "userId")

	mgs.AssertExpectations(t)
	req.NoError(err)
	req.Equal("userId", inst.State.Players[0].userId)

	// Don't allow the same user to join a game twice
	mgs.On("GetInstanceById", "instanceId").Return(inst, nil).Once()

	err = int.JoinGame("instanceId", "userId")
	req.Error(err, "The user should not be able to join the same game twice.")
}
