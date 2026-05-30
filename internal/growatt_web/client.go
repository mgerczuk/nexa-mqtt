package growatt_web

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

func (c *Client) postForm(url string, data url.Values, responseBody any) error {
	err := c.client.postForm(url, data, responseBody)
	if err != nil {
		notLoggedIn := strings.Contains(err.Error(), "invalid character '<' looking for beginning of value")
		if notLoggedIn {
			slog.Warn("re-login (web)", slog.String("error", err.Error()))
			if err := c.Login(); err != nil {
				slog.Error("could not re-login", slog.String("error", err.Error()))
				misc.Panic(err)
			}
			return c.postForm(url, data, responseBody)
		} else {
			return err
		}
	}

	return nil
}

func (c *Client) Login() error {
	var result GrowattResult
	if err := c.postForm(c.serverUrl+"/login", url.Values{
		"account":  {c.username},
		"password": {c.password},
	}, &result); err != nil {
		return err
	}
	if result.Result < 0 {
		return errors.New(result.Msg)
	}
	return nil
}

func (c *Client) GetPlantList() ([]GrowattPlant, error) {
	var result []GrowattPlant
	if err := c.postForm(c.serverUrl+"/index/getPlantListTitle", url.Values{}, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetPlantDevices(plantId string) (*GrowattPlantDevices, error) {
	var result GrowattPlantDevices
	if err := c.postForm(c.serverUrl+"/panel/getDevicesByPlantList", url.Values{
		"plantId":  {plantId},
		"currPage": {"1"},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetNoahList(plantId int) (*GrowattNoahList, error) {
	var result GrowattNoahList
	if err := c.postForm(c.serverUrl+"/device/getNoahList", url.Values{
		"plantId":  {fmt.Sprintf("%d", plantId)},
		"currPage": {"1"},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetNoahDetails(plantId int, serial string) (*GrowattNoahList, error) {
	var result GrowattNoahList
	if err := c.postForm(c.serverUrl+"/device/getNoahList", url.Values{
		"plantId":  {fmt.Sprintf("%d", plantId)},
		"deviceSn": {serial},
		"currPage": {"1"},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetNoahHistory(serial string, startDate string, endDate string) (*GrowattNoahHistory, error) {
	if startDate == "" {
		startDate = time.Now().Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}
	var result GrowattNoahHistory
	if err := c.postForm(c.serverUrl+"/device/getNoahHistory", url.Values{
		"deviceSn":  {serial},
		"start":     {"0"},
		"startDate": {startDate},
		"endDate":   {endDate},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetNoahStatus(plantId int, serial string) (*GrowattNoahStatus, error) {
	var result GrowattNoahStatus
	if err := c.postForm(fmt.Sprintf(c.serverUrl+"/panel/noah/getNoahStatusData?plantId=%d", plantId), url.Values{
		"deviceSn": {serial},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetNoahTotals(plantId int, serial string) (*GrowattNoahTotals, error) {
	var result GrowattNoahTotals
	if err := c.postForm(fmt.Sprintf(c.serverUrl+"/panel/noah/getNoahTotalData?plantId=%d", plantId), url.Values{
		"deviceSn": {serial},
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type SetResponse struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

func (c *Client) tcpset(serialNumber string, _type string, params ...string) error {
	body := url.Values{
		"action":    {"noahSet"},
		"serialNum": {serialNumber},
		"type":      {_type},
	}
	for i, param := range params {
		body.Set(fmt.Sprintf("param%d", i+1), param)
	}
	var result SetResponse
	if err := c.postForm(c.serverUrl+"/tcpSet.do", body, &result); err != nil {
		return err
	}
	if !result.Success {
		return errors.New(result.Msg)
	}
	return nil
}

// "Charge Upper Limit SOC" setting
func (h *Client) SetChargingSocHighLimit(serialNumber string, chargingLimit float64) error {
	val := math.Max(70, math.Min(100, chargingLimit))
	return h.tcpset(serialNumber, "charging_soc_high_limit", fmt.Sprintf("%.0f", val))
}

// "Discharge Lower Limit SOC" setting
func (c *Client) SetChargingSocLowLimit(serialNumber string, dischargeLimit float64) error {
	val := math.Max(0, math.Min(30, dischargeLimit))
	return c.tcpset(serialNumber, "charging_soc_low_limit", fmt.Sprintf("%.0f", val))
}

// "Set Exportlimit" setting
func (h *Client) SetAntiBackflowSetting(serialNumber string, enableLimit int, powerSettingPercent float64) error {
	val := math.Max(0, math.Min(100, powerSettingPercent))
	return h.tcpset(serialNumber, "anti_back_flow_setting", fmt.Sprintf("%d", enableLimit), fmt.Sprintf("%.0f", val))
}

// "Power+ Function" setting
func (c *Client) SetACCouplePowerControl(serialNumber string, _1000WEnable int) error {
	return c.tcpset(serialNumber, "ac_couple_power_control", fmt.Sprintf("%d", _1000WEnable))
}

// "Off-Grid Enable" setting
func (c *Client) SetGridConnectionControl(serialNumber string, offlineEnable int) error {
	return c.tcpset(serialNumber, "grid_connection_control", fmt.Sprintf("%d", offlineEnable))
}

// "AC Always On" setting
func (c *Client) SetLightLoadEnable(serialNumber string, enable int) error {
	return c.tcpset(serialNumber, "light_load_enable", fmt.Sprintf("%d", enable))
}

// "System Default Output Power" setting
func (c *Client) SetSystemOutputPower(serialNumber string, mode int, power float64) error {
	val := math.Max(0, math.Min(1000, power))
	return c.tcpset(serialNumber, "system_out_put_power", fmt.Sprintf("%d", mode), fmt.Sprintf("%.0f", val))
}

// "AC Couple Enable" setting - not used
func (c *Client) SetAcCoupleEnable(serialNumber string, enabled int) error {
	return c.tcpset(serialNumber, "ac_couple_enable", fmt.Sprintf("%d", enabled))
}

// "Draw power from the grid" setting
func (c *Client) SetAllowGridCharging(serialNumber string, allow int) error {
	return c.tcpset(serialNumber, "allow_grid_charging", fmt.Sprintf("%d", allow))
}

// "Always On" setting
func (h *Client) SetNeverPowerOff(serialNumber string, enable int) error {
	return h.tcpset(serialNumber, "never_power_off", fmt.Sprintf("%d", enable))
}
