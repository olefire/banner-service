package http

import (
	"banner-service/internal/models"
	"banner-service/internal/repository"
	authservice "banner-service/internal/service/auth"
	bannerservice "banner-service/internal/service/banner"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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
	publicKey string
}

func NewController(as AuthService, bs BannerService, pk string) *Controller {
	return &Controller{
		AuthService:   as,
		BannerService: bs,
		publicKey:     pk,
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

func (ctr *Controller) GetBanner(w http.ResponseWriter, r *http.Request) {
	tagId, err := strconv.ParseUint(r.URL.Query().Get("tag_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	featureId, err := strconv.ParseUint(r.URL.Query().Get("feature_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	role, ok := r.Context().Value("role").(models.UserRole)
	if !ok {
		http.Error(w, "Role not found in context", http.StatusUnauthorized)
		return
	}

	content, err := ctr.BannerService.GetBanner(r.Context(), tagId, featureId, role)
	if errors.Is(err, repository.ErrAccessDenied) {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(content))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

//ToDo: add methods for banner
