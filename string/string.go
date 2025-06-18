package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// func main() {
// 	// fmt.Println(os.Args)
// 	// keyword := os.Args[1]
// 	// test(keyword)
// 	//testPoint()
// 	//getPartDir("阳光电影www ygdy8.com.长津湖之水门桥.2022.BD.1080P.国语中字_20230811114016.mp4")
// 	// fmt.Println(TranslateStr("green"))
// 	// fmt.Println(TranslateSliceStr("[\"red\",\"blue\"]"))
// 	// fmt.Println(TranslateTableStr("[[\"red\",\"blue\",\"green\"],[\"big\",\"small\"]]"))
// 	//fmt.Println(GetCNPrizeLevel(10))
// }

func testPoint() {
	var a *string
	var b *int
	fmt.Println(&a, &b)
	//fmt.Println(*a,&b)
}

func getPartDir(filename string) string {
	fmt.Println("1111111111111111111", strings.SplitN(filename, ".", 2))
	fmt.Println(strings.SplitN(filename, ".", 2)[0])
	return strings.SplitN(filename, ".", 2)[0]
}

func test(keyword string) {
	keyword = strings.ReplaceAll(keyword, "\\", "\\\\")
	keyword = strings.ReplaceAll(keyword, "%", "\\%")
	keyword = strings.ReplaceAll(keyword, "'", "''")
	fmt.Println(keyword)
}

/*
*
判断是否为字母： unicode.IsLetter(v)
判断是否为十进制数字： unicode.IsDigit(v)
判断是否为数字： unicode.IsNumber(v)
判断是否为空白符号： unicode.IsSpace(v)
判断是否为Unicode标点字符 :unicode.IsPunct(v)
判断是否为中文：unicode.Han(v)
*/
func SpecialLetters(letter rune) (bool, []rune) {
	if unicode.IsPunct(letter) || unicode.IsSymbol(letter) {
		var chars []rune
		chars = append(chars, '\\', '\\', letter)
		return true, chars
	}
	return false, nil
}

// TranslateStr 翻译字符串 "green" -> "绿色"
func TranslateStr(str string) (string, error) {
	// "green" -> "{"green":"绿色"}"
	resp := TranslateData(str)
	// "{"green":"绿色"}" -> "绿色"
	var data map[string]string
	err := json.Unmarshal([]byte(resp), &data)
	if err != nil {
		return "", err
	}

	return data[str], nil
}

// TranslateSliceStr 翻译切片字符串 "["red","blue"]" -> ["红色","蓝色"]
func TranslateSliceStr(source string) (string, error) {
	var slice []string
	err := json.Unmarshal([]byte(source), &slice)
	if err != nil {
		fmt.Println("这里报错")
		return "", err
	}

	data, err := TranslateSlice(slice)
	if err != nil {
		fmt.Println("这里报错2")
		return "", err
	}

	jsonString, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return "", err
	}
	return string(jsonString), nil
}

// TranslateSliceStrByMap 通过Map翻译切片字符串 "["red","blue"]" -> ["红色","蓝色"]
func TranslateSliceStrByMap(source string, transMap map[string]string) (string, error) {
	var slice []string
	err := json.Unmarshal([]byte(source), &slice)
	if err != nil {
		return "", err
	}
	for k, v := range slice {
		if transMap[v] != "" {
			slice[k] = transMap[v]
		}
	}
	return fmt.Sprintf("%q", slice), nil
}

// TranslateTableStr 翻译二维数组字符串 "[["red","blue","green"],["big","small"]]" -> "[["红色","蓝色","绿色"],["大","小"]]"
func TranslateTableStr(source string) (string, map[string]string, error) {
	var table [][]string
	err := json.Unmarshal([]byte(source), &table)
	if err != nil {
		return "", nil, err
	}
	fmt.Println(table)

	res := make([][]string, 0)
	optionTransMap := make(map[string]string)
	for _, slice := range table {
		data, err := TranslateSlice(slice)
		if err != nil {
			return "", nil, err
		}
		if len(data) == len(slice) {
			for i := 0; i < len(slice); i++ {
				optionTransMap[slice[i]] = data[i]
			}
		}
		res = append(res, data)
	}

	jsonString, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return "", nil, err
	}
	return string(jsonString), optionTransMap, nil
}

// TranslateSlice 翻译切片 ["red","blue"] -> "["红色","蓝色"]"
func TranslateSlice(slice []string) ([]string, error) {
	// 将切片转换为 JSON 字符串 ["red","blue"] -> "{"red":"red","blue":"blue"}"
	mapFromSlice := make(map[string]string)
	for _, item := range slice {
		mapFromSlice[item] = item
	}
	fmt.Println(mapFromSlice)
	jsData, err := json.Marshal(mapFromSlice)
	if err != nil {
		return nil, err
	}
	jsStr := string(jsData)
	fmt.Println("kk", jsStr)
	// 翻译 JSON 字符串 "{"red":"red","blue":"blue"}" -> "{"red":"红色","blue":"蓝色"}"
	resp := TranslateData(jsStr)
	fmt.Println("翻译后的数据", resp)
	// 将翻译后的json字符串转成切片 "{"red":"红色","blue":"蓝色"}" -> ["红色","蓝色"]
	mapFromResp := make(map[string]string)
	err = json.Unmarshal([]byte(resp), &mapFromResp)
	if err != nil {
		return nil, err
	}
	var data []string
	for _, v := range slice {
		data = append(data, mapFromResp[v])
	}
	return data, nil
}

