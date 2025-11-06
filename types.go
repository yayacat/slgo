package slgo

import (
	"github.com/yayacat/slgo/handlers"
)

type handler interface {
	Callback(*handlers.SeedLinkClient, handlers.SeedLinkProvider, handlers.SeedLinkConsumer, ...string) error
	Fallback(*handlers.SeedLinkClient, handlers.SeedLinkProvider, handlers.SeedLinkConsumer, ...string)
}

type SeedLinkCommand struct {
	Handler handler
	HasArgs bool
}

type SeedLinkServer struct {
	hooks    handlers.SeedLinkHooks
	Provider handlers.SeedLinkProvider
	Consumer handlers.SeedLinkConsumer
}
