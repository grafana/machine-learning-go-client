package mlapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		require.NoError(t, err)
	}))
	defer s.Close()

	c, err := New(s.URL, Config{})
	require.NoError(t, err)
	ctx := context.Background()
	err = c.request(ctx, "GET", "/", nil, nil, nil)
	require.NoError(t, err)
}

func TestBasicAuth(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		if user != "hello" {
			http.Error(w, "bad username", http.StatusUnauthorized)
			return
		}
		if pass != "world" {
			http.Error(w, "bad password", http.StatusUnauthorized)
			return
		}
		_, err := w.Write([]byte("OK"))
		require.NoError(t, err)
	}))
	defer s.Close()

	c, err := New(s.URL, Config{
		BasicAuth: url.UserPassword("hello", "world"),
	})
	require.NoError(t, err)
	ctx := context.Background()
	err = c.request(ctx, "GET", "/", nil, nil, nil)
	require.NoError(t, err)
}

func TestBearerAuth(t *testing.T) {
	token := "my_token"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authorization") != "Bearer "+token {
			http.Error(w, "bad authorization header", http.StatusUnauthorized)
			return
		}
		_, err := w.Write([]byte("OK"))
		require.NoError(t, err)
	}))
	defer s.Close()

	c, err := New(s.URL, Config{
		BearerToken: token,
	})
	require.NoError(t, err)
	ctx := context.Background()
	err = c.request(ctx, "GET", "/", nil, nil, nil)
	require.NoError(t, err)
}

func TestFailures(t *testing.T) {
	message := "OK"
	status := http.StatusOK
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, message, status)
	}))
	defer s.Close()
	c, err := New(s.URL, Config{})
	require.NoError(t, err)

	for _, tc := range []struct {
		message string
		status  int
	}{
		{"not found", http.StatusNotFound},
		{"bad request", http.StatusBadRequest},
	} {
		message = tc.message
		status = tc.status
		err = c.request(context.Background(), "GET", "/", nil, nil, nil)
		require.EqualError(t, err, fmt.Sprintf("status: %d, body: %s\n", tc.status, tc.message))
	}
}

func TestRetries(t *testing.T) {
	i := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the body was sent in every retried request.
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, []byte("hello"), body)
		if i == 0 {
			i++
			http.Error(w, "failure!", http.StatusInternalServerError)
			return
		}
		_, err = w.Write([]byte("OK"))
		require.NoError(t, err)
	}))
	defer s.Close()
	c, err := New(s.URL, Config{
		NumRetries: 2,
	})
	require.NoError(t, err)

	reqBody := bytes.NewReader([]byte("hello"))
	err = c.request(context.Background(), "GET", "/", nil, reqBody, nil)
	assert.NoError(t, err)
}
