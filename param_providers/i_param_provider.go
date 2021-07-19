package param_providers

type ParamProvider interface {
	Init()
	Changed() bool
}
