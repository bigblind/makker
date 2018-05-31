package games

import (
	"math/rand"
	"time"
	"context"
	"fmt"
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
	UserId string

	PrivateState interface{}
	PublicState  interface{}
	Score        int32
}

type MetaState int8

const (
	WaitingForPlayers MetaState = iota
	InProgress
	GameOver
)

type GameState struct {
	Players     []PlayerState
	SharedState interface{}
}

type Move struct {
	Player int8
	Data   interface{}
	Time   time.Time
}

type Game interface {
	Info() GameInfo

	InitializeState(state *GameState)
	HandleUpdate(g *GameState, m Move) (error)
	CanPlayerMove(playerIndex int, g *GameState) bool
	IsGameOver(g *GameState) bool
}

type GameInstance struct {
	Id          string
	GameName    string
	GameVersion int
	State       GameState
	Moves       []Move
	MetaState   MetaState
	AdminUserId string
}

func NewInstance(g Game, adminUserId string) *GameInstance {
	info := g.Info()
	instance := GameInstance{
		GameName:    info.Name,
		GameVersion: info.Version,
		AdminUserId: adminUserId,

		Moves:     make([]Move, 2),
		MetaState: WaitingForPlayers,
		State: GameState{
			Players: make([]PlayerState, 0, info.MaxPlayers),
		},
	}

	return &instance
}

func (i *GameInstance) Game() Game {
	g, err := Registry.GetGame(i.GameName, i.GameVersion)
	if err != nil {
		panic(err)
	}
	return g
}

func (gs *GameState) AddPlayer(userId string) {
	gs.Players = append(gs.Players, PlayerState{
		UserId: userId,
	})
}

func (i *GameInstance) AddPlayer(userId string) {i.State.AddPlayer(userId)}

func (gs *GameState) HasPlayer(userId string) bool {
	return gs.GetPlayerIndex(userId) >= 0
}

func (i *GameInstance) HasPlayer(userId string) bool {return i.State.HasPlayer(userId)}

func (gs *GameState) GetPlayerIndex(userId string) int16 {
	for i, p := range gs.Players {
		if p.UserId == userId {
			return int16(i)
		}
	}

	return -1
}

func (i *GameInstance) GetPlayerIndex(userId string) int16 {return i.State.GetPlayerIndex(userId)}

// ShufflePlayers randomly reorders the players, so they're not playing in the order they joined
func (i *GameInstance) ShufflePlayers() {
	p := i.State.Players
	n := len(p)
	for j := 0; j < n; j++ {
		k := j + rand.Intn(n-j)
		p[j], p[k] = p[k], p[j]
	}
}

type InstanceChannels struct {
	Public, Private string
}

func (i *GameInstance) Channels(userId string) InstanceChannels {
	return InstanceChannels{
		Public: i.Id,
		Private: fmt.Sprintf("%v-%v", i.Id, userId),
	}
}

type GameStore interface {
	SaveInstance(instance *GameInstance) error
	GetInstanceById(id string) (*GameInstance, error)
	GetInstancesByGame(gameName string) (*[]GameInstance, error)
	GetInstancesByGameVersion(game Game) (*[]GameInstance, error)
	DeleteGameInstance(instance *GameInstance) error
}

type StoreConstructor func(ctx context.Context) GameStore
