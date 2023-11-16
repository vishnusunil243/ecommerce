package helperStruct

type UserReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Mobile   string `json:"mobile" binding:"required"`
	Password string `json:"password"`
	OTP      string `json:"OTP"`
}

type LoginReq struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password" binding:"required"`
}

type OTPData struct {
	PhoneNumber string `json:"phoneNumber,omitempty" validate:"required"`
}

type VerifyOtp struct {
	User *OTPData `json:"user,omitempty" validate:"required"`
	Code string   `json:"code,omitempty" validate:"required"`
}

type Address struct {
	House_number string `json:"house_number" `
	Street       string `json:"street" `
	City         string `json:"city" `
	District     string `json:"district" `
	Landmark     string `json:"landmark" `
	Pincode      int    `json:"pincode" `
	IsDefault    bool   `json:"isdefault" `
}

type UpdatePassword struct {
	OldPassword string `json:"oldpassword" `
	NewPassword string `json:"newpassword" `
}
type ForgotPassword struct {
	Email       string `json:"email"`
	NewPassword string `json:"newpassword"`
	OTP         string `json:"otp"`
}
type UpdateMobile struct {
	Mobile string `json:"mobile"`
}
