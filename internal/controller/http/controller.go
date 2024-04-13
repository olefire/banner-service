package http

import (
	"banner-service/internal/models"
	"banner-service/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type AuthProvider struct {
	AuthManagement
}

type BannerService struct {
	BannerManagement
}

type Controller struct {
	AuthProvider
	BannerService
	publicKey string
}

func NewController(as AuthProvider, bs BannerService, pk string) *Controller {
	return &Controller{
		AuthProvider:  as,
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

	err = ctr.AuthProvider.SignUp(r.Context(), &user)
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
	var signInInput models.User
	err := json.NewDecoder(r.Body).Decode(&signInInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := ctr.AuthProvider.SignIn(r.Context(), &signInInput)
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

	//role, ok := r.Context().Value("role").(models.UserRole)
	//if !ok {
	//	http.Error(w, "Role not found in context", http.StatusUnauthorized)
	//	return
	//}
	role := models.UserRole("admin")
	content, err := ctr.BannerService.GetBanner(r.Context(), tagId, featureId, role)
	if errors.Is(err, repository.ErrBannerInactive) {
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

func (ctr *Controller) GetFilteredBanners(w http.ResponseWriter, r *http.Request) {
	var filter models.FilterBanner

	if featureIdStr := r.URL.Query().Get("feature_id"); featureIdStr != "" {
		featureId, err := strconv.ParseUint(featureIdStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid feature_id", http.StatusBadRequest)
			return
		}
		filter.FeatureId = &featureId
	}

	if tagIdStr := r.URL.Query().Get("tag_id"); tagIdStr != "" {
		tagId, err := strconv.ParseUint(tagIdStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid tag_id", http.StatusBadRequest)
			return
		}
		filter.TagId = &tagId
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		filter.Limit = limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.ParseUint(offsetStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
		filter.Offset = offset
	}

	banners, err := ctr.BannerService.GetFilteredBanners(r.Context(), &filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bannersJSON, err := json.Marshal(banners)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bannersJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ctr *Controller) CreateBanner(w http.ResponseWriter, r *http.Request) {
	var banner *models.Banner
	err := json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bannerId, err := ctr.BannerService.CreateBanner(r.Context(), banner)
	if err != nil {
		http.Error(w, fmt.Sprintf("CreateBanner error: %v ", err), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(strconv.FormatUint(bannerId, 10)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ctr *Controller) PartialUpdateBanner(w http.ResponseWriter, r *http.Request) {
	bannerId, err := strconv.ParseUint(chi.URLParam(r, "banner_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var banner models.PatchBanner
	err = json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ctr.BannerService.PartialUpdateBanner(r.Context(), bannerId, &banner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ctr *Controller) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	bannerId, err := strconv.ParseUint(chi.URLParam(r, "banner_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ctr.BannerService.DeleteBanner(r.Context(), bannerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ctr *Controller) GetListOfVersions(w http.ResponseWriter, r *http.Request) {
	bannerId, err := strconv.ParseUint(chi.URLParam(r, "banner_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	banners, err := ctr.BannerService.GetListOfVersions(r.Context(), bannerId)
	if err != nil {
		//ToDo: add error handling
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bannersJSON, err := json.Marshal(banners)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bannersJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (ctr *Controller) ChooseBannerVersion(w http.ResponseWriter, r *http.Request) {
	bannerId, err := strconv.ParseUint(chi.URLParam(r, "version"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	version, err := strconv.ParseUint(chi.URLParam(r, "version"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ctr.BannerService.ChooseVersion(r.Context(), bannerId, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

//ToDo: add methods for banner
