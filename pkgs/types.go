package mapper

type FilePath string

type Ethnic string

func IsValidEthnic(ethnic string) bool {
	validEthnic := EthnicSet.Contains(Ethnic(ethnic))
	return validEthnic
}

type PlayerID string

type Player struct {
	ID     PlayerID
	Ethnic Ethnic
}
