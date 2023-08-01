package http

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/requiemofthesouls/logger"

	"github.com/requiemofthesouls/svc-http/server"
)

type (
	Manager interface {
		Start(name string) error
		StartAll(ctx context.Context)
		Stop(name string) error
	}

	manager struct {
		servers   map[string]server.Server
		stopChans map[string]chan struct{}
		l         logger.Wrapper
	}
)

func New(servers map[string]server.Server, l logger.Wrapper) Manager {
	return &manager{servers: servers, stopChans: make(map[string]chan struct{}), l: l}
}

func (m *manager) Start(name string) error {
	var (
		s  server.Server
		ok bool
	)
	if s, ok = m.servers[name]; !ok {
		return fmt.Errorf("unknown server '%s'", name)
	}

	m.stopChans[name] = make(chan struct{})
	go func(s server.Server, stopChan chan struct{}, l logger.Wrapper) {
		l.Info(fmt.Sprintf("Start HTTP server '%s'", name))

		if err := s.Start(); err != nil {
			l.Error("Error start HTTP server", logger.Error(err))
		}

		stopChan <- struct{}{}
	}(s, m.stopChans[name], m.l)

	return nil
}

func (m *manager) StartAll(ctx context.Context) {
	for name := range m.servers {
		if err := m.Start(name); err != nil {
			m.l.Error("Error start HTTP server", logger.Error(err))
		}
	}

	<-ctx.Done()

	for name := range m.servers {
		if err := m.Stop(name); err != nil {
			m.l.Error("Error stop HTTP server", logger.Error(err))
		}
	}
}

func (m *manager) Stop(name string) error {
	var (
		s  server.Server
		ok bool
	)
	if s, ok = m.servers[name]; !ok {
		return fmt.Errorf("unknown server '%s'", name)
	}

	if !s.IsStarted() {
		return nil
	}

	m.l.Info(fmt.Sprintf("Stop HTTP server '%s'", name))
	if err := s.Stop(); err != nil {
		return err
	}

	select {
	case <-time.After(time.Second * 5):
		return errors.New("couldn't stop http server within the specified timeout (5 sec)")
	case <-m.stopChans[name]:
	}

	return nil
}
