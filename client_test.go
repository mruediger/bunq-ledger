package bunq_ledger

import (
	"net/http"
	"testing"
)

const sandboxUrl = "https://public-api.sandbox.bunq.com"

func TestCreateURL(t *testing.T) {
	c, err := NewClient(sandboxUrl)
	if err != nil {
		t.Fatalf("%v", err)
	}

	c.URL.Path = "/v1/installation"
	if c.URL.String() != "https://public-api.sandbox.bunq.com/v1/installation" {
		t.Errorf("error")
	}
}

func TestRegisterDevice(t *testing.T) {
	c, err := NewClient(sandboxUrl)
	if err != nil {
		t.Fatalf("%v", err)
	}
	c.URL.Path = "/v1/installation"

	r, err := http.NewRequest(http.MethodPost, c.URL.String(), nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	r.Header.Add("Cache-Control", "no-cache")
	r.Header.Add("User-Agent", "bunq-ledger-test")

	t.Logf("%v", r)

	client := &http.Client{}
	resp, err := client.Do(r)

	t.Logf("%v", resp)
	t.Logf("%v", err)
}
