package growatt_web

type GrowattResult struct {
	Result int    `json:"result"`
	Msg    string `json:"msg"`
}

type GrowattPlant struct {
	PlantId   string `json:"id"`
	PlantName string `json:"plantName"`
}

type GrowattPlantDevices struct {
	Result int `json:"result"`
	Obj    struct {
		CurrPage int `json:"currPage"`
		Pages    int `json:"pages"`
		PageSize int `json:"pageSize"`
		Count    int `json:"count"`
		Ind      int `json:"ind"`
		Datas    []struct {
			DeviceType      string `json:"deviceType"`
			PtoStatus       string `json:"ptoStatus"`
			ShellyDeviceSn  string `json:"shellyDeviceSn"`
			TimeServer      string `json:"timeServer"`
			AccountName     string `json:"accountName"`
			Timezone        string `json:"timezone"`
			PlantID         string `json:"plantId"`
			DeviceTypeName  string `json:"deviceTypeName"`
			NominalPower    string `json:"nominalPower"`
			BdcStatus       string `json:"bdcStatus"`
			EToday          string `json:"eToday"`
			EMonth          string `json:"eMonth"`
			DatalogTypeTest string `json:"datalogTypeTest"`
			ETotal          string `json:"eTotal"`
			Pac             string `json:"pac"`
			DatalogSn       string `json:"datalogSn"`
			Alias           string `json:"alias"`
			Location        string `json:"location"`
			DeviceModel     string `json:"deviceModel"`
			Sn              string `json:"sn"`
			PlantName       string `json:"plantName"`
			Status          string `json:"status"`
			LastUpdateTime  string `json:"lastUpdateTime"`
		} `json:"datas"`
		NotPager bool `json:"notPager"`
	} `json:"obj"`
}

type GrowattNoahListData struct {
	AcCouplePowerControl  string `json:"acCouplePowerControl"`
	Alias                 string `json:"alias"`
	AllowGridCharging     string `json:"allowGridCharging"`
	ChargingSocHighLimit  string `json:"chargingSocHighLimit"`
	ChargingSocLowLimit   string `json:"chargingSocLowLimit"`
	DefaultACCouplePower  string `json:"defaultACCouplePower"`
	DefaultMode           string `json:"defaultMode"`
	DeviceModel           string `json:"deviceModel"`
	GridConnectionControl string `json:"gridConnectionControl"`
	LightLoadEnable       string `json:"lightLoadEnable"`
	NeverPowerOff         string `json:"neverPowerOff"`
	PlantID               string `json:"plantId"`
	Sn                    string `json:"sn"`
	Version               string `json:"version"`
}

type GrowattNoahList struct {
	CurrPage int                   `json:"currPage"`
	Pages    int                   `json:"pages"`
	PageSize int                   `json:"pageSize"`
	Count    int                   `json:"count"`
	Ind      int                   `json:"ind"`
	Datas    []GrowattNoahListData `json:"datas"`
	NotPager bool                  `json:"notPager"`
}

