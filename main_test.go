package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

// test health, nohealth

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize()

	code := m.Run()

	os.Exit(code)
}

// test the health URL returns something proper
func TestHealth(t *testing.T) {
	health_req, _ := http.NewRequest("GET", "/health", nil)
	health_rsp := executeRequest(health_req)

	checkResponseCode(t, http.StatusOK, health_rsp.Code)

	if body := health_rsp.Body.String(); body != `"ok"` {
		t.Errorf("Expected an 'ok' message. Got %s", body)
	} else {
		fmt.Println("- Test OK: health check")
	}
}
