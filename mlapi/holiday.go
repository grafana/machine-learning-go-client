package mlapi

import (
	"bytes"
	"context"
	"encoding/json"
	"time"
)

// Holiday is a collection of time periods indicating where a time series
// will behave differently than normal.
//
// A holiday may be specified using either an iCal (in which case both
// the ICalURL and ICalTimeZone fields must be provided), or using a set
// of custom time periods directly (in which case the CustomPeriods field
// must be provided).
//
// Holidays can be linked to jobs in three ways:
// - at holiday creation time, if the job already exists, using the Jobs field
// - at job creation time, if the holiday already exists, using the Holidays field
// - using the LinkHolidaysToJob method on a Client.
type Holiday struct {
	ID string `json:"id,omitempty"`
	// Name is a human readable name for the holiday.
	Name string `json:"name"`
	// Description is a human readable description for the holiday.
	Description string `json:"description"`

	// ICalURL is the URL to an iCal file containing all occurrences
	// of the holiday.
	ICalURL *string `json:"iCalUrl,omitempty"`
	// ICalTimeZone is the timezone to use for 'All Day' events on the iCal in ICalURL,
	// if present. This is required because All Day events don't come with time zone information by default,
	// and there is no well-adhered-to standard for timezones of entire iCal files.
	ICalTimeZone *string `json:"iCalTimeZone,omitempty"`

	// CustomPeriods are holiday periods that are specified explicitly.
	CustomPeriods CustomPeriods `json:"customPeriods,omitempty"`

	// Jobs is a slice of IDs of Jobs that are using this holiday.
	// Requests may specify either IDs or names. Responses will always contain IDs.
	Jobs []string `json:"jobs"`
}

// CustomPeriods is a slice of CustomPeriods representing all occurrences
// of a holiday.
type CustomPeriods []CustomPeriod

// CustomPeriod is a single period to be included in a holiday.
type CustomPeriod struct {
	// Name is the name of this period.
	Name string `json:"name"`
	// StartTime is the (inclusive) start of this period.
	StartTime time.Time `json:"startTime"`
	// EndTime is the (exclusive) end of this period.
	EndTime time.Time `json:"endTime"`
}

// NewHoliday creates a new holiday.
func (c *Client) NewHoliday(ctx context.Context, holiday Holiday) (Holiday, error) {
	data, err := json.Marshal(holiday)
	if err != nil {
		return Holiday{}, err
	}
	result := responseWrapper[Holiday]{}
	err = c.request(ctx, "POST", "/manage/api/v1/holidays", nil, bytes.NewReader(data), &result)
	if err != nil {
		return Holiday{}, err
	}
	return result.Data, err
}

// Holidays fetches all existing holidays.
func (c *Client) Holidays(ctx context.Context) ([]Holiday, error) {
	result := responseWrapper[[]Holiday]{}
	err := c.request(ctx, "GET", "/manage/api/v1/holidays", nil, nil, &result)
	if err != nil {
		return []Holiday{}, err
	}
	return result.Data, nil
}

// Holiday fetches an existing holiday.
func (c *Client) Holiday(ctx context.Context, id string) (Holiday, error) {
	result := responseWrapper[Holiday]{}
	err := c.request(ctx, "GET", "/manage/api/v1/holidays/"+id, nil, nil, &result)
	if err != nil {
		return Holiday{}, err
	}
	return result.Data, err
}

// UpdateHoliday updates an existing holiday.
func (c *Client) UpdateHoliday(ctx context.Context, holiday Holiday) (Holiday, error) {
	id := holiday.ID
	// Clear the ID before sending otherwise validation fails.
	holiday.ID = ""
	data, err := json.Marshal(holiday)
	if err != nil {
		return Holiday{}, err
	}

	result := responseWrapper[Holiday]{}
	err = c.request(ctx, "POST", "/manage/api/v1/holidays/"+id, nil, bytes.NewReader(data), &result)
	if err != nil {
		return Holiday{}, err
	}
	return result.Data, err
}

// DeleteHoliday deletes an existing holiday.
func (c *Client) DeleteHoliday(ctx context.Context, id string) error {
	return c.request(ctx, "DELETE", "/manage/api/v1/holidays/"+id, nil, nil, nil)
}
