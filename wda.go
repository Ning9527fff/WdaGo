package WdaGo

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Ning9527fff/MyLog"
	"github.com/tidwall/gjson"
)

const (
	UserAgent       = "Go-HTTP-Client/1.0"
	ContentTypeJson = "application/json"
	PicturePath     = "screenShot/"
	LinkText        = 1
	PartialLinkText = 2
	ClassName       = 3
	Path            = 4
	ClassChain      = 5
)

func GetWdaSession(url string) *WdaSession {

	header := map[string]string{
		"User-Agent":   UserAgent,
		"content-type": ContentTypeJson,
	}

	session := &WdaSession{
		url:     url,
		headers: header,
		client:  NewHTTPClient(0),
	}
	return session
}

// GetStatus 获取当前iphone上的wda状态
func (session *WdaSession) GetStatus() (*PhoneStatus, error) {

	api := session.url + "/status"

	body, err := session.client.GetRequest(api, session.headers)

	if err != nil {
		return nil, err
	}

	deviceStatus := &PhoneStatus{
		Device:       gjson.Get(string(body), "value.device").String(),
		DeviceIP:     gjson.Get(string(body), "value.ios.ip").String(),
		AgentVersion: gjson.Get(string(body), "value.build.version").String(),
		OsName:       gjson.Get(string(body), "value.os.name").String(),
		OsVersion:    gjson.Get(string(body), "value.os.version").String(),
		SdkVersion:   gjson.Get(string(body), "value.os.sdkVersion").String(),
		State:        gjson.Get(string(body), "value.os.state").String(),
		IsReady:      gjson.Get(string(body), "value.ready").Bool(),
	}

	return deviceStatus, nil
}

func (session *WdaSession) GetSession(bundleId string) error {

	api := session.url + "/session"

	data := SessionRequest{
		Capabilities: Capabilities{
			BundleId: bundleId,
		},
	}

	body, err := session.client.PostRequest(api, data, session.headers)
	log.DebugF("Response body: %v", string(body))

	if err != nil {
		return err
	}

	session.sessionId = gjson.Get(string(body), "value.sessionId").String()
	return nil
}

// CloseSession 关闭session
func (session *WdaSession) CloseSession() error {
	if session.sessionId == "" {
		return fmt.Errorf(" No session can be closed.")
	}

	err := session.CloseSession()
	if err != nil {
		return fmt.Errorf(" Close session failed:  %v", err)
	} else {
		session.sessionId = ""
		return nil
	}
}

func (session *WdaSession) CheckSession() (bool, error) {

	api := session.url + "/session/" + session.sessionId

	body, err := session.client.GetRequest(api, session.headers)
	if err != nil {
		return false, err
	}

	if gjson.Get(string(body), "sessionId").String() == session.sessionId {
		return true, nil
	} else {
		return false, nil
	}

}

func (session *WdaSession) DeleteSession() error {

	api := session.url + "/session/" + session.sessionId

	body, err := session.client.DeleteRequest(api, session.headers)
	if err != nil {
		return err
	}
	log.DebugF("response body is %v ", string(body))

	if gjson.Get(string(body), "sessionId").String() == "" {
		return nil
	} else {
		return fmt.Errorf(" Delete session failed ")
	}

}

// GetDeviceInfo 获取设备当前的状态
func (session *WdaSession) GetDeviceInfo() (*DeviceInfo, error) {

	api := session.url + "/session" + session.sessionId + "/wda/device/info"
	body, err := session.client.GetRequest(api, session.headers)
	if err != nil {
		return nil, err
	}

	value := gjson.Get(string(body), "value")
	if !value.Exists() {
		return nil, fmt.Errorf(" Get device info failed, No device info ")
	}

	result := value.Value().(map[string]interface{})

	Device := DeviceInfo{
		TimeZone:           GetStringFromValueInterface(result, "timeZone"),
		CurrentLocale:      GetStringFromValueInterface(result, "currentLocale"),
		Model:              GetStringFromValueInterface(result, "model"),
		Uuid:               GetStringFromValueInterface(result, "uuid"),
		ThermalState:       GetStringFromValueInterface(result, "thermalState"),
		UserInterfaceIdiom: GetNumFromValueInterface(result, "userInterfaceIdiom"),
		UserInterfaceStyle: GetStringFromValueInterface(result, "userInterfaceStyle"),
		Name:               GetStringFromValueInterface(result, "name"),
		IsSimulator:        GetBoolFromValueInterface(result, "isSimulator"),
	}
	return &Device, nil
}

