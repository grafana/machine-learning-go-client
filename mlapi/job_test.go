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

func TestNewJob(t *testing.T) {
	id := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	job := Job{}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/manage/api/v1/jobs" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		parsedJob := Job{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&parsedJob)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		assert.Equal(t, job, parsedJob)
		parsedJob.ID = id
		enc := json.NewEncoder(w)
		_ = enc.Encode(jobResponseWrapper{Data: parsedJob})
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	returnedJob, err := c.NewJob(ctx, job)
	require.NoError(t, err)
	assert.Equal(t, id, returnedJob.ID)
}

func TestJob(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/manage/api/v1/jobs/8b154ff8-3d64-4b79-8b26-02b4baeb44e4" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		_, err := w.Write([]byte(
			`{"status":"success","data":{"id":"8b154ff8-3d64-4b79-8b26-02b4baeb44e4","created":"2022-01-05T15:48:48.647Z","modified":"2022-01-05T15:48:48.647Z","createdBy":"a_user","modifiedBy":null,"name":"Test Job","metric":"test_job","description":"","grafanaUrl":"http://localhost:3000/","grafanaApiKey":"\u003credacted\u003e","datasourceId":10,"datasourceType":"prometheus","queryParams":{"exemplar":true,"expr":"sum(up)","interval":"","legendFormat":"","refId":"A"},"interval":300,"algorithm":"grafana_prophet_1_0_1","hyperParams":{"changepoint_prior_scale":0.05,"growth":"linear","holidays_prior_scale":10,"interval_width":0.95,"seasonality_mode":"additive","seasonality_prior_scale":10},"trainingWindow":7776000,"trainingFrequency":86400,"status":"pending","nextTrainingAt":"2022-01-05T15:48:48.638971435Z","trainingScheduledAt":null,"trainingCompletedAt":null,"lastTrainingStatus":null,"trainingResult":"Pending","trainingFailures":0}}`,
		))
		require.NoError(t, err)
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	id := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	job, err := c.Job(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, job.ID)
	assert.Equal(t, "Test Job", job.Name)
	assert.Equal(t, "test_job", job.Metric)
	assert.NotEmpty(t, job.QueryParams)
	assert.NotEmpty(t, job.HyperParams)
}

func TestUpdateJob(t *testing.T) {
	id := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	job := Job{
		ID: id,
	}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/manage/api/v1/jobs/"+id {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		parsedJob := Job{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&parsedJob)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		assert.Equal(t, job, parsedJob)
		parsedJob.ID = id
		enc := json.NewEncoder(w)
		_ = enc.Encode(jobResponseWrapper{Data: parsedJob})
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	returnedJob, err := c.UpdateJob(ctx, job)
	require.NoError(t, err)
	assert.Equal(t, job, returnedJob)
}

func TestDeleteJob(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/manage/api/v1/jobs/8b154ff8-3d64-4b79-8b26-02b4baeb44e4" {
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

	err = c.DeleteJob(ctx, "8b154ff8-3d64-4b79-8b26-02b4baeb44e4")
	require.NoError(t, err)
}
