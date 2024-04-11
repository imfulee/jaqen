package mapper

type IMapping interface {
	Exist(PlayerID) bool
	MapToImage(PlayerID, FilePath)
	Save() error
	// Write(path string)
}

type IImages interface {
	Random(ethnic Ethnic) FilePath
}

func Map(getPlayers func() ([]Player, error), mapping IMapping, imagePool IImages, preserve bool) error {
	players, err := getPlayers()
	if err != nil {
		return err
	}

	for _, player := range players {
		if preserve && mapping.Exist(player.id) {
			continue
		}

		imgFilename := imagePool.Random(player.ethnic)
		mapping.MapToImage(player.id, imgFilename)
	}

	if err := mapping.Save(); err != nil {
		return err
	}

	return nil
}
