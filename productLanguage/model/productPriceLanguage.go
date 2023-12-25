package model

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
