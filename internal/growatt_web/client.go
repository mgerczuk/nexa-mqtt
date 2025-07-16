package growatt_web

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"nexa-mqtt/internal/misc"
	"strings"
	"time"
)

type Client struct {
	client    HttpClient
	serverUrl string
	username  string
	password  string
	jar       *cookiejar.Jar
}

func newClient(serverUrl string, username string, password string) *Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		slog.Error("could not create cookie jar", slog.String("error", err.Error()))
		misc.Panic(err)
	}

	slog.Info("setting server url (web)", slog.String("url", serverUrl))

	return &Client{
		client: &httpClient{
			client: &http.Client{
				Transport:     nil,
				CheckRedirect: nil,
				Jar:           jar,
				Timeout:       10 * time.Second,
			},
		},
		serverUrl: serverUrl,
		username:  username,
		password:  password,
		jar:       jar,
	}
}

func (h *Client) postForm(url string, data url.Values, responseBody any) error {
	err := h.client.postForm(url, data, responseBody)
	if err != nil {
		notLoggedIn := strings.Contains(err.Error(), "invalid character '<' looking for beginning of value")
		if notLoggedIn {
			slog.Warn("re-login (web)", slog.String("error", err.Error()))
			if err := h.Login(); err != nil {
				slog.Error("could not re-login", slog.String("error", err.Error()))
				misc.Panic(err)
			}
			return h.postForm(url, data, responseBody)
		} else {
			return err
		}
	}

	return nil
}

func (h *Client) Login() error {
	var result GrowattResult
	if err := h.postForm(h.serverUrl+"/login", url.Values{
		"account":  {h.username},
		"password": {h.password},
	}, &result); err != nil {
		return err
	}
	if result.Result < 0 {
		return errors.New(result.Msg)
	}
	return nil
}

func (h *Client) GetPlantList() ([]GrowattPlant, error) {
	var result []GrowattPlant
	if err := h.postForm(h.serverUrl+"/index/getPlantListTitle", url.Values{}, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (h *Client) GetPlantDevices(plantId string) (*GrowattPlantDevices, error) {
	var result GrowattPlantDevices
	if err := h.postForm(h.serverUrl+"/panel/getDevicesByPlantList", url.Values{
		"plantId":  {plantId},
		"currPage": {"1"},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (h *Client) GetNoahList(plantId int) (*GrowattNoahList, error) {
	var result GrowattNoahList
	if err := h.postForm(h.serverUrl+"/device/getNoahList", url.Values{
		"plantId":  {fmt.Sprintf("%d", plantId)},
		"currPage": {"1"},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (h *Client) GetNoahDetails(plantId int, serial string) (*GrowattNoahList, error) {
	var result GrowattNoahList
	if err := h.postForm(h.serverUrl+"/device/getNoahList", url.Values{
		"plantId":  {fmt.Sprintf("%d", plantId)},
		"deviceSn": {serial},
		"currPage": {"1"},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (h *Client) GetNoahHistory(serial string, startDate string, endDate string) (*GrowattNoahHistory, error) {
	if startDate == "" {
		startDate = time.Now().Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}
	var result GrowattNoahHistory
	if err := h.postForm(h.serverUrl+"/device/getNoahHistory", url.Values{
		"deviceSn":  {serial},
		"start":     {"0"},
		"startDate": {startDate},
		"endDate":   {endDate},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (h *Client) GetNoahStatus(plantId int, serial string) (*GrowattNoahStatus, error) {
	var result GrowattNoahStatus
	if err := h.postForm(fmt.Sprintf(h.serverUrl+"/panel/noah/getNoahStatusData?plantId=%d", plantId), url.Values{
		"deviceSn": {serial},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (h *Client) GetNoahTotals(plantId int, serial string) (*GrowattNoahTotals, error) {
	var result GrowattNoahTotals
	if err := h.postForm(fmt.Sprintf(h.serverUrl+"/panel/noah/getNoahTotalData?plantId=%d", plantId), url.Values{
		"deviceSn": {serial},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
