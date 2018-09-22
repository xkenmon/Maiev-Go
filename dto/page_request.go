package dto

type PageRequest struct {
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
	Order string `form:"order"`
	Sort  string `form:"sort"`
}
