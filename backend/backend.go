package backend

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := initBackend()
	if err != nil {
		return nil, err
	}
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

type backend struct {
	*framework.Backend
}

func initBackend() (*backend, error) {
	var b backend
	b.Backend = &framework.Backend{
		Help: "The atcrypto secret backend provides a secure way to interact with did:plc rotation keys",
		Paths: framework.PathAppend(
			b.paths(),
		),
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"keys/",
			},
		},
		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}
	return &b, nil
}
