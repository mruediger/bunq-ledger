package bunq_ledger

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
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

func (c *Client) postInstallation() error {
	url := c.URL
	url.Path = "/v1/installation"

	pubKey, err := x509.MarshalPKIXPublicKey(c.private_key.Public())
	if err != nil {
		return fmt.Errorf("foo")
	}

	pubKeyPemBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubKey,
	}

	return nil
}
