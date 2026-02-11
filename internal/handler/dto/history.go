package dto

import "github.com/ArtemChadaev/GoCreateHistory/internal/domain"

type CreateHistoryRequest struct {
	Description string `json:"description" validate:"required"`
	ChapterSize *int   `json:"chapter_size" default:"3" validate:"min=1,max=10"`
	ImageSize   *int   `json:"image_size" default:"1" validate:"min=0,max=2"`
	Save        *bool  `json:"save" default:"true"`
}

func (r *CreateHistoryRequest) ToDomain(userID int) domain.UserRequest {
	return domain.UserRequest{
		UserID:      userID,
		Description: r.Description,
		ChapterSize: *r.ChapterSize,
		ImageSize:   *r.ImageSize,
		Save:        *r.Save,
	}
}
