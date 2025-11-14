package mlapi

import "context"

type TenantInfo struct {
	MaxSeriesPerJob     uint `json:"maxSeriesPerJob"`
	MaxSeriesPerOutlier uint `json:"maxSeriesPerOutlier"`
}

// TenantInfo returns the per forecast/outlier limits for the authenticated tenant.
func (c *Client) TenantInfo(ctx context.Context) (TenantInfo, error) {
	result := responseWrapper[TenantInfo]{}
	err := c.request(ctx, "GET", "/tenant/api/v1/info", nil, nil, &result)
	if err != nil {
		return TenantInfo{}, err
	}
	return result.Data, nil
}
