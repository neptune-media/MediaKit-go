package tools

type Check interface {
	Validate() error
}

type Runner interface {
	Do() error
	GetCommandString() string
	GetOutput() []byte
}
