package smpp

import (
	"context"
	"os"
	"testing"

	"github.com/tsocial/logger"
)

var testConfigs []*Config

// test connect to the same server
var (
	testEventHandler   *SampleEventHandler
	testClient         *Client
	testSameServerPool *ClientPool
	testPoolSize       = 2
)

// test connect to different servers
var (
	testDifferentServersPool    *ClientPool
	testDifferentServersHandler *SampleEventHandler
)

func TestMain(m *testing.M) {
	logger.SetupLogger("smpp")

	testConfigs = LoadConfigsFromEnvironment()
	testEventHandler = NewSampleEventHandler()
	testDifferentServersHandler = NewSampleEventHandler()

	ctx := context.Background()
	testClient, err := NewClient(ctx, testConfigs[0], testEventHandler)
	if err != nil {
		panic(err)
	}

	sameServerConfigs := []*Config{testConfigs[0], testConfigs[0]}
	testSameServerPool, err = NewClientPool(ctx, sameServerConfigs, testEventHandler)
	if err != nil {
		panic(err)
	}

	testDifferentServersPool, err = NewClientPool(ctx, testConfigs[1:], testDifferentServersHandler)
	if err != nil {
		panic(err)
	}

	defer func() {
		if testClient != nil {
			testClient.Disconnect(ctx)
		}

		if testSameServerPool != nil {
			testSameServerPool.Disconnect(ctx)
		}

		if testDifferentServersPool != nil {
			testDifferentServersPool.Disconnect(ctx)
		}
	}()

	code := m.Run()
	os.Exit(code) //nolint:gocritic
}
