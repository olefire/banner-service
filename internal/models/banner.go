package models

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/tomi77/go-sqlx"
	"time"
)

type Banner struct {
	BannerId  uuid.UUID              `db:"banner_id" json:"banner_id"`
	FeatureId uuid.UUID              `db:"feature_id" json:"feature_id"`
	TagIds    []uuid.UUID            `db:"tag_ids" json:"tag_ids"`
	Content   map[string]interface{} `db:"content" json:"content"`
	IsActive  bool                   `db:"is_active" json:"is_active"`
	CreatedAt time.Time              `db:"created_at" json:"created_at"`
	UpdatedAt time.Time              `db:"updated_at" json:"updated_at"`
}

type PatchBanner struct {
	FeatureId sqlx.NullUint   `db:"feature_id" json:"feature_id"`
	TagIds    []sqlx.NullUint `db:"tag_ids" json:"tag_ids"`
	Content   map[string]interface{}
	IsActive  sql.NullBool
}

type FilterBanner struct {
	FeatureId uint `db:"feature_id" json:"featureId"`
	TagId     uint `db:"tag_id" json:"tag_id"`
	Limit     uint
	Offset    uint
}
