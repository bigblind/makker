package games

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGameRegistry_Register(t *testing.T) {
	gr := newRegistry()
	g1 := NewMockGame("myGame", 1)
	g2 := NewMockGame("mySecondGame", 1)
	g3 := NewMockGame("mySecondGame", 2)

	gr.Register(g1)
	gr.Register(g2)
	gr.Register(g3)
}

func TestGameRegistry_GetGame(t *testing.T) {
	req := require.New(t)
	gr := newRegistry()
	g1 := NewMockGame("myGame", 1)
	g2 := NewMockGame("mySecondGame", 1)
	g3 := NewMockGame("mySecondGame", 2)

	gr.Register(g1)
	gr.Register(g2)
	gr.Register(g3)

	res1, e1 := gr.GetGame("myGame", 1)
	req.NoError(e1)
	req.Equal(g1, res1)

	res2, e2 := gr.GetGame("mySecondGame", 1)
	req.NoError(e2)
	req.Equal(g2, res2)

	res3, e3 := gr.GetGame("mySecondGame", 2)
	req.NoError(e3)
	req.Equal(g3, res3)

	_, e4 := gr.GetGame("gameDoesNotExist", 1)
	req.Error(e4, "The game doesn't exist.")

	_, e5 := gr.GetGame("myGame", 2)
	req.Error(e5, "The game's version does not exist")
}

func TestGameRegistry_GetGameLatestVersion(t *testing.T) {
	req := require.New(t)
	gr := newRegistry()
	g1 := NewMockGame("myGame", 1)
	g2 := NewMockGame("myGame", 2)
	g3 := NewMockGame("mySecondGame", 1)
	g4 := NewMockGame("mySecondGame", 2)

	gr.Register(g1)
	gr.Register(g2)
	// Make sure that the latest version is also returned when games are registered in reverse order
	gr.Register(g4)
	gr.Register(g3)

	res1, e1 := gr.GetGameLatestVersion("myGame")
	req.NoError(e1)
	req.Equal(g2, res1)

	res2, e2 := gr.GetGameLatestVersion("mySecondGame")
	req.NoError(e2)
	req.Equal(g4, res2)

	_, e3 := gr.GetGameLatestVersion("doesNotExist")
	req.Error(e3)
}
