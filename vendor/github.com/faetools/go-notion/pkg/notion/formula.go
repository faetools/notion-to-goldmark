package notion

import "encoding/json"

type formula Formula

// MarshalJSON fulfils json.Marshaler.
// Its purpose is to write `null` for empty fields that correspond to the type
// but not write it for fields that don't correspond to the type.
func (v Formula) MarshalJSON() ([]byte, error) {
	switch v.Type {
	case FormulaTypeDate:
		return json.Marshal(date{
			Type: string(v.Type),
			Date: v.Date,
		})
	case FormulaTypeNumber:
		return json.Marshal(number{
			Type:   string(v.Type),
			Number: v.Number,
		})
	default:
		return json.Marshal(formula(v))
	}
}
