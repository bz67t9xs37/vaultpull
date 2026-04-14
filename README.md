# vaultpull

> CLI tool to sync HashiCorp Vault secrets into local `.env` files with diff-aware updates.

---

## Installation

```bash
go install github.com/yourname/vaultpull@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourname/vaultpull/releases).

---

## Usage

```bash
# Authenticate and pull secrets from a Vault path into a local .env file
vaultpull --addr https://vault.example.com \
          --path secret/data/myapp \
          --output .env

# Preview changes without writing to disk
vaultpull --path secret/data/myapp --output .env --dry-run
```

vaultpull will compare the remote secrets against your existing `.env` file and display a diff before applying any changes, so you always know what's being updated.

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--addr` | Vault server address | `$VAULT_ADDR` |
| `--token` | Vault token | `$VAULT_TOKEN` |
| `--path` | Secret path in Vault | *(required)* |
| `--output` | Output `.env` file path | `.env` |
| `--dry-run` | Preview diff without writing | `false` |

---

## Authentication

vaultpull respects standard Vault environment variables:

```bash
export VAULT_ADDR=https://vault.example.com
export VAULT_TOKEN=s.xxxxxxxxxxxxxxxx
```

---

## License

[MIT](LICENSE) © yourname