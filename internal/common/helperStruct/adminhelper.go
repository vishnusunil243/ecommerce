package helperStruct

type CreateAdmin struct {
	Name     string ` json:"name" validate:"required"`
	Email    string ` json:"email" validate:"required" binding:"email"`
	Password string ` json:"password" validate:"required"`
}
type SuperLoginReq struct {
	Email    string
	Password string
}
type BlockData struct {
	UserId uint ` json:"userid" validate:"required"`
}
type Dashboard struct {
	StartDate string `json:"date1"`
	EndDate   string `json:"date2"`
	Month     int    `json:"month"`
	Day       int    `json:"day"`
	Year      int    `json:"year"`
}
