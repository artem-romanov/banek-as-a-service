package banekloader

import (
	"context"

	"baneks.com/internal/model"
)

type BaneksLoader interface {
	GetRandomBanek(context.Context) (model.Banek, error)
}
