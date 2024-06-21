package main

import (
	"net/http"
	"snippetbox.rahilganatra.net/internal/assert"
	"testing"
)

func TestPing(t *testing.T) {

	//indicates its ok to run test concurrently alongside other tests. These tests will be run
	//parallel with tests marked with Parallel() only
	t.Parallel()

	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	code, _, body := ts.get(t, "/ping")
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}