func TranslateData(str string) string {
	if str == "{\"red\":\"red\",\"blue\":\"blue\"}" {
		return "{\"red\":\"红色\",\"blue\":\"蓝色\"}"
	}
	if str == "{\"blue\":\"blue\",\"red\":\"red\"}" {
		return "{\"red\":\"红色\",\"blue\":\"蓝色\"}"
	}
	if str == "green" {
		return "{\"green\":\"绿色\"}"
	}
	if str == "{\"big\":\"big\",\"small\":\"small\"}" {
		return "{\"big\":\"大的\",\"small\":\"小的\"}"
	}
	if str == "{\"red\":\"red\",\"blue\":\"blue\",\"green\":\"green\"}" {
		return "{\"red\":\"红色\",\"blue\":\"蓝色\",\"green\":\"绿色\"}"
	}
	if str == "{\"blue\":\"blue\",\"green\":\"green\",\"red\":\"red\"}" {
		return "{\"red\":\"红色\",\"blue\":\"蓝色\",\"green\":\"绿色\"}"
	}
	return ""
}

func numberToChinese(num int) string {
	if num == 0 {
		return "零"
	}

	chineseDigits := "零一二三四五六七八九"
	chineseUnits := []string{"", "十", "百", "千"}
	chineseBigUnits := []string{"", "万", "亿", "兆"}

	var result strings.Builder
	unitIndex := 0
	bigUnitIndex := 0
	zeroFlag := false

	for num > 0 {
		part := num % 10
		if part != 0 {
			if zeroFlag {
				result.WriteString("零")
				zeroFlag = false
			}
			result.WriteString(string(chineseDigits[part]))
			result.WriteString(chineseUnits[unitIndex])
		} else {
			zeroFlag = true
		}

		num /= 10
		unitIndex++
		if unitIndex == 4 {
			unitIndex = 0
			result.WriteString(chineseBigUnits[bigUnitIndex])
			bigUnitIndex++
		}
	}

	// 反转结果
	resultString := result.String()
	resultString = reverse(resultString)

	// 处理末尾可能的零
	if strings.HasSuffix(resultString, "零") {
		resultString = resultString[:len(resultString)-1]
	}

	return resultString
}

// 反转字符串
func reverse(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		result[len(s)-1-i] = r
	}
	return string(result)
}

// 根据数字返回对应的"几等奖"
func getPrizeRank(num int) string {
	chineseNum := numberToChinese(num)
	return chineseNum + "等奖"
}

// ConvertToChineseNumber 将阿拉伯数字转换为中文数字
func ConvertToChineseNumber(number int) string {
	chineseDigits := []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
	numStr := ""

	strNumber := strconv.Itoa(number)
	for _, digit := range strNumber {
		numStr += chineseDigits[digit-'0'] + "等奖"
	}

	return numStr
}

// GetPrizeLevel 根据数字返回对应的中文奖等级表达形式
func GetCNPrizeLevel(number int) string {
	chineseDigits := []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
	if number <= 0 || number > 9 {
		return "无效的奖项"
	}

	numStr := ""

	strNumber := strconv.Itoa(number)
	for _, digit := range strNumber {
		numStr += chineseDigits[digit-'0']
	}
	numStr += "等奖"
	return numStr
}

type GQLShopInfo struct {
	Id              string `json:"id"`
	MyshopifyDomain string `json:"myshopifyDomain"`
	Name            string `json:"name"`
	PrimaryDomain   struct {
		Host string `json:"host"`
		Id   string `json:"id"`
		Kkk  struct {
			Host string `json:"host"`
			Id   string `json:"id"`
		} `json:"kkk"`
	} `json:"primaryDomain"`
}

// 生成 GraphQL fragment 字符串
// 生成 GraphQL fragment 字符串
func generateFragment(fragmentName string, fragmentType string, model interface{}) string {
	fields := generateFields(model)

	fragment := fmt.Sprintf(`
		fragment %s on %s {
			%s
		}
	`, fragmentName, fragmentType, fields)

	return fragment
}

// 递归生成嵌套字段
func generateFields(model interface{}) string {
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var fields []string
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		fieldType := field.Type
		if fieldType.Kind() == reflect.Struct {
			// 嵌套结构体，递归处理字段
			nestedFields := generateFields(val.Field(i).Interface())
			fields = append(fields, fmt.Sprintf("%s {\n%s\n}", jsonTag, nestedFields))
		} else {
			// 普通字段
			fields = append(fields, jsonTag)
		}
	}

	return strings.Join(fields, "\n")
}

func main() {
	// shopInfo := GQLShopInfo{}
	// fragment := generateFragment("shopFields", "Shop", shopInfo)
	// fmt.Println(fragment)
	testString()
}

// 小坑，空字符串strings.Split后的长度是1.
func testString() {
	input := "1,"
	elements := strings.Split(input, ",")
	fmt.Println(input == "", len(elements) == 0)
	fmt.Println(len(elements), elements)
}
