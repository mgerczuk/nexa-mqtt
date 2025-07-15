package growatt_web

import (
	"errors"
	"net/http/cookiejar"
	"net/url"
	"testing"
	"time"

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

	mockHttpClient.OnLogin("user", "secret", GrowattResult{}, nil)

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

	mockHttpClient.OnLogin("user", "secret", GrowattResult{}, errors.New("login failed"))

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

	mockHttpClient.OnLogin("user", "secret", GrowattResult{}, nil)

	err := client.Login()

	assert.Nil(t, err)
}

func TestLogin_ErrResult(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnLogin("user", "secret", GrowattResult{Result: -1}, nil)

	err := client.Login()

	assert.NotNil(t, err)
}

func TestLogin_Fails(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnLogin("user", "secret", GrowattResult{}, errors.New("login fails"))

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

func TestGetPlantDevices_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetPlantDevices("pid", GrowattPlantDevices{}, nil)

	result, err := client.GetPlantDevices("pid")

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestGetPlantDevices_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetPlantDevices("pid", GrowattPlantDevices{}, errors.New("GetPlantDevices fails"))

	result, err := client.GetPlantDevices("pid")

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestGetNoahList_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahList(12, GrowattNoahList{}, nil)

	result, err := client.GetNoahList(12)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestGetNoahList_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahList(12, GrowattNoahList{}, errors.New("GetNoahList fails"))

	result, err := client.GetNoahList(12)

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestGetNoahDetails_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahDetails(12, "Serial123", GrowattNoahList{}, nil)

	result, err := client.GetNoahDetails(12, "Serial123")

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestGetNoahDetails_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahDetails(12, "Serial123", GrowattNoahList{}, errors.New("GetNoahDetails fails"))

	result, err := client.GetNoahDetails(12, "Serial123")

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestGetNoahHistory_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahHistory("Serial123", "2025-07-01", "2025-07-02", GrowattNoahHistory{}, nil)

	result, err := client.GetNoahHistory("Serial123", "2025-07-01", "2025-07-02")

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestGetNoahHistory_EmptyDate(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	expected := time.Now().Format("2006-01-02")
	mockHttpClient.OnGetNoahHistory("Serial123", expected, expected, GrowattNoahHistory{}, nil)

	result, err := client.GetNoahHistory("Serial123", "", "")

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestGetNoahHistory_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahHistory("Serial123", "2025-07-01", "2025-07-02", GrowattNoahHistory{}, errors.New("GetNoahHistory fails"))

	result, err := client.GetNoahHistory("Serial123", "2025-07-01", "2025-07-02")

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestGetNoahStatus_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahStatus(42, "Serial123", GrowattNoahStatus{}, nil)

	result, err := client.GetNoahStatus(42, "Serial123")

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestGetNoahStatus_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahStatus(42, "Serial123", GrowattNoahStatus{}, errors.New("GetNoahStatus fails"))

	result, err := client.GetNoahStatus(42, "Serial123")

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestGetNoahTotals_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahTotals(42, "Serial123", GrowattNoahTotals{}, nil)

	result, err := client.GetNoahTotals(42, "Serial123")

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestGetNoahTotals_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahTotals(42, "Serial123", GrowattNoahTotals{}, errors.New("GetNoahTotals fails"))

	result, err := client.GetNoahTotals(42, "Serial123")

	assert.Nil(t, result)
	assert.NotNil(t, err)
}
