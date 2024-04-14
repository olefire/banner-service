package http

import (
	"banner-service/internal/auth"
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

func (ctr *Controller) SignUpEndpoint(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := ctr.AuthProvider.SignUp(r.Context(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(token))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctr *Controller) SignInEndpoint(w http.ResponseWriter, r *http.Request) {
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

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(token))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctr *Controller) GetBannerEndpoint(w http.ResponseWriter, r *http.Request) {
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

	useLastRevision := r.URL.Query().Get("use_last_revision") == "true"

	role := auth.GetRole(r.Context())

	content, err := ctr.BannerService.GetBanner(r.Context(), tagId, featureId, role, useLastRevision)
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

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(content))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctr *Controller) GetFilteredBannersEndpoint(w http.ResponseWriter, r *http.Request) {
	var filter models.FilterBanner

	featureIdStr := r.URL.Query().Get("feature_id")
	featureId, err := strconv.ParseUint(featureIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid feature_id", http.StatusBadRequest)
		return
	}
	filter.FeatureId = featureId

	tagIdStr := r.URL.Query().Get("tag_id")
	tagId, err := strconv.ParseUint(tagIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid tag_id", http.StatusBadRequest)
		return
	}
	filter.TagId = tagId

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}
	filter.Limit = limit

	offsetStr := r.URL.Query().Get("offset")
	offset, err := strconv.ParseUint(offsetStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid offset", http.StatusBadRequest)
		return
	}
	filter.Offset = offset

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

type CreateDTO struct {
	TagIds    []uint64        `json:"tag_ids"`
	FeatureId uint64          `json:"feature_id"`
	Content   json.RawMessage `json:"content"`
	IsActive  bool            `json:"is_active"`
}

func (ctr *Controller) CreateBannerEndpoint(w http.ResponseWriter, r *http.Request) {
	var banner *CreateDTO
	err := json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bannerId, err := ctr.BannerService.CreateBanner(r.Context(), &models.Banner{
		TagIds:    banner.TagIds,
		FeatureId: banner.FeatureId,
		Content:   banner.Content,
		IsActive:  banner.IsActive,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("CreateBanner error: %v ", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(strconv.FormatUint(bannerId, 10)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctr *Controller) PartialUpdateBannerEndpoint(w http.ResponseWriter, r *http.Request) {
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
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ctr *Controller) DeleteBannerEndpoint(w http.ResponseWriter, r *http.Request) {
	bannerId, err := strconv.ParseUint(chi.URLParam(r, "banner_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = ctr.BannerService.DeleteBanner(r.Context(), bannerId); errors.Is(err, repository.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ctr *Controller) GetListOfVersionsEndpoint(w http.ResponseWriter, r *http.Request) {
	bannerId, err := strconv.ParseUint(chi.URLParam(r, "banner_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	banners, err := ctr.BannerService.GetListOfVersions(r.Context(), bannerId)
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bannersJSON, err := json.Marshal(banners)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bannersJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctr *Controller) ChooseBannerVersionEndpoint(w http.ResponseWriter, r *http.Request) {
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

	err = ctr.BannerService.ChooseBannerVersion(r.Context(), bannerId, version)
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ctr *Controller) MarkBannerAsDeletedEndpoint(w http.ResponseWriter, r *http.Request) {
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

	err = ctr.BannerService.MarkBannerAsDeleted(r.Context(), tagId, featureId)
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
