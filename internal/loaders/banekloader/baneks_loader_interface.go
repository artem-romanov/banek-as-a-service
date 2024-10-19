package banekloader

import "baneks.com/internal/model"

type BaneksLoader interface {
	GetRandomBanek() (model.Banek, error)
}
