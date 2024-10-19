package banek_loader

import (
	"sync"
)

type Loader uint

const (
	SITE_LOADER = iota
	RU_LOADER
)

var mutex sync.Mutex

type BanekBalancer struct {
	currentLoader Loader
}

var instance *BanekBalancer
var once sync.Once

func GetBalancer() *BanekBalancer {
	once.Do(func() {
		instance = &BanekBalancer{
			currentLoader: SITE_LOADER,
		}
	})
	return instance
}

func (balancer *BanekBalancer) GetLoader() BaneksLoader {
	// thread safety for gorutines
	// because it's a round robin approach,
	// we cant allow them to read and write balancer's info all in once
	mutex.Lock()
	defer mutex.Unlock()
	defer balancer.toggleBalancer()

	switch balancer.currentLoader {
	case SITE_LOADER:
		return NewBaneksSiteLoader()
	case RU_LOADER:
		return NewBanekRuLoader()
	default:
		return NewBaneksSiteLoader()
	}
}

func (balancer *BanekBalancer) toggleBalancer() {
	// not thread safe
	switch balancer.currentLoader {
	case SITE_LOADER:
		balancer.currentLoader = RU_LOADER
	case RU_LOADER:
		balancer.currentLoader = SITE_LOADER
	default:
		balancer.currentLoader = SITE_LOADER
	}
}
