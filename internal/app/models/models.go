package models

// Ad ...
type Ad struct {
	ID          int      `json:"id,omitempty"`
	Name        string   `json:"name" binding:"required,min=0,max=200"`
	Description string   `json:"description,omitempty" binding:"required,min=0,max=1000"`
	Links       []string `json:"photos" binding:"required,max=3,dive,url"`
	Price       float64  `json:"price" binding:"min=0"`
	Link        string   `json:"-"`
}
