package mapper

type Ethnic string

type PlayerID string

type Player struct {
	id     PlayerID
	ethnic Ethnic
}

type FilePath string
