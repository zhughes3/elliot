package secret

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

// AzureKeyVaultClientFunctions are the azure golang sdk functions we use for azure key vault client
type AzureKeyVaultClientFunctions interface {
	SetSecret(ctx context.Context, name string, parameters azsecrets.SetSecretParameters, options *azsecrets.SetSecretOptions) (azsecrets.SetSecretResponse, error)
	GetSecret(ctx context.Context, name string, version string, options *azsecrets.GetSecretOptions) (azsecrets.GetSecretResponse, error)
	DeleteSecret(ctx context.Context, name string, options *azsecrets.DeleteSecretOptions) (azsecrets.DeleteSecretResponse, error)
}

// AzureKeyVault represents a client connection to Azure Key Vault
type AzureKeyVault struct {
	client AzureKeyVaultClientFunctions
}

// NewAzureKeyVault attempts to authenticate with Azure and creates an Azure Key Vault client
func NewAzureKeyVault(uri string) (AzureKeyVault, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return AzureKeyVault{}, fmt.Errorf("problem getting default Azure credential")
	}
	client, err := azsecrets.NewClient(uri, cred, nil)
	if err != nil {
		return AzureKeyVault{}, fmt.Errorf("problem creating Azure Key Vault client")
	}

	return AzureKeyVault{client: client}, nil
}

// CreateSecret attempts to create a new secret in Azure Key Vault
func (a AzureKeyVault) CreateSecret(ctx context.Context, name, value string) error {
	params := azsecrets.SetSecretParameters{Value: &value}
	_, err := a.client.SetSecret(ctx, name, params, nil)
	if err != nil {
		return fmt.Errorf("problem creating Azure Key Vault secret")
	}

	return nil
}

// ReadSecret fetches an Azure Key Vault secret by name
func (a AzureKeyVault) ReadSecret(ctx context.Context, name string, version string) (*string, error) {
	resp, err := a.client.GetSecret(ctx, name, version, nil)
	if err != nil {
		return new(string), fmt.Errorf("problem reading Azure Key Vault secret")
	}

	return resp.Value, nil
}

// DeleteSecret attempts to delete an Azure Key Vault secret by identifier
func (a AzureKeyVault) DeleteSecret(ctx context.Context, name string) error {
	_, err := a.client.DeleteSecret(ctx, name, nil)
	if err != nil {
		return fmt.Errorf("problem deleting Azure Key Vault secret")
	}

	return nil
}
