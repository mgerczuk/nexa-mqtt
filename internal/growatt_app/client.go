package growatt_app

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"nexa-mqtt/internal/misc"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	client    HttpClient
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
		password:  hashPassword(password),
		jar:       jar,
	}
}

func (h *Client) postForm(url string, data url.Values, responseBody any) error {
	err := h.client.postForm(url, h.token, data, responseBody)
	if err != nil {
		notLoggedIn := strings.Contains(err.Error(), "Dear user, you have not login to the system") ||
			strings.Contains(err.Error(), "invalid character '<' looking for beginning of value")

		if notLoggedIn {
			slog.Warn("re-login", slog.String("error", err.Error()))
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

func (h *Client) loginGetToken() error {
	var data TokenResponse
	if err := h.postForm("https://evcharge.growatt.com/ocpp/user", url.Values{
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
	if err := h.postForm(h.serverUrl+"/newTwoLoginAPIV2.do", url.Values{
		"userName":          {h.username},
		"password":          {h.password},
		"newLogin":          {"1"},
		"phoneType":         {"android"},
		"shinephoneVersion": {"8.3.6.0"},
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
	if err := h.postForm(h.serverUrl+"/newTwoPlantAPI.do?op=getAllPlantListTwo", url.Values{
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
	if err := h.postForm(h.serverUrl+"/noahDeviceApi/noah/isPlantNoahSystem", url.Values{
		"plantId": {plantId},
	}, &data); err != nil {
		return nil, err
	}

	if !data.Obj.IsPlantHaveNexa {
		return nil, errors.New("No NEXA device")
	}

	return &data, nil
}

func (h *Client) GetSystemStatus(serialNumber string) (*NoahStatus, error) {
	var data NoahStatus
	if err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/getSystemStatus", url.Values{
		"deviceSn": {serialNumber},
	}, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (h *Client) GetNexaInfoBySn(serialNumber string) (*NexaInfo, error) {
	var data NexaInfo
	if err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/getNexaInfoBySn", url.Values{
		"language": {"1"},
		"deviceSn": {serialNumber},
	}, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// No longer used in app 8.3.1.0 and higher
func (h *Client) GetBatteryData(serialNumber string) (*BatteryInfo, error) {
	var data BatteryInfo
	if err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/getBatteryData", url.Values{
		"deviceSn": {serialNumber},
	}, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (h *Client) nexaSet(serialNumber string, settingType string, params ...string) error {
	body := url.Values{
		"serialNum": {serialNumber},
		"type":      {settingType},
	}
	for i, param := range params {
		body.Set(fmt.Sprintf("param%d", i+1), param)
	}
	var data SetResponse
	if err := h.postForm(h.serverUrl+"/noahDeviceApi/nexa/set", body, &data); err != nil {
		return err
	}
	if data.Result <= 0 {
		return errors.New(data.Msg)
	}

	return nil
}

func (h *Client) SetSystemOutputPower(serialNumber string, mode int, power float64) error {
	p := math.Max(0, math.Min(1000, power))
	return h.nexaSet(serialNumber, "system_out_put_power", fmt.Sprintf("%d", mode), fmt.Sprintf("%.0f", p))
}

func (h *Client) SetChargingSoc(serialNumber string, chargingLimit float64, dischargeLimit float64) error {
	c := math.Max(70, math.Min(100, chargingLimit))
	d := math.Max(0, math.Min(30, dischargeLimit))
	return h.nexaSet(serialNumber, "charging_soc", fmt.Sprintf("%.0f", c), fmt.Sprintf("%.0f", d))
}

func (h *Client) SetAllowGridCharging(serialNumber string, allow int) error {
	return h.nexaSet(serialNumber, "allow_grid_charging", fmt.Sprintf("%d", allow))
}

func (h *Client) SetGridConnectionControl(serialNumber string, offlineEnable int) error {
	return h.nexaSet(serialNumber, "grid_connection_control", fmt.Sprintf("%d", offlineEnable))
}

// "Power+ Function" setting
func (h *Client) SetACCouplePowerControl(serialNumber string, _1000WEnable int) error {
	return h.nexaSet(serialNumber, "ac_couple_power_control", fmt.Sprintf("%d", _1000WEnable))
}

// "AC Always On" setting
func (h *Client) SetLightLoadEnable(serialNumber string, enable int) error {
	return h.nexaSet(serialNumber, "light_load_enable", fmt.Sprintf("%d", enable))
}

// "Always On" setting
func (h *Client) SetNeverPowerOff(serialNumber string, enable int) error {
	return h.nexaSet(serialNumber, "never_power_off", fmt.Sprintf("%d", enable))
}

// "Export Limitation" setting
func (h *Client) SetBackflow(serialNumber string, enableLimit int, powerSettingPercent float64) error {
	return h.nexaSet(serialNumber, "backflow_setting", fmt.Sprintf("%d", enableLimit), fmt.Sprintf("%.0f", powerSettingPercent))
}
