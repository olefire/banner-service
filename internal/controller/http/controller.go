package http

import (
	"banner-service/internal/models"
	authservice "banner-service/internal/service/auth"
	bannerservice "banner-service/internal/service/banner"
	"encoding/json"
	"net/http"
)

type AuthService struct {
	authservice.AuthManagement
}

type BannerService struct {
	bannerservice.BannerManagement
}

type Controller struct {
	AuthService
	BannerService
}

func NewController(as AuthService, bs BannerService) *Controller {
	return &Controller{
		AuthService:   as,
		BannerService: bs,
	}
}

func (ctr *Controller) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ctr.AuthService.SignUp(r.Context(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte("User created"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (ctr *Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	var signInInput *models.User
	err := json.NewDecoder(r.Body).Decode(&signInInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := ctr.AuthService.SignIn(r.Context(), signInInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(token))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}

//ToDo: add methods for banner
