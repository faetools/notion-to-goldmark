package notion

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
