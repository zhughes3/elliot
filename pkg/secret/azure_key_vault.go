package secret

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

// azureKeyVaultCLient are the azure golang sdk functions we use for azure key vault client
type AzureKeyVaultClient interface {
	SetSecret(ctx context.Context, name string, parameters azsecrets.SetSecretParameters, options *azsecrets.SetSecretOptions) (azsecrets.SetSecretResponse, error)
	GetSecret(ctx context.Context, name string, version string, options *azsecrets.GetSecretOptions) (azsecrets.GetSecretResponse, error)
	DeleteSecret(ctx context.Context, name string, options *azsecrets.DeleteSecretOptions) (azsecrets.DeleteSecretResponse, error)
}

// AzureKeyVault represents a client connection to Azure Key Vault
type azureKeyVault struct {
	client AzureKeyVaultClient
}

// NewAzureKeyVault attempts to authenticate with Azure and creates an Azure Key Vault client
func NewAzureKeyVaultClient(uri string) (AzureKeyVaultClient, error) {
	if len(strings.TrimSpace(uri)) == 0 {
		return nil, errors.New("the uri cannot be empty")
	}
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("problem getting default Azure credential: %w", err)
	}
	client, err := azsecrets.NewClient(uri, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("problem creating Azure Key Vault client: %w", err)
	}

	return client, nil
}

// NewAzureKeyVault returns a KeyVault that uses the given AzureKeyVaultClient.
func NewAzureKeyVault(client AzureKeyVaultClient) KeyVault {
	return azureKeyVault{client: client}
}

// CreateSecret attempts to create a new secret in Azure Key Vault
func (a azureKeyVault) StoreSecret(ctx context.Context, name, value string) error {
	return a.doStoreSecret(ctx, name, azsecrets.SetSecretParameters{Value: &value})
}

func (a azureKeyVault) doStoreSecret(ctx context.Context, name string, params azsecrets.SetSecretParameters) error {
	_, err := a.client.SetSecret(ctx, name, params, nil)
	if err != nil {
		return fmt.Errorf("problem setting secret: %w", err)
	}

	return nil
}

// ReadSecret fetches an Azure Key Vault secret by name
func (a azureKeyVault) ReadSecret(ctx context.Context, name string) (string, bool, error) {
	resp, found, err := a.doReadSecret(ctx, name)
	if err != nil || !found {
		return "", found, err
	}

	return stringValueFor(resp.Value), true, nil
}

func (a azureKeyVault) doReadSecret(ctx context.Context, name string) (azsecrets.GetSecretResponse, bool, error) {
	resp, err := a.client.GetSecret(ctx, name, "", nil)
	if err != nil {
		if isSecretNotFound(err) {
			return azsecrets.GetSecretResponse{}, false, nil
		}
		return azsecrets.GetSecretResponse{}, false, fmt.Errorf("problem reading secret: %w", err)
	}

	return resp, true, nil
}

// DeleteSecret attempts to delete an Azure Key Vault secret by identifier
func (a azureKeyVault) DeleteSecret(ctx context.Context, name string) (bool, error) {
	_, err := a.client.DeleteSecret(ctx, name, nil)
	if err != nil {
		if isSecretNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("problem deleting secret: %w", err)
	}

	return true, nil
}

// isSecretNotFound returns true if the given error is a ResponseError from
// Azure with a StatusCode of 404 (not found), otherwise false. More details here:
// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore#ResponseError
func isSecretNotFound(err error) bool {
	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) {
		if respErr.StatusCode == http.StatusNotFound {
			return true
		}
	}
	return false
}

func stringValueFor(in *string) string {
	if in != nil {
		return *in
	}
	return ""
}

func stringPtrFor(in string) *string {
	return &in
}
