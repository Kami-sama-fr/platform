package models

import "time"

type Contact struct {
	Common
	UserID      string     `gorm:"column:user_id;type:text;index;not null" json:"userId"`
	Email       string     `gorm:"column:email;type:text;not null" json:"email"`
	DisplayName string     `gorm:"column:display_name;type:text;not null" json:"displayName"`
	AvatarUrl   *string    `gorm:"column:avatar_url;type:text" json:"avatarUrl,omitempty"`
	Phone       *string    `gorm:"column:phone;type:text" json:"phone,omitempty"`
	Company     *string    `gorm:"column:company;type:text" json:"company,omitempty"`
	JobTitle    *string    `gorm:"column:job_title;type:text" json:"jobTitle,omitempty"`
	IsFavorite  bool       `gorm:"column:is_favorite;not null;default:false" json:"isFavorite"`
	Notes       *string    `gorm:"column:notes;type:text" json:"notes,omitempty"`
	Tags        *string    `gorm:"column:tags;type:text" json:"tags,omitempty"`
}

func (Contact) TableName() string { return "contacts" }

type ContactGroup struct {
	Common
	WorkspaceId string  `gorm:"column:workspace_id;type:text;index;not null" json:"workspaceId"`
	Name        string  `gorm:"column:name;type:text;not null" json:"name"`
	Description *string `gorm:"column:description;type:text" json:"description,omitempty"`
	Color       *string `gorm:"column:color;type:text" json:"color,omitempty"`
}

func (ContactGroup) TableName() string { return "contact_groups" }

type FamilyMember struct {
	Common
	UserID       string     `gorm:"column:user_id;type:text;index;not null" json:"userId"`
	Email        string     `gorm:"column:email;type:text;not null" json:"email"`
	DisplayName  string     `gorm:"column:display_name;type:text;not null" json:"displayName"`
	Role         string     `gorm:"column:role;type:text;not null" json:"role"`
	Status       string     `gorm:"column:status;type:text;not null" json:"status"`
	InvitedBy    string     `gorm:"column:invited_by;type:text;not null" json:"invitedBy"`
	InvitedAt    time.Time  `gorm:"column:invited_at;not null" json:"invitedAt"`
	JoinedAt     *time.Time `gorm:"column:joined_at" json:"joinedAt,omitempty"`
	AvatarUrl    *string    `gorm:"column:avatar_url;type:text" json:"avatarUrl,omitempty"`
	Phone        *string    `gorm:"column:phone;type:text" json:"phone,omitempty"`
	Relationship *string    `gorm:"column:relationship;type:text" json:"relationship,omitempty"`
}

func (FamilyMember) TableName() string { return "family_members" }

type Invitation struct {
	Common
	Email      string     `gorm:"column:email;type:text;not null" json:"email"`
	Token      string     `gorm:"column:token;type:text;uniqueIndex;not null" json:"-"`
	Type       string     `gorm:"column:type;type:text;not null" json:"type"`
	Role       string     `gorm:"column:role;type:text;not null" json:"role"`
	Message    *string    `gorm:"column:message;type:text" json:"message,omitempty"`
	Status     string     `gorm:"column:status;type:text;not null;default:'pending'" json:"status"`
	InvitedBy  string     `gorm:"column:invited_by;type:text;not null" json:"invitedBy"`
	TargetId   *string    `gorm:"column:target_id;type:text" json:"targetId,omitempty"`
	AcceptedAt *time.Time `gorm:"column:accepted_at" json:"acceptedAt,omitempty"`
	ExpiresAt  time.Time  `gorm:"column:expires_at;not null" json:"expiresAt"`
	Metadata   *string    `gorm:"column:metadata;type:text" json:"metadata,omitempty"`
}

func (Invitation) TableName() string { return "invitations" }