// GetLocation 用于获取iphone的经纬度，授权状态等数据
func (session *WdaSession) GetLocation() (error, *Location) {
	api := session.url + "/session/" + session.sessionId + "/wda/location"

	body, err := session.client.GetRequest(api, session.headers)
	if err != nil {
		return fmt.Errorf(" Get location from api failed: %v ", err), nil
	}

	data, err := GetDataFromRespBody(body)
	if err != nil {
		return fmt.Errorf(" Get location failed %v ", err), nil
	}

	return nil, &Location{
		Latitude:            GetNumFromValueInterface(data, "latitude"),
		AuthorizationStatus: GetNumFromValueInterface(data, "authorizationStatus"),
		Longitude:           GetNumFromValueInterface(data, "longitude"),
		Altitude:            GetNumFromValueInterface(data, "altitude"),
	}
}

// GetBatteryInfo 获取电池信息
func (session *WdaSession) GetBatteryInfo() (*BatteryInfo, error) {

	api := session.url + "/session/" + session.sessionId + "/wda/batteryInfo"
	body, err := session.client.GetRequest(api, session.headers)
	if err != nil {
		return nil, fmt.Errorf(" Get battery info failed from api : %v ", err)
	}

	data, err := GetDataFromRespBody(body)
	if err != nil {
		return nil, fmt.Errorf(" Get battery info failed : %v ", err)
	}

	return &BatteryInfo{
		Level: GetNumFromValueInterface(data, "level"),
		State: GetNumFromValueInterface(data, "state"),
	}, nil

}

// BackToHomePage 返回home页
func (session *WdaSession) BackToHomePage() error {
	api := session.url + "/wda/homescreen"

	body, err := session.client.PostRequest(api, nil, session.headers)
	if err != nil {
		return err
	}

	if gjson.Get(string(body), "sessionId").String() == session.sessionId &&
		gjson.Get(string(body), "value").String() == "" {
		return nil
	} else {
		return fmt.Errorf(" Back to home page failed ")
	}
}

// CurrentScreenShot 当前页面截屏, 不置顶文件后缀，默认为.png
func (session *WdaSession) CurrentScreenShot(picturePath, pictureName string) (string, error) {
	api := session.url + "/screenshot"

	body, err := session.client.GetRequest(api, session.headers)
	if err != nil {
		return StringNull, err
	}

	pictureData := gjson.Get(string(body), "value")

	if !pictureData.Exists() {
		return StringNull, fmt.Errorf(" Current screenshot failed, there is no vaild data ")
	}

	//base64 decode
	imageDataByte, err := base64.StdEncoding.DecodeString(pictureData.String())
	if err != nil {
		return StringNull, fmt.Errorf(" Decode picture data with base64 failed ")
	}

	// 如果传入未指定文件后缀，则默认为png
	if !(strings.Contains(pictureName, ".png") ||
		strings.Contains(pictureName, ".jpg") ||
		strings.Contains(pictureName, ".jpeg")) {
		pictureName = pictureName + ".png"
	}

	// err = os.WriteFile(PicturePath+session.sessionId + ".png", imageDataByte, 0644)
	imagePath := filepath.Join(picturePath, pictureName)
	err = os.WriteFile(imagePath, imageDataByte, 0644)
	if err != nil {
		return StringNull, fmt.Errorf(" Write image file failed %v", err)
	} else {
		return imagePath, nil
	}
}

// GetAkaTree 获取当前页面树🌲
func (session *WdaSession) GetAkaTree() error {

	api := session.url + "/source"

	body, err := session.client.GetRequest(api, session.headers)
	if err != nil {
		return fmt.Errorf(" Get Aka Tree failed %v", err)
	}

	xmlFlow := gjson.Get(string(body), "value").String()

	log.DebugF("%v", string(body))
	log.DebugF("%v", xmlFlow)

	err = os.WriteFile("test.xml", []byte(xmlFlow), 0644)

	return nil
}

