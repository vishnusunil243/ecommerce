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
