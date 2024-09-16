package cmd

type JaqenConfig struct {
	Preserve        *bool              `field:"preserve"`
	XMLPath         *string            `field:"xml_path"`
	RTFPath         *string            `field:"rtf_path"`
	IMGPath         *string            `field:"img_path"`
	FMVersion       *string            `field:"fm_version"`
	MappingOverride *map[string]string `field:"mapping_override"`
}
