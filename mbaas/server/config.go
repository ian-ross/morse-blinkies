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
}