// SearchElement 以不同方式搜索元素
func (session *WdaSession) SearchElement(searchType int, Parms string) (string, error) {

	api := session.url + "/session/" + session.sessionId + "/elements"

	eleReq := ElementSearchRequest{}
	switch searchType {
	case LinkText:
		eleReq.Using = "link text"
	case PartialLinkText:
		eleReq.Using = "partial link text"
	case ClassName:
		eleReq.Using = "class name"
	case ClassChain:
		eleReq.Using = "xpath "
	case Path:
		eleReq.Using = "xpath"
	default:
		return "", fmt.Errorf(" Not supported search type now ")
	}

	eleReq.Value = Parms

	body, err := session.client.PostRequest(api, eleReq, session.headers)
	if err != nil {
		return "", fmt.Errorf(" Search element failed %v", err)
	}
	//这里要加一个返回多个element的逻辑
	return gjson.Get(string(body), "value.0.ELEMENT").String(), nil

}

func (session *WdaSession) ClickElement(elementId string) error {
	api := session.url + "/session/" + session.sessionId + "/element/" + elementId + "/click"

	body, err := session.client.PostRequest(api, nil, session.headers)
	if err != nil {
		return fmt.Errorf(" Click element failed %v", err)
	}

	if gjson.Get(string(body), "value").String() == "" &&
		gjson.Get(string(body), "sessionId").String() == session.sessionId {
		return nil
	} else {
		return fmt.Errorf(" Click element failed ")
	}

}

func (session *WdaSession) TypingText(elementId string, Text string) error {
	api := session.url + "/session/" + session.sessionId + "/element/" + elementId + "/value"

	typingReq := TypingRequest{
		Value: []byte(Text),
	}

	body, err := session.client.PostRequest(api, typingReq, session.headers)
	if err != nil {
		return fmt.Errorf(" Typing text failed %v", err)
	}

	if gjson.Get(string(body), "value").String() == "" &&
		gjson.Get(string(body), "sessionId").String() == session.sessionId {
		return nil
	} else {
		return fmt.Errorf(" Typing text failed ")
	}
}

func (session *WdaSession) ClearText(elementId string) error {
	api := session.url + "/session/" + session.sessionId + "/element/" + elementId + "/clear"

	body, err := session.client.PostRequest(api, nil, session.headers)
	if err != nil {
		return fmt.Errorf(" Clear text failed %v", err)
	}

	if gjson.Get(string(body), "value").String() == "" &&
		gjson.Get(string(body), "sessionId").String() == session.sessionId {
		return nil
	} else {
		return fmt.Errorf(" Click element failed ")
	}
}

func (session *WdaSession) AlertGet(client *HTTPClient) error {
	//curl -X GET $JSON_HEADER $DEVICE_URL/session/$SESSION_ID/alert/text
	return nil
}

func (session *WdaSession) AlertAccept(client *HTTPClient) error {
	//curl -X POST $JSON_HEADER -d "" $DEVICE_URL/session/$SESSION_ID/alert/accept
	return nil
}

func (session *WdaSession) AlertDismiss(client *HTTPClient) error {
	//curl -X POST $JSON_HEADER -d "" $DEVICE_URL/session/$SESSION_ID/alert/dismiss
	return nil
}

// GetWindowSize 获取当前窗口大小
func (session *WdaSession) GetWindowSize() (*WindowSize, error) {

	api := session.url + "/session/" + session.sessionId + "/windows/size"

	body, err := session.client.GetRequest(api, session.headers)
	if err != nil {
		return nil, fmt.Errorf(" Get WindowSize failed from api :%v", err)
	}

	data, err := GetDataFromRespBody(body)
	if err != nil {
		return nil, fmt.Errorf(" Get WindowSize failed :%v", err)
	}

	return &WindowSize{
		Width:  GetNumFromValueInterface(data, "width"),
		Height: GetNumFromValueInterface(data, "height"),
	}, err
}

// GetScreenSize 获取设备屏幕的点长和点宽，返回换算系数和ScreenSize
func (session *WdaSession) GetScreenSize() (*ScreenSizeResponse, error) {
	api := session.url + "/session/" + session.sessionId + "/wda/screen"
	body, err := session.client.GetRequest(api, session.headers)
	if err != nil {
		return nil, fmt.Errorf(" Get Screen Size failed from api :%v", err)
	}

	scrSize := &ScreenSizeResponse{}

	err = json.Unmarshal(body, &scrSize)
	if err != nil {
		return nil, fmt.Errorf(" Parse Screen Size failed :%v", err)
	}
	return scrSize, nil
}

