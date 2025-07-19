package dto

type CreateProductInput struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// UpdateProductInput represents the fields allowed when updating a product.
// Only the Name and Price can be modified.
type UpdateProductInput struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetJWTInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UpdateOwnProfileInput struct {
	Email       string `json:"email" binding:"required" validate:"required,email"`
	Password    string `json:"password" binding:"required" validate:"required"`
	Name        string `json:"name" binding:"required" validate:"required"`
	NewPassword string `json:"new_password"`
}
type GetJWTOutput struct {
	AccessToken string `json:"access_token"`
}
