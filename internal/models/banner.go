package models

import (
	"database/sql"
	"github.com/tomi77/go-sqlx"
	"time"
)

type Banner struct {
	BannerId  uint64    `db:"banner_id" json:"banner_id"`
	FeatureId uint64    `db:"feature_id" json:"feature_id"`
	TagIds    []uint64  `db:"tag_ids" json:"tag_ids"`
	Content   string    `db:"content" json:"content"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type PatchBanner struct {
	FeatureId sqlx.NullUint64   `db:"feature_id" json:"feature_id"`
	TagIds    []sqlx.NullUint64 `db:"tag_ids" json:"tag_ids"`
	Content   map[string]interface{}
	IsActive  sql.NullBool
}

type FilterBanner struct {
	FeatureId uint64 `db:"feature_id" json:"featureId"`
	TagId     uint64 `db:"tag_id" json:"tag_id"`
	Limit     uint64
	Offset    uint64
}
