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
	"strings"
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

func TMPTestRegisterDevice(t *testing.T) {
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

func TestParseResponseBody(t *testing.T) {
	body := strings.NewReader("{\"Response\":[{\"Id\":{\"id\":4570524}},{\"Token\":{\"id\":18186516,\"created\":\"2022-04-19 14:24:12.594108\",\"updated\":\"2022-04-19 14:24:12.594108\",\"token\":\"1b7b362557436e270519c0ec2a6714024da32fb642e375ac5ccbda4f9ee879c0\"}},{\"ServerPublicKey\":{\"server_public_key\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3\\/0m0tvJ58mP6itSQ11e\\nEsfGgsGg\\/RVtWz+fl7IYJKhUOTHhqgqJp6notPbk82LwGvpF8LuWYBgFOI7ou2Dd\\nuQaSpZcz+KQaBdxDjpPqWYUMB+lZYR8DhJOyUPVf4i4knjs7Op0mMQdwDR+fjZKz\\n\\/MMw1PR3fPUQmUyJWhUDMWI76s+zQDRLbkiZAClTJ3+jZ8576TxlG4DkTdmdLney\\nkjzyF70bXGarVNVBAt5VHIu6h8Ir3KMvgLpGdCznxqYaRMTsDjJlTPYzF0xH2hBj\\nOgtzl28dlvPO94CE7ek3611In7CBms8ml1mb1LvTz1eVKF2XnG4YgKBxhSD6PaOH\\n5QIDAQAB\\n-----END PUBLIC KEY-----\\n\"}}]}")

	type installation struct {
		Token           string `json:"token"`
		ServerPublicKey string `json:"server_public_key"`
	}

	var instresp map[string][]map[string]json.RawMessage

	err := json.NewDecoder(body).Decode(&instresp)
	if err != nil {
		t.Fatalf("%v", err)
	}

	var inst installation

	for _, thing := range instresp["Response"] {
		for _, v := range thing {
			json.Unmarshal(v, &inst)
		}
	}

	t.Logf("%+v", inst)
}
