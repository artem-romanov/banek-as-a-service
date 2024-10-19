package banek_loader

import (
	"math/rand/v2"
)

func GetRandomLoader() BaneksLoader {
	siteLoader := NewBaneksSiteLoader()
	ruLoader := NewBanekRuLoader()
	loaders := []BaneksLoader{siteLoader, ruLoader}
	randomLoaderIndex := rand.IntN(len(loaders))
	return loaders[randomLoaderIndex]
}