type GrowattNoahHistoryData struct {
	FunctionCode                   int     `json:"functionCode"`
	DeviceSn                       string  `json:"deviceSn"`
	DatalogSn                      string  `json:"datalogSn"`
	Time                           string  `json:"time"`
	IsAgain                        int     `json:"isAgain"`
	Dtc                            int     `json:"dtc"`
	Status                         int     `json:"status"`
	MpptProtectStatus              int     `json:"mpptProtectStatus"`
	PdWarnStatus                   int     `json:"pdWarnStatus"`
	Pac                            float64 `json:"pac"`
	EacToday                       float64 `json:"eacToday"`
	EacMonth                       float64 `json:"eacMonth"`
	EacYear                        float64 `json:"eacYear"`
	EacTotal                       float64 `json:"eacTotal"`
	Ppv                            float64 `json:"ppv"`
	WorkMode                       int     `json:"workMode"`
	TotalBatteryPackChargingStatus int     `json:"totalBatteryPackChargingStatus"`
	TotalBatteryPackChargingPower  int     `json:"totalBatteryPackChargingPower"`
	BatteryPackageQuantity         int     `json:"batteryPackageQuantity"`
	TotalBatteryPackSoc            int     `json:"totalBatteryPackSoc"`
	HeatingStatus                  int     `json:"heatingStatus"`
	FaultStatus                    int     `json:"faultStatus"`
	Battery1SerialNum              string  `json:"battery1SerialNum"`
	Battery1Soc                    int     `json:"battery1Soc"`
	Battery1Temp                   float64 `json:"battery1Temp"`
	Battery1WarnStatus             int     `json:"battery1WarnStatus"`
	Battery1ProtectStatus          int     `json:"battery1ProtectStatus"`
	Battery2SerialNum              string  `json:"battery2SerialNum"`
	Battery2Soc                    int     `json:"battery2Soc"`
	Battery2Temp                   float64 `json:"battery2Temp"`
	Battery2WarnStatus             int     `json:"battery2WarnStatus"`
	Battery2ProtectStatus          int     `json:"battery2ProtectStatus"`
	Battery3SerialNum              string  `json:"battery3SerialNum"`
	Battery3Soc                    int     `json:"battery3Soc"`
	Battery3Temp                   float64 `json:"battery3Temp"`
	Battery3WarnStatus             int     `json:"battery3WarnStatus"`
	Battery3ProtectStatus          int     `json:"battery3ProtectStatus"`
	Battery4SerialNum              string  `json:"battery4SerialNum"`
	Battery4Soc                    int     `json:"battery4Soc"`
	Battery4Temp                   float64 `json:"battery4Temp"`
	Battery4WarnStatus             int     `json:"battery4WarnStatus"`
	Battery4ProtectStatus          int     `json:"battery4ProtectStatus"`
	SettableTimePeriod             int     `json:"settableTimePeriod"`
	AcCoupleWarnStatus             int     `json:"acCoupleWarnStatus"`
	AcCoupleProtectStatus          int     `json:"acCoupleProtectStatus"`
	CtFlag                         int     `json:"ctFlag"`
	TotalHouseholdLoad             float64 `json:"totalHouseholdLoad"`
	HouseholdLoadApartFromGroplug  float64 `json:"householdLoadApartFromGroplug"`
	OnOffGrid                      int     `json:"onOffGrid"`
	CtSelfPower                    float64 `json:"ctSelfPower"`
	PlugPower                      float64 `json:"plugPower"`
	SocketPower                    int     `json:"socketPower"`
	ChargeSocLimit                 int     `json:"chargeSocLimit"`
	DischargeSocLimit              int     `json:"dischargeSocLimit"`
	Pv1Voltage                     float64 `json:"pv1Voltage"`
	Pv1Current                     float64 `json:"pv1Current"`
	Pv1Temp                        float64 `json:"pv1Temp"`
	Pv2Voltage                     float64 `json:"pv2Voltage"`
	Pv2Current                     float64 `json:"pv2Current"`
	Pv2Temp                        float64 `json:"pv2Temp"`
	SystemTemp                     float64 `json:"systemTemp"`
	MaxCellVoltage                 float64 `json:"maxCellVoltage"`
	MinCellVoltage                 float64 `json:"minCellVoltage"`
	BatteryCycles                  int     `json:"batteryCycles"`
	BatterySoh                     int     `json:"batterySoh"`
	Pv3Voltage                     float64 `json:"pv3Voltage"`
	Pv3Current                     float64 `json:"pv3Current"`
	Pv3Temp                        float64 `json:"pv3Temp"`
	Pv4Voltage                     float64 `json:"pv4Voltage"`
	Pv4Current                     float64 `json:"pv4Current"`
	Pv4Temp                        float64 `json:"pv4Temp"`
	TotalReactivePower             float64 `json:"totalReactivePower"`
	TotalPowerFactor               float64 `json:"totalPowerFactor"`
	ForwardReactiveEnergy          float64 `json:"forwardReactiveEnergy"`
	ReverseReactiveEnergy          float64 `json:"reverseReactiveEnergy"`
	ApparentEnergy                 float64 `json:"apparentEnergy"`
	TotalActiveEnergy              float64 `json:"totalActiveEnergy"`
	TotalReactiveEnergy            float64 `json:"totalReactiveEnergy"`
	Frequency                      float64 `json:"frequency"`
	PhaseAReactivePower            float64 `json:"phaseAReactivePower"`
	PhaseBReactivePower            float64 `json:"phaseBReactivePower"`
	PhaseCReactivePower            float64 `json:"phaseCReactivePower"`
	VoltageAB                      float64 `json:"voltageAB"`
	VoltageBC                      float64 `json:"voltageBC"`
	VoltageCA                      float64 `json:"voltageCA"`
	PhaseAActiveEnergy             float64 `json:"phaseAActiveEnergy"`
	PhaseAReactiveEnergy           float64 `json:"phaseAReactiveEnergy"`
	PhaseBActiveEnergy             float64 `json:"phaseBActiveEnergy"`
	PhaseBReactiveEnergy           float64 `json:"phaseBReactiveEnergy"`
	PhaseCActiveEnergy             float64 `json:"phaseCActiveEnergy"`
	PhaseCReactiveEnergy           float64 `json:"phaseCReactiveEnergy"`
	EastronFlag                    int     `json:"eastronFlag"`
	EastronStatus                  int     `json:"eastronStatus"`
	ForwardTotalActiveEnergy       float64 `json:"forwardTotalActiveEnergy"`
	ReverseTotalActiveEnergy       float64 `json:"reverseTotalActiveEnergy"`
	TotalActivePower               float64 `json:"totalActivePower"`
	TotalApparentPower             float64 `json:"totalApparentPower"`
	PhaseAActivePower              float64 `json:"phaseAActivePower"`
	PhaseAApparentPower            float64 `json:"phaseAApparentPower"`
	PhaseAVoltage                  float64 `json:"phaseAVoltage"`
	PhaseACurrent                  float64 `json:"phaseACurrent"`
	PhaseAPowerFactor              float64 `json:"phaseAPowerFactor"`
	PhaseBActivePower              float64 `json:"phaseBActivePower"`
	PhaseBApparentPower            float64 `json:"phaseBApparentPower"`
	PhaseBVoltage                  float64 `json:"phaseBVoltage"`
	PhaseBCurrent                  float64 `json:"phaseBCurrent"`
	PhaseBPowerFactor              float64 `json:"phaseBPowerFactor"`
	PhaseCActivePower              float64 `json:"phaseCActivePower"`
	PhaseCApparentPower            float64 `json:"phaseCApparentPower"`
	PhaseCVoltage                  float64 `json:"phaseCVoltage"`
	PhaseCCurrent                  float64 `json:"phaseCCurrent"`
	PhaseCPowerFactor              float64 `json:"phaseCPowerFactor"`
}

