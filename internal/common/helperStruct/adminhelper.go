package helperStruct

type CreateAdmin struct {
	Name     string ` json:"name" validate:"required"`
	Email    string ` json:"email" validate:"required" binding:"email"`
	Password string ` json:"password" validate:"required"`
	IsSuper  bool   `json:"isSuper" validate:"required"`
}

type BlockData struct {
	UserId uint ` json:"userid" validate:"required"`
}
