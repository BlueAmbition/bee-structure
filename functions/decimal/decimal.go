package decimal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//保留float小数位
func NumFormat(floatStr interface{}, precision uint64) string {
	var (
		tempFloat float64
	)
	switch floatStr.(type) {
	case float64:
		tempFloat = floatStr.(float64)
		break
	case string:
		tempFloat, _ = strconv.ParseFloat(floatStr.(string), 64)
		break
	}
	if tempFloat == 0 {
		return "0"
	}
	originStr := fmt.Sprintf("%."+strconv.FormatUint(precision, 10)+"f", tempFloat)
	if !strings.Contains(originStr, ".") {
		return originStr
	}
	pattern := "0*$"
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return originStr
	}
	//将匹配到的部分替换为"##.#"
	str := reg.ReplaceAllString(originStr, "")
	if len(str) > 0 {
		if str[len(str)-1:] == "." {
			str = strings.Replace(str, ".", "", 1)
		}
	}
	return str
}
