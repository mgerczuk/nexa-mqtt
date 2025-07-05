package growatt_web

import (
	"errors"
	"net/http/cookiejar"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupMocks(t *testing.T) (*MockHttpClient, *Client) {
	mockHttpClient := MockHttpClient{}
	jar, err := cookiejar.New(nil)
	assert.Nil(t, err)

	client := Client{
		client:    &mockHttpClient,
		serverUrl: "https://openapi.growatt.com",
		username:  "user",
		password:  "secret",
		jar:       jar,
	}

	return &mockHttpClient, &client
}

func Test_postForm_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"http://someurl",
		url.Values{},
		nil,
	).Return(nil)

	err := client.postForm("http://someurl", url.Values{}, nil)
	assert.NoError(t, err)

	mockHttpClient.AssertExpectations(t)
}

func Test_postForm_AnyError(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"http://someurl",
		url.Values{},
		nil,
	).Return(errors.New("some error"))

	err := client.postForm("http://someurl", url.Values{}, nil)
	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func Test_postForm_RetryLogin(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"http://someurl",
		url.Values{},
		nil,
	).Return(errors.New("invalid character '<' looking for beginning of value")).Once()

	mockHttpClient.OnLogin("user", "secret", nil)

	mockHttpClient.On(
		"postForm",
		"http://someurl",
		url.Values{},
		nil,
	).Return(nil)

	err := client.postForm("http://someurl", url.Values{}, nil)
	assert.NoError(t, err)

	mockHttpClient.AssertExpectations(t)
}

func Test_postForm_RetryLoginFails(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"http://someurl",
		url.Values{},
		nil,
	).Return(errors.New("invalid character '<' looking for beginning of value"))

	mockHttpClient.OnLogin("user", "secret", errors.New("login failed"))

	defer func() {
		if r := recover(); r != nil {
			mockHttpClient.AssertExpectations(t)
		}
	}()

	client.postForm("http://someurl", url.Values{}, nil)
	t.Errorf("Test failed, panic was expected")
}

func TestLogin_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnLogin("user", "secret", nil)

	err := client.Login()

	assert.Nil(t, err)
}

func TestLogin_Fails(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnLogin("user", "secret", errors.New("login fails"))

	err := client.Login()

	assert.NotNil(t, err)
}

func TestGetPlantList_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	expected := []GrowattPlant{
		{PlantId: "1", PlantName: "Test1"},
		{PlantId: "2", PlantName: "Test2"},
	}
	mockHttpClient.OnGetPlantList(expected, nil)

	result, err := client.GetPlantList()

	assert.Equal(t, expected, result)

	assert.Nil(t, err)
}

func TestGetPlantList_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	expected := []GrowattPlant{}
	mockHttpClient.OnGetPlantList(expected, errors.New("GetPlantList fails"))

	result, err := client.GetPlantList()

	assert.Nil(t, result)
	assert.NotNil(t, err)
}
