package games

import "fmt"

type gameVersions map[int]Game

type gameRegistry struct {
	games map[string]gameVersions
}

func (gr *gameRegistry) Register(g Game) {
	i := g.Info()
	if _, ok := gr.games[i.Name]; !ok {
		gr.games[i.Name] = make(map[int]Game)
	}

	gr.games[i.Name][i.Version] = g
}

func (gr *gameRegistry) GetGame(name string, version int) (Game, error) {
	if versions, ok := gr.games[name]; ok {
		if game, ok := versions[version]; ok {
			return game, nil
		}

		return nil, fmt.Errorf("game %v has no version %v", name, version)
	}

	return nil, fmt.Errorf("game %v not found", name)
}

func newRegistry() *gameRegistry {
	gr := gameRegistry{games: make(map[string]gameVersions)}
	return &gr
}

var Registry = newRegistry()