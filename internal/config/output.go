package config

// OutputFormat controls how sync results are printed.
type OutputFormat string

const (
	OutputFormatText OutputFormat = "text"
	OutputFormatJSON OutputFormat = "json"
)

// OutputConfig holds configuration for the output/printer subsystem.
type OutputConfig struct {
	Format  OutputFormat `yaml:"format"`
	Verbose bool         `yaml:"verbose"`
	NoColor bool         `yaml:"no_color"`
	MaskSensitive bool  `yaml:"mask_sensitive"`
}

// DefaultOutputConfig returns sensible defaults for output.
func DefaultOutputConfig() *OutputConfig {
	return &OutputConfig{
		Format:        OutputFormatText,
		Verbose:       false,
		NoColor:       false,
		MaskSensitive: true,
	}
}

// ApplyOutputDefaults fills zero-value fields with defaults.
func ApplyOutputDefaults(c *OutputConfig) {
	if c == nil {
		return
	}
	if c.Format == "" {
		c.Format = OutputFormatText
	}
}

// IsJSON returns true when JSON output format is selected.
func (c *OutputConfig) IsJSON() bool {
	return c != nil && c.Format == OutputFormatJSON
}

// IsVerbose returns true when verbose mode is enabled.
func (c *OutputConfig) IsVerbose() bool {
	return c != nil && c.Verbose
}
