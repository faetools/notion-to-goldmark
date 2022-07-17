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
	End      *string `json:"end"`
	Start    string  `json:"start"`
	TimeZone *string `json:"time_zone"`
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

	loc, err := loadLocation(tmp.TimeZone)
	if err != nil {
		return err
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

func loadLocation(tz *string) (*time.Location, error) {
	if tz == nil {
		return nil, nil
	}

	return time.LoadLocation(*tz)
}

// MarshalJSON fulfils json.Marshaler.
func (d Date) MarshalJSON() ([]byte, error) {
	loc, err := loadLocation(d.TimeZone)
	if err != nil {
		return nil, err
	}

	tmp := tmpDate{
		Start:    formatTime(d.Start, loc),
		TimeZone: d.TimeZone,
	}

	if d.End != nil {
		endStr := formatTime(*d.End, loc)
		tmp.End = &endStr
	}

	return json.Marshal(tmp)
}

func parseTimeOrDate(ts string, loc *time.Location) (time.Time, error) {
	if lenLayoutDate == len(ts) {
		return time.Parse(layoutDate, ts)
	}

	if loc == nil {
		return time.Parse(time.RFC3339, ts)
	}

	t, err := time.ParseInLocation(time.RFC3339, ts, loc)
	if err != nil || loc == nil {
		return t, err
	}

	sec := t.Second()

	t = t.In(loc)

	// adjust seconds offset
	// see issue: https://github.com/golang/go/issues/53919
	// and fix: https://github.com/golang/go/pull/53920
	diff := sec - t.Second()
	if diff > 0 {
		diff -= 60
	}

	return t.Add(time.Duration(diff) * time.Second), nil
}

func (d Date) String() string {
	if d.End == nil {
		return formatTime(d.Start, nil)
	}

	return fmt.Sprintf("%s - %s",
		formatTime(d.Start, nil),
		formatTime(*d.End, nil))
}

func formatTime(t time.Time, loc *time.Location) string {
	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
		return t.Format(layoutDate)
	}

	if loc != nil {
		t = t.In(loc)
	}

	return t.Format(time.RFC3339)
}
