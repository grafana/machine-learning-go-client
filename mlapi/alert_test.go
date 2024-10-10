package mlapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJobAlert(t *testing.T) {
	jobID := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	alertID := "5218f38f-569b-448f-b81d-578173412195"
	alert := Alert{}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/manage/api/v1/jobs/"+jobID+"/alerts" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		parsedAlert := Alert{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&parsedAlert)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		assert.Equal(t, alert, parsedAlert)
		parsedAlert.ID = alertID
		enc := json.NewEncoder(w)
		_ = enc.Encode(responseWrapper[Alert]{Data: parsedAlert})
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	returnedAlert, err := c.NewJobAlert(ctx, jobID, alert)
	require.NoError(t, err)
	assert.Equal(t, alertID, returnedAlert.ID)
}

func TestJobAlert(t *testing.T) {
	jobID := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	alertID := "5218f38f-569b-448f-b81d-578173412195"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != fmt.Sprintf("/manage/api/v1/jobs/%s/alerts/%s", jobID, alertID) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		_, err := w.Write([]byte(
			`{"status":"success","data":{"id":"5218f38f-569b-448f-b81d-578173412195","created":"2024-07-02T16:19:39.992Z","modified":"2024-07-02T16:19:39.992Z","createdBy":null,"modifiedBy":null,"title":"test job alert","anomalyCondition":"any","for":"5m","window":"0s","labels":{"foo":"bar"},"annotations":{"description":"Anomaly detected","summary":"Anomaly detected"},"noDataCondition":""}}`,
		))
		require.NoError(t, err)
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	alert, err := c.JobAlert(ctx, jobID, alertID)
	require.NoError(t, err)
	assert.Equal(t, alertID, alert.ID)
	assert.Equal(t, "test job alert", alert.Title)
	assert.Equal(t, AnomalyConditionAny, alert.AnomalyCondition)
	assert.Equal(t, model.Duration(5*time.Minute), alert.For)
	assert.Equal(t, map[string]string{
		"foo": "bar",
	}, alert.Labels)
	assert.Equal(t, map[string]string{
		"description": "Anomaly detected",
		"summary":     "Anomaly detected",
	}, alert.Annotations)
}

func TestUpdateJobAlert(t *testing.T) {
	jobID := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	alertID := "5218f38f-569b-448f-b81d-578173412195"
	alert := Alert{
		ID: alertID,
	}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != fmt.Sprintf("/manage/api/v1/jobs/%s/alerts/%s", jobID, alertID) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		parsedAlert := Alert{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&parsedAlert)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if parsedAlert.ID != "" {
			http.Error(w, "id should be empty when updating", http.StatusBadRequest)
			return
		}
		parsedAlert.ID = alertID
		assert.Equal(t, alert, parsedAlert)
		enc := json.NewEncoder(w)
		_ = enc.Encode(responseWrapper[Alert]{Data: parsedAlert})
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	returnedAlert, err := c.UpdateJobAlert(ctx, jobID, alert)
	require.NoError(t, err)
	assert.Equal(t, alert, returnedAlert)
}

func TestDeleteJobAlert(t *testing.T) {
	jobID := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	alertID := "5218f38f-569b-448f-b81d-578173412195"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != fmt.Sprintf("/manage/api/v1/jobs/%s/alerts/%s", jobID, alertID) {
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

	err = c.DeleteJobAlert(ctx, jobID, alertID)
	require.NoError(t, err)
}

func TestNewOutlierAlert(t *testing.T) {
	outlierID := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	alertID := "903a57c2-04cf-4a03-b4aa-54567e981ac0"
	alert := Alert{}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/manage/api/v1/outliers/"+outlierID+"/alerts" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		parsedAlert := Alert{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&parsedAlert)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		assert.Equal(t, alert, parsedAlert)
		parsedAlert.ID = alertID
		enc := json.NewEncoder(w)
		_ = enc.Encode(responseWrapper[Alert]{Data: parsedAlert})
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	returnedAlert, err := c.NewOutlierAlert(ctx, outlierID, alert)
	require.NoError(t, err)
	assert.Equal(t, alertID, returnedAlert.ID)
}

func TestOutlierAlert(t *testing.T) {
	jobID := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	alertID := "903a57c2-04cf-4a03-b4aa-54567e981ac0"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != fmt.Sprintf("/manage/api/v1/outliers/%s/alerts/%s", jobID, alertID) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		_, err := w.Write([]byte(
			`{"status":"success","data":{"id":"903a57c2-04cf-4a03-b4aa-54567e981ac0","created":"2024-07-02T19:04:56.604Z","modified":"2024-07-02T19:04:56.604Z","createdBy":null,"modifiedBy":null,"title":"test outlier alert","for":"5m","window":"1h","labels":{"foo":"bar"},"annotations":{"description":"Outlier detected","summary":"Outlier detected"},"noDataCondition":""}}`,
		))
		require.NoError(t, err)
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	alert, err := c.OutlierAlert(ctx, jobID, alertID)
	require.NoError(t, err)
	assert.Equal(t, alertID, alert.ID)
	assert.Equal(t, "test outlier alert", alert.Title)
	assert.EqualValues(t, "", alert.AnomalyCondition)
	assert.EqualValues(t, 5*time.Minute, alert.For)
	assert.EqualValues(t, time.Hour, alert.Window)
	assert.Equal(t, map[string]string{
		"foo": "bar",
	}, alert.Labels)
	assert.Equal(t, map[string]string{
		"description": "Outlier detected",
		"summary":     "Outlier detected",
	}, alert.Annotations)
}

func TestUpdateOutlierAlert(t *testing.T) {
	jobID := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	alertID := "903a57c2-04cf-4a03-b4aa-54567e981ac0"
	alert := Alert{
		ID: alertID,
	}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != fmt.Sprintf("/manage/api/v1/outliers/%s/alerts/%s", jobID, alertID) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		parsedAlert := Alert{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&parsedAlert)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if parsedAlert.ID != "" {
			http.Error(w, "id should be empty when updating", http.StatusBadRequest)
			return
		}
		parsedAlert.ID = alertID
		assert.Equal(t, alert, parsedAlert)
		enc := json.NewEncoder(w)
		_ = enc.Encode(responseWrapper[Alert]{Data: parsedAlert})
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	returnedAlert, err := c.UpdateOutlierAlert(ctx, jobID, alert)
	require.NoError(t, err)
	assert.Equal(t, alert, returnedAlert)
}

func TestDeleteOutlierAlert(t *testing.T) {
	jobID := "8b154ff8-3d64-4b79-8b26-02b4baeb44e4"
	alertID := "903a57c2-04cf-4a03-b4aa-54567e981ac0"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != fmt.Sprintf("/manage/api/v1/outliers/%s/alerts/%s", jobID, alertID) {
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

	err = c.DeleteOutlierAlert(ctx, jobID, alertID)
	require.NoError(t, err)
}
