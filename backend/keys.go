package backend

import (
	"context"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"

	"github.com/bluesky-social/indigo/atproto/atcrypto"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type key struct {
	Exportable bool   `json:"exportable"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

func (k *key) GetPrivateKey() (atcrypto.PrivateKey, error) {
	return atcrypto.ParsePrivateMultibase(k.PrivateKey)
}

func (k *key) GetPublicKey() (atcrypto.PublicKey, error) {
	return atcrypto.ParsePublicMultibase(k.PublicKey)
}

func (b *backend) listKeys(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "keys/")
	if err != nil {
		return nil, err
	}

	entries = slices.DeleteFunc(entries, func(s string) bool {
		if strings.HasSuffix(s, "/") {
			return true
		}

		return false
	})

	return logical.ListResponse(entries), nil
}

func (b *backend) createKey(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	keyName := data.Get("name").(string)
	keyInput := data.Get("privateKey").(string)
	exportable := data.Get("exportable").(bool)
	algorithm := data.Get("algorithm").(string)
	var privKey atcrypto.PrivateKeyExportable
	var err error

	if keyInput != "" {
		keyBytes, err := hex.DecodeString(keyInput)
		if err != nil {
			return nil, err
		}
		if len(keyBytes) != 32 {
			return nil, fmt.Errorf("invalid private key length: %d", len(keyBytes))
		}

		switch algorithm {
		case "secp256k1":
			privKey, err = atcrypto.ParsePrivateBytesK256(keyBytes)
			if err != nil {
				return nil, err
			}
		case "p256":
			privKey, err = atcrypto.ParsePrivateBytesP256(keyBytes)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
		}
	} else {
		switch algorithm {
		case "secp256k1":
			privKey, err = atcrypto.GeneratePrivateKeyK256()
			if err != nil {
				return nil, err
			}
		case "p256":
			privKey, err = atcrypto.GeneratePrivateKeyP256()
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
		}
	}

	pubKey, err := privKey.PublicKey()
	if err != nil {
		return nil, err
	}

	key := &key{
		Exportable: exportable,
		PrivateKey: privKey.Multibase(),
		PublicKey:  pubKey.Multibase(),
	}

	entry, err := logical.StorageEntryJSON(fmt.Sprintf("keys/%s", keyName), key)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"public_key": key.PublicKey,
		},
	}, nil
}

func (b *backend) deleteKey(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	keyName := data.Get("name").(string)
	err := req.Storage.Delete(ctx, fmt.Sprintf("keys/%s", keyName))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) readKey(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	keyName := data.Get("name").(string)
	key, err := b.internalReadKey(ctx, req, keyName)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"exportable": key.Exportable,
			"public_key": key.PublicKey,
		},
	}, nil
}

func (b *backend) exportKey(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	keyName := data.Get("name").(string)
	key, err := b.internalReadKey(ctx, req, keyName)
	if err != nil {
		return nil, err
	}
	if !key.Exportable {
		return nil, fmt.Errorf("key is not exportable")
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"public_key":  key.PublicKey,
			"private_key": key.PrivateKey,
		},
	}, nil
}

func (b *backend) sign(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	keyName := data.Get("name").(string)
	message := data.Get("message").(string)
	key, err := b.internalReadKey(ctx, req, keyName)
	if err != nil {
		return nil, err
	}
	privKey, err := key.GetPrivateKey()
	if err != nil {
		return nil, err
	}
	messageBytes, err := hex.DecodeString(message)
	if err != nil {
		return nil, err
	}
	signature, err := privKey.HashAndSign(messageBytes)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"signature": hex.EncodeToString(signature),
		},
	}, nil
}

func (b *backend) internalReadKey(ctx context.Context, req *logical.Request, keyName string) (*key, error) {
	entry, err := req.Storage.Get(ctx, fmt.Sprintf("keys/%s", keyName))
	if err != nil {
		b.Logger().Error("Failed to read key", "error", err)
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf("key not found")
	}
	var keyData key
	err = entry.DecodeJSON(&keyData)
	if err != nil {
		b.Logger().Error("Failed to decode key", "error", err)
		return nil, err
	}
	return &keyData, nil
}
