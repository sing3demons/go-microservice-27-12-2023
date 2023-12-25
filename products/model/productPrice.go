package model

type ProductPrice struct {
	Type               string                 `json:"@type" validate:"required" bson:"@type"`
	ID                 string                 `json:"id" validate:"required" bson:"id" binding:"required"`
	Href               string                 `json:"href,omitempty" bson:"href,omitempty"`
	Name               string                 `json:"name,omitempty" bson:"name,omitempty"`
	Version            string                 `json:"version,omitempty" bson:"version,omitempty"`
	LastUpdate         string                 `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	ValidFor           *ValidFor              `json:"validFor,omitempty" bson:"validFor,omitempty"`
	Price              float64                `json:"price,omitempty" bson:"price,omitempty"`
	Tax                float64                `json:"tax,omitempty" bson:"tax,omitempty"`
	PopRelationship    []PopRelationship      `json:"popRelationship,omitempty" bson:"popRelationship,omitempty"`
	SupportingLanguage []ProductPriceLanguage `json:"supportingLanguage,omitempty" bson:"supportingLanguage,omitempty"`
}

type PopRelationship struct {
	ID   string `json:"id,omitempty" bson:"id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

type ProductPriceLanguage struct {
	Type         string    `json:"@type" validate:"required" bson:"@type"`
	ID           string    `json:"id" bson:"id,omitempty"`
	Href         string    `json:"href,omitempty" bson:"href,omitempty"`
	LanguageCode string    `json:"languageCode,omitempty" bson:"languageCode,omitempty"`
	Name         string    `json:"name,omitempty" bson:"name,omitempty"`
	Version      string    `json:"version,omitempty" bson:"version,omitempty"`
	LastUpdate   string    `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	ValidFor     *ValidFor `json:"validFor,omitempty" bson:"validFor,omitempty"`
	Price        Price     `json:"price,omitempty" bson:"price,omitempty"`
	Tax          Tax       `json:"tax,omitempty" bson:"tax,omitempty"`
}

type Tax struct {
	Unit  string  `json:"unit,omitempty" bson:"unit,omitempty"`
	Value float64 `json:"value,omitempty" bson:"value,omitempty"`
}

type Price struct {
	Unit  string  `json:"unit,omitempty" bson:"unit,omitempty"`
	Value float64 `json:"value,omitempty" bson:"value,omitempty"`
}
