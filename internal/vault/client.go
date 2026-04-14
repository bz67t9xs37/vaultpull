package vault

import (
	"context"
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	api    *vaultapi.Client
	Mount  string
}

// Config holds the configuration needed to connect to Vault.
type Config struct {
	Address string
	Token   string
	Mount   string
}

// NewClient creates a new Vault client from the given config.
// Falls back to VAULT_ADDR and VAULT_TOKEN environment variables if not set.
func NewClient(cfg Config) (*Client, error) {
	vcfg := vaultapi.DefaultConfig()

	if cfg.Address != "" {
		vcfg.Address = cfg.Address
	} else if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		vcfg.Address = addr
	} else {
		vcfg.Address = "http://127.0.0.1:8200"
	}

	client, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required (set VAULT_TOKEN or pass --token)")
	}
	client.SetToken(token)

	mount := cfg.Mount
	if mount == "" {
		mount = "secret"
	}

	return &Client{api: client, Mount: mount}, nil
}

// GetSecrets reads a KV v2 secret at the given path and returns a map of key→value.
func (c *Client) GetSecrets(ctx context.Context, path string) (map[string]string, error) {
	secretPath := fmt.Sprintf("%s/data/%s", c.Mount, path)

	secret, err := c.api.Logical().ReadWithContext(ctx, secretPath)
	if err != nil {
		return nil, fmt.Errorf("reading vault path %q: %w", secretPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no secret found at path %q", secretPath)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected secret format at path %q", secretPath)
	}

	result := make(map[string]string, len(data))
	for k, v := range data {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}
