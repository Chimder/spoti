package user

type CreateUserReq struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Image         string `json:"image"`
	Followers     uint32 `json:"followers"`
	PremiumStatus bool   `json:"premium_status"`
}

// func (cu *CreateUserReq) AddHashPass(new string) {
// 	cu.HashPassword = new
// }
