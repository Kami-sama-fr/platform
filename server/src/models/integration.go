package models

import (
	"time"

	"gorm.io/datatypes"
)

type Integration struct {
	Common
	Name           string         `gorm:"column:name;type:text;not null" json:"name"`
	Provider       string         `gorm:"column:provider;type:text;not null" json:"provider"`
	Type           string         `gorm:"column:type;type:text;not null" json:"type"`
	Status         string         `gorm:"column:status;type:text;not null;default:'inactive'" json:"status"`
	Config         datatypes.JSON `gorm:"column:config;type:jsonb" json:"config,omitempty"`
	Secrets        datatypes.JSON `gorm:"column:secrets;type:jsonb" json:"-"`
	Metadata       datatypes.JSON `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
	LastSyncedAt   *time.Time     `gorm:"column:last_synced_at" json:"lastSyncedAt,omitempty"`
	LastError      string         `gorm:"column:last_error;type:text" json:"lastError,omitempty"`
}

func (Integration) TableName() string { return "integrations" }

type ApiKey struct {
	Common
	Name        string     `gorm:"column:name;type:text;not null" json:"name"`
	KeyHash     string     `gorm:"column:key_hash;type:text;uniqueIndex;not null" json:"-"`
	KeyPrefix   string     `gorm:"column:key_prefix;type:text;not null" json:"keyPrefix"`
	UserID      string     `gorm:"column:user_id;type:text;index;not null" json:"userId"`
	Permissions datatypes.JSON `gorm:"column:permissions;type:jsonb" json:"permissions,omitempty"`
	LastUsedAt  *time.Time `gorm:"column:last_used_at" json:"lastUsedAt,omitempty"`
	ExpiresAt   *time.Time `gorm:"column:expires_at" json:"expiresAt,omitempty"`
	RevokedAt   *time.Time `gorm:"column:revoked_at" json:"revokedAt,omitempty"`
	UsageCount  int64      `gorm:"column:usage_count;default:0" json:"usageCount"`
}

func (ApiKey) TableName() string { return "api_keys" }

type DomainConfig struct {
	Common
	Domain     string     `gorm:"column:domain;type:text;uniqueIndex;not null" json:"domain"`
	IsVerified bool       `gorm:"column:is_verified;default:false" json:"isVerified"`
	IsSSL      bool       `gorm:"column:is_ssl;default:false" json:"isSsl"`
	VerifiedAt *time.Time `gorm:"column:verified_at" json:"verifiedAt,omitempty"`
	SSLAt      *time.Time `gorm:"column:ssl_at" json:"sslAt,omitempty"`
}

func (DomainConfig) TableName() string { return "domain_configs" }
