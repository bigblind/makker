package games

import "fmt"

type GamesInteractor struct {
	store GameStore
}

func (inter GamesInteractor) CreateInstance(gameName, userId string) (GameInstance, error) {
	g, err := Registry.GetGameLatestVersion(gameName)

	if err != nil {
		return GameInstance{}, err
	}

	inst := NewInstance(g, userId)
	inst.AddPlayer(userId)
	err = inter.store.SaveInstance(inst)
	return *inst, err
}

func (inter GamesInteractor) JoinGame(instanceId, userId string) error {
	inst, err := inter.store.GetInstanceById(instanceId)
	if err != nil {
		return err
	}

	if inst.HasPlayer(userId) {
		return fmt.Errorf("%v is already in the game.", userId)
	}

	inst.AddPlayer(userId)

	return inter.store.SaveInstance(inst)
}

func (inter GamesInteractor) StartGame(instanceId, userId string) error {
	inst, err := inter.store.GetInstanceById(instanceId)
	if err != nil {
		return err
	}

	if inst.AdminUserId != userId {
		return fmt.Errorf("you're not the admin of this game")
	}

	inst.ShufflePlayers()
	inst.MetaState = InProgress

	return inter.store.SaveInstance(inst)
}