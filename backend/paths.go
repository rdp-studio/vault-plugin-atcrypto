package backend

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathKeysList() *framework.Path {
	return &framework.Path{
		Pattern: "keys/?$",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.listKeys,
		},
		HelpSynopsis: "List keys maintained by the plugin backend.",
	}
}

func (b *backend) pathKeysOperation() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name"),
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.createKey,
			logical.DeleteOperation: b.deleteKey,
			logical.ReadOperation:   b.readKey,
		},
		HelpSynopsis: "Create a new key, import an existing key, read a key or delete a key.",
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the key.",
			},
			"exportable": {
				Type:        framework.TypeBool,
				Description: "Enables key to be exportable.",
				Default:     false,
			},
			"privateKey": {
				Type:        framework.TypeString,
				Description: "Hexidecimal string of the private key (32-byte or 64-char long). If present, the request will import the given key instead of generating a new key.",
				Default:     "",
			},
			"algorithm": {
				Type:        framework.TypeString,
				Description: "Algorithm of the key (e.g., secp256k1).",
				Default:     "secp256k1",
			},
		},
	}
}

func (b *backend) pathKeysExport() *framework.Path {
	return &framework.Path{
		Pattern: "export/keys/" + framework.GenericNameRegex("name"),
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.exportKey,
		},
		HelpSynopsis: "Export a key.",
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the key.",
			},
		},
	}
}

func (b *backend) pathKeysSign() *framework.Path {
	return &framework.Path{
		Pattern: "sign/keys/" + framework.GenericNameRegex("name"),
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.sign,
		},
		HelpSynopsis: "Sign a message with a key.",
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the key.",
			},
			"message": {
				Type:        framework.TypeString,
				Description: "Message to sign. Should be hexidecimal string of the message.",
			},
		},
	}
}

func (b *backend) paths() []*framework.Path {
	return []*framework.Path{
		b.pathKeysList(),
		b.pathKeysOperation(),
		b.pathKeysExport(),
		b.pathKeysSign(),
	}
}
