package xtime

var (
	RemoteServerAddr = ""
)

type Option func(*Options)

type Options struct {
	ErrHandlers []func(error)
}

func WithErrHandler(cbs ...func(error)) Option {
	return func(o *Options) {
		o.ErrHandlers = append(o.ErrHandlers, cbs...)
	}
}
