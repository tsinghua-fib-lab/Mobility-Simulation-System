package client

import (
	"context"
	"flag"
	"sync"

	"git.fiblab.net/sim/sidecar/core"
	routingv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/routing/v2"
)

var (
	// uriHints
	routingHint = flag.String("routing", "localhost:52101", "routing service uri hint")
)

type RoutingServiceClient struct {
	*core.BaseServiceClient
	client routingv2.RoutingServiceClient

	wg sync.WaitGroup
}

func NewRoutingServiceClient() *RoutingServiceClient {
	s := &RoutingServiceClient{
		BaseServiceClient: core.NewBaseServiceClient("", "routing", nil, *routingHint),
	}
	s.client = routingv2.NewRoutingServiceClient(s.Conn)
	return s
}

func (s *RoutingServiceClient) GetRoute(
	ctx context.Context, in *routingv2.GetRouteRequest, process func(res *routingv2.GetRouteResponse)) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.BaseServiceClient.Wait()
		if res, err := s.client.GetRoute(ctx, in); err != nil {
			log.Panic(err)
		} else {
			process(res)
		}
	}()
}

func (s *RoutingServiceClient) Wait() {
	s.wg.Wait()
}
