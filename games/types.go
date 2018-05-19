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
	State GameState
	Moves []Move
	MetaState MetaState
	AdminUserId string
}

func NewInstance(g Game, adminUserId string) *GameInstance {
	info := g.Info()
	instance := GameInstance{
		GameName: info.Name,
		GameVersion: info.Version,
		AdminUserId: adminUserId,

		Moves: make([]Move, 2),
		MetaState: WaitingForPlayers,
		State: GameState{
			Players: make([]PlayerState, 0, info.MaxPlayers),
		},
	}

	return &instance
}

func (i *GameInstance) AddPlayer(userId string) {
	i.State.Players = append(i.State.Players, PlayerState{
		userId: userId,
	})
}

type GameStore interface{
	SaveInstance(instance *GameInstance) (GameInstance, error)
	GetInstanceById(id string) (*GameInstance, error)
	GetInstancesByGame(gameName string) (*[]GameInstance, error)
	GetInstancesByGameVersion(game Game) (*[]GameInstance, error)
	DeleteGameInstance(instance *GameInstance) error
}

