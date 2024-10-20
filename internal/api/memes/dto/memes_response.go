package dto

type MemeResponse struct {
	ImageUri        string `json:"image_uri"`
	OriginalPostUri string `json:"post_uri"`
}

type MemesResponse struct {
	Memes []MemeResponse `json:"memes"`
}
