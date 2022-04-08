package bunq_ledger

import (
	"fmt"
	"net/url"
)

type Client struct {
	URL url.URL
}

func NewClient(address string) (*Client, error) {
	addressURL, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("bunq-ledger: could not parse entpoint address %q: %w", address, err)
	}

	c := &Client{
		URL: *addressURL,
	}

	return c, nil
}
