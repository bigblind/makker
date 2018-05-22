package games

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

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

func (mgs MockGamesStore) GetInstancesByGame(gameName string) (*[]GameInstance, error) {
	args := mgs.Called(gameName)
	return args.Get(0).(*[]GameInstance), args.Error(1)
}

func (mgs MockGamesStore) GetInstancesByGameVersion(game Game) (*[]GameInstance, error) {
	args := mgs.Called(game)
	return args.Get(0).(*[]GameInstance), args.Error(1)
}

func (mgs MockGamesStore) DeleteGameInstance(instance *GameInstance) error {
	args := mgs.Called(instance)
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

	inst, err := int.CreateInstance("myGame", "UserId")

	req.NoError(err)
	req.Equal("myGame", inst.GameName)
	req.Equal(2, inst.GameVersion)
	req.Equal("UserId", inst.AdminUserId)
	req.Equal(1, len(inst.State.Players))
	req.Equal("UserId", inst.State.Players[0].UserId)

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
	g := makeGame("myGame", 1)
	Registry.Register(g)
	inst := NewInstance(g, "adminId")
	mgs.On("GetInstanceById", "instanceId").Return(inst, nil).Once()
	mgs.On("SaveInstance", inst).Return(nil).Once()

	err := int.JoinGame("instanceId", "UserId")

	mgs.AssertExpectations(t)
	req.NoError(err)
	req.Equal("UserId", inst.State.Players[0].UserId)

	// Don't allow the same user to join a game twice
	mgs.On("GetInstanceById", "instanceId").Return(inst, nil).Once()

	err = int.JoinGame("instanceId", "UserId")
	req.Error(err, "The user should not be able to join the same game twice.")
}

func TestGamesInteractor_StartGame(t *testing.T) {
	req := require.New(t)
	int, mgs := createInteractor()
	g := makeGame("myGame", 1)
	inst := NewInstance(g, "UserId")
	mgs.On("GetInstanceById", "instanceId").Return(inst, nil).Once()
	mgs.On("SaveInstance", inst).Return(nil).Once()
	g.On("InitializeState", &inst.State).Return()
	Registry.Register(g)

	err := int.StartGame("instanceId", "UserId")

	req.NoError(err)
	req.Equal(InProgress, inst.MetaState)
	mgs.AssertExpectations(t)
	g.AssertExpectations(t)
}

func TestGamesInteractor_GetInstance(t *testing.T) {
	req := require.New(t)
	int, mgs := createInteractor()
	g := makeGame("myGame", 1)
	inst := NewInstance(g, "adminId")
	inst.AddPlayer("player1")
	inst.AddPlayer("player2")
	inst.AddPlayer("player3")
	for i, _ := range inst.State.Players {
		inst.State.Players[i].PrivateState = "private"
		inst.State.Players[i].PublicState = "public"
	}
	inst.MetaState = InProgress
	fmt.Println(inst)
	mgs.On("GetInstanceById", "instanceId").Return(inst, nil)

	inst2, err := int.GetInstance("instanceId", "player2")
	req.NoError(err)
	req.Equal("myGame", inst2.GameName)
	req.Equal(1, inst2.GameVersion)

	ids := make([]string, 3)
	for i, p := range inst2.State.Players {
		fmt.Println(p)
		ids[i] = p.UserId
		req.Equal("public", p.PublicState)
		if p.UserId == "player2" {
			req.Equal("private", p.PrivateState)
		} else {
			req.Nil(p.PrivateState)
		}

		mgs.AssertExpectations(t)
	}
	req.Equal([]string{"player2", "player3", "player1"}, ids)
}