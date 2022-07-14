package notion

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	layoutDate    = "2006-01-02"
	lenLayoutDate = len(layoutDate)
)

// tmpDate is used to unmarshall into this so we can properly adjust the times to be in the right time zone.
type tmpDate struct {
	End      *string `json:"end,omitempty"`
	Start    string  `json:"start"`
	TimeZone *string `json:"time_zone,omitempty"`
}

// UnmarshalJSON fulfils json.Unmarshaller.
func (d *Date) UnmarshalJSON(b []byte) error {
	tmp := &tmpDate{}

	err := json.Unmarshal(b, tmp)
	if err != nil {
		return err
	}

	// Time zone information for start and end. Possible values are extracted from the IANA database and they are based on the time zones from Moment.js.
	//
	// When time zone is provided, start and end should not have any UTC offset. In addition, when time zone is provided, start and end cannot be dates without time information.
	//
	// If null, time zone information will be contained in UTC offsets in start and end.
	d.TimeZone = tmp.TimeZone

	loc := time.UTC
	if tmp.TimeZone != nil {
		loc, err = time.LoadLocation(*d.TimeZone)
		if err != nil {
			return err
		}
	}

	d.Start, err = parseTimeOrDate(tmp.Start, loc)
	if err != nil {
		return err
	}

	if tmp.End != nil {
		end, err := parseTimeOrDate(*tmp.End, loc)
		if err != nil {
			return err
		}

		d.End = &end
	}

	return nil
}

// MarshalJSON fulfils json.Marshaler.
func (d Date) MarshalJSON() ([]byte, error) {
	tmp := tmpDate{
		Start:    formatTime(d.Start, true),
		TimeZone: d.TimeZone,
	}

	if d.End != nil {
		end := formatTime(*d.End, true)
		tmp.End = &end
	}

	return json.Marshal(tmp)
}

func parseTimeOrDate(ts string, loc *time.Location) (time.Time, error) {
	if lenLayoutDate == len(ts) {
		return time.ParseInLocation(layoutDate, ts, loc)
	}

	t, err := time.ParseInLocation(time.RFC3339, ts, loc)
	if err != nil || loc == nil {
		return t, err
	}

	return t.In(loc), nil
}

func (d Date) String() string {
	if d.End == nil {
		return formatTime(d.Start, false)
	}

	return fmt.Sprintf("%s - %s",
		formatTime(d.Start, false), formatTime(*d.End, false))
}

func formatTime(t time.Time, inUTC bool) string {
	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
		return t.Format(layoutDate)
	}

	if inUTC {
		t = t.In(time.UTC)
	}

	return t.Format(time.RFC3339)
}
