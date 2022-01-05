package mlapi

import (
	"bytes"
	"context"
	"encoding/json"
)

// Job is a job that will be scheduled.
type Job struct {
	ID                string                 `json:"id,omitempty"`
	Name              string                 `json:"name"`
	Metric            string                 `json:"metric"`
	Description       string                 `json:"description"`
	GrafanaURL        string                 `json:"grafanaUrl"`
	DatasourceID      uint                   `json:"datasourceId"`
	DatasourceType    string                 `json:"datasourceType"`
	QueryParams       map[string]interface{} `json:"queryParams"`
	Interval          uint                   `json:"interval"`
	Algorithm         string                 `json:"algorithm"`
	HyperParams       map[string]interface{} `json:"hyperParams"`
	TrainingWindow    uint                   `json:"trainingWindow"`
	TrainingFrequency uint                   `json:"trainingFrequency"`
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
