package notion

import (
	"encoding/json"
)

type (
	propertyValue PropertyValue

	date struct {
		ID   string `json:"id,omitempty"`
		Type string `json:"type"`
		Date *Date  `json:"date"`
	}

	number struct {
		ID     string   `json:"id,omitempty"`
		Type   string   `json:"type"`
		Number *float32 `json:"number"`
	}

	propertyValueSelect struct {
		ID     string       `json:"id"`
		Type   PropertyType `json:"type"`
		Select *SelectValue `json:"select"`
	}

	propertyValueStatus struct {
		ID     string       `json:"id"`
		Type   PropertyType `json:"type"`
		Status *SelectValue `json:"status"`
	}

	propertyValueURL struct {
		ID   string       `json:"id"`
		Type PropertyType `json:"type"`
		URL  *string      `json:"url"`
	}

	propertyValuePhone struct {
		ID          string       `json:"id"`
		Type        PropertyType `json:"type"`
		PhoneNumber *string      `json:"phone_number"`
	}

	propertyValueEmail struct {
		ID    string       `json:"id"`
		Type  PropertyType `json:"type"`
		Email *string      `json:"email"`
	}
)

// MarshalJSON fulfils json.Marshaler.
func (v PropertyValue) MarshalJSON() ([]byte, error) {
	switch v.Type {
	case PropertyTypeDate:
		return json.Marshal(date{
			ID:   v.Id,
			Type: string(v.Type),
			Date: v.Date,
		})
	case PropertyTypeNumber:
		return json.Marshal(number{
			ID:     v.Id,
			Type:   string(v.Type),
			Number: v.Number,
		})
	case PropertyTypeSelect:
		return json.Marshal(propertyValueSelect{
			ID:     v.Id,
			Type:   v.Type,
			Select: v.Select,
		})
	case PropertyTypeStatus:
		return json.Marshal(propertyValueStatus{
			ID:     v.Id,
			Type:   v.Type,
			Status: v.Status,
		})
	case PropertyTypeUrl:
		return json.Marshal(propertyValueURL{
			ID:   v.Id,
			Type: v.Type,
			URL:  v.Url,
		})
	case PropertyTypePhoneNumber:
		return json.Marshal(propertyValuePhone{
			ID:          v.Id,
			Type:        v.Type,
			PhoneNumber: v.PhoneNumber,
		})
	case PropertyTypeEmail:
		return json.Marshal(propertyValueEmail{
			ID:    v.Id,
			Type:  v.Type,
			Email: v.Email,
		})
	default:
		return json.Marshal(propertyValue(v))
	}
}

// GetMultiSelect returns the multiselect value.
func (v PropertyValue) GetMultiSelect() PropertyOptions {
	if v.MultiSelect == nil {
		return PropertyOptions{}
	}

	return *v.MultiSelect
}

// GetCheckbox returns the checkbox value.
func (v PropertyValue) GetCheckbox() bool {
	return v.Checkbox != nil && *v.Checkbox
}

// GetDate returns the date value.
func (v PropertyValue) GetDate() Date {
	if v.Date == nil {
		return Date{}
	}

	return *v.Date
}

// GetFiles returns the files value.
func (v PropertyValue) GetFiles() Files {
	if v.Files == nil {
		return nil
	}

	return *v.Files
}

// GetNumber returns the number value.
func (v PropertyValue) GetNumber() float32 {
	if v.Number == nil {
		return 0
	}

	return *v.Number
}

// GetRichText returns the rich text value.
func (v PropertyValue) GetRichText() RichTexts {
	if v.RichText == nil {
		return nil
	}

	return *v.RichText
}

// GetSelect returns the value that was selected.
func (v PropertyValue) GetSelect() SelectValue {
	if v.Select == nil {
		return SelectValue{}
	}

	return *v.Select
}

// GetRelation returns the relation value.
func (v PropertyValue) GetRelation() References {
	if v.Relation == nil {
		return nil
	}

	return *v.Relation
}

// GetTitle returns the title value.
func (v PropertyValue) GetTitle() RichTexts {
	if v.Title == nil {
		return nil
	}

	return *v.Title
}
