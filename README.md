# vault-plugin-atcrypto

A HashiCorp Vault plugin that provides secure key management and cryptographic operations for AT Protocol, supporting secp256k1 and P-256 curves.

## Features

- **Key Management**: Create, import, read, and delete cryptographic keys
- **Supported Algorithms**: secp256k1 and P-256 (NIST)
- **Key Export**: Configurable exportable keys for secure backup
- **Digital Signatures**: Sign messages using stored keys
- **Secure Storage**: Keys are stored using Vault's seal wrap encryption
- **AT Protocol Integration**: Uses Bluesky's atcrypto library for AT Protocol compatibility

## Requirements

- Go 1.25.0 or later
- HashiCorp Vault 1.15.0 or later

## Installation

### Build the Plugin

```bash
go build -o vault-plugin-atcrypto
```

### Register the Plugin with Vault

1. Move the compiled binary to your Vault plugin directory:
```bash
mv vault-plugin-atcrypto /path/to/vault/plugins/
```

2. Calculate the SHA256 checksum:
```bash
sha256sum vault-plugin-atcrypto
```

3. Register the plugin with Vault:
```bash
vault plugin register \
  -sha256=<SHA256_CHECKSUM> \
  secret vault-plugin-atcrypto
```

4. Enable the secrets engine:
```bash
vault secrets enable -path=atcrypto vault-plugin-atcrypto
```

## Usage

### Create a New Key

Generate a new secp256k1 key:

```bash
vault write atcrypto/keys/my-key \
  algorithm=secp256k1 \
  exportable=false
```

Generate a new P-256 key:

```bash
vault write atcrypto/keys/my-key \
  algorithm=p256 \
  exportable=false
```

### Import an Existing Key

Import a key from hex-encoded private key bytes (32 bytes):

```bash
vault write atcrypto/keys/my-imported-key \
  algorithm=secp256k1 \
  privateKey=<HEX_PRIVATE_KEY> \
  exportable=true
```

### Read Key Information

```bash
vault read atcrypto/keys/my-key
```

Response:
```json
{
  "data": {
    "exportable": false,
    "public_key": "z6Mk..."
  }
}
```

### List All Keys

```bash
vault list atcrypto/keys
```

### Export a Key

Only available if the key was created with `exportable=true`:

```bash
vault read atcrypto/export/keys/my-key
```

Response:
```json
{
  "data": {
    "public_key": "z6Mk...",
    "private_key": "z6Mk..."
  }
}
```

### Sign a Message

Sign a hex-encoded message using a stored key:

```bash
vault write atcrypto/sign/keys/my-key \
  message=<HEX_MESSAGE>
```

Response:
```json
{
  "data": {
    "signature": "<HEX_SIGNATURE>"
  }
}
```

### Delete a Key

```bash
vault delete atcrypto/keys/my-key
```

## API Reference

| Path | Method | Description |
|------|--------|-------------|
| `atcrypto/keys` | LIST | List all keys |
| `atcrypto/keys/:name` | POST | Create or import a key |
| `atcrypto/keys/:name` | GET | Read key information |
| `atcrypto/keys/:name` | DELETE | Delete a key |
| `atcrypto/export/keys/:name` | GET | Export a key (if exportable) |
| `atcrypto/sign/keys/:name` | POST | Sign a message |

## Development

### Project Structure

```
vault-plugin-atcrypto/
├── backend/
│   ├── backend.go    # Backend initialization and factory
│   ├── keys.go       # Key management operations
│   └── paths.go      # API path definitions
├── devenv/           # Development environment configuration
├── main.go           # Plugin entry point
├── go.mod            # Go module definition
└── go.sum            # Go module checksums
```

### Local Development

1. Start a development Vault instance:

```bash
vault server -dev -dev-root-token-id=root -dev-plugin-dir=./devenv/plugins
```

2. Build and register the plugin:

```bash
go build -o devenv/plugins/vault-plugin-atcrypto
vault plugin register -sha256=$(sha256sum devenv/plugins/vault-plugin-atcrypto | cut -d' ' -f1) secret vault-plugin-atcrypto
vault secrets enable -path=atcrypto vault-plugin-atcrypto
```

## Dependencies

- [github.com/hashicorp/vault/api](https://github.com/hashicorp/vault) - Vault API client
- [github.com/hashicorp/vault/sdk](https://github.com/hashicorp/vault) - Vault SDK for plugin development
- [github.com/bluesky-social/indigo](https://github.com/bluesky-social/indigo) - AT Protocol cryptographic library

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
