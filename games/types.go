package games

import "time"

type GameInfo struct {
	// The name of the game
	Name string

	// The minimum number of players
	MinPlayers int

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

type GameState struct {
	Players []PlayerState
	SharedState interface{}
}


type Game interface {
	Info() GameInfo

	GetInitialStat(players []PlayerState)
	HandleUpdate(g GameState) (GameState, error)
	CanPlayerMove(playerIndex int, g GameState) bool
	IsGameOver(g GameState)
}
