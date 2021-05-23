package tools

type Runner interface {
	Do() error
	GetCommandString() string
	GetOutput() []byte
}
