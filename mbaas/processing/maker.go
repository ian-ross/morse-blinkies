package processing

// BlinkyMaker is a structure used to manage information needed for
// running the Python blinky maker script.
type BlinkyMaker struct {
	Script      string
	WorkDir     string
	OutputDir   string
	TemplateDir string
}
