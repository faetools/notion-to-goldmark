package notion

import (
	"encoding/json"
)

type (
	rollup Rollup

	rollupNumber struct {
		Type     RollupType `json:"type"`
		Number   *float32   `json:"number"`
		Function string     `json:"function"`
	}

	rollupDate struct {
		Type     RollupType `json:"type"`
		Date     *Date      `json:"date"`
		Function string     `json:"function"`
	}

	rollupString struct {
		Type     RollupType `json:"type"`
		String   *string    `json:"string"`
		Function string     `json:"function"`
	}

	rollupArrayItem RollupArrayItem

	rollupArrayItemDate struct {
		Type RollupArrayItemType `json:"type"`
		Date *Date               `json:"date"`
	}
)

// MarshalJSON fulfils json.Marshaler.
func (r Rollup) MarshalJSON() ([]byte, error) {
	switch r.Type {
	case RollupTypeNumber:
		return json.Marshal(rollupNumber{
			Type:     r.Type,
			Number:   r.Number,
			Function: r.Function,
		})
	case RollupTypeDate:
		return json.Marshal(rollupDate{
			Type:     r.Type,
			Date:     r.Date,
			Function: r.Function,
		})
	case RollupTypeString:
		return json.Marshal(rollupString{
			Type:     r.Type,
			String:   r.String,
			Function: r.Function,
		})
	default:
		return json.Marshal(rollup(r))
	}
}

func (r RollupArrayItem) MarshalJSON() ([]byte, error) {
	switch r.Type {
	case RollupArrayItemTypeDate:
		return json.Marshal(rollupArrayItemDate{
			Type: r.Type,
			Date: r.Date,
		})
	default:
		return json.Marshal(rollupArrayItem(r))
	}
}
