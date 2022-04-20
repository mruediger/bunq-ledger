package bunq_ledger

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	URL         url.URL
	private_key rsa.PrivateKey
}

func NewClient(address string) (*Client, error) {
	addressURL, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("bunq-ledger: could not parse entpoint address %q: %w", address, err)
	}

	private_key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("bunq-ledger: could not create the keypair %w", err)
	}

	c := &Client{
		URL:         *addressURL,
		private_key: *private_key,
	}

	return c, nil
}

func (c *Client) postInstallation() (string, error) {
	url := c.URL
	url.Path = "/v1/installation"

	pubKey, err := x509.MarshalPKIXPublicKey(c.private_key.Public())
	if err != nil {
		return "", fmt.Errorf("bunq-ledger: cannot marshall public key %w", err)
	}

	pubKeyPemBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubKey,
	}

	pubKeyString := string(pem.EncodeToMemory(&pubKeyPemBlock))

	postBody, err := json.Marshal(map[string]string{
		"client_public_key": pubKeyString,
	})
	if err != nil {
		return "", fmt.Errorf("bunq-ledger: cannot marshal post body %w", err)
	}

	r, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(postBody))
	if err != nil {
		return "", fmt.Errorf("bunq-ledger: error creating request %w", err)
	}
	r.Header.Add("Cache-Control", "no-cache")
	r.Header.Add("User-Agent", "bunq-ledger-test")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return "", fmt.Errorf("bunq-ledger: error performing request %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	sb := string(body)
	fmt.Printf("Body: %q", sb)

	parseInstallationResponseBody(resp.Body)

	return "foo", nil
}

type InstallationResponse struct {
	Token           string `json:"token"`
	Created         string `json:"created"`
	Updated         string `json:"updated"`
	ServerPublicKey string `json:"server_public_key"`
}

func parseInstallationResponseBody(body io.Reader) (InstallationResponse, error) {
	var response map[string][]map[string]json.RawMessage
	err := json.NewDecoder(body).Decode(&response)
	if err != nil {
		return InstallationResponse{}, err
	}

	var result InstallationResponse
	for _, thing := range response["Response"] {
		for _, v := range thing {
			err := json.Unmarshal(v, &result)
			if err != nil {
				return result, err
			}
		}
	}
	return result, nil
}

func (c *Client) postDeviceServer(token string, apiKey string) error {
	url := c.URL
	url.Path = "/v1/device-server"

	postBody, err := json.Marshal(map[string]string{
		"secret": apiKey,
	})
	if err != nil {
		return fmt.Errorf("bunq-ledger: cannot marshal post body %w", err)
	}
	r, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(postBody))
	if err != nil {
		return fmt.Errorf("bunq-ledger: error creating request %w", err)
	}

	r.Header.Add("Cache-Control", "no-cache")
	r.Header.Add("User-Agent", "bunq-ledger-test")
	r.Header.Add("X-Bunq-Client-Authentication", token)

	return nil
}

func (c *Client) postSessionServer(token string, apiKey string) error {
	return nil
}
