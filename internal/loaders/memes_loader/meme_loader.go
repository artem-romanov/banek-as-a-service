package memesloader

import "baneks.com/internal/model"

type RandomMemesConfig struct {
	Year int
}
type MemeLoader interface {
	GetRandomMemes() ([]model.Meme, error)
	GetRandomMemesWithConfig(RandomMemesConfig) ([]model.Meme, error)
}
