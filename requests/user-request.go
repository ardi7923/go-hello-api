package requests

type UserCreateRequest struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required" validate:"min:6"`
}

type UserUpdateRequest struct {
	Id       uint64 `json:"id" form:"id" binding:"required"`
	Name     string `json:"name" form:"name" binding:"required"`
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required" validate:"min:6"`
}
