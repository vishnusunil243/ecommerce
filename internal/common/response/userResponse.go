package response

type UserData struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	ReportCount int
}

type UserDetails struct {
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
	Name    string `json:"name"`
	Email   string `json:"email"`
	Mobile  string `json:"mobile"`
	Address `gorm:"embedded" json:"address"`
}
type Address struct {
	House_number string `json:"house_number" `
	Street       string `json:"street" `
	City         string `json:"city" `
	District     string `json:"district" `
	Landmark     string `json:"landmark" `
	Pincode      int    `json:"pincode" `
	IsDefault    bool   `json:"isdefault"`
}
