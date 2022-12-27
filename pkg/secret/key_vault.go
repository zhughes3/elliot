package secret

import "context"

// KeyVault is an interface representing the functions we expose to CRUD secrets from key vault
type KeyVault interface {
	StoreSecret(ctx context.Context, name, value string) error
	ReadSecret(ctx context.Context, name string) (string, bool, error)
	DeleteSecret(ctx context.Context, name string) (bool, error)
}
