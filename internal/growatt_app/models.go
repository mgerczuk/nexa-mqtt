package growatt_app

type TokenResponse struct {
	Code  int    `json:"code"`
	Data  string `json:"data"`
	Token string `json:"token"`
}

type LoginResult struct {
	Back struct {
		Msg     string `json:"msg"`
		Success bool   `json:"success"`
		User    struct {
			ID int `json:"id"`
		} `json:"user"`
	} `json:"back"`
}

type PlantListV2 struct {
	PlantList []struct {
		ID int `json:"id"`
	} `json:"PlantList"`
}

type ResponseContainerV2[T any] struct {
	Msg    string `json:"msg"`
	Result int    `json:"result"`
	Obj    T      `json:"obj"`
}

type NoahPlantInfoObj struct {
	IsPlantNoahSystem bool   `json:"isPlantNoahSystem"`
	PlantID           string `json:"plantId"`
	IsPlantHaveNoah   bool   `json:"isPlantHaveNoah"`
	IsPlantHaveNexa   bool   `json:"isPlantHaveNexa"`
	DeviceSn          string `json:"deviceSn"`
	PlantName         string `json:"plantName"`
}

type NoahPlantInfo struct {
	ResponseContainerV2[NoahPlantInfoObj]
}

type NoahStatusObj struct {
	LoadPower       string `json:"loadPower"` // new
	GridPower       string `json:"gridPower"` // new
	ChargePower     string `json:"chargePower"`
	GroplugPower    string `json:"groplugPower"` // new
	WorkMode        string `json:"workMode"`
	Soc             string `json:"soc"`
	EastronStatus   string `json:"eastronStatus"` // new
	AssociatedInvSn string `json:"associatedInvSn"`
	BatteryNum      string `json:"batteryNum"`
	ProfitToday     string `json:"profitToday"`
	PlantID         string `json:"plantId"`
	DisChargePower  string `json:"disChargePower"`
	EacTotal        string `json:"eacTotal"`
	EacToday        string `json:"eacToday"`
	IsHaveCt        string `json:"isHaveCt"`  // new
	OnOffGrid       string `json:"onOffGrid"` // new
	Pac             string `json:"pac"`
	Ppv             string `json:"ppv"`
	Alias           string `json:"alias"`
	ProfitTotal     string `json:"profitTotal"`
	MoneyUnit       string `json:"moneyUnit"`
	GroplugNum      string `json:"groplugNum"` // new
	OtherPower      string `json:"otherPower"` // new
	Status          string `json:"status"`     // 1 = online, -1 = offline, 5 = heating
}

type NoahStatus struct {
	ResponseContainerV2[NoahStatusObj]
}

type NexaInfoObj struct {
	Noah struct {
		AcCouple                    string              `json:"acCouple"`
		AcCoupleEnable              string              `json:"acCoupleEnable"`
		AcCouplePowerControl        string              `json:"acCouplePowerControl"`
		Alias                       string              `json:"alias"`
		AllowGridCharging           string              `json:"allowGridCharging"`
		AmmeterModel                string              `json:"ammeterModel"`
		AmmeterSn                   string              `json:"ammeterSn"`
		AntiBackflowEnable          string              `json:"antiBackflowEnable"`
		AntiBackflowPowerPercentage string              `json:"antiBackflowPowerPercentage"`
		BatSns                      []string            `json:"batSns"`
		ChargingSocHighLimit        string              `json:"chargingSocHighLimit"`
		ChargingSocLowLimit         string              `json:"chargingSocLowLimit"`
		CtType                      string              `json:"ctType"`
		DefaultACCouplePower        string              `json:"defaultACCouplePower"`
		DefaultMode                 string              `json:"defaultMode"`
		DeviceSn                    string              `json:"deviceSn"`
		FormulaMoney                string              `json:"formulaMoney"`
		GridConnectionControl       string              `json:"gridConnectionControl"`
		GridSet                     string              `json:"gridSet"`
		LightLoadEnable             string              `json:"light_load_enable"`
		Model                       string              `json:"model"`
		MoneyUnitText               string              `json:"moneyUnitText"`
		NeverPowerOff               string              `json:"never_power_off"`
		NeverPowerOffSet            string              `json:"neverPowerOffSet"`
		PlantID                     string              `json:"plantId"`
		PlantName                   string              `json:"plantName"`
		Safety                      int                 `json:"safety"`
		SafetyEnable                string              `json:"safetyEnable"`
		ShellyList                  []interface{}       `json:"shellyList"`
		SmartPlan                   string              `json:"smartPlan"`
		TempType                    string              `json:"tempType"`
		TimeSegment                 []map[string]string `json:"time_segment"`
		Version                     string              `json:"version"`
		WorkMode                    string              `json:"workMode"`
	} `json:"noah"`
	PlantList []struct {
		PlantID      string      `json:"plantId"`
		PlantImgName interface{} `json:"plantImgName"`
		PlantName    string      `json:"plantName"`
	} `json:"plantList"`
	UnitList map[string]string `json:"unitList"` // new
}

type NexaInfo struct {
	ResponseContainerV2[NexaInfoObj]
}

type BatteryInfoObj struct {
	Batter   []BatteryDetails `json:"batter"`
	TempType string           `json:"tempType"`
	Time     string           `json:"time"`
}

type BatteryInfo struct {
	ResponseContainerV2[BatteryInfoObj]
}

type BatteryDetails struct {
	SerialNum string `json:"serialNum"`
	Soc       string `json:"soc"`
	Temp      string `json:"temp"`
}

type SetResponse struct {
	ResponseContainerV2[any]
}
