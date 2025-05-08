package models

// ProductAttribute represents an attribute with its possible values
type ProductAttribute struct {
	AttributeID int64
	Name        string
	Options     []AttributeOption
}

// AttributeOption represents a possible value for an attribute
type AttributeOption struct {
	OptionID int64
	Value    string
}
