package time

type Option func(*Options)

type Options struct {
	UseSystemTime bool

	RemoteServerAddr string

	RemoteServerToken string

	Ns string
}

func DefaultOptions() *Options {
	return &Options{
		UseSystemTime:    true,
		RemoteServerAddr: "",
	}
}

func WithGo() Option {
	return func(o *Options) {
		o.UseSystemTime = true
	}
}

func WithLocal() Option {
	return func(o *Options) {
		o.RemoteServerAddr = ""
		o.UseSystemTime = false
	}
}

func WithRemote(addr, token string) Option {
	return func(o *Options) {
		o.RemoteServerAddr = addr
		o.RemoteServerToken = token
		o.UseSystemTime = false
	}
}

func WithNamespace(ns string) Option {
	return func(o *Options) {
		o.Ns = ns
	}
}
