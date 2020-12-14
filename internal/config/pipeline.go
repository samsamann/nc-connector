package config

// FilePipelineConfig struct holds all configurations for one file pipeline.
type FilePipelineConfig struct {
	Import FileImporterConfig
}

// FileImporterConfig struct holds all configurations for the importer.
type FileImporterConfig struct {
	Name    string
	Options map[string]interface{}
}
