package internal

type JaqenConfig struct {
	Preserve        *bool              `field:"preserve" toml:"preserve"`
	XMLPath         *string            `field:"xml_path" toml:"xml_path"`
	RTFPath         *string            `field:"rtf_path" toml:"rtf_path"`
	IMGPath         *string            `field:"img_path" toml:"img_path"`
	FMVersion       *string            `field:"fm_version" toml:"fm_version"`
	MappingOverride *map[string]string `field:"mapping_override" toml:"mapping_override"`
}
