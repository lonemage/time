package time

type Option func(*Options)

type Options struct {
	RemoteServerAddr  string
	RemoteServerToken string

	Ns string
}

func DefaultOptions() *Options {
	return &Options{
		RemoteServerAddr: "",
	}
}

func WithRemote(addr string) Option {
	return func(o *Options) {
		o.RemoteServerAddr = addr
	}
}

func WithNamespace(ns string) Option {
	return func(o *Options) {
		o.Ns = ns
	}
}