func (session *WdaSession) GetActiveAppInfo() (*AppInfo, error) {
	api := session.url + "/session/" + session.sessionId + "/wda/activeAppInfo"

	body, err := session.client.GetRequest(api, session.headers)
	if err != nil {
		return nil, fmt.Errorf(" Get Active App info failed from api :%v", err)
	}

	var appInfo AppInfo
	if err = json.Unmarshal(body, &appInfo); err != nil {
		return nil, fmt.Errorf(" Get Active App info failed from response body :%v", err)
	}

	return &appInfo, nil
}

func (session *WdaSession) GetAppList() (*[]AppBaseInfo, error) {

	api := session.url + "/session/" + session.sessionId + "/wda/apps/list"

	body, err := session.client.GetRequest(api, session.headers)
	if err != nil {
		return nil, fmt.Errorf(" Get App list failed from api :%v", err)
	}

	var appList AppList
	if err = json.Unmarshal(body, &appList); err != nil {
		return nil, fmt.Errorf(" Get App list failed from response body :%v", err)
	}

	return &appList.Value, nil
}

func (session *WdaSession) GetAppState(bundleIdString string) (int64, error) {
	api := session.url + "/session/" + session.sessionId + "/wda/apps/state"

	bundleId := BundleIdRequest{BundleId: bundleIdString}

	body, err := session.client.PostRequest(api, bundleId, session.headers)
	if err != nil {
		return 0, fmt.Errorf(" Get App state failed from api :%v", err)
	}

	return gjson.Get(string(body), "value").Int(), nil
}

// IsLocked 是否锁屏
func (session *WdaSession) IsLocked() (bool, error) {
	api := session.url + "/session/" + session.sessionId + "/wda/locked"

	body, err := session.client.GetRequest(api, session.headers)
	if err != nil {
		return false, fmt.Errorf(" Get Locked status failed from api :%v", err)
	}

	return gjson.Get(string(body), "value").Bool(), nil
}

// UnlockedDevice 解锁设备
func (session *WdaSession) UnlockedDevice() error {
	api := session.url + "/session/" + session.sessionId + "/wda/unlock"

	body, err := session.client.PostRequest(api, nil, session.headers)
	if err != nil {
		return fmt.Errorf(" Unlocked device failed from api :%v", err)
	}

	if gjson.Get(string(body), "value").String() == "" &&
		gjson.Get(string(body), "sessionId").String() == session.sessionId {
		return nil
	} else {
		return fmt.Errorf(" Unlocked device failed ")
	}
}

func (session *WdaSession) LockedDevice() error {
	api := session.url + "/session/" + session.sessionId + "/wda/lock"

	body, err := session.client.PostRequest(api, nil, session.headers)
	if err != nil {
		return fmt.Errorf(" Lock device failed from api :%v", err)
	}

	if gjson.Get(string(body), "value").String() == "" &&
		gjson.Get(string(body), "sessionId").String() == session.sessionId {
		return nil
	} else {
		return fmt.Errorf(" Lock device failed ")
	}
}

func (session *WdaSession) LaunchApp(bundleId string) error {

	api := session.url + "/session/" + session.sessionId + "/wda/apps/launch"
	bundleIdReq := BundleIdRequest{
		BundleId: bundleId,
	}

	body, err := session.client.PostRequest(api, bundleIdReq, session.headers)
	if err != nil {
		return fmt.Errorf(" Launch App failed from api :%v", err)
	}

	if JudgeResponseCorrect(body, session.sessionId) {
		return nil
	} else {
		return fmt.Errorf(" Launch App failed ")
	}

}

// LaunchAppWithoutSession 不需要指定session来启动app
func (session *WdaSession) LaunchAppWithoutSession(bundleId string) error {
	api := session.url + "/wda/apps/launchUnattached"
	bundleIdReq := BundleIdRequest{
		BundleId: bundleId,
	}

	body, err := session.client.PostRequest(api, bundleIdReq, session.headers)
	if err != nil {
		return fmt.Errorf(" Launch App without session failed from api :%v", err)
	}

	if JudgeResponseCorrect(body, session.sessionId) {
		return nil
	} else {
		return fmt.Errorf(" Launch App without session failed ")
	}
}

// TerminateApp 关闭app
func (session *WdaSession) TerminateApp(bundleId string) error {
	api := session.url + "/session/" + session.sessionId + "/wda/apps/terminate"
	bundleIdReq := BundleIdRequest{
		BundleId: bundleId,
	}

	body, err := session.client.PostRequest(api, bundleIdReq, session.headers)
	if err != nil {
		return fmt.Errorf(" Terminate App failed from api :%v", err)
	}

	if gjson.Get(string(body), "value").Bool() {
		return nil
	} else {
		return fmt.Errorf(" Terminate App failed ")
	}
}

