package mlapi

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// Job is a job that will be scheduled.
type Job struct {
	ID string `json:"id,omitempty"`
	// Name is a human readable name for the job.
	Name string `json:"name"`
	// Metric is the metric name used to query the job. Must match Prometheus
	// naming requirements:
	// https://prometheus.io/docs/concepts/data_model/#metric-names-and-labels.
	Metric      string `json:"metric"`
	Description string `json:"description"`

	// GrafanaURL is the full URL to the Grafana instance. For example,
	// https://myinstance.grafana.net/.
	GrafanaURL string `json:"grafanaUrl"`
	// DatasourceID is the numeric ID of the datasource to query when training
	// data.
	DatasourceID uint `json:"datasourceId"`
	// DatasourceUID is the string UID of the datasource to query when training
	// data.
	DatasourceUID  string                 `json:"datasourceUid"`
	DatasourceType string                 `json:"datasourceType"`
	CustomLabels   map[string]interface{} `json:"customLabels"`
	QueryParams    map[string]interface{} `json:"queryParams"`
	// Interval is the data resolution in seconds.
	Interval uint `json:"interval"`
	// TrainingWindow is the lookback window to train on in seconds.
	TrainingWindow uint `json:"trainingWindow"`
	// TrainingFrequency is how often to re-train a model in seconds.
	TrainingFrequency uint `json:"trainingFrequency"`

	// Algorithm is the algorithm to use for machine learning.
	// https://grafana.com/docs/grafana-cloud/machine-learning/models/ contains
	// information on all supported algorithms.
	Algorithm string `json:"algorithm"`
	// HyperParams are the hyperparameters that can be specified. See
	// https://grafana.com/docs/grafana-cloud/machine-learning/models/ for the
	// various hyperparameters that can be changed.
	HyperParams map[string]interface{} `json:"hyperParams"`

	// Holidays is a slice of IDs or names of Holidays to be linked to this job.
	// Requests may specify either IDs or names. Responses will always contain IDs.
	Holidays []string `json:"holidays"`

	// ManagedBy is used to identify who controls system forecasts. It is
	// required when creating system forecasts and must not be set otherwise.
	ManagedBy string `json:"managedBy,omitempty"`
}

// NewJob creates a machine learning job and schedules a training.
func (c *Client) NewJob(ctx context.Context, job Job) (Job, error) {
	return c.newJob(ctx, job, "/manage/api/v1/jobs")
}

// NewSystemJob creates a system machine learning job and schedules a training.
func (c *Client) NewSystemJob(ctx context.Context, job Job) (Job, error) {
	return c.newJob(ctx, job, "/manage/api/v1/system-jobs")
}

func (c *Client) newJob(ctx context.Context, job Job, path string) (Job, error) {
	data, err := json.Marshal(job)
	if err != nil {
		return Job{}, err
	}

	result := responseWrapper[Job]{}
	err = c.request(ctx, "POST", path, nil, bytes.NewReader(data), &result)
	if err != nil {
		return Job{}, err
	}
	return result.Data, nil
}

// Jobs fetches all existing machine learning jobs.
func (c *Client) Jobs(ctx context.Context) ([]Job, error) {
	result := responseWrapper[[]Job]{}
	err := c.request(ctx, "GET", "/manage/api/v1/jobs", nil, nil, &result)
	if err != nil {
		return []Job{}, err
	}
	return result.Data, nil
}

// Job fetches an existing machine learning job.
func (c *Client) Job(ctx context.Context, id string) (Job, error) {
	result := responseWrapper[Job]{}
	err := c.request(ctx, "GET", "/manage/api/v1/jobs/"+id, nil, nil, &result)
	if err != nil {
		return Job{}, err
	}
	return result.Data, nil
}

// UpdateJob updates a machine learning job. A new training will be scheduled as part of updating.
func (c *Client) UpdateJob(ctx context.Context, job Job) (Job, error) {
	return c.updateJob(ctx, job, "/manage/api/v1/jobs/")
}

// UpdateSystemJob updates a system machine learning job and schedules a new
// training. It can also be used to change a user job into a system job if
// necessary.
func (c *Client) UpdateSystemJob(ctx context.Context, job Job) (Job, error) {
	return c.updateJob(ctx, job, "/manage/api/v1/system-jobs/")
}

func (c *Client) updateJob(ctx context.Context, job Job, path string) (Job, error) {
	id := job.ID
	// Clear the ID before sending otherwise validation fails.
	job.ID = ""
	data, err := json.Marshal(job)
	if err != nil {
		return Job{}, err
	}

	result := responseWrapper[Job]{}
	err = c.request(ctx, "POST", path+id, nil, bytes.NewReader(data), &result)
	if err != nil {
		return Job{}, err
	}
	return result.Data, nil
}

// DeleteJob deletes a machine learning job.
func (c *Client) DeleteJob(ctx context.Context, id string) error {
	return c.request(ctx, "DELETE", "/manage/api/v1/jobs/"+id, nil, nil, nil)
}

// DeleteJob deletes a system machine learning job.
func (c *Client) DeleteSystemJob(ctx context.Context, id string) error {
	return c.request(ctx, "DELETE", "/manage/api/v1/system-jobs/"+id, nil, nil, nil)
}

// LinkHolidaysToJob links a job to a set of holidays.
// Only the ID and Holidays fields of the Job struct are used.
func (c *Client) LinkHolidaysToJob(ctx context.Context, jobID string, holidayIDs []string) (Job, error) {
	job := Job{
		Holidays: holidayIDs,
	}

	data, err := json.Marshal(job)
	if err != nil {
		return Job{}, err
	}

	result := responseWrapper[Job]{}
	err = c.request(ctx, "PUT", "/manage/api/v1/jobs/"+jobID+"/holidays", nil, bytes.NewReader(data), &result)
	if err != nil {
		return Job{}, err
	}
	return result.Data, err
}

// ForecastParams are parameters used in an ephemeral job prediction.
// Note this doesn't include `series` because users will not know
// the hash of the training's series in advance. Instead we enforce
// that ephemeral forecasts can only include a single series.
type ForecastParams struct {
	// Start is the start time of the forecast.
	Start time.Time `json:"start"`
	// End is the end time of the forecast.
	End time.Time `json:"end"`
	// Interval is the interval of the forecast.
	Interval uint `json:"interval"`
}

// ForecastRequest is a request to run an ephemeral forecast.
type ForecastRequest struct {
	// Job is the specification of a job to run the forecast for.
	Job Job `json:"job"`
	// ForecastParams specify the start, end, and interval of the forecast.
	ForecastParams ForecastParams `json:"forecastParams"`
}

// ForecastJob returns a forecast for a job definition and time range.
//
// This is a convenience API to avoid having to create a job, wait for it to train,
// and then query the forecast. It is designed for exploratory usage rather than
// to be called regularly; if you want to regularly query a forecast, create a job
// and query it using the `grafanacloud-ml-metrics` datasource.
// Jobs specified in the ForecastRequest must have a single series, or this
// will return an error.
// This function may be slow the first time it is called, but the result will be
// cached for 24 hours after that.
func (c *Client) ForecastJob(ctx context.Context, spec ForecastRequest) (backend.QueryDataResponse, error) {
	data, err := json.Marshal(spec)
	if err != nil {
		return backend.QueryDataResponse{}, err
	}

	result := responseWrapper[backend.QueryDataResponse]{}
	err = c.request(ctx, "POST", "/predict/api/v1/forecast", nil, bytes.NewReader(data), &result)
	if err != nil {
		return backend.QueryDataResponse{}, err
	}
	return result.Data, nil
}
