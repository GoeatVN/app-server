package userdto

type AddUserRequest struct {
	Username  string `json:"username" binding:"required" validate:"required" message:"Username is required."`
	Password  string `json:"password" binding:"required" validate:"required" message:"Password is required."`
	Email     string `json:"email" binding:"required" validate:"required,email" message:"Email is required and must be valid."`
	Phone     string `json:"phone" binding:"required" validate:"required" message:"Phone number is required."`
	FullName  string `json:"fullName" binding:"required" validate:"required" message:"Full name is required."`
	Status    string `json:"status" binding:"required" validate:"required" message:"Status is required."`
	CreatedBy string `json:"-"`
}
