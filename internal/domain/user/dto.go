package user

type CreateUserReq struct {
	Name          string
	Email         string
	Password      string
	HashPassword  string
	Image         string
	Followers     uint32
	PremiumStatus bool
}

// func (cu *CreateUserReq) AddHashPass(new string) {
// 	cu.HashPassword = new
// }
