package memesloader

import (
	"context"

	"baneks.com/internal/model"
)

type RandomMemesConfig struct {
	Year int
}
type MemeLoader interface {
	GetRandomMemes(context.Context) ([]model.Meme, error)
	GetRandomMemesWithConfig(context.Context, RandomMemesConfig) ([]model.Meme, error)
}
