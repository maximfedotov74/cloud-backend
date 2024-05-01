package dto

type CreateUser struct {
	Email    string `json:"email" example:"makc-dgek@mail.ru" validate:"required,email"`
	Password string `json:"password" example:"1234567890" validate:"required,min=6"`
}

type UpdateUser struct {
	AvatarPath *string `json:"avatar_path" validate:"omitempty,filepath"`
}

type ConfirmChangePassword struct {
	Code string `json:"code" validate:"required,min=6,max=6" example:"123456"`
}

type ChangePassword struct {
	OldPassword string `json:"old_password" validate:"required,min=6,max=100" example:"sdfsdfs222"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=100" example:"sdfsdfs222"`
}
