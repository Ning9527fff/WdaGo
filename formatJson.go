package WdaGo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// GetDataFromRespBody 用于处理wda标准返回数据中的value结构的数据，将其转为map
func GetDataFromRespBody(body []byte) (map[string]interface{}, error) {
	result := gjson.Get(string(body), "value")
	if !result.Exists() {
		return nil, fmt.Errorf(" Get value failed from result, because result is not valid ")
	}

	return result.Value().(map[string]interface{}), nil
}

func GetStringFromValueInterface(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok && val != nil {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func GetNumFromValueInterface(data map[string]interface{}, key string) int64 {
	if val, ok := data[key]; ok && val != nil {
		switch v := val.(type) {
		case int64:
			return v
		case int:
			return int64(v)
		case float64:
			return int64(v)
		case string:
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				return i
			}
		}
	}
	return 0
}

func GetBoolFromValueInterface(data map[string]interface{}, key string) bool {
	if val, ok := data[key]; ok && val != nil {
		switch v := val.(type) {
		case bool:
			return v
		case string:
			return strings.ToLower(v) == "true"
		case int, int64:
			return v != 0
		}
	}
	return false
}

// JudgeResponseCorrect 判断wda请求返回结果是否正确, 正确为true，错误为false
func JudgeResponseCorrect(body []byte, sessionId string) bool {
	if gjson.Get(string(body), "value").String() == "" &&
		gjson.Get(string(body), "sessionId").String() == sessionId {
		return true
	} else {
		return false
	}
}
