package time

type Option func(*Options)

type Options struct {
	ErrHandlers []func(error)

	Location *Location

	RemoteServerAddr  string
	RemoteServerForce bool
	RemoteServerToken string
}

func DefaultOptions() *Options {
	return &Options{
		Location:         Local,
		RemoteServerAddr: "",
	}
}

func WithErrHandler(cbs ...func(error)) Option {
	return func(o *Options) {
		o.ErrHandlers = append(o.ErrHandlers, cbs...)
	}
}

func WithRemote(addr string, force bool) Option {
	return func(o *Options) {
		o.RemoteServerAddr = addr
		o.RemoteServerForce = force
	}
}

func WithLocation(loc *Location) Option {
	return func(o *Options) {
		o.Location = loc
	}
}
