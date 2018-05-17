package storage

import "github.com/bigblind/makker/games"

type GameStore interface{
	SaveInstance(instance games.GameInstance) (games.GameInstance, error)
	GetInstanceById(id string) (games.GameInstance, error)
	GetInstancesByGame(gameName string) ([]games.GameInstance, error)
	GetInstancesByGameVersion(game games.Game) ([]games.GameInstance, error)
	DeleteGameInstance(instance games.GameInstance) error
}
