package games

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