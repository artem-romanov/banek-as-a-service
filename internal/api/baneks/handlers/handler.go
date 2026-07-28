package handlers

import "baneks.com/internal/loaders/banekloader"

type Handlers struct {
	banekBalancer *banekloader.BanekBalancer
	banekSiteLoader *banekloader.BaneksSiteLoader
}

func NewHandlers(
	banekBalancer *banekloader.BanekBalancer,
	banekSiteLoader *banekloader.BaneksSiteLoader,
) *Handlers {
	return &Handlers{
		banekBalancer: banekBalancer,
		banekSiteLoader: banekSiteLoader,
	}
}
