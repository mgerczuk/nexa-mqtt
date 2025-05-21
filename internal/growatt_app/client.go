package growatt_app

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"noah-mqtt/internal/misc"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	client    *http.Client
	serverUrl string
	username  string
	password  string
	userAgent string
	userId    string
	token     string
	jar       *cookiejar.Jar
}

func newClient(serverUrl string, username string, password string) *Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		slog.Error("could not create cookie jar", slog.String("error", err.Error()))
		misc.Panic(err)
	}

	slog.Info("setting server url (app)", slog.String("url", serverUrl))

	return &Client{
		client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           jar,
			Timeout:       10 * time.Second,
		},
		serverUrl: serverUrl,
		username:  username,
		password:  hashPassword(password),
		jar:       jar,
	}
}

func (h *Client) loginGetToken() error {
	var data TokenResponse
	if _, err := h.postForm("https://evcharge.growatt.com/ocpp/user", url.Values{
		"cmd":      {"shineLogin"},
		"userId":   {fmt.Sprintf("SHINE%s", h.username)},
		"password": {h.password},
		"lan":      {"1"},
	}, &data); err != nil {
		return err
	}

	h.token = data.Token
	return nil
}

func (h *Client) Login() error {
	if err := h.loginGetToken(); err != nil {
		return err
	}

	var data LoginResult
	if _, err := h.postForm(h.serverUrl+"/newTwoLoginAPIV2.do", url.Values{
		"userName":          {h.username},
		"password":          {h.password},
		"newLogin":          {"1"},
		"phoneType":         {"android"},
		"shinephoneVersion": {"8.3.0.2"},
		"phoneSn":           {uuid.New().String()},
		"ipvcpc":            {ipvcpc(h.username)},
		"language":          {"1"},
		"systemVersion":     {"15"},
		"phoneModel":        {"Mi A1"},
		"loginTime":         {time.Now().Format(time.DateTime)},
		"appType":           {"ShinePhone"},
		"timestamp":         {timestamp()},
	}, &data); err != nil {
		return err
	}

	if !data.Back.Success {
		return fmt.Errorf("login failed: %s", data.Back.Msg)
	}

	h.userId = fmt.Sprintf("%d", data.Back.User.ID)
	return nil
}

func (h *Client) GetPlantList() (*PlantListV2, error) {
	var data PlantListV2
	if _, err := h.postForm(h.serverUrl+"/newTwoPlantAPI.do?op=getAllPlantListTwo", url.Values{
		"plantStatus": {""},
		"pageSize":    {"20"},
		"language":    {"1"},
		"toPageNum":   {"1"},
		"order":       {"1"},
	}, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (h *Client) GetNoahPlantInfo(plantId string) (*NoahPlantInfo, error) {
	var data NoahPlantInfo
	if _, err := h.postForm(h.serverUrl+"/noahDeviceApi/noah/isPlantNoahSystem", url.Values{
		"plantId": {plantId},
	}, &data); err != nil {
		return nil, err
	}

	if !data.Obj.IsPlantHaveNexa {
		err := errors.New("No NEXA device")
		slog.Error(err.Error())
		misc.Panic(err)
	}

	return &data, nil
}

func (h *Client) GetNoahStatus(serialNumber string) (*NoahStatus, error) {
	var data NoahStatus
	if _, err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/getSystemStatus", url.Values{
		"deviceSn": {serialNumber},
	}, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (h *Client) GetNoahInfo(serialNumber string) (*NexaInfo, error) {
	var data NexaInfo
	if _, err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/getNexaInfoBySn", url.Values{
		"deviceSn": {serialNumber},
	}, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (h *Client) GetBatteryData(serialNumber string) (*BatteryInfo, error) {
	var data BatteryInfo
	if _, err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/getBatteryData", url.Values{
		"deviceSn": {serialNumber},
	}, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (h *Client) SetSystemOutputPower(serialNumber string, mode int, power float64) error {
	p := math.Max(0, math.Min(800, power))
	var data map[string]any
	if _, err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/set", url.Values{
		"serialNum": {serialNumber},
		"type":      {"system_out_put_power"},
		"param1":    {fmt.Sprintf("%d", mode)},
		"param2":    {fmt.Sprintf("%.0f", p)},
	}, &data); err != nil {
		return err
	}

	return nil
}

func (h *Client) SetSocLimit(serialNumber string, chargingLimit float64, dischargeLimit float64) error {
	c := math.Max(70, math.Min(100, chargingLimit))
	d := math.Max(0, math.Min(30, dischargeLimit))
	var data map[string]any
	if _, err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/set", url.Values{
		"serialNum": {serialNumber},
		"type":      {"charging_soc"},
		"param1":    {fmt.Sprintf("%.0f", c)},
		"param2":    {fmt.Sprintf("%.0f", d)},
	}, &data); err != nil {
		return err
	}

	return nil
}
