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
	AcCoupleEnable              string `json:"acCoupleEnable"`
	AcCouplePowerControl        string `json:"acCouplePowerControl"`
	AccountName                 string `json:"accountName"`
	Address                     string `json:"address"`
	Alias                       string `json:"alias"`
	AllowGridCharging           string `json:"allowGridCharging"`
	AntiBackflowEnable          string `json:"antiBackflowEnable"`
	AntiBackflowPowerPercentage string `json:"antiBackflowPowerPercentage"`
	AssociatedInvManAndModel    string `json:"associatedInvManAndModel"`
	ChargingSocHighLimit        string `json:"chargingSocHighLimit"`
	ChargingSocLowLimit         string `json:"chargingSocLowLimit"`
	ComVersion                  string `json:"comVersion"`
	CountryAndArea              string `json:"countryAndArea"`
	DatalogSn                   string `json:"datalogSn"`
	DatalogType                 string `json:"datalogType"`
	DefaultACCouplePower        string `json:"defaultACCouplePower"`
	DefaultMode                 string `json:"defaultMode"`
	DeviceModel                 string `json:"deviceModel"`
	DeviceType                  string `json:"deviceType"`
	Dtc                         string `json:"dtc"`
	EMonth                      string `json:"eMonth"`
	EToday                      string `json:"eToday"`
	ETotal                      string `json:"eTotal"`
	FwVersion                   string `json:"fwVersion"`
	GridConnectionControl       string `json:"gridConnectionControl"`
	HidePowerFromGrid           string `json:"hidePowerFromGrid"`
	IsHaveCT                    string `json:"isHaveCT"`
	LastUpdateTime              string `json:"lastUpdateTime"`
	LightLoadEnable             string `json:"lightLoadEnable"`
	Lost                        string `json:"lost"`
	ManAddress                  string `json:"manAddress"`
	ManName                     string `json:"manName"`
	Model                       string `json:"model"`
	NeverPowerOff               string `json:"neverPowerOff"`
	NominalPower                string `json:"nominalPower"`
	Pac                         string `json:"pac"`
	PlantID                     string `json:"plantId"`
	PlantName                   string `json:"plantName"`
	Ppv                         string `json:"ppv"`
	SafetyCorrespondNum         string `json:"safetyCorrespondNum"`
	ShellyDeviceSn              string `json:"shellyDeviceSn"`
	ShellyFlag                  string `json:"shellyFlag"`
	Sn                          string `json:"sn"`
	Soc                         string `json:"soc"`
	Status                      string `json:"status"`
	SysTime                     string `json:"sysTime"`
	Time1Enable                 string `json:"time1Enable"`
	Time1End                    string `json:"time1End"`
	Time1Mode                   string `json:"time1Mode"`
	Time1Power                  string `json:"time1Power"`
	Time1Repeat                 string `json:"time1Repeat"`
	Time1Start                  string `json:"time1Start"`
	Time2Enable                 string `json:"time2Enable"`
	Time2End                    string `json:"time2End"`
	Time2Mode                   string `json:"time2Mode"`
	Time2Power                  string `json:"time2Power"`
	Time2Repeat                 string `json:"time2Repeat"`
	Time2Start                  string `json:"time2Start"`
	Time3Enable                 string `json:"time3Enable"`
	Time3End                    string `json:"time3End"`
	Time3Mode                   string `json:"time3Mode"`
	Time3Power                  string `json:"time3Power"`
	Time3Repeat                 string `json:"time3Repeat"`
	Time3Start                  string `json:"time3Start"`
	Time4Enable                 string `json:"time4Enable"`
	Time4End                    string `json:"time4End"`
	Time4Mode                   string `json:"time4Mode"`
	Time4Power                  string `json:"time4Power"`
	Time4Repeat                 string `json:"time4Repeat"`
	Time4Start                  string `json:"time4Start"`
	Time5Enable                 string `json:"time5Enable"`
	Time5End                    string `json:"time5End"`
	Time5Mode                   string `json:"time5Mode"`
	Time5Power                  string `json:"time5Power"`
	Time5Repeat                 string `json:"time5Repeat"`
	Time5Start                  string `json:"time5Start"`
	Time6Enable                 string `json:"time6Enable"`
	Time6End                    string `json:"time6End"`
	Time6Mode                   string `json:"time6Mode"`
	Time6Power                  string `json:"time6Power"`
	Time6Repeat                 string `json:"time6Repeat"`
	Time6Start                  string `json:"time6Start"`
	Time7Enable                 string `json:"time7Enable"`
	Time7End                    string `json:"time7End"`
	Time7Mode                   string `json:"time7Mode"`
	Time7Power                  string `json:"time7Power"`
	Time7Repeat                 string `json:"time7Repeat"`
	Time7Start                  string `json:"time7Start"`
	Time8Enable                 string `json:"time8Enable"`
	Time8End                    string `json:"time8End"`
	Time8Mode                   string `json:"time8Mode"`
	Time8Power                  string `json:"time8Power"`
	Time8Repeat                 string `json:"time8Repeat"`
	Time8Start                  string `json:"time8Start"`
	Time9Enable                 string `json:"time9Enable"`
	Time9End                    string `json:"time9End"`
	Time9Mode                   string `json:"time9Mode"`
	Time9Power                  string `json:"time9Power"`
	Time9Repeat                 string `json:"time9Repeat"`
	Time9Start                  string `json:"time9Start"`
	Timezone                    string `json:"timezone"`
	Version                     string `json:"version"`
	WorkMode                    string `json:"workMode"`
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
