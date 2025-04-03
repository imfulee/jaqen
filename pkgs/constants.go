package mapper

import (
	mapset "github.com/deckarep/golang-set/v2"
)

const (
	African                     Ethnic = "African"
	Asian                       Ethnic = "Asian"
	Caucasian                   Ethnic = "Caucasian"
	CentralEuropean             Ethnic = "Central European"
	EasternEuropeanCentralAsian Ethnic = "EECA"
	ItalianMediterranean        Ethnic = "Italmed"
	MiddleEastNorthAfrican      Ethnic = "MENA"
	MiddleEastSouthAsian        Ethnic = "MESA"
	SouthAmericanMediterranean  Ethnic = "SAMed"
	Scandinavian                Ethnic = "Scandinavian"
	SouthEastAsian              Ethnic = "Seasian"
	SouthAmerican               Ethnic = "South American"
	SpanishMediterranean        Ethnic = "SpanMed"
	YugoslavGreek               Ethnic = "YugoGreek"
)

var Ethnicities = [...]Ethnic{
	African,
	Asian,
	Caucasian,
	CentralEuropean,
	EasternEuropeanCentralAsian,
	ItalianMediterranean,
	MiddleEastNorthAfrican,
	MiddleEastSouthAsian,
	SouthAmericanMediterranean,
	Scandinavian,
	SouthEastAsian,
	SouthAmerican,
	SpanishMediterranean,
	YugoslavGreek,
}

var EthnicSet = mapset.NewSet[Ethnic]()

