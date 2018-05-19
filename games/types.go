package games

import (
	"time"
)

type GameInfo struct {
	// The name of the game
	Name string

	// The minimum number of players
	MinPlayers int

	// The version of this game
	Version int

	// The maximum number of players
	MaxPlayers int

	// After not receiving a state update from a player for this duration, they're considered "gone", and lose the game
	PlayerTimeout time.Duration

	// Whether the player automatically loses the game when they disconnect
	LoseOnDisconnect bool
}

type PlayerState struct {
	userId string

	PrivateState interface{}
	PublicState interface{}
	Score int32
}

type MetaState uint8

const (
	WaitingForPlayers MetaState = iota
	InProgress
	GameOver
)

type GameState struct {
	Players []PlayerState
	SharedState interface{}
}

type Move struct {
	Player uint8
	Data interface{}
	Time time.Time
}

type Game interface {
	Info() GameInfo

	GetInitialStat(players []PlayerState)
	HandleUpdate(g GameState, m Move) (GameState, error)
	CanPlayerMove(playerIndex int, g GameState) bool
	IsGameOver(g GameState)
}

type GameInstance struct {
	Id string
	GameName string
	GameVersion int
	States GameState
	Moves []Move
	MetaState MetaState
	AdminUserId string
}

type GameStore interface{
	SaveInstance(instance GameInstance) (GameInstance, error)
	GetInstanceById(id string) (GameInstance, error)
	GetInstancesByGame(gameName string) ([]GameInstance, error)
	GetInstancesByGameVersion(game Game) ([]GameInstance, error)
	DeleteGameInstance(instance GameInstance) error
}

