package mlapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/prometheus/common/model"
)

// AnomalyCondition describes what sort of anomaly an alert should be generated for.
//
// It is only used for forecast alerts and is not supported for outlier alerts.
type AnomalyCondition string

const (
	// AnomalyConditionHigh alerts when the actual value of a metric is higher than the expected range.
	AnomalyConditionHigh AnomalyCondition = "high"
	// AnomalyConditionLow alerts when the actual value of a metric is lower than the expected range.
	AnomalyConditionLow AnomalyCondition = "low"
	// AnomalyConditionAny alerts when any anomalous condition is present (too high or too low).
	AnomalyConditionAny AnomalyCondition = "any"
)

func (d *AnomalyCondition) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	switch AnomalyCondition(value) {
	case AnomalyConditionHigh, AnomalyConditionLow, AnomalyConditionAny:
		*d = AnomalyCondition(value)
		return nil
	default:
		return fmt.Errorf("unrecognized anomalyCondition: %s", value)
	}
}

type NoDataState string

const (
	NoDataStateOK       = "OK"
	NoDataStateAlerting = "Alerting"
	NoDataStateNoData   = "NoData"
)

func (d *NoDataState) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	switch AnomalyCondition(value) {
	case NoDataStateOK, NoDataStateAlerting, NoDataStateNoData:
		*d = NoDataState(value)
		return nil
	case "":
		*d = NoDataState("")
		return nil
	default:
		return fmt.Errorf("unrecognized noDataState: %s", value)
	}
}

type Alert struct {
	ID string `json:"id,omitempty"`

	// Ttile of the alert, at most 190 characters.
	Title string `json:"title"`
	// Whether we look for anomalies that are too low, too high, or both.
	AnomalyCondition AnomalyCondition `json:"anomalyCondition,omitempty"`
	// The Prometheus style for clause of an alert.
	For model.Duration `json:"for"`
	// Allows a comparison of if the average value over the window is
	// anomalous, for example `>0.7` would alert if more than 70% of the points
	// in the window are anomalous according to the anomaly condition.
	Threshold string `json:"threshold,omitempty"`
	// Specifying a window will issue a range query instead of an instant query
	// and average values over that range. Maximum of 12h.
	Window model.Duration `json:"window"`
	// Additional labels to add to the alert.
	Labels map[string]string `json:"labels"`
	// Annotations to include on an alert, such as severity or a description.
	Annotations map[string]string `json:"annotations"`
	// NoDataState allows alerting if no data is found in the query. Empty will
	// default to OK to match Prometheus behavior.
	NoDataState NoDataState `json:"noDataCondition"`
	// CustomQuery [Experimental] allows specifying a custom query to use for
	// the alert. Alerts will still be triggered for values that are not zero,
	// or in the case a threshold is defined, meeting the threshold. This field
	// is experimental and may change or be removed at any time.
	CustomQuery string `json:"customQuery"`

	// NotificationSettings are overrides to how notifications for an alert are
	// specified. The field is passed to the alert without modification.
	NotificationSettings *models.AlertRuleNotificationSettings `json:"notificationSettings"`

	SyncError string `json:"syncError,omitempty"`
}

// NewJobAlert creates an alert for a job.
func (c *Client) NewJobAlert(ctx context.Context, jobID string, alert Alert) (Alert, error) {
	data, err := json.Marshal(alert)
	if err != nil {
		return Alert{}, err
	}

	result := responseWrapper[Alert]{}
	err = c.request(ctx, "POST", fmt.Sprintf("/manage/api/v1/jobs/%s/alerts", jobID), nil, bytes.NewReader(data), &result)
	if err != nil {
		return Alert{}, err
	}
	return result.Data, nil
}

