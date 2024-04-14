package models

import (
	"encoding/json"
	"time"
)

type Banner struct {
	BannerId  uint64          `db:"banner_id" json:"banner_id"`
	FeatureId uint64          `db:"feature_id" json:"feature_id"`
	TagIds    []uint64        `db:"tag_ids" json:"tag_ids"`
	Content   json.RawMessage `db:"content" json:"content"`
	IsActive  bool            `db:"is_active" json:"is_active"`
	Version   uint64          `db:"version" json:"version"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}

type PatchBanner struct {
	FeatureId *uint64         `db:"feature_id" json:"feature_id"`
	TagIds    []uint64        `db:"tag_ids" json:"tag_ids"`
	Content   json.RawMessage `db:"content" json:"content"`
	IsActive  *bool           `db:"is_active" json:"is_active"`
}

type FilterBanner struct {
	FeatureId uint64 `db:"feature_id" json:"featureId"`
	TagId     uint64 `db:"tag_id" json:"tag_id"`
	Limit     uint64
	Offset    uint64
}

type FeatureTag struct {
	FeatureId uint64 `db:"feature_id" json:"feature_id"`
	TagId     uint64 `db:"tag_id" json:"tag_id"`
}

type BannerContent struct {
	Content  string `db:"content"`
	IsActive bool   `db:"is_active"`
}
