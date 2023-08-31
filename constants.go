package main

const (
	African                     string = "African"
	Asian                       string = "Asian"
	Caucasion                   string = "Caucasion"
	CentralEuropean             string = "CentralEuropean"
	EasternEuropeanCentralAsian string = "EECA"
	ItalianMediterranean        string = "ItaMed"
	MiddleEastNorthAfrican      string = "MENA"
	MiddleEastSouthAmerican     string = "MESA"
	SouthAmericanMediterranean  string = "SAMed"
	Scandinavian                string = "Scandinavian"
	SouthEastAsian              string = "Seasian"
	SouthAmerican               string = "SouthAmerican"
	SpanishMediterranean        string = "SpanMed"
	YugoslavGreek               string = "YugoGreek"
)

type iModes struct {
	Overwrite string
	Preserve  string
}

var Modes = iModes{
	Overwrite: "Overwrite",
	Preserve:  "Preserve",
}