var NationEthnicMapping = map[string]Ethnic{
	"AFG": MiddleEastSouthAsian,
	"AIA": African,
	"ALB": YugoslavGreek,
	"ALG": MiddleEastNorthAfrican,
	"AND": SpanishMediterranean,
	"ANG": African,
	"ARG": SouthAmericanMediterranean,
	"ARM": EasternEuropeanCentralAsian,
	"ARU": African,
	"ASA": African,
	"ATG": African,
	"AUS": CentralEuropean,
	"AUT": CentralEuropean,
	"AXL": Scandinavian,
	"AZE": EasternEuropeanCentralAsian,
	"BAH": African,
	"BAN": MiddleEastSouthAsian,
	"BAS": SpanishMediterranean,
	"BDI": African,
	"BEL": CentralEuropean,
	"BEN": African,
	"BER": African,
	"BFA": African,
	"BHR": MiddleEastSouthAsian,
	"BHU": Asian,
	"BIH": YugoslavGreek,
	"BLM": Caucasian,
	"BLR": EasternEuropeanCentralAsian,
	"BLZ": SouthAmerican,
	"BOE": African,
	"BOL": SouthAmerican,
	"BOT": African,
	"BRA": SouthAmerican,
	"BRB": African,
	"BRU": SouthEastAsian,
	"BUL": EasternEuropeanCentralAsian,
	"CAM": SouthEastAsian,
	"CAN": Caucasian,
	"CAY": African,
	"CGO": African,
	"CHA": African,
	"CHI": SouthAmerican,
	"CHN": Asian,
	"CIV": African,
	"CMR": African,
	"COD": African,
	"COK": African,
	"COL": SouthAmerican,
	"COM": African,
	"CPV": African,
	"CRC": SouthAmerican,
	"CRO": YugoslavGreek,
	"CTA": African,
	"CUB": SouthAmerican,
	"CUW": African,
	"CYP": MiddleEastNorthAfrican,
	"CZE": EasternEuropeanCentralAsian,
	"DEN": Scandinavian,
	"DJI": African,
	"DMA": African,
	"DOM": SouthAmerican,
	"ECU": SouthAmerican,
	"EGY": MiddleEastNorthAfrican,
	"ENG": Caucasian,
	"EQG": African,
	"ERI": African,
	"ESP": SpanishMediterranean,
	"EST": EasternEuropeanCentralAsian,
	"ESW": African,
	"ETH": African,
	"FIJ": African,
	"FIN": Scandinavian,
	"FRA": CentralEuropean,
	"FRO": Scandinavian,
	"FSM": African,
	"GAB": African,
	"GAM": African,
	"GBR": Caucasian,
	"GEO": EasternEuropeanCentralAsian,
	"GER": CentralEuropean,
	"GHA": African,
	"GIB": Caucasian,
	"GLP": African,
	"GNB": African,
	"GRE": YugoslavGreek,
	"GRL": Caucasian,
	"GRN": African,
	"GUA": SouthAmerican,
	"GUF": African,
	"GUI": African,
	"GUM": African,
	"GUY": African,
	"HAI": African,
	"HKG": Asian,
	"HON": SouthAmerican,
	"HUN": CentralEuropean,
	"IDN": SouthEastAsian,
	"IND": MiddleEastSouthAsian,
	"IRL": Caucasian,
	"IRN": MiddleEastSouthAsian,
	"IRQ": MiddleEastSouthAsian,
	"ISL": Scandinavian,
	"ISR": MiddleEastNorthAfrican,
	"ITA": ItalianMediterranean,
	"JAM": African,
	"JOR": MiddleEastSouthAsian,
	"JPN": Asian,
	"KAZ": EasternEuropeanCentralAsian,
	"KEN": African,
	"KGZ": EasternEuropeanCentralAsian,
	"KIR": African,
	"KOR": Asian,
	"KOS": YugoslavGreek,
	"KSA": MiddleEastSouthAsian,
	"KUW": MiddleEastSouthAsian,
	"KVX": YugoslavGreek,
	"LAO": SouthEastAsian,
	"LBN": MiddleEastNorthAfrican,
	"LBR": African,
	"LBY": African,
	"LCA": African,
	"LES": African,
	"LIB": MiddleEastNorthAfrican,
	"LIE": CentralEuropean,
	"LTU": EasternEuropeanCentralAsian,
	"LUX": CentralEuropean,
	"LVA": EasternEuropeanCentralAsian,
	"MAC": Asian,
	"MAD": African,
	"MAR": MiddleEastNorthAfrican,
	"MAS": SouthEastAsian,
	"MAY": African,
	"MDA": EasternEuropeanCentralAsian,
	"MDV": African,
	"MEX": SouthAmerican,
	"MGL": Asian,
	"MKD": EasternEuropeanCentralAsian,
	"MLI": African,
	"MLT": ItalianMediterranean,
	"MNE": YugoslavGreek,
	"MNG": Asian,
	"MON": ItalianMediterranean,
	"MOZ": African,
	"MRI": African,
	"MSR": African,
	"MTN": African,
	"MTQ": African,
	"MWI": African,
	"MYA": SouthEastAsian,
	"NAM": African,
	"NCA": SouthAmerican,
	"NCL": African,
	"NED": CentralEuropean,
	"NEP": MiddleEastSouthAsian,
	"NGA": African,
	"NIG": African,
	"NIR": Caucasian,
	"NIU": African,
	"NMI": African,
	"NOR": Scandinavian,
	"NZL": Caucasian,
	"OMA": MiddleEastSouthAsian,
	"PAK": MiddleEastSouthAsian,
	"PAN": SouthAmerican,
	"PAR": SouthAmerican,
	"PER": SouthAmerican,
	"PHI": SouthEastAsian,
	"PLE": MiddleEastSouthAsian,
	"PLW": African,
	"PNG": African,
	"POL": CentralEuropean,
	"POR": SpanishMediterranean,
	"PRK": Asian,
	"PUR": SouthAmerican,
	"QAT": MiddleEastSouthAsian,
	"REU": African,
	"ROU": EasternEuropeanCentralAsian,
	"RSA": African,
	"RUS": EasternEuropeanCentralAsian,
	"RWA": African,
	"SAM": African,
	"SCO": Caucasian,
	"SDN": African,
	"SEN": African,
	"SEY": African,
	"SGP": Asian,
	"SIN": SouthEastAsian,
	"SKN": African,
	"SLE": African,
	"SLV": SouthAmerican,
	"SMA": African,
	"SMN": African,
	"SMR": ItalianMediterranean,
	"SOL": African,
	"SOM": African,
	"SPM": Caucasian,
	"SRB": YugoslavGreek,
	"SRI": African,
	"SSD": African,
	"STP": African,
	"SUD": MiddleEastNorthAfrican,
	"SUI": CentralEuropean,
	"SUR": African,
	"SVK": CentralEuropean,
	"SVN": YugoslavGreek,
	"SWE": Scandinavian,
	"SWZ": African,
	"SYR": MiddleEastSouthAsian,
	"TAH": African,
	"TAN": African,
	"TCA": African,
	"TGA": African,
	"THA": SouthEastAsian,
	"TJK": EasternEuropeanCentralAsian,
	"TKM": EasternEuropeanCentralAsian,
	"TLS": African,
	"TOG": African,
	"TPE": Asian,
	"TRI": African,
	"TUN": MiddleEastNorthAfrican,
	"TUR": MiddleEastNorthAfrican,
	"TUV": African,
	"UAE": MiddleEastSouthAsian,
	"UGA": African,
	"UKR": EasternEuropeanCentralAsian,
	"URU": SouthAmericanMediterranean,
	"USA": Caucasian,
	"UZB": EasternEuropeanCentralAsian,
	"VAN": African,
	"VAT": Caucasian,
	"VEN": SouthAmerican,
	"VGB": African,
	"VIE": SouthEastAsian,
	"VIN": African,
	"VIR": African,
	"WAL": Caucasian,
	"WFI": African,
	"YEM": MiddleEastSouthAsian,
	"ZAM": African,
	"ZAN": African,
	"ZIM": African,
}

const (
	FMVersion2020 = "2020"
	FMVersion2021 = "2021"
	FMVersion2022 = "2022"
	FMVersion2023 = "2023"
	FMVersion2024 = "2024"
)

var FMVersions = []string{
	FMVersion2020,
	FMVersion2021,
	FMVersion2022,
	FMVersion2023,
	FMVersion2024,
}

func init() {
	for _, ethnic := range Ethnicities {
		EthnicSet.Add(ethnic)
	}
}
