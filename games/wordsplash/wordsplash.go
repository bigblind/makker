package wordsplash

import (
	"github.com/bigblind/makker/games"
	"time"
	"fmt"
)

type wordSplash struct {}

type gameState struct {
	Round int `json:"round"`
	Letters string `json:"letters"`
	Stage string `json:"stage"`
	Ready []bool `json:"ready"`
	RoundStarted time.Time `json:"round_started"`
	Submissions []string `json:"submissions"`
	SubmissionExists []bool `json:"submission_exists"`
}

type privateState struct {
	Submission string
	Exists bool
}

func (wordSplash) Info() games.GameInfo {
	return games.GameInfo{
		Name: "WordSplash",
		Version: 1,
		MinPlayers: 2,
		MaxPlayers: 4,
		LoseOnDisconnect: true,
	}
}

func (wordSplash) InitializeState(state *games.GameState) {
	state.SharedState = gameState{
		Round: 0,
		Letters: "",
		Ready: make([]bool, len(state.Players)),
		Submissions: make([]string, len(state.Players)),
		SubmissionExists: make([]bool, len(state.Players)),
		Stage: "picking",
	}
}

func (wordSplash) HandleUpdate(g *games.GameState, m games.Move) (error) {
	state := (g.SharedState.(gameState))
	np := len(g.Players)

	action, ok := m.Data.(string)
	if !ok {
		return fmt.Errorf("move isn't a string")
	}

	if state.Stage == "picking" {
		picker := state.Round % np
		if m.Player != uint8(picker) {
			return fmt.Errorf("you're not the picker")
		}

		state.Letters = addLetter(state.Letters, action)

		if len(state.Letters) == 9 {
			state.Stage = "game"
			state.RoundStarted = time.Now()
			state.Ready = make([]bool, len(g.Players))
		}
	}

	if state.Stage == "game" {
		if !state.Ready[m.Player] {
			exists := wordExists(action)
			g.Players[m.Player].PrivateState = privateState{
				Submission: action,
				Exists: exists,
			}
			state.Ready[m.Player] = true
		}

		if all(state.Ready) {
			state.Stage = "result"
			for i, p := range g.Players {
				state.Submissions[i] = p.PrivateState.(privateState).Submission
				state.SubmissionExists[i] = p.PrivateState.(privateState).Exists
				state.Ready[i] = false
			}
		}
	}

	if state.Stage == "result" {
		state.Ready[m.Player] = true
		if all(state.Ready) {
			state.Stage = "picking"
			state.Round += 1
		}
	}

	g.SharedState = state

	return nil
}

func all(bs []bool) bool {
	for _, b := range bs {
		if !b {
			return false
		}
	}

	return true
}

func (wordSplash) CanPlayerMove(playerIndex int, g *games.GameState) bool {
	if g.SharedState.(gameState).Stage == "picking" {
		return playerIndex == g.SharedState.(gameState).Round % len(g.Players)
	}
	return !g.SharedState.(gameState).Ready[playerIndex]
}

func (wordSplash) IsGameOver(g *games.GameState) bool {
	return g.SharedState.(gameState).Round == len(g.Players) * 2
}

func init() {
	games.Registry.Register(wordSplash{})
}


