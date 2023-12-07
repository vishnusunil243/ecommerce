package response

type UserData struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	ReportCount int
}

type UserDetails struct {
	Id                uint
	Name              string
	Email             string
	Mobile            string
	IsBlocked         bool
	ReportCount       int
	BlockedAt         string `json:",omitempty"`
	BlockedBy         uint   `json:",omitempty"`
	ReasonForBlocking string `json:",omitempty"`
}
type UserProfile struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
	ReferralId string `json:"ReferralId,omitempty"`
	Address    `gorm:"embedded" json:"address,omitempty"`
}
type Address struct {
	House_number string `json:"house_number,omitempty" `
	Street       string `json:"street,omitempty" `
	City         string `json:"city,omitempty" `
	District     string `json:"district,omitempty" `
	Landmark     string `json:"landmark,omitempty" `
	Pincode      int    `json:"pincode,omitempty" `
	IsDefault    bool   `json:"isdefault,omitempty"`
}
