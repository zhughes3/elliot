package secret

import (
	"context"
	"errors"
	"testing"

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
	wrapper    AzureKeyVault
	mockClient *MockazureKeyVaultClient
	ctx        context.Context
}

func (suite *AzureKeyVaultTestSuite) SetupTest() {
	ctx := context.Background()
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	mockClient := NewMockazureKeyVaultClient(ctrl)
	wrapper := AzureKeyVault{client: mockClient}
	suite.wrapper = wrapper
	suite.mockClient = mockClient
	suite.ctx = ctx
}

func (suite *AzureKeyVaultTestSuite) TestStoreSecret() {
	val := secretValue
	azParams := azsecrets.SetSecretParameters{Value: &val}
	suite.mockClient.EXPECT().
		SetSecret(gomock.Eq(suite.ctx), gomock.Eq(secretName), gomock.Eq(azParams), gomock.Nil()).
		Return(azsecrets.SetSecretResponse{}, nil)
	err := suite.wrapper.StoreSecret(suite.ctx, secretName, secretValue)
	assert.Nil(suite.T(), err)
}

func (suite *AzureKeyVaultTestSuite) TestStoreSecretWithErrs() {
	val := secretValue
	azParams := azsecrets.SetSecretParameters{Value: &val}
	suite.mockClient.EXPECT().
		SetSecret(gomock.Eq(suite.ctx), gomock.Eq(secretName), gomock.Eq(azParams), gomock.Nil()).
		Return(azsecrets.SetSecretResponse{}, errors.New(errorMessage))

	err := suite.wrapper.StoreSecret(suite.ctx, secretName, secretValue)
	assert.Error(suite.T(), err)
	assert.ErrorContains(suite.T(), err, "problem creating Azure Key Vault secret")
}

func (suite *AzureKeyVaultTestSuite) TestGetSecret() {
	val := secretValue
	suite.mockClient.EXPECT().
		GetSecret(gomock.Eq(suite.ctx), gomock.Eq(secretName), gomock.Eq(""), gomock.Nil()).
		Return(
			azsecrets.GetSecretResponse{
				SecretBundle: azsecrets.SecretBundle{
					Value: &val,
				}}, nil)

	secret, err := suite.wrapper.ReadSecret(suite.ctx, secretName)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), *secret, val)
}

func (suite *AzureKeyVaultTestSuite) TestGetSecretWithErrs() {
	suite.mockClient.EXPECT().
		GetSecret(gomock.Eq(suite.ctx), gomock.Eq(secretName), gomock.Eq(""), gomock.Nil()).
		Return(azsecrets.GetSecretResponse{}, errors.New(errorMessage))

	val, err := suite.wrapper.ReadSecret(suite.ctx, secretName)
	assert.NotNil(suite.T(), err)
	assert.Len(suite.T(), *val, 0)
}

func (suite *AzureKeyVaultTestSuite) TestDeleteSecret() {
	suite.mockClient.EXPECT().
		DeleteSecret(gomock.Eq(suite.ctx), gomock.Eq(secretName), gomock.Nil()).
		Return(azsecrets.DeleteSecretResponse{}, nil)
	err := suite.wrapper.DeleteSecret(suite.ctx, secretName)
	assert.Nil(suite.T(), err)
}

func (suite *AzureKeyVaultTestSuite) TestDeleteSecretWithErrs() {
	suite.mockClient.EXPECT().
		DeleteSecret(gomock.Eq(suite.ctx), gomock.Eq(secretName), gomock.Nil()).
		Return(azsecrets.DeleteSecretResponse{}, errors.New(errorMessage))
	err := suite.wrapper.DeleteSecret(suite.ctx, secretName)
	assert.NotNil(suite.T(), err)
}

func TestAzureKeyVaultTestSuite(t *testing.T) {
	suite.Run(t, new(AzureKeyVaultTestSuite))
}
