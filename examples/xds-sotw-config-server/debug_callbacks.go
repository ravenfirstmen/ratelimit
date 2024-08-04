package xdsconfig

import (
	"context"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	log "github.com/sirupsen/logrus"
	"sync"
)

type DebugCallback struct {
	Logger *log.Entry
	//
	fetches        int
	requests       int
	responses      int
	deltaRequests  int
	deltaResponses int
	mu             sync.Mutex
}

var _ server.Callbacks = &DebugCallback{}

func NewDebugCallback(logger *log.Entry) *DebugCallback {
	return &DebugCallback{Logger: logger}
}

func (cb *DebugCallback) Report() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	log.Printf("server callbacks fetches=%d requests=%d responses=%d\n", cb.fetches, cb.requests, cb.responses)
}

func (cb *DebugCallback) OnStreamOpen(_ context.Context, id int64, typ string) error {
	log.Printf("stream %d open for %s\n", id, typ)
	return nil
}

func (cb *DebugCallback) OnStreamClosed(id int64, node *core.Node) {
	log.Printf("stream %d of node %s closed\n", id, node.Id)
}

func (cb *DebugCallback) OnDeltaStreamOpen(_ context.Context, id int64, typ string) error {
	log.Printf("delta stream %d open for %s\n", id, typ)
	return nil
}
func (cb *DebugCallback) OnDeltaStreamClosed(id int64, node *core.Node) {
	log.Printf("delta stream %d of node %s closed\n", id, node.Id)
}

func (cb *DebugCallback) OnStreamRequest(id int64, req *discovery.DiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.requests++
	log.Printf("received request for %s on stream %d: %v:%v", req.GetTypeUrl(), id, req.VersionInfo, req.ResourceNames)

	return nil
}

func (cb *DebugCallback) OnStreamResponse(ctx context.Context, id int64, req *discovery.DiscoveryRequest, res *discovery.DiscoveryResponse) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.responses++
	log.Printf("responding to request for %s on stream %d", req.GetTypeUrl(), id)
}

func (cb *DebugCallback) OnStreamDeltaResponse(id int64, req *discovery.DeltaDiscoveryRequest, res *discovery.DeltaDiscoveryResponse) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.deltaResponses++
}
func (cb *DebugCallback) OnStreamDeltaRequest(int64, *discovery.DeltaDiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.deltaRequests++

	return nil
}
func (cb *DebugCallback) OnFetchRequest(context.Context, *discovery.DiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.fetches++
	return nil
}
func (cb *DebugCallback) OnFetchResponse(*discovery.DiscoveryRequest, *discovery.DiscoveryResponse) {}
