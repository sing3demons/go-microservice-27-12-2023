package model

type Category struct {
	Type            string     `json:"@type,omitempty" validate:"required" bson:"@type,omitempty"`
	ID              string     `json:"id" validate:"required" bson:"id"`
	Href            string     `json:"href,omitempty" bson:"href,omitempty"`
	Name            string     `json:"name,omitempty" bson:"name,omitempty"`
	Version         string     `json:"version,omitempty" bson:"version,omitempty"`
	LastUpdate      string     `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	ValidFor        *ValidFor  `json:"validFor,omitempty" bson:"validFor,omitempty"`
	LifecycleStatus string     `json:"lifecycleStatus,omitempty" bson:"lifecycleStatus,omitempty"`
	Products        []Products `json:"products,omitempty" bson:"products,omitempty"`
}
