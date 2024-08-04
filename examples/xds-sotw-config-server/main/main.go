package main

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	xdwconfig "github.com/envoyproxy/ratelimit/examples/xds-sotw-config-server"
)

var (
	port   uint
	nodeID string
)

func init() {
	flag.UintVar(&port, "port", 18000, "xDS management server port")
	flag.StringVar(&nodeID, "nodeID", "test-node-id", "Node ID")
}

func initLogger() *log.Logger {
	logger := log.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(&log.TextFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.DebugLevel)

	return logger
}

func run(l *log.Logger) int {
	flag.Parse()

	logger := l.WithFields(log.Fields{
		"nodeID": nodeID,
		"port":   port,
	})
	logger.Info("Starting up...")
	defer logger.Info("Ending up...")

	// Create a cache
	snapshotCache := cache.NewSnapshotCache(false, cache.IDHash{}, logger)

	// Create the snapshot that we'll serve to Envoy
	snapshot := xdwconfig.GenerateSnapshot()
	if err := snapshot.Consistent(); err != nil {
		logger.Errorf("Snapshot is inconsistent: %+v\n%+v", snapshot, err)
		return 1
	}
	logger.Debugf("Will serve snapshot %+v", snapshot)

	ctx := context.Background()

	// Add the snapshot to the cache
	if err := snapshotCache.SetSnapshot(ctx, nodeID, snapshot); err != nil {
		logger.Errorf("Snapshot error %q for %+v", err, snapshot)
		return 1
	}

	// Run the xDS server
	cb := xdwconfig.NewDebugCallback(logger)
	srv := server.NewServer(ctx, snapshotCache, cb)
	xdwconfig.RunServer(logger, srv, port)

	return 0
}

func main() {
	os.Exit(run(initLogger()))
}