type GrowattNoahHistoryObj struct {
	EndDate  string                   `json:"endDate"`
	Datas    []GrowattNoahHistoryData `json:"datas"`
	Start    int                      `json:"start"`
	HaveNext bool                     `json:"haveNext"`
}

type GrowattNoahHistory struct {
	Result int                   `json:"result"`
	Obj    GrowattNoahHistoryObj `json:"obj"`
}

type GrowattNoahStatusObj struct {
	SmartSocketPower              string `json:"smartSocketPower"`
	CtSelfPower                   string `json:"ctSelfPower"`
	GroplugFlag                   string `json:"groplugFlag"`
	HouseholdLoadApartFromGroplug string `json:"householdLoadApartFromGroplug"`
	ShellyFlag                    string `json:"shellyFlag"`
	TotalHouseholdLoad            string `json:"totalHouseholdLoad"`
	TotalBatteryPackSoc           string `json:"totalBatteryPackSoc"`
	Pac                           string `json:"pac"`
	WorkMode                      string `json:"workMode"`
	EastronFlag                   string `json:"eastronFlag"`
	BatteryPackageQuantity        string `json:"batteryPackageQuantity"`
	Ppv                           string `json:"ppv"`
	GroplugNum                    string `json:"groplugNum"`
	TotalBatteryPackChargingPower string `json:"totalBatteryPackChargingPower"`
	OtherPower                    string `json:"otherPower"`
	Status                        string `json:"status"`
}

type GrowattNoahStatus struct {
	Result  int                  `json:"result"`
	Msg     interface{}          `json:"msg"`
	Obj     GrowattNoahStatusObj `json:"obj"`
	Request interface{}          `json:"request"`
}

type GrowattNoahTotalsObj struct {
	MTotal    string `json:"mTotal"`
	MUnitText string `json:"mUnitText"`
	MToday    string `json:"mToday"`
	EacTotal  string `json:"eacTotal"`
	EacToday  string `json:"eacToday"`
}

type GrowattNoahTotals struct {
	Result  int                  `json:"result"`
	Msg     interface{}          `json:"msg"`
	Obj     GrowattNoahTotalsObj `json:"obj"`
	Request interface{}          `json:"request"`
}
