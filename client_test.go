package bunq_ledger

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
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

	apiKey := "foo"

	token, err := c.postInstallation()
	if err != nil {
		t.Fatalf("%v", err)
	}

	c.postDeviceServer(token, apiKey)
	c.postSessionServer(token, apiKey)
}

func OldTestRegisterDevice(t *testing.T) {
	c, err := NewClient(sandboxUrl)
	if err != nil {
		t.Fatalf("%v", err)
	}
	c.URL.Path = "/v1/installation"

	private_key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Logf("%v", err)
	}

	pubKey, err := x509.MarshalPKIXPublicKey(private_key.Public())
	if err != nil {
		t.Logf("%v", err)
	}

	pubKeyPemBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubKey,
	}

	pubKeyString := string(pem.EncodeToMemory(&pubKeyPemBlock))

	t.Logf("PublicKey: %v", pubKeyString)

	postBody, _ := json.Marshal(map[string]string{
		"client_public_key": pubKeyString,
	})

	r, err := http.NewRequest(http.MethodPost, c.URL.String(), bytes.NewBuffer(postBody))
	if err != nil {
		t.Fatalf("%v", err)
	}
	r.Header.Add("Cache-Control", "no-cache")
	r.Header.Add("User-Agent", "bunq-ledger-test")

	t.Logf("%v", r)

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		t.Logf("%v", err)
	}

	t.Logf("Resp: %v", resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Logf("%v", err)
	}

	sb := string(body)
	t.Logf("Body: %s", sb)
}
