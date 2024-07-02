package mlapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOutlierDetector(t *testing.T) {
	id := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	outlier := OutlierDetector{}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/manage/api/v1/outliers" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		parsedOutlierDetector := OutlierDetector{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&parsedOutlierDetector)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		assert.Equal(t, outlier, parsedOutlierDetector)
		parsedOutlierDetector.ID = id
		enc := json.NewEncoder(w)
		_ = enc.Encode(responseWrapper[OutlierDetector]{Data: parsedOutlierDetector})
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	returnedOutlierDetector, err := c.NewOutlierDetector(ctx, outlier)
	require.NoError(t, err)
	assert.Equal(t, id, returnedOutlierDetector.ID)
}

func TestOutlierDetector(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/manage/api/v1/outliers/8b154ff8-3d64-4b79-8b26-02b4baeb44e4" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		_, err := w.Write([]byte(
			`{"status":"success","data":{"id":"8b154ff8-3d64-4b79-8b26-02b4baeb44e4","created":"2022-01-05T15:48:48.647Z","modified":"2022-01-05T15:48:48.647Z","createdBy":"a_user","modifiedBy":null,"name":"Test OutlierDetector","metric":"test_outlier","description":"","grafanaUrl":"http://localhost:3000/","grafanaApiKey":"\u003credacted\u003e","datasourceId":10,"datasourceUid":"abcd1234","datasourceType":"prometheus","queryParams":{"exemplar":true,"expr":"sum(up)","interval":"","legendFormat":"","refId":"A"},"interval":300,"algorithm":{"name":"dbscan","config":{"eps":0.5}}}}`,
		))
		require.NoError(t, err)
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	id := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	outlier, err := c.OutlierDetector(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, outlier.ID)
	assert.Equal(t, "Test OutlierDetector", outlier.Name)
	assert.Equal(t, "test_outlier", outlier.Metric)
	assert.NotEmpty(t, outlier.QueryParams)
	assert.NotEmpty(t, outlier.Algorithm)
}

func TestUpdateOutlierDetector(t *testing.T) {
	id := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	outlier := OutlierDetector{
		ID: id,
	}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/manage/api/v1/outliers/"+id {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		parsedOutlierDetector := OutlierDetector{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&parsedOutlierDetector)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if parsedOutlierDetector.ID != "" {
			http.Error(w, "id should be empty when updating", http.StatusBadRequest)
			return
		}
		parsedOutlierDetector.ID = id
		assert.Equal(t, outlier, parsedOutlierDetector)
		enc := json.NewEncoder(w)
		_ = enc.Encode(responseWrapper[OutlierDetector]{Data: parsedOutlierDetector})
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	returnedOutlierDetector, err := c.UpdateOutlierDetector(ctx, outlier)
	require.NoError(t, err)
	assert.Equal(t, outlier, returnedOutlierDetector)
}

func TestDeleteOutlierDetector(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/manage/api/v1/outliers/8b154ff8-3d64-4b79-8b26-02b4baeb44e4" {
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

	err = c.DeleteOutlierDetector(ctx, "8b154ff8-3d64-4b79-8b26-02b4baeb44e4")
	require.NoError(t, err)
}
