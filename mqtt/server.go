package mqtt

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/tx7do/kratos-transport/common"
	"net/url"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

// Server is a mqtt server wrapper.
type Server struct {
	common.Broker

	subscribers map[string]common.Subscriber

	baseCtx context.Context
	log     *log.Helper
	err     error
}

// NewServer creates a mqtt server by options.
func NewServer(opts ...common.Option) *Server {
	srv := &Server{
		baseCtx:     context.Background(),
		log:         log.NewHelper(log.GetLogger()),
		Broker:      NewBroker(opts...),
		subscribers: map[string]common.Subscriber{},
	}

	return srv
}

// Endpoint return a real address to registry endpoint.
// examples:
//   tcp://emqx:public@broker.emqx.io:1883
func (s *Server) Endpoint() (*url.URL, error) {
	if s.err != nil {
		return nil, s.err
	}
	return url.Parse(s.Address())
}

// Start the mqtt server.
func (s *Server) Start(ctx context.Context) error {
	if s.err != nil {
		return s.err
	}

	s.baseCtx = ctx

	s.log.Infof("[mqtt] server listening on: %s", s.Address())

	s.err = s.Connect()
	if s.err != nil {
		return s.err
	}

	return nil
}

// Stop the mqtt server.
func (s *Server) Stop(_ context.Context) error {
	s.log.Info("[mqtt] server stopping")
	return s.Disconnect()
}

// RegisterSubscriber is syntactic sugar for registering a subscriber
func (s *Server) RegisterSubscriber(topic string, h common.Handler, opts ...common.SubscribeOption) error {
	sub, err := s.Subscribe(topic, h, opts...)
	if err != nil {
		return err
	}

	s.subscribers[topic] = sub

	return nil
}
