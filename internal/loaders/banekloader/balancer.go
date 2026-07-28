package banekloader

import (
	"net/http"
	"sync"
)

type BanekBalancer struct {
	current int
	mu sync.Mutex
	loaders []BaneksLoader
}

func NewBanekBalancer(client *http.Client) *BanekBalancer {
	instance := &BanekBalancer{
		loaders: []BaneksLoader{
			NewBanekRuLoader(client),
			NewBaneksSiteLoader(client),
		},
	}
	return instance
}

func (b *BanekBalancer) GetLoader() BaneksLoader {
	// thread safety for gorutines
	// because it's a round robin approach,
	// we cant allow them to read and write balancer's info all in once
	b.mu.Lock()
	defer b.mu.Unlock()

	// на случай если кто-то решит создать балансер не через New
	if len(b.loaders) == 0 {
		return nil
	}

	loader := b.loaders[b.current]
	b.current++
	if b.current >= len(b.loaders) {
		b.current = 0
	}

	return loader
}
