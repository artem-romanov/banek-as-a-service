package memesloader

import "baneks.com/internal/model"

type MemeLoader interface {
	GetRandomMemes() ([]model.Meme, error)
}