// ActivateApp 激活app？与启动有何区别暂时没搞清楚
func (session *WdaSession) ActivateApp(bundleId string) error {
	api := session.url + "/session/" + session.sessionId + "/wda/apps/activate"
	bundleIdReq := BundleIdRequest{
		BundleId: bundleId,
	}
	body, err := session.client.PostRequest(api, bundleIdReq, session.headers)
	if err != nil {
		return fmt.Errorf(" Activate App failed from api :%v", err)
	}
	if JudgeResponseCorrect(body, session.sessionId) {
		return nil
	} else {
		return fmt.Errorf(" Activate App failed ")
	}
}

// DeactivateApp 让app处于后台状态指定时间
func (session *WdaSession) DeactivateApp(time int) error {
	api := session.url + "/session/" + session.sessionId + "/wda/deactivateApp"

	dura := PauseTime{
		Duration: time,
	}

	body, err := session.client.PostRequest(api, dura, session.headers)
	if err != nil {
		return fmt.Errorf(" Deactivate app failed %v", err)
	}

	if JudgeResponseCorrect(body, session.sessionId) {
		return nil
	} else {
		return fmt.Errorf(" Deactivate app failed ")
	}
}

// ResetAppAuth 重置app auth，暂时不清楚如何使用，先实现
func (session *WdaSession) ResetAppAuth(resource string) error {
	api := session.url + "/session/" + session.sessionId + "/wda/resetAppAuth"
	sourceReq := SourceRequest{
		Resource: resource,
	}

	body, err := session.client.PostRequest(api, sourceReq, session.headers)
	if err != nil {
		return fmt.Errorf(" Reset App Auth failed from api :%v", err)
	}
	if JudgeResponseCorrect(body, session.sessionId) {
		return nil
	} else {
		return fmt.Errorf(" Reset App Auth failed ")
	}
}

// TapWithLocation  使用坐标点击
func (session *WdaSession) TapWithLocation(location ElementLocation) error {
	api := session.url + "/session/" + session.sessionId + "/wda/tap"

	body, err := session.client.PostRequest(api, ElementLocation{
		X: location.X,
		Y: location.Y,
	}, session.headers)
	if err != nil {
		return fmt.Errorf(" Tap With Location failed from api :%v", err)
	}

	if JudgeResponseCorrect(body, session.sessionId) {
		return nil
	} else {
		return fmt.Errorf(" Tap With Location failed ")
	}

}

// DoubleTapWithLocation 使用坐标双击
func (session *WdaSession) DoubleTapWithLocation(x, y float64) error {
	api := session.url + "/session/" + session.sessionId + "/wda/doubleTap"

	body, err := session.client.PostRequest(api, ElementLocation{
		X: x,
		Y: y,
	}, session.headers)
	if err != nil {
		return fmt.Errorf(" Tap With Location failed from api :%v", err)
	}

	if JudgeResponseCorrect(body, session.sessionId) {
		return nil
	} else {
		return fmt.Errorf(" Double Tap With Location failed ")
	}
}

// TouchAndHoldWithLocation 对指定坐标长按
func (session *WdaSession) TouchAndHoldWithLocation(x, y, duration float64) error {
	api := session.url + "/session/" + session.sessionId + "/wda/touchAndHold"

	body, err := session.client.PostRequest(api, HoldRequest{
		elementLocation: ElementLocation{
			X: x,
			Y: y,
		},
		Duration: duration,
	}, session.headers)
	if err != nil {
		return fmt.Errorf(" TouchAndHold With Location failed from api :%v", err)
	}

	if JudgeResponseCorrect(body, session.sessionId) {
		return nil
	} else {
		return fmt.Errorf(" TouchAndHold With Location failed ")
	}
}

// DragWithLocation 拖动操作 swipe操作与该操作本纸上为同一个
func (session *WdaSession) DragWithLocation(xBefore, yBefore, xLater, yLater float64) error {
	api := session.url + "/session/" + session.sessionId + "/wda/dragfromtoforduration"

	body, err := session.client.PostRequest(api, DragOption{
		FromX: xBefore,
		FromY: yBefore,
		ToX:   xLater,
		ToY:   yLater,
	}, session.headers)
	if err != nil {
		return fmt.Errorf(" Drag With Location failed from api :%v", err)
	}

	if JudgeResponseCorrect(body, session.sessionId) {
		return nil
	} else {
		return fmt.Errorf(" Drag With Location failed ")
	}
}

