package def

import (
	cfgDef "github.com/requiemofthesouls/config/def"
	"github.com/requiemofthesouls/container"
	logDef "github.com/requiemofthesouls/logger/def"

	http "github.com/requiemofthesouls/svc-http"
	"github.com/requiemofthesouls/svc-http/server"
	serverDef "github.com/requiemofthesouls/svc-http/server/def"
)

const DIServerManager = "http.server_manager"

func init() {
	container.Register(func(builder *container.Builder, _ map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DIServerManager,
			Build: func(cont container.Container) (interface{}, error) {
				var cfg cfgDef.Wrapper
				if err := cont.Fill(cfgDef.DIWrapper, &cfg); err != nil {
					return nil, err
				}

				var serversConfig []server.Config
				if err := cfg.UnmarshalKey("httpServers", &serversConfig); err != nil {
					return nil, err
				}

				var servers = make(map[string]server.Server, len(serversConfig))
				for _, serverConfig := range serversConfig {
					var serverFactory serverDef.ServerFactory
					if err := cont.Fill(serverDef.DIServerFactory, &serverFactory); err != nil {
						return nil, err
					}

					var (
						srv server.Server
						err error
					)
					if srv, err = serverFactory(serverConfig); err != nil {
						return nil, err
					}

					servers[serverConfig.Name] = srv
				}

				var l logDef.Wrapper
				if err := cont.Fill(logDef.DIWrapper, &l); err != nil {
					return nil, err
				}

				return http.New(servers, l), nil
			},
		})
	})
}
