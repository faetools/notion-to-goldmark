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

	// 	Date   *time.Time `json:"date"`
	// String *string    `json:"string"`
)

// MarshalJSON fulfils json.Marshaler.
func (v Rollup) MarshalJSON() ([]byte, error) {
	switch v.Type {
	case RollupTypeNumber:
		return json.Marshal(rollupNumber{
			Type:     v.Type,
			Number:   v.Number,
			Function: v.Function,
		})
	default:
		return json.Marshal(rollup(v))
	}
}
