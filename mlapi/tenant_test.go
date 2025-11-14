package mlapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTenantInfo(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/tenant/api/v1/info" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		_, err := w.Write([]byte(`{"status":"success","data":{"id":1,"features":{"AugursMSTL":true,"SIFTDev":true,"SIFTExperimental":true,"SiftOnCallSummaries":true,"SIFTPreview":true,"SIFTAgent":true},"maxSeriesPerJob":1000,"maxSeriesPerOutlier":5000}}`))
		require.NoError(t, err)
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()

	tenantInfo, err := c.TenantInfo(ctx)
	require.NoError(t, err)
	assert.Equal(t, TenantInfo{
		MaxSeriesPerJob:     1000,
		MaxSeriesPerOutlier: 5000,
	}, tenantInfo)
}
