package mlapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHoliday(t *testing.T) {
	id := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	holiday := Holiday{}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/manage/api/v1/holidays" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		parsedHoliday := Holiday{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&parsedHoliday)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		assert.Equal(t, holiday, parsedHoliday)
		parsedHoliday.ID = id
		enc := json.NewEncoder(w)
		_ = enc.Encode(responseWrapper[Holiday]{Data: parsedHoliday})
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	returnedHoliday, err := c.NewHoliday(ctx, holiday)
	require.NoError(t, err)
	assert.Equal(t, id, returnedHoliday.ID)
}

func TestHolidayICal(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/manage/api/v1/holidays/8b154ff8-3d64-4b79-8b26-02b4baeb44e4" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		_, err := w.Write([]byte(
			`{"status":"success","data":{"id":"8b154ff8-3d64-4b79-8b26-02b4baeb44e4","created":"2023-01-04T10:09:47.254Z","modified":"2023-01-04T10:09:47.254Z","tenantId":"0","createdBy":null,"modifiedBy":null,"name":"Test Holiday","description":null,"iCalUrl":"https://calendar.google.com/calendar/ical/en.uk%23holiday%40group.v.calendar.google.com/public/basic.ics","iCalTimeZone":"Europe/London","customPeriods":null,"resolvedDates":[{"name":"New Year's Day","startTime":"2022-01-01T00:00:00Z","endTime":"2022-01-02T00:00:00Z"},{"name":"New Year's Day observed","startTime":"2022-01-03T00:00:00Z","endTime":"2022-01-04T00:00:00Z"},{"name":"2nd January (substitute day) (Scotland)","startTime":"2022-01-04T00:00:00Z","endTime":"2022-01-05T00:00:00Z"},{"name":"Twelfth Night","startTime":"2022-01-05T00:00:00Z","endTime":"2022-01-06T00:00:00Z"},{"name":"Valentine's Day","startTime":"2022-02-14T00:00:00Z","endTime":"2022-02-15T00:00:00Z"},{"name":"Carnival / Shrove Tuesday / Pancake Day","startTime":"2022-03-01T00:00:00Z","endTime":"2022-03-02T00:00:00Z"},{"name":"St Patrick's Day (Northern Ireland)","startTime":"2022-03-17T00:00:00Z","endTime":"2022-03-18T00:00:00Z"},{"name":"Mother's Day","startTime":"2022-03-27T00:00:00Z","endTime":"2022-03-27T23:00:00Z"},{"name":"Daylight Saving Time starts","startTime":"2022-03-27T00:00:00Z","endTime":"2022-03-27T23:00:00Z"},{"name":"Good Friday","startTime":"2022-04-14T23:00:00Z","endTime":"2022-04-15T23:00:00Z"},{"name":"Easter Sunday","startTime":"2022-04-16T23:00:00Z","endTime":"2022-04-17T23:00:00Z"},{"name":"Easter Monday (regional holiday)","startTime":"2022-04-17T23:00:00Z","endTime":"2022-04-18T23:00:00Z"},{"name":"St. George's Day","startTime":"2022-04-22T23:00:00Z","endTime":"2022-04-23T23:00:00Z"},{"name":"Early May Bank Holiday","startTime":"2022-05-01T23:00:00Z","endTime":"2022-05-02T23:00:00Z"},{"name":"Spring Bank Holiday","startTime":"2022-06-01T23:00:00Z","endTime":"2022-06-02T23:00:00Z"},{"name":"Queen Elizabeth II's Platinum Jubilee","startTime":"2022-06-02T23:00:00Z","endTime":"2022-06-03T23:00:00Z"},{"name":"Queen Elizabeth II's Birthday","startTime":"2022-06-10T23:00:00Z","endTime":"2022-06-11T23:00:00Z"},{"name":"Father's Day","startTime":"2022-06-18T23:00:00Z","endTime":"2022-06-19T23:00:00Z"},{"name":"Battle of the Boyne (Northern Ireland)","startTime":"2022-07-11T23:00:00Z","endTime":"2022-07-12T23:00:00Z"},{"name":"Summer Bank Holiday (Scotland)","startTime":"2022-07-31T23:00:00Z","endTime":"2022-08-01T23:00:00Z"},{"name":"Summer Bank Holiday (regional holiday)","startTime":"2022-08-28T23:00:00Z","endTime":"2022-08-29T23:00:00Z"},{"name":"State Funeral of Queen Elizabeth II","startTime":"2022-09-18T23:00:00Z","endTime":"2022-09-19T23:00:00Z"},{"name":"Daylight Saving Time ends","startTime":"2022-10-29T23:00:00Z","endTime":"2022-10-31T00:00:00Z"},{"name":"Halloween","startTime":"2022-10-31T00:00:00Z","endTime":"2022-11-01T00:00:00Z"},{"name":"Guy Fawkes Day","startTime":"2022-11-05T00:00:00Z","endTime":"2022-11-06T00:00:00Z"},{"name":"Remembrance Sunday","startTime":"2022-11-13T00:00:00Z","endTime":"2022-11-14T00:00:00Z"},{"name":"St Andrew's Day (Scotland)","startTime":"2022-11-30T00:00:00Z","endTime":"2022-12-01T00:00:00Z"},{"name":"Christmas Eve","startTime":"2022-12-24T00:00:00Z","endTime":"2022-12-25T00:00:00Z"},{"name":"Christmas Day","startTime":"2022-12-25T00:00:00Z","endTime":"2022-12-26T00:00:00Z"},{"name":"Boxing Day","startTime":"2022-12-26T00:00:00Z","endTime":"2022-12-27T00:00:00Z"},{"name":"Substitute Bank Holiday for Christmas Day","startTime":"2022-12-27T00:00:00Z","endTime":"2022-12-28T00:00:00Z"},{"name":"New Year's Eve","startTime":"2022-12-31T00:00:00Z","endTime":"2023-01-01T00:00:00Z"},{"name":"New Year's Day","startTime":"2023-01-01T00:00:00Z","endTime":"2023-01-02T00:00:00Z"},{"name":"New Year's Day observed","startTime":"2023-01-02T00:00:00Z","endTime":"2023-01-03T00:00:00Z"},{"name":"2nd January (substitute day) (Scotland)","startTime":"2023-01-03T00:00:00Z","endTime":"2023-01-04T00:00:00Z"},{"name":"Twelfth Night","startTime":"2023-01-05T00:00:00Z","endTime":"2023-01-06T00:00:00Z"},{"name":"Valentine's Day","startTime":"2023-02-14T00:00:00Z","endTime":"2023-02-15T00:00:00Z"},{"name":"Carnival / Shrove Tuesday / Pancake Day","startTime":"2023-02-21T00:00:00Z","endTime":"2023-02-22T00:00:00Z"},{"name":"St Patrick's Day (Northern Ireland)","startTime":"2023-03-17T00:00:00Z","endTime":"2023-03-18T00:00:00Z"},{"name":"Mother's Day","startTime":"2023-03-19T00:00:00Z","endTime":"2023-03-20T00:00:00Z"},{"name":"Daylight Saving Time starts","startTime":"2023-03-26T00:00:00Z","endTime":"2023-03-26T23:00:00Z"},{"name":"Good Friday","startTime":"2023-04-06T23:00:00Z","endTime":"2023-04-07T23:00:00Z"},{"name":"Easter Sunday","startTime":"2023-04-08T23:00:00Z","endTime":"2023-04-09T23:00:00Z"},{"name":"Easter Monday (regional holiday)","startTime":"2023-04-09T23:00:00Z","endTime":"2023-04-10T23:00:00Z"},{"name":"St. George's Day","startTime":"2023-04-22T23:00:00Z","endTime":"2023-04-23T23:00:00Z"},{"name":"Early May Bank Holiday","startTime":"2023-04-30T23:00:00Z","endTime":"2023-05-01T23:00:00Z"},{"name":"The Coronation of King Charles III","startTime":"2023-05-05T23:00:00Z","endTime":"2023-05-06T23:00:00Z"},{"name":"Bank Holiday for the Coronation of King Charles III","startTime":"2023-05-07T23:00:00Z","endTime":"2023-05-08T23:00:00Z"},{"name":"Spring Bank Holiday","startTime":"2023-05-28T23:00:00Z","endTime":"2023-05-29T23:00:00Z"},{"name":"King's Birthday","startTime":"2023-06-16T23:00:00Z","endTime":"2023-06-17T23:00:00Z"},{"name":"Father's Day","startTime":"2023-06-17T23:00:00Z","endTime":"2023-06-18T23:00:00Z"},{"name":"Battle of the Boyne (Northern Ireland)","startTime":"2023-07-11T23:00:00Z","endTime":"2023-07-12T23:00:00Z"},{"name":"Summer Bank Holiday (Scotland)","startTime":"2023-08-06T23:00:00Z","endTime":"2023-08-07T23:00:00Z"},{"name":"Summer Bank Holiday (regional holiday)","startTime":"2023-08-27T23:00:00Z","endTime":"2023-08-28T23:00:00Z"},{"name":"Daylight Saving Time ends","startTime":"2023-10-28T23:00:00Z","endTime":"2023-10-30T00:00:00Z"},{"name":"Halloween","startTime":"2023-10-31T00:00:00Z","endTime":"2023-11-01T00:00:00Z"},{"name":"Guy Fawkes Day","startTime":"2023-11-05T00:00:00Z","endTime":"2023-11-06T00:00:00Z"},{"name":"Remembrance Sunday","startTime":"2023-11-12T00:00:00Z","endTime":"2023-11-13T00:00:00Z"},{"name":"St Andrew's Day (Scotland)","startTime":"2023-11-30T00:00:00Z","endTime":"2023-12-01T00:00:00Z"},{"name":"Christmas Eve","startTime":"2023-12-24T00:00:00Z","endTime":"2023-12-25T00:00:00Z"},{"name":"Christmas Day","startTime":"2023-12-25T00:00:00Z","endTime":"2023-12-26T00:00:00Z"},{"name":"Boxing Day","startTime":"2023-12-26T00:00:00Z","endTime":"2023-12-27T00:00:00Z"},{"name":"New Year's Eve","startTime":"2023-12-31T00:00:00Z","endTime":"2024-01-01T00:00:00Z"},{"name":"New Year's Day","startTime":"2024-01-01T00:00:00Z","endTime":"2024-01-02T00:00:00Z"},{"name":"2nd January (Scotland)","startTime":"2024-01-02T00:00:00Z","endTime":"2024-01-03T00:00:00Z"},{"name":"Twelfth Night","startTime":"2024-01-05T00:00:00Z","endTime":"2024-01-06T00:00:00Z"},{"name":"Carnival / Shrove Tuesday / Pancake Day","startTime":"2024-02-13T00:00:00Z","endTime":"2024-02-14T00:00:00Z"},{"name":"Valentine's Day","startTime":"2024-02-14T00:00:00Z","endTime":"2024-02-15T00:00:00Z"},{"name":"Mother's Day","startTime":"2024-03-10T00:00:00Z","endTime":"2024-03-11T00:00:00Z"},{"name":"St Patrick's Day (Northern Ireland)","startTime":"2024-03-17T00:00:00Z","endTime":"2024-03-18T00:00:00Z"},{"name":"Day off for St Patrick's Day (Northern Ireland)","startTime":"2024-03-18T00:00:00Z","endTime":"2024-03-19T00:00:00Z"},{"name":"Good Friday","startTime":"2024-03-29T00:00:00Z","endTime":"2024-03-30T00:00:00Z"},{"name":"Daylight Saving Time starts","startTime":"2024-03-31T00:00:00Z","endTime":"2024-03-31T23:00:00Z"},{"name":"Easter Sunday","startTime":"2024-03-31T00:00:00Z","endTime":"2024-03-31T23:00:00Z"},{"name":"Easter Monday (regional holiday)","startTime":"2024-03-31T23:00:00Z","endTime":"2024-04-01T23:00:00Z"},{"name":"St. George's Day","startTime":"2024-04-22T23:00:00Z","endTime":"2024-04-23T23:00:00Z"},{"name":"Early May Bank Holiday","startTime":"2024-05-05T23:00:00Z","endTime":"2024-05-06T23:00:00Z"},{"name":"Spring Bank Holiday","startTime":"2024-05-26T23:00:00Z","endTime":"2024-05-27T23:00:00Z"},{"name":"King's Birthday","startTime":"2024-06-14T23:00:00Z","endTime":"2024-06-15T23:00:00Z"},{"name":"Father's Day","startTime":"2024-06-15T23:00:00Z","endTime":"2024-06-16T23:00:00Z"},{"name":"Battle of the Boyne (Northern Ireland)","startTime":"2024-07-11T23:00:00Z","endTime":"2024-07-12T23:00:00Z"},{"name":"Summer Bank Holiday (Scotland)","startTime":"2024-08-04T23:00:00Z","endTime":"2024-08-05T23:00:00Z"},{"name":"Summer Bank Holiday (regional holiday)","startTime":"2024-08-25T23:00:00Z","endTime":"2024-08-26T23:00:00Z"},{"name":"Daylight Saving Time ends","startTime":"2024-10-26T23:00:00Z","endTime":"2024-10-28T00:00:00Z"},{"name":"Halloween","startTime":"2024-10-31T00:00:00Z","endTime":"2024-11-01T00:00:00Z"},{"name":"Guy Fawkes Day","startTime":"2024-11-05T00:00:00Z","endTime":"2024-11-06T00:00:00Z"},{"name":"Remembrance Sunday","startTime":"2024-11-10T00:00:00Z","endTime":"2024-11-11T00:00:00Z"},{"name":"St Andrew's Day (Scotland)","startTime":"2024-11-30T00:00:00Z","endTime":"2024-12-01T00:00:00Z"},{"name":"St Andrew's Day observed (Scotland)","startTime":"2024-12-02T00:00:00Z","endTime":"2024-12-03T00:00:00Z"},{"name":"Christmas Eve","startTime":"2024-12-24T00:00:00Z","endTime":"2024-12-25T00:00:00Z"},{"name":"Christmas Day","startTime":"2024-12-25T00:00:00Z","endTime":"2024-12-26T00:00:00Z"},{"name":"Boxing Day","startTime":"2024-12-26T00:00:00Z","endTime":"2024-12-27T00:00:00Z"},{"name":"New Year's Eve","startTime":"2024-12-31T00:00:00Z","endTime":"2025-01-01T00:00:00Z"}],"jobs":[]}}`,
		))
		require.NoError(t, err)
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	id := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	holiday, err := c.Holiday(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, holiday.ID)
	assert.Equal(t, "Test Holiday", holiday.Name)
	assert.Equal(t, *holiday.ICalURL, "https://calendar.google.com/calendar/ical/en.uk%23holiday%40group.v.calendar.google.com/public/basic.ics")
	assert.Equal(t, *holiday.ICalTimeZone, "Europe/London")
	assert.Empty(t, holiday.CustomPeriods)
}

