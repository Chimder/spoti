package user

type CreateUserReq struct {
	Name          string
	Email         string
	Image         string
	Followers     uint32
	PremiumStatus bool
}
