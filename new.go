package slgo

import "github.com/yayacat/slgo/handlers"

func New(provider handlers.SeedLinkProvider, consumer handlers.SeedLinkConsumer, hooks handlers.SeedLinkHooks) SeedLinkServer {
	return SeedLinkServer{
		hooks:    hooks,
		Provider: provider,
		Consumer: consumer,
	}
}
