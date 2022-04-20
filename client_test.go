package bunq_ledger

import (
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

func TestParseResponseBody(t *testing.T) {
	body := strings.NewReader("{\"Response\":[{\"Id\":{\"id\":4570524}},{\"Token\":{\"id\":18186516,\"created\":\"2022-04-19 14:24:12.594108\",\"updated\":\"2022-04-19 14:24:12.594108\",\"token\":\"1b7b362557436e270519c0ec2a6714024da32fb642e375ac5ccbda4f9ee879c0\"}},{\"ServerPublicKey\":{\"server_public_key\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3\\/0m0tvJ58mP6itSQ11e\\nEsfGgsGg\\/RVtWz+fl7IYJKhUOTHhqgqJp6notPbk82LwGvpF8LuWYBgFOI7ou2Dd\\nuQaSpZcz+KQaBdxDjpPqWYUMB+lZYR8DhJOyUPVf4i4knjs7Op0mMQdwDR+fjZKz\\n\\/MMw1PR3fPUQmUyJWhUDMWI76s+zQDRLbkiZAClTJ3+jZ8576TxlG4DkTdmdLney\\nkjzyF70bXGarVNVBAt5VHIu6h8Ir3KMvgLpGdCznxqYaRMTsDjJlTPYzF0xH2hBj\\nOgtzl28dlvPO94CE7ek3611In7CBms8ml1mb1LvTz1eVKF2XnG4YgKBxhSD6PaOH\\n5QIDAQAB\\n-----END PUBLIC KEY-----\\n\"}}]}")

	result, err := parseInstallationResponseBody(body)
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Logf("%+v", result)
}
