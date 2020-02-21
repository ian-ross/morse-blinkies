package server

// Config contains the configuration information needed to start the
// status server process.
type Config struct {
	// DevMode is a development mode flag: if true, logging is in a
	// human-readable format; if false, logging is in a JSON format.
	DevMode bool `env:"DEV_MODE,default=false"`

	// Port is the port to run the HTTP server on.
	Port int `env:"PORT,default=8080"`

	// CSRFSecret is a secret string used for generating CSRF tokens.
	CSRFSecret string `env:"CSRF_SECRET"`

	// BlinkyScript is the Python script for blinky creation.
	BlinkyScript string `env:"BLINKY_SCRIPT"`

	// WorkDir is the working directory for blinky creation.
	WorkDir string `env:"WORKDIR"`

	// OutputDir is the directory where blinky ZIP and HTML explanation
	// pages are stored.
	OutputDir string `env:"OUTDIR"`

	// Library directory containing default rules file.
	LibDir string `env:"LIBDIR"`

	// TemplateDir is the directory holding the blinky project template.
	TemplateDir string `env:"TEMPLATEDIR"`
}
