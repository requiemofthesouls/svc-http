package def

import (
	"github.com/requiemofthesouls/container"
	logDef "github.com/requiemofthesouls/logger/def"

	"github.com/requiemofthesouls/svc-http/server"
)

const (
	DIServerFactory = "http.server_factory"
	DIHandlerPrefix = "http.handler."
)

type ServerFactory func(cfg server.Config) (server.Server, error)

func init() {
	container.Register(func(builder *container.Builder, _ map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DIServerFactory,
			Build: func(cont container.Container) (interface{}, error) {
				var l logDef.Wrapper
				if err := cont.Fill(logDef.DIWrapper, &l); err != nil {
					return nil, err
				}

				return func(cfg server.Config) (server.Server, error) {
					var handler server.Handler
					if err := cont.Fill(DIHandlerPrefix+cfg.Name, &handler); err != nil {
						return nil, err
					}

					return server.New(cfg, handler), nil
				}, nil
			},
		})
	})
}
