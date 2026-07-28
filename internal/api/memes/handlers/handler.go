package handlers

import memesloader "baneks.com/internal/loaders/memes_loader"


type Handlers struct {
	memeLoader memesloader.MemeLoader
}

func NewHandlers(memeLoader memesloader.MemeLoader) *Handlers {
	return &Handlers{
		memeLoader: memeLoader,
	}
}
