package dto

type CreateUserRequest struct {
	Name      string   `json:"name" validate:"required"`
	Neighbors []string `json:"neighbors" validate:"min=1,dive,required"`
}

type UpdateUserRequest struct {
	Name      *string   `json:"name" validate:"omitempty"`
	Neighbors *[]string `json:"neighbors" validate:"omitempty,min=1,dive,required"`
}

type UserResponse struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Neighbors []string `json:"neighbors"`
}
