package games

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"github.com/bigblind/makker/channels"
)

func createInteractor() (GamesInteractor, *MockGamesStore, *channels.MockChannelProvider) {
	mgs := new(MockGamesStore)
	mcp := new(channels.MockChannelProvider)
	return GamesInteractor{mgs, mcp}, mgs, mcp
}

func TestGamesInteractor_CreateInstance(t *testing.T) {
	req := require.New(t)
	int, mgs, mcp := createInteractor()
	game := makeGame("myGame", 1)
	game2 := makeGame("myGame", 2)
	Registry.Register(game)
	Registry.Register(game2)

	pubc := channels.NewMockChannel("games", "publicId", true)
	privc := channels.NewMockChannel("games", "privateId", false)

	mgs.On("SaveInstance", mock.AnythingOfType("*games.GameInstance")).Return(nil).Once()
	mcp.On("NewChannel", nil, "games", "", true).Return(pubc)
	mcp.On("NewChannel", nil, "games", ";userId", false).Return(privc)
	pubc.On("ClientId").Return("pubId")
	privc.On("ClientId").Return("privId")

	inst, err := int.CreateInstance("myGame", "userId")

	req.NoError(err)
	req.Equal("myGame", inst.GameInfo.Name)
	req.Equal(2, inst.GameInfo.Version)
	req.Equal(1, len(inst.Players))
	req.Equal("userId", inst.Players[0].UserId)
	req.Equal("pubId", inst.PublicChannel)
	req.Equal("privId", inst.PrivateChannel)

	// error cases
	// The game does not exist
	_, err = int.CreateInstance("nonExistentGame", "foo")
	req.Error(err, "Should throw an error when there's no game with the given name")

	// The GameStore returned an error
	mgs.On("SaveInstance", mock.AnythingOfType("*games.GameInstance")).Return(fmt.Errorf("foo"))
	_, err = int.CreateInstance("myGame", "foo")
	mgs.AssertExpectations(t)
	req.Error(err, "Should return an error when the GameStore returns an error")
}

func TestGamesInteractor_JoinGame(t *testing.T) {
	req := require.New(t)
	int, mgs, _ := createInteractor()
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
	int, mgs, _ := createInteractor()
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
	int, mgs, mcp := createInteractor()
	g := makeGame("myGame", 1)
	inst := NewInstance(g, "adminId")
	inst.AddPlayer("player1")
	inst.AddPlayer("player2")
	inst.AddPlayer("player3")

	pubc := channels.NewMockChannel("games", "publicId", true)
	privc := channels.NewMockChannel("games", "privateId", false)

	inst.MetaState = InProgress
	mgs.On("GetInstanceById", "instanceId").Return(inst, nil)
	mcp.On("NewChannel", nil, "games", "", true).Return(pubc)
	mcp.On("NewChannel", nil, "games", ";", false).Return(privc)
	pubc.On("ClientId").Return("pubId")
	privc.On("ClientId").Return("privId")

	inst2, err := int.GetInstance("instanceId")
	req.NoError(err)
	req.Equal("myGame", inst2.GameInfo.Name)
	req.Equal(1, inst2.GameInfo.Version)
	req.Equal("pubId", inst2.PublicChannel)
	req.Equal("privId", inst2.PrivateChannel)
}

func TestGamesInteractor_MakeMove(t *testing.T) {
	req := require.New(t)
	int, mgs, mcp := createInteractor()
	g := makeGame("myGame", 1)
	inst := NewInstance(g, "adminId")
	inst.AddPlayer("player1")
	inst.AddPlayer("player2")

	pubc := channels.NewMockChannel("games", "publicId", true)
	privc := channels.NewMockChannel("games", "privateId", false)

	mgs.On("GetInstanceById", "instanceId").Return(inst, nil)
	mgs.On("SaveInstance", inst).Return(nil)
	mcp.On("NewChannel", nil, "games", "", true).Return(pubc)
	mcp.On("NewChannel", nil, "games", ";player2", false).Return(privc)
	pubc.On("ClientId").Return("pubId")
	privc.On("ClientId").Return("privId")

	var move Move
	g.On("HandleUpdate", &inst.State, mock.AnythingOfType("games.Move")).Run(func(args mock.Arguments) {
		move = args.Get(1).(Move)
		req.Equal(int8(1), move.Player)
		req.Equal("MoveData", move.Data)
	}).Return(nil)
	g.On("CanPlayerMove", 1, &inst.State).Return(true)
	g.On("IsGameOver", &inst.State).Return((true))
	Registry.Register(g)

	_, err := int.MakeMove("instanceId", "player2", "MoveData")
	req.NoError(err)
	g.AssertExpectations(t)
	mgs.AssertExpectations(t)
}
