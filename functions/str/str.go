package str

import (
	"math/rand"
	"strconv"
	"time"
)

//随机数字码
func RandNumCode(codeLen int) string {
	nums := ""
	rand.Seed(time.Now().Unix())
	for i := 0; i < codeLen; i++ {
		t := rand.Intn(9)
		nums += strconv.Itoa(t)
	}
	return nums
}

//随机数字字符串码
func RandMixCode(codeLen int) string {
	mixArr := [36]string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	mixsLen := len(mixArr) - 1
	codes := ""
	for i := 0; i < codeLen; i++ {
		t := rand.Intn(mixsLen)
		codes += mixArr[t]
	}
	return codes
}

//任何字符串（中英文）按一个个算长度
func StringLen(str string) int {
	runes := []rune(str)
	return len(runes)
}

//截断文本
func SubStr(str string, begin int, end int) string {
	if begin < 0 || begin > end {
		return str
	}
	runes := []rune(str)
	if len(runes) >= end {
		runes = runes[begin:end]
		return string(runes)
	}
	return str
}

//截断文本
func ShortTxt(str string, shortLen int) string {
	runes := []rune(str)
	if len(runes) > shortLen {
		runes = runes[:shortLen+1]
		return string(runes) + "......"
	}
	return string(runes)
}