// JobAlerts fetches all alerts for a given Job.
func (c *Client) JobAlerts(ctx context.Context, jobID string) ([]Alert, error) {
	result := responseWrapper[[]Alert]{}
	err := c.request(ctx, "GET", fmt.Sprintf("/manage/api/v1/jobs/%s/alerts", jobID), nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

// JobAlert fetches an existing alert for the given machine learning job.
func (c *Client) JobAlert(ctx context.Context, jobID, alertID string) (Alert, error) {
	result := responseWrapper[Alert]{}
	err := c.request(ctx, "GET", fmt.Sprintf("/manage/api/v1/jobs/%s/alerts/%s", jobID, alertID), nil, nil, &result)
	if err != nil {
		return Alert{}, err
	}
	return result.Data, nil
}

// UpdateJobAlert updates the alert for a machine learning job.
func (c *Client) UpdateJobAlert(ctx context.Context, jobID string, alert Alert) (Alert, error) {
	alertID := alert.ID
	// Clear the ID before sending otherwise validation fails.
	alert.ID = ""
	data, err := json.Marshal(alert)
	if err != nil {
		return Alert{}, err
	}

	result := responseWrapper[Alert]{}
	err = c.request(ctx, "POST", fmt.Sprintf("/manage/api/v1/jobs/%s/alerts/%s", jobID, alertID), nil, bytes.NewReader(data), &result)
	if err != nil {
		return Alert{}, err
	}
	return result.Data, nil
}

// DeleteJobAlert deletes an alert on a job.
func (c *Client) DeleteJobAlert(ctx context.Context, jobID, alertID string) error {
	return c.request(ctx, "DELETE", fmt.Sprintf("/manage/api/v1/jobs/%s/alerts/%s", jobID, alertID), nil, nil, nil)
}

// NewOutlierAlert creates an alert for an outlier detector.
func (c *Client) NewOutlierAlert(ctx context.Context, outlierID string, alert Alert) (Alert, error) {
	data, err := json.Marshal(alert)
	if err != nil {
		return Alert{}, err
	}

	result := responseWrapper[Alert]{}
	err = c.request(ctx, "POST", fmt.Sprintf("/manage/api/v1/outliers/%s/alerts", outlierID), nil, bytes.NewReader(data), &result)
	if err != nil {
		return Alert{}, err
	}
	return result.Data, nil
}

// OutlierAlerts fetches all alerts for a given Job.
func (c *Client) OutlierAlerts(ctx context.Context, outlierID string) ([]Alert, error) {
	result := responseWrapper[[]Alert]{}
	err := c.request(ctx, "GET", fmt.Sprintf("/manage/api/v1/outliers/%s/alerts", outlierID), nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

// JobAlert fetches an existing alert for the given outlier detector.
func (c *Client) OutlierAlert(ctx context.Context, outlierID, alertID string) (Alert, error) {
	result := responseWrapper[Alert]{}
	err := c.request(ctx, "GET", fmt.Sprintf("/manage/api/v1/outliers/%s/alerts/%s", outlierID, alertID), nil, nil, &result)
	if err != nil {
		return Alert{}, err
	}
	return result.Data, nil
}

// UpdateJobAlert updates the alert for an outlier detector.
func (c *Client) UpdateOutlierAlert(ctx context.Context, outlierID string, alert Alert) (Alert, error) {
	alertID := alert.ID
	// Clear the ID before sending otherwise validation fails.
	alert.ID = ""
	data, err := json.Marshal(alert)
	if err != nil {
		return Alert{}, err
	}

	result := responseWrapper[Alert]{}
	err = c.request(ctx, "POST", fmt.Sprintf("/manage/api/v1/outliers/%s/alerts/%s", outlierID, alertID), nil, bytes.NewReader(data), &result)
	if err != nil {
		return Alert{}, err
	}
	return result.Data, nil
}

// DeleteOutlierAlert deletes an alert on an outlier detector.
func (c *Client) DeleteOutlierAlert(ctx context.Context, outlierID, alertID string) error {
	return c.request(ctx, "DELETE", fmt.Sprintf("/manage/api/v1/outliers/%s/alerts/%s", outlierID, alertID), nil, nil, nil)
}
