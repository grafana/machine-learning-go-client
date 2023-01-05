package mlapi

import (
	"bytes"
	"context"
	"encoding/json"
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
}

type jobResponseWrapper struct {
	Status   string   `json:"status"`
	Data     Job      `json:"data"`
	Warnings []string `json:"warnings"`
	Error    string   `json:"error"`
}

// NewJob creates a machine learning job and schedules a training.
func (c *Client) NewJob(ctx context.Context, job Job) (Job, error) {
	data, err := json.Marshal(job)
	if err != nil {
		return Job{}, err
	}

	result := jobResponseWrapper{}
	err = c.request(ctx, "POST", "/manage/api/v1/jobs", nil, bytes.NewBuffer(data), &result)
	if err != nil {
		return Job{}, err
	}
	return result.Data, nil
}

// Job fetches an existing machine learning job.
func (c *Client) Job(ctx context.Context, id string) (Job, error) {
	result := jobResponseWrapper{}
	err := c.request(ctx, "GET", "/manage/api/v1/jobs/"+id, nil, nil, &result)
	if err != nil {
		return Job{}, err
	}
	return result.Data, nil
}

// UpdateJob updates a machine learning job. A new training will be scheduled as part of updating.
func (c *Client) UpdateJob(ctx context.Context, job Job) (Job, error) {
	id := job.ID
	// Clear the ID before sending otherwise validation fails.
	job.ID = ""
	data, err := json.Marshal(job)
	if err != nil {
		return Job{}, err
	}

	result := jobResponseWrapper{}
	err = c.request(ctx, "POST", "/manage/api/v1/jobs/"+id, nil, bytes.NewBuffer(data), &result)
	if err != nil {
		return Job{}, err
	}
	return result.Data, nil
}

// DeleteJob deletes a machine learning job.
func (c *Client) DeleteJob(ctx context.Context, id string) error {
	return c.request(ctx, "DELETE", "/manage/api/v1/jobs/"+id, nil, nil, nil)
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

	result := jobResponseWrapper{}
	err = c.request(ctx, "PUT", "/manage/api/v1/jobs/"+jobID+"/holidays", nil, bytes.NewBuffer(data), &result)
	if err != nil {
		return Job{}, err
	}
	return result.Data, err
}
