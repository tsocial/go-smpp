package smpp

import (
	"context"
	"sync"

	"github.com/tsocial/logger"
)

// ClientPool is a pool of clients
type ClientPool struct {
	clients []*Client
	size    int
	index   int
	lock    sync.RWMutex
}

// NewClientPool initializes the pool with n clients that connect to different SMSC servers.
func NewClientPool(ctx context.Context, servers []*Config, eventHandler IEventHandler) (*ClientPool, error) {
	pool := ClientPool{size: len(servers)}

	for _, server := range servers {
		client, err := NewClient(ctx, server, eventHandler)
		if err != nil {
			return nil, err
		}
		pool.clients = append(pool.clients, client)
	}

	return &pool, nil
}

// Disconnect disconnects all clients in the pool
func (pool *ClientPool) Disconnect(ctx context.Context) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	for _, client := range pool.clients {
		client.Disconnect(ctx)
	}
	pool.clients = nil
}

// Size returns the number of clients in the pool
func (pool *ClientPool) Size() int {
	pool.lock.RLock()
	defer pool.lock.RUnlock()
	return pool.size
}

// Send sends SMS using one client from the pool
func (pool *ClientPool) Send(ctx context.Context, sms *SMS) *SubmitReport {
	client := pool.NextClient()
	report := client.Send(ctx, sms)

	if report.Error == nil {
		return report
	}

	logger.PrintError(ctx, "retrying", report.Error)

	if newClient := pool.GetAnotherClient(client); newClient != nil {
		return newClient.Send(ctx, sms)
	}

	return report
}

// NextClient returns the next client from the pool using round-robin mechanism
func (pool *ClientPool) NextClient() *Client {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	pool.index = (pool.index + 1) % pool.size
	return pool.clients[pool.index]
}

// GetAnotherClient returns a client that is different with the given client
func (pool *ClientPool) GetAnotherClient(selected *Client) *Client {
	pool.lock.RLock()
	defer pool.lock.RUnlock()

	for _, client := range pool.clients {
		if client != selected {
			return client
		}
	}

	return nil
}
