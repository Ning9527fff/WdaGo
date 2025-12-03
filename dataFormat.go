package WdaGo

type WdaSession struct {
	url       string
	sessionId string
	headers   map[string]string
	client    *HTTPClient
}

type PhoneStatus struct {
	Device       string
	DeviceIP     string
	AgentVersion string
	OsName       string
	OsVersion    string
	SdkVersion   string
	State        string
	IsReady      bool
}

const (
	VolumeUp               = 1
	VolumeDown             = 2
	Home                   = 3
	NotificationTypePlain  = "plain"
	NotificationTypeDarwin = "darwin"
	StringNull             = ""
)

type Capabilities struct {
	BundleId string `json:"bundleId"`
}
type SessionRequest struct {
	Capabilities Capabilities `json:"capabilities"`
}

type PauseTime struct {
	Duration int `json:"duration"`
}

type ElementSearchRequest struct {
	Using string `json:"using"`
	Value string `json:"value"`
}

type TypingRequest struct {
	Value []byte `json:"value"`
}

type DeviceInfo struct {
	TimeZone           string `json:"timeZone"`
	CurrentLocale      string `json:"currentLocale"`
	Model              string `json:"model"`
	Uuid               string `json:"uuid"`
	ThermalState       string `json:"thermalState"`
	UserInterfaceIdiom int64  `json:"userInterfaceIdiom"`
	UserInterfaceStyle string `json:"userInterfaceStyle"`
	Name               string `json:"name"`
	IsSimulator        bool   `json:"isSimulator"`
}

type Location struct {
	Latitude            int64 `json:"latitude"`
	AuthorizationStatus int64 `json:"authorizationStatus"`
	Longitude           int64 `json:"longitude"`
	Altitude            int64 `json:"altitude"`
}

type BatteryInfo struct {
	Level int64 `json:"level"`
	State int64 `json:"state"`
}

type WindowSize struct {
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
}

type ScreenSize struct {
	StatusBarSize WindowSize `json:"statusBarSize"`
	Scale         int64      `json:"scale"`
	ScreenSize    WindowSize `json:"screenSize"`
}

type ScreenSizeResponse struct {
	Value     ScreenSize `json:"value"`
	SessionId string     `json:"sessionId"`
}

type AppInfo struct {
	Value struct {
		ProcessArguments struct {
			Env  interface{}   `json:"env"`
			Args []interface{} `json:"args"`
		} `json:"processArguments"`
		Name     string `json:"name"`
		Pid      int    `json:"pid"`
		BundleId string `json:"bundleId"`
	} `json:"value"`
	SessionId string `json:"sessionId"`
}

type AppBaseInfo struct {
	Pid      int64  `json:"pid"`
	BundleId string `json:"bundleId"`
}

type AppList struct {
	Value     []AppBaseInfo `json:"value"`
	SessionId string        `json:"sessionId"`
}

type BundleIdRequest struct {
	BundleId string `json:"bundleId"`
}

type SourceRequest struct {
	Resource string `json:"resource"`
}

type ElementLocation struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type HoldRequest struct {
	elementLocation ElementLocation
	Duration        float64 `json:"duration"`
}

type DragOption struct {
	FromX float64 `json:"fromX"`
	FromY float64 `json:"fromY"`
	ToX   float64 `json:"toX"`
	ToY   float64 `json:"toY"`
}

type ButtonName struct {
	Name string `json:"name"`
}

type NotificationExpect struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Timeout int64  `json:"timeout"`
}

type TextRequest struct {
	Text string `json:"text"`
}

type UrlBody struct {
	Url string `json:"url"`
}

type Scale struct {
	ScaleX float64
	ScaleY float64
}
