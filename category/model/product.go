package model

import "time"

type Products struct {
	Type               string            `json:"@type,omitempty" bson:"@type,omitempty"`
	ID                 string            `json:"id" bson:"id,omitempty"`
	Name               string            `json:"name,omitempty" bson:"name,omitempty"`
	Href               string            `json:"href,omitempty" bson:"href,omitempty"`
	LifecycleStatus    string            `json:"lifecycleStatus,omitempty" bson:"lifecycleStatus,omitempty"`
	Version            string            `json:"version,omitempty" bson:"version,omitempty"`
	LastUpdate         string            `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	ValidFor           *ValidFor         `json:"validFor,omitempty" bson:"validFor,omitempty"`
	Category           []Category        `json:"category,omitempty" bson:"category,omitempty"`
	SupportingLanguage []ProductLanguage `json:"supportingLanguage,omitempty" bson:"supportingLanguage,omitempty"`
}

type ProductLanguage struct {
	Type         string       `json:"@type" validate:"required" bson:"@type"`
	ID           string       `json:"id,omitempty" validate:"required" bson:"id,omitempty" binding:"required"`
	Href         string       `json:"href,omitempty" bson:"href,omitempty"`
	LanguageCode string       `json:"languageCode,omitempty" bson:"languageCode,omitempty" binding:"required"`
	Name         string       `json:"name,omitempty" bson:"name,omitempty"`
	Version      string       `json:"version,omitempty" bson:"version,omitempty"`
	LastUpdate   string       `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	ValidFor     *ValidFor    `json:"validFor,omitempty" bson:"validFor,omitempty"`
	Attachment   []Attachment `json:"attachment,omitempty" bson:"attachment,omitempty"`
}

type Attachment struct {
	Type        string      `json:"@type" validate:"required" bson:"@type"`
	ID          string      `json:"id" validate:"required" bson:"id" binding:"required"`
	Href        string      `json:"href,omitempty" bson:"href,omitempty"`
	Description string      `json:"description,omitempty" bson:"description,omitempty"`
	MimeType    string      `json:"mimeType,omitempty" bson:"mimeType,omitempty"`
	Name        string      `json:"name,omitempty" bson:"name,omitempty"`
	Url         string      `json:"url,omitempty" bson:"url,omitempty"`
	ValidFor    *ValidFor   `json:"validFor,omitempty" bson:"validFor,omitempty"`
	RedirectUrl string      `json:"redirectUrl,omitempty" bson:"redirectUrl,omitempty"`
	DisplayInfo DisplayInfo `json:"displayInfo,omitempty" bson:"displayInfo,omitempty"`
}

type DisplayInfo struct {
	ValueType string   `json:"valueType,omitempty" bson:"valueType,omitempty"`
	Value     []string `json:"value,omitempty" bson:"value,omitempty"`
}

type ValidFor struct {
	StartDateTime time.Time `json:"startDateTime,omitempty" bson:"startDateTime,omitempty"`
	EndDateTime   time.Time `json:"endDateTime,omitempty" bson:"endDateTime,omitempty"`
}
