package games

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockGamesStore struct {
	mock.Mock
}

func (mgs MockGamesStore) SaveInstance(ctx context.Context, instance *GameInstance) error {
	args := mgs.Called(ctx, instance)
	return args.Error(0)
}

func (mgs MockGamesStore) GetInstanceById(ctx context.Context, id string) (*GameInstance, error) {
	args := mgs.Called(ctx, id)
	return args.Get(0).(*GameInstance), args.Error(1)
}

func (mgs MockGamesStore) GetInstancesByGame(ctx context.Context, gameName string, state ...MetaState) (*[]GameInstance, error) {
	args := mgs.Called(ctx, gameName)
	return args.Get(0).(*[]GameInstance), args.Error(1)
}

func (mgs MockGamesStore) GetInstancesByGameVersion(ctx context.Context, game Game, state ...MetaState) (*[]GameInstance, error) {
	args := mgs.Called(ctx, game)
	return args.Get(0).(*[]GameInstance), args.Error(1)
}

func (mgs MockGamesStore) DeleteGameInstance(ctx context.Context, instance *GameInstance) error {
	args := mgs.Called(ctx, instance)
	return args.Error(0)
}
