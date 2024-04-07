package http

type AuthService struct {
}

type BannerService struct {
}

type Controller struct {
	AuthService
	BannerService
	privateKey string
	publicKey  string
}

func NewController(as AuthService, bs BannerService, privateKey, publicKey string) *Controller {
	return &Controller{
		AuthService:   as,
		BannerService: bs,
		privateKey:    privateKey,
		publicKey:     publicKey,
	}
}