func TestHolidayCustomPeriods(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/manage/api/v1/holidays/8b154ff8-3d64-4b79-8b26-02b4baeb44e4" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		_, err := w.Write([]byte(
			`{"status":"success","data":{"id":"8b154ff8-3d64-4b79-8b26-02b4baeb44e4","created":"2023-01-04T10:09:47.254Z","modified":"2023-01-04T10:09:47.254Z","tenantId":"0","createdBy":null,"modifiedBy":null,"name":"Test Holiday","description":null,"iCalUrl":null,"iCalTimeZone":null,"customPeriods":[{"startTime":"2023-01-01T00:00:00Z","endTime":"2023-01-01T23:59:59Z"}],"resolvedDates":[{"name":"New Year's Day","startTime":"2022-01-01T00:00:00Z","endTime":"2022-01-02T00:00:00Z"},{"name":"New Year's Day observed","startTime":"2022-01-03T00:00:00Z","endTime":"2022-01-04T00:00:00Z"},{"name":"2nd January (substitute day) (Scotland)","startTime":"2022-01-04T00:00:00Z","endTime":"2022-01-05T00:00:00Z"},{"name":"Twelfth Night","startTime":"2022-01-05T00:00:00Z","endTime":"2022-01-06T00:00:00Z"},{"name":"Valentine's Day","startTime":"2022-02-14T00:00:00Z","endTime":"2022-02-15T00:00:00Z"},{"name":"Carnival / Shrove Tuesday / Pancake Day","startTime":"2022-03-01T00:00:00Z","endTime":"2022-03-02T00:00:00Z"},{"name":"St Patrick's Day (Northern Ireland)","startTime":"2022-03-17T00:00:00Z","endTime":"2022-03-18T00:00:00Z"},{"name":"Mother's Day","startTime":"2022-03-27T00:00:00Z","endTime":"2022-03-27T23:00:00Z"},{"name":"Daylight Saving Time starts","startTime":"2022-03-27T00:00:00Z","endTime":"2022-03-27T23:00:00Z"},{"name":"Good Friday","startTime":"2022-04-14T23:00:00Z","endTime":"2022-04-15T23:00:00Z"},{"name":"Easter Sunday","startTime":"2022-04-16T23:00:00Z","endTime":"2022-04-17T23:00:00Z"},{"name":"Easter Monday (regional holiday)","startTime":"2022-04-17T23:00:00Z","endTime":"2022-04-18T23:00:00Z"},{"name":"St. George's Day","startTime":"2022-04-22T23:00:00Z","endTime":"2022-04-23T23:00:00Z"},{"name":"Early May Bank Holiday","startTime":"2022-05-01T23:00:00Z","endTime":"2022-05-02T23:00:00Z"},{"name":"Spring Bank Holiday","startTime":"2022-06-01T23:00:00Z","endTime":"2022-06-02T23:00:00Z"},{"name":"Queen Elizabeth II's Platinum Jubilee","startTime":"2022-06-02T23:00:00Z","endTime":"2022-06-03T23:00:00Z"},{"name":"Queen Elizabeth II's Birthday","startTime":"2022-06-10T23:00:00Z","endTime":"2022-06-11T23:00:00Z"},{"name":"Father's Day","startTime":"2022-06-18T23:00:00Z","endTime":"2022-06-19T23:00:00Z"},{"name":"Battle of the Boyne (Northern Ireland)","startTime":"2022-07-11T23:00:00Z","endTime":"2022-07-12T23:00:00Z"},{"name":"Summer Bank Holiday (Scotland)","startTime":"2022-07-31T23:00:00Z","endTime":"2022-08-01T23:00:00Z"},{"name":"Summer Bank Holiday (regional holiday)","startTime":"2022-08-28T23:00:00Z","endTime":"2022-08-29T23:00:00Z"},{"name":"State Funeral of Queen Elizabeth II","startTime":"2022-09-18T23:00:00Z","endTime":"2022-09-19T23:00:00Z"},{"name":"Daylight Saving Time ends","startTime":"2022-10-29T23:00:00Z","endTime":"2022-10-31T00:00:00Z"},{"name":"Halloween","startTime":"2022-10-31T00:00:00Z","endTime":"2022-11-01T00:00:00Z"},{"name":"Guy Fawkes Day","startTime":"2022-11-05T00:00:00Z","endTime":"2022-11-06T00:00:00Z"},{"name":"Remembrance Sunday","startTime":"2022-11-13T00:00:00Z","endTime":"2022-11-14T00:00:00Z"},{"name":"St Andrew's Day (Scotland)","startTime":"2022-11-30T00:00:00Z","endTime":"2022-12-01T00:00:00Z"},{"name":"Christmas Eve","startTime":"2022-12-24T00:00:00Z","endTime":"2022-12-25T00:00:00Z"},{"name":"Christmas Day","startTime":"2022-12-25T00:00:00Z","endTime":"2022-12-26T00:00:00Z"},{"name":"Boxing Day","startTime":"2022-12-26T00:00:00Z","endTime":"2022-12-27T00:00:00Z"},{"name":"Substitute Bank Holiday for Christmas Day","startTime":"2022-12-27T00:00:00Z","endTime":"2022-12-28T00:00:00Z"},{"name":"New Year's Eve","startTime":"2022-12-31T00:00:00Z","endTime":"2023-01-01T00:00:00Z"},{"name":"New Year's Day","startTime":"2023-01-01T00:00:00Z","endTime":"2023-01-02T00:00:00Z"},{"name":"New Year's Day observed","startTime":"2023-01-02T00:00:00Z","endTime":"2023-01-03T00:00:00Z"},{"name":"2nd January (substitute day) (Scotland)","startTime":"2023-01-03T00:00:00Z","endTime":"2023-01-04T00:00:00Z"},{"name":"Twelfth Night","startTime":"2023-01-05T00:00:00Z","endTime":"2023-01-06T00:00:00Z"},{"name":"Valentine's Day","startTime":"2023-02-14T00:00:00Z","endTime":"2023-02-15T00:00:00Z"},{"name":"Carnival / Shrove Tuesday / Pancake Day","startTime":"2023-02-21T00:00:00Z","endTime":"2023-02-22T00:00:00Z"},{"name":"St Patrick's Day (Northern Ireland)","startTime":"2023-03-17T00:00:00Z","endTime":"2023-03-18T00:00:00Z"},{"name":"Mother's Day","startTime":"2023-03-19T00:00:00Z","endTime":"2023-03-20T00:00:00Z"},{"name":"Daylight Saving Time starts","startTime":"2023-03-26T00:00:00Z","endTime":"2023-03-26T23:00:00Z"},{"name":"Good Friday","startTime":"2023-04-06T23:00:00Z","endTime":"2023-04-07T23:00:00Z"},{"name":"Easter Sunday","startTime":"2023-04-08T23:00:00Z","endTime":"2023-04-09T23:00:00Z"},{"name":"Easter Monday (regional holiday)","startTime":"2023-04-09T23:00:00Z","endTime":"2023-04-10T23:00:00Z"},{"name":"St. George's Day","startTime":"2023-04-22T23:00:00Z","endTime":"2023-04-23T23:00:00Z"},{"name":"Early May Bank Holiday","startTime":"2023-04-30T23:00:00Z","endTime":"2023-05-01T23:00:00Z"},{"name":"The Coronation of King Charles III","startTime":"2023-05-05T23:00:00Z","endTime":"2023-05-06T23:00:00Z"},{"name":"Bank Holiday for the Coronation of King Charles III","startTime":"2023-05-07T23:00:00Z","endTime":"2023-05-08T23:00:00Z"},{"name":"Spring Bank Holiday","startTime":"2023-05-28T23:00:00Z","endTime":"2023-05-29T23:00:00Z"},{"name":"King's Birthday","startTime":"2023-06-16T23:00:00Z","endTime":"2023-06-17T23:00:00Z"},{"name":"Father's Day","startTime":"2023-06-17T23:00:00Z","endTime":"2023-06-18T23:00:00Z"},{"name":"Battle of the Boyne (Northern Ireland)","startTime":"2023-07-11T23:00:00Z","endTime":"2023-07-12T23:00:00Z"},{"name":"Summer Bank Holiday (Scotland)","startTime":"2023-08-06T23:00:00Z","endTime":"2023-08-07T23:00:00Z"},{"name":"Summer Bank Holiday (regional holiday)","startTime":"2023-08-27T23:00:00Z","endTime":"2023-08-28T23:00:00Z"},{"name":"Daylight Saving Time ends","startTime":"2023-10-28T23:00:00Z","endTime":"2023-10-30T00:00:00Z"},{"name":"Halloween","startTime":"2023-10-31T00:00:00Z","endTime":"2023-11-01T00:00:00Z"},{"name":"Guy Fawkes Day","startTime":"2023-11-05T00:00:00Z","endTime":"2023-11-06T00:00:00Z"},{"name":"Remembrance Sunday","startTime":"2023-11-12T00:00:00Z","endTime":"2023-11-13T00:00:00Z"},{"name":"St Andrew's Day (Scotland)","startTime":"2023-11-30T00:00:00Z","endTime":"2023-12-01T00:00:00Z"},{"name":"Christmas Eve","startTime":"2023-12-24T00:00:00Z","endTime":"2023-12-25T00:00:00Z"},{"name":"Christmas Day","startTime":"2023-12-25T00:00:00Z","endTime":"2023-12-26T00:00:00Z"},{"name":"Boxing Day","startTime":"2023-12-26T00:00:00Z","endTime":"2023-12-27T00:00:00Z"},{"name":"New Year's Eve","startTime":"2023-12-31T00:00:00Z","endTime":"2024-01-01T00:00:00Z"},{"name":"New Year's Day","startTime":"2024-01-01T00:00:00Z","endTime":"2024-01-02T00:00:00Z"},{"name":"2nd January (Scotland)","startTime":"2024-01-02T00:00:00Z","endTime":"2024-01-03T00:00:00Z"},{"name":"Twelfth Night","startTime":"2024-01-05T00:00:00Z","endTime":"2024-01-06T00:00:00Z"},{"name":"Carnival / Shrove Tuesday / Pancake Day","startTime":"2024-02-13T00:00:00Z","endTime":"2024-02-14T00:00:00Z"},{"name":"Valentine's Day","startTime":"2024-02-14T00:00:00Z","endTime":"2024-02-15T00:00:00Z"},{"name":"Mother's Day","startTime":"2024-03-10T00:00:00Z","endTime":"2024-03-11T00:00:00Z"},{"name":"St Patrick's Day (Northern Ireland)","startTime":"2024-03-17T00:00:00Z","endTime":"2024-03-18T00:00:00Z"},{"name":"Day off for St Patrick's Day (Northern Ireland)","startTime":"2024-03-18T00:00:00Z","endTime":"2024-03-19T00:00:00Z"},{"name":"Good Friday","startTime":"2024-03-29T00:00:00Z","endTime":"2024-03-30T00:00:00Z"},{"name":"Daylight Saving Time starts","startTime":"2024-03-31T00:00:00Z","endTime":"2024-03-31T23:00:00Z"},{"name":"Easter Sunday","startTime":"2024-03-31T00:00:00Z","endTime":"2024-03-31T23:00:00Z"},{"name":"Easter Monday (regional holiday)","startTime":"2024-03-31T23:00:00Z","endTime":"2024-04-01T23:00:00Z"},{"name":"St. George's Day","startTime":"2024-04-22T23:00:00Z","endTime":"2024-04-23T23:00:00Z"},{"name":"Early May Bank Holiday","startTime":"2024-05-05T23:00:00Z","endTime":"2024-05-06T23:00:00Z"},{"name":"Spring Bank Holiday","startTime":"2024-05-26T23:00:00Z","endTime":"2024-05-27T23:00:00Z"},{"name":"King's Birthday","startTime":"2024-06-14T23:00:00Z","endTime":"2024-06-15T23:00:00Z"},{"name":"Father's Day","startTime":"2024-06-15T23:00:00Z","endTime":"2024-06-16T23:00:00Z"},{"name":"Battle of the Boyne (Northern Ireland)","startTime":"2024-07-11T23:00:00Z","endTime":"2024-07-12T23:00:00Z"},{"name":"Summer Bank Holiday (Scotland)","startTime":"2024-08-04T23:00:00Z","endTime":"2024-08-05T23:00:00Z"},{"name":"Summer Bank Holiday (regional holiday)","startTime":"2024-08-25T23:00:00Z","endTime":"2024-08-26T23:00:00Z"},{"name":"Daylight Saving Time ends","startTime":"2024-10-26T23:00:00Z","endTime":"2024-10-28T00:00:00Z"},{"name":"Halloween","startTime":"2024-10-31T00:00:00Z","endTime":"2024-11-01T00:00:00Z"},{"name":"Guy Fawkes Day","startTime":"2024-11-05T00:00:00Z","endTime":"2024-11-06T00:00:00Z"},{"name":"Remembrance Sunday","startTime":"2024-11-10T00:00:00Z","endTime":"2024-11-11T00:00:00Z"},{"name":"St Andrew's Day (Scotland)","startTime":"2024-11-30T00:00:00Z","endTime":"2024-12-01T00:00:00Z"},{"name":"St Andrew's Day observed (Scotland)","startTime":"2024-12-02T00:00:00Z","endTime":"2024-12-03T00:00:00Z"},{"name":"Christmas Eve","startTime":"2024-12-24T00:00:00Z","endTime":"2024-12-25T00:00:00Z"},{"name":"Christmas Day","startTime":"2024-12-25T00:00:00Z","endTime":"2024-12-26T00:00:00Z"},{"name":"Boxing Day","startTime":"2024-12-26T00:00:00Z","endTime":"2024-12-27T00:00:00Z"},{"name":"New Year's Eve","startTime":"2024-12-31T00:00:00Z","endTime":"2025-01-01T00:00:00Z"}],"jobs":[]}}`,
		))
		require.NoError(t, err)
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	id := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	holiday, err := c.Holiday(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, holiday.ID)
	assert.Equal(t, "Test Holiday", holiday.Name)
	assert.Equal(t, holiday.CustomPeriods, CustomPeriods([]CustomPeriod{
		{
			StartTime: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2023, time.January, 1, 23, 59, 59, 0, time.UTC),
		},
	}))
	assert.Nil(t, holiday.ICalURL)
	assert.Nil(t, holiday.ICalTimeZone)
}

func TestUpdateHoliday(t *testing.T) {
	id := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	holiday := Holiday{
		ID: id,
	}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/manage/api/v1/holidays/"+id {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		parsedHoliday := Holiday{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&parsedHoliday)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if parsedHoliday.ID != "" {
			http.Error(w, "id should be empty when updating", http.StatusBadRequest)
			return
		}
		parsedHoliday.ID = id
		assert.Equal(t, holiday, parsedHoliday)
		enc := json.NewEncoder(w)
		_ = enc.Encode(responseWrapper[Holiday]{Data: parsedHoliday})
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	returnedHoliday, err := c.UpdateHoliday(ctx, holiday)
	require.NoError(t, err)
	assert.Equal(t, holiday, returnedHoliday)
}

func TestDeleteHoliday(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/manage/api/v1/holidays/8b154ff8-3d64-4b79-8b26-02b4baeb44e4" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		_, err := w.Write([]byte("successfully deleted"))
		require.NoError(t, err)
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	err = c.DeleteHoliday(ctx, "8b154ff8-3d64-4b79-8b26-02b4baeb44e4")
	require.NoError(t, err)
}
