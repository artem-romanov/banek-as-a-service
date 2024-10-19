package dto

import "baneks.com/internal/models"

// Yeah, I know this struct is stupid.
// But it's better for future changes between DTO and original entity
type BanekResponse struct {
	Text  string `json:"text"`
	Likes int    `json:"likes"`
}

func BanekToResponse(banek *models.Banek) BanekResponse {
	return BanekResponse{
		Text:  banek.Text,
		Likes: banek.Likes,
	}
}
