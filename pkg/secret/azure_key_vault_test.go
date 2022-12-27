package secret

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	secretName   = "secret-identifier"
	secretValue  = "secret-value"
	errorMessage = "arbitrary error from azsecrets"
)

type AzureKeyVaultTestSuite struct {
	suite.Suite
	vault  azureKeyVault
	client *MockAzureKeyVaultClient
	ctx    context.Context
}

func TestAzureKeyVault(t *testing.T) {
	suite.Run(t, new(AzureKeyVaultTestSuite))
}

func (s *AzureKeyVaultTestSuite) SetupTest() {
	s.client = NewMockAzureKeyVaultClient(gomock.NewController(s.T()))
	s.vault = azureKeyVault{s.client}
	s.ctx = context.Background()
}

func (s *AzureKeyVaultTestSuite) TestStoreSecret() {
	azParams := azsecrets.SetSecretParameters{Value: stringPtrFor(secretValue)}
	s.client.EXPECT().
		SetSecret(gomock.Eq(s.ctx), gomock.Eq(secretName), gomock.Eq(azParams), gomock.Nil()).
		Return(azsecrets.SetSecretResponse{}, nil)
	err := s.vault.StoreSecret(s.ctx, secretName, secretValue)
	s.Nil(err)
}

func (s *AzureKeyVaultTestSuite) TestStoreSecretWithErrs() {
	azParams := azsecrets.SetSecretParameters{Value: stringPtrFor(secretValue)}
	s.client.EXPECT().
		SetSecret(gomock.Eq(s.ctx), gomock.Eq(secretName), gomock.Eq(azParams), gomock.Nil()).
		Return(azsecrets.SetSecretResponse{}, errors.New(errorMessage))

	err := s.vault.StoreSecret(s.ctx, secretName, secretValue)
	s.Error(err)
	s.Equal(err.Error(), "problem setting secret: "+errorMessage)
}

func (s *AzureKeyVaultTestSuite) TestGetSecret() {
	s.client.EXPECT().
		GetSecret(gomock.Eq(s.ctx), gomock.Eq(secretName), gomock.Eq(""), gomock.Nil()).
		Return(
			azsecrets.GetSecretResponse{
				SecretBundle: azsecrets.SecretBundle{
					Value: stringPtrFor(secretValue),
				}}, nil)

	secret, found, err := s.vault.ReadSecret(s.ctx, secretName)
	s.Nil(err)
	s.True(found)
	s.Equal(secret, secretValue)
}

func (s *AzureKeyVaultTestSuite) TestGetSecretWithSecretNotFoundErr() {
	err := newSecretNotFoundError()
	s.client.EXPECT().
		GetSecret(gomock.Eq(s.ctx), gomock.Eq(secretName), gomock.Eq(""), gomock.Nil()).
		Return(azsecrets.GetSecretResponse{}, err)

	_, found, err := s.vault.ReadSecret(s.ctx, secretName)
	s.Nil(err)
	s.False(found)
}

func (s *AzureKeyVaultTestSuite) TestGetSecretWithErrs() {
	s.client.EXPECT().
		GetSecret(gomock.Eq(s.ctx), gomock.Eq(secretName), gomock.Eq(""), gomock.Nil()).
		Return(azsecrets.GetSecretResponse{}, errors.New(errorMessage))

	val, found, err := s.vault.ReadSecret(s.ctx, secretName)
	s.NotNil(err)
	s.Equal(err.Error(), "problem reading secret: "+errorMessage)
	s.Len(val, 0)
	s.False(found)
}

func (s *AzureKeyVaultTestSuite) TestDeleteSecret() {
	s.client.EXPECT().
		DeleteSecret(gomock.Eq(s.ctx), gomock.Eq(secretName), gomock.Nil()).
		Return(azsecrets.DeleteSecretResponse{}, nil)
	found, err := s.vault.DeleteSecret(s.ctx, secretName)
	s.Nil(err)
	s.True(found)
}

func (s *AzureKeyVaultTestSuite) TestDeleteSecretWithSecretNotFoundErr() {
	s.client.EXPECT().
		DeleteSecret(gomock.Eq(s.ctx), gomock.Eq(secretName), gomock.Nil()).
		Return(azsecrets.DeleteSecretResponse{}, newSecretNotFoundError())
	found, err := s.vault.DeleteSecret(s.ctx, secretName)
	s.Nil(err)
	s.False(found)
}

func (s *AzureKeyVaultTestSuite) TestDeleteSecretWithErrs() {
	s.client.EXPECT().
		DeleteSecret(gomock.Eq(s.ctx), gomock.Eq(secretName), gomock.Nil()).
		Return(azsecrets.DeleteSecretResponse{}, errors.New(errorMessage))
	found, err := s.vault.DeleteSecret(s.ctx, secretName)

	s.NotNil(err)
	s.False(found)
}

func TestIsSecretNotFoundFalseCases(t *testing.T) {
	testCases := []error{
		&azcore.ResponseError{StatusCode: 401},
		&azcore.ResponseError{StatusCode: 403},
		&azcore.ResponseError{StatusCode: 500},
	}

	for _, test := range testCases {
		respValue := isSecretNotFound(test)
		assert.False(t, respValue)
	}
}
func TestIsSecretNotFoundError(t *testing.T) {
	respValue := isSecretNotFound(newSecretNotFoundError())
	assert.True(t, respValue)
}

func newSecretNotFoundError() error {
	return &azcore.ResponseError{
		StatusCode: 404,
	}
}
