package mlapi

import (
	"bytes"
	"context"
	"encoding/json"
)

type OutlierAlgorithmConfig struct {
	Epsilon float64 `json:"epsilon"`
}

type OutlierAlgorithm struct {
	Name        string                  `json:"name"`
	Sensitivity float64                 `json:"sensitivity"` // used by MAD
	Config      *OutlierAlgorithmConfig `json:"config"`      // used by DBSCAN
}

// OutlierDetector defines an outlier detector instance
type OutlierDetector struct {
	ID string `json:"id,omitempty"`
	// Name is a human readable name for the outlier detector.
	Name string `json:"name"`
	// Metric is the metric name used to query the outlier detector. Must match Prometheus
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

	// Algorithm specifies the algorithm to use and its configuration. See
	// https://grafana.com/docs/grafana-cloud/machine-learning/outlier-detection/ for the
	// options.
	Algorithm OutlierAlgorithm `json:"algorithm"`
}

// NewOutlierDetector creates an outlier detector.
func (c *Client) NewOutlierDetector(ctx context.Context, outlier OutlierDetector) (OutlierDetector, error) {
	data, err := json.Marshal(outlier)
	if err != nil {
		return OutlierDetector{}, err
	}

	result := responseWrapper[OutlierDetector]{}
	err = c.request(ctx, "POST", "/manage/api/v1/outliers", nil, bytes.NewReader(data), &result)
	if err != nil {
		return OutlierDetector{}, err
	}
	return result.Data, nil
}

// OutlierDetectors fetches all existing outlier detectors.
func (c *Client) OutlierDetectors(ctx context.Context) ([]OutlierDetector, error) {
	result := responseWrapper[[]OutlierDetector]{}
	err := c.request(ctx, "GET", "/manage/api/v1/outliers", nil, nil, &result)
	if err != nil {
		return []OutlierDetector{}, err
	}
	return result.Data, nil
}

// OutlierDetector fetches an existing outlier detector.
func (c *Client) OutlierDetector(ctx context.Context, id string) (OutlierDetector, error) {
	result := responseWrapper[OutlierDetector]{}
	err := c.request(ctx, "GET", "/manage/api/v1/outliers/"+id, nil, nil, &result)
	if err != nil {
		return OutlierDetector{}, err
	}
	return result.Data, nil
}

// UpdateOutlierDetector updates an outlier detector.
func (c *Client) UpdateOutlierDetector(ctx context.Context, outlier OutlierDetector) (OutlierDetector, error) {
	id := outlier.ID
	// Clear the ID before sending otherwise validation fails.
	outlier.ID = ""
	data, err := json.Marshal(outlier)
	if err != nil {
		return OutlierDetector{}, err
	}

	result := responseWrapper[OutlierDetector]{}
	err = c.request(ctx, "POST", "/manage/api/v1/outliers/"+id, nil, bytes.NewReader(data), &result)
	if err != nil {
		return OutlierDetector{}, err
	}
	return result.Data, nil
}

// DeleteOutlierDetector deletes an outlier detector.
func (c *Client) DeleteOutlierDetector(ctx context.Context, id string) error {
	return c.request(ctx, "DELETE", "/manage/api/v1/outliers/"+id, nil, nil, nil)
}
