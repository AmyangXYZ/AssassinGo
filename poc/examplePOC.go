package poc

// ExamplePOC .
type ExamplePOC struct {
	target string
}

// NewExamplePOC returns ExamplePOC.
func NewExamplePOC() *ExamplePOC {
	return &ExamplePOC{}
}

// Run impletements POC interface.
func (e *ExamplePOC) Run(target string) {
	e.target = target
}
