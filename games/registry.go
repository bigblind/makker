package games

import "fmt"

type gameVersions map[int]Game

type gameRegistry struct {
	games map[string]gameVersions
	latestVersion map[string]*Game
}

func (gr *gameRegistry) Register(g Game) {
	i := g.Info()
	if _, ok := gr.games[i.Name]; !ok {
		gr.games[i.Name] = make(map[int]Game)
	}

	gr.games[i.Name][i.Version] = g

	if latestGame, ok := gr.latestVersion[i.Name]; !ok || (*latestGame).Info().Version < i.Version {
		gr.latestVersion[i.Name] = &g
	}
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

func (gr *gameRegistry) GetGameLatestVersion(name string) (Game, error) {
	if g, ok := gr.latestVersion[name]; ok {
		return *g, nil
	}

	return nil, fmt.Errorf("No game with name %v found.", name)
}

func newRegistry() *gameRegistry {
	gr := gameRegistry{
			games: make(map[string]gameVersions),
			latestVersion: make(map[string]*Game),
		}
	return &gr
}

var Registry = newRegistry()