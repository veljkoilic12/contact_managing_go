package main

import (
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {

	// spin up the server and defer server closing till the end of the test
	app := newTestApp()
	ts := newTestServer(app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/v1/healthcheck")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	expResp := `{
	"status": "available",
	"system_info": {
		"environment": "testing",
		"version": "1.0.0"
	}
}
`
	if string(body) != expResp {
		t.Errorf("want %q; got %q", expResp, string(body))
	}
}