//W3C操作暂不实现

// PressButton 点击按钮，此处按钮指的是iphone的硬件按钮，硬件按钮名如下：
//
//	home,volumeUp,volumeDown
func (session *WdaSession) PressButton(buttonType int) error {

	var button ButtonName
	switch buttonType {
	case VolumeUp:
		button.Name = "volumeUp"
	case VolumeDown:
		button.Name = "volumeDown"
	case Home:
		button.Name = "home"
	default:
		return fmt.Errorf(" Error: UnKnown Button ")
	}

	api := session.url + "/session/" + session.sessionId + "/wda/pressButton"

	body, err := session.client.PostRequest(api, button, session.headers)
	if err != nil {
		return fmt.Errorf(" PressButton failed from api :%v", err)
	}

	if JudgeResponseCorrect(body, session.sessionId) {
		return nil
	} else {
		return fmt.Errorf(" PressButton failed ")
	}
}

// ExpectedNotification 判断是否出现一个预期中的notification
func (session *WdaSession) ExpectedNotification(notificationName string, notificationType string, timeOut int64) error {
	api := session.url + "/session/" + session.sessionId + "/wda/expectedNotification"

	body, err := session.client.PostRequest(api, NotificationExpect{
		Name:    notificationName,
		Type:    notificationType,
		Timeout: timeOut,
	}, session.headers)
	if err != nil {
		return fmt.Errorf(" Get Expected Notification failed from api :%v", err)
	}
	if JudgeResponseCorrect(body, session.sessionId) {
		return nil
	} else {
		return fmt.Errorf(" No Expected Notification found ")
	}
}

// ActiveSiri 启动siri,输入指定文本
func (session *WdaSession) ActiveSiri(text string) error {
	api := session.url + "/session/" + session.sessionId + "/wda/siri/activate"

	body, err := session.client.PostRequest(api, TextRequest{
		Text: text,
	}, session.headers)
	if err != nil {
		return fmt.Errorf(" Active Siri failed from api :%v", err)
	}

	if JudgeResponseCorrect(body, session.sessionId) {
		return nil
	} else {
		return fmt.Errorf(" Active Siri failed ")
	}

}

// LetSiriOpenUrl 让siri打开一个指定的url
// 传入的url必须是绝对url，即带https或者http
func (session *WdaSession) LetSiriOpenUrl(RawUrl string) error {
	api := session.url + "/session/" + session.sessionId + "/url"

	realUrl, err := url.Parse(RawUrl)
	if err != nil {
		return fmt.Errorf("Parse Url Failed  :%v", err)
	}

	if !realUrl.IsAbs() {
		return fmt.Errorf(" Url is not a absolutly url  ")
	}

	body, err := session.client.PostRequest(api, UrlBody{Url: realUrl.String()}, session.headers)
	if err != nil {
		return fmt.Errorf(" Siri Open Url failed from api :%v", err)
	}
	if JudgeResponseCorrect(body, session.sessionId) {
		return nil
	} else {
		return fmt.Errorf(" Siri Open Url failed : %v ", err)
	}
}

// GetOrientation 获取当前屏幕方向
func (session *WdaSession) GetOrientation() (string, error) {
	api := session.url + "/session/" + session.sessionId + "/orientation"

	body, err := session.client.GetRequest(api, session.headers)
	if err != nil {
		return "", fmt.Errorf(" Get Orientation failed from api :%v", err)
	}

	return gjson.Get(string(body), "value").String(), nil
}

// GetRotation 获取当前设备的旋转坐标，暂时没用，先不实现
func GetRotation() error {
	return nil
}

// ShutDownWda 关闭wda
func (session *WdaSession) ShutDownWda() error {
	api := session.url + "wda/shutDown"
	body, err := session.client.GetRequest(api, session.headers)
	if err != nil {
		return fmt.Errorf(" ShutDownWda failed from api :%v", err)
	}
	if JudgeResponseCorrect(body, session.sessionId) {
		return nil
	} else {
		return fmt.Errorf(" ShutDown Wda failed ")
	}
}
