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
	str := `Boa noite
Como posso proceder à devolução destas calças que encomendei?
Cumprimentos
susana féria

Em seg., 28 de abr. de 2025 às 16:05, buddhastoneshop &lt;
contact@mail.buddhastoneshop.com&gt; escreveu:

&gt; The below items from your order are now out for delivery.
&gt;
&gt; ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌ ͏‌
&gt;
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AADL4ylXwwRH547t%2FWSfUUxjddL2aSNdzz0RArhQ%2BHbu6xZ2MeRm14ooSiClWsOG6G0XnC1ax%2FLG%2BqU3PKnMJfrDlqcfM7FaDi%2FA3Kf3K6uGJqndTXEBITbaWQxOmfNMFmWeNcqE4dzIJLR1nh%2FUoUq9serOtP6zdMDyQh%2B5AcyFAMTS6wboZEDHw1WbixL6rK%2BXmW3Bg%2BwheHZRYeg2RL7yO%2BR5UcfyDDaksFpqowzM64YhR%2B7KzNuIN5OHMqBbfuMj2RjiInyPD7GcMjP1Zn2gKAG3y66DWhc2fBan0Zbm9RYhQsdHgUqjC%2BSfEciNYV%2BoQbtdNjiBlza4ZjswpdOEgttvafXk9y0%2FJV4pISAG8we7dmm1%2B04OGgklxWRmfSg%3D&amp;c=AAArTUIEk%2BBtaSM8xJy9q18jtCOF5yB3A21fgKrQqEIgupwodb6v7JmUtxhPNBJBXlyg9X1oZmSWG717koqaK9PBqASI4RDpySmYxDGgxJSNGB1yX%2BgRt32wAqbzesvrbyahXP%2FmaO0n1CACVJYwFvPgBVvNlQ9JaYLXg%2F5eIn65RkMkzFRkriRP1Ip03zQdUOt7SnY8EBIO5AVFIbHJnv%2FZf1uc4ac7Tab8Zv9%2BCEus1EbLLwlg04ruL3bp1wHsmZjr9SSBwrHNUDlOWwxRBB9y5YZbFvDJvBfHjMhtPLhoICyGZRqxrFmfAuv8G4wfJd0noxWw46zypJLBxcRPM9ZT4xhdTyAT%2FIpFSOR9Ske0IiDYPEkBt3DWpE2EXL%2FsVPEuOoRmBHqf3xroNSQ4GnWXHRv60iO1ZUdTkjTKjmvDpT%2FZZ6lTuBM3ZeLZUxzqzu9QcBBO8unIflYeUZSAjawIWUfy2LAhHhxzgjgFHw%3D%3D">
&gt; Out for Delivery
&gt;
&gt; Hey Susana,
&gt;
&gt; The below items from your order are now out for delivery.
&gt; You can follow the status of your shipment by clicking the button below:
&gt; Track Package
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AAAyZVn6cY5Sg2tVVqCH5iw0pGE6BoBiR%2B6V2vEKm3B7b2Y2nCc0FaqJqjTXGDz2%2FvFRYzdHZyIKfYm2EdMoP9WnNV9SbTFyxXh2nTs3SAE77IUZbxsP6ZMkxFERe0p5bZKZGTBzPHQaQ7HCse9RXI18hLR4tlAA8SYeaxPaf51OAkQEplUUn1XBZolLeX5%2Fbl2imihq2p0NqFePWmiRs0c3W7zcqHiC1yAbhVrUiHhRzox7Zh37hlRYQFmOGUHx1nVdXFlo8jht1d2p6cIv4LPC47li8c7VGOsK3%2Fnw4o5sIHweencbVHjHEli8XOl4URNgU4ZwRii3mEJ7tgFfva565gSEgTkML4eR6ivBDUWZZYPsIOrhG3e5EpmDp%2BquhVf2nekC3963Nd3muhd1Q0%2FrEm5SdIxf2lADbvalMT%2F1YJEdxd07Ae3e4cWAxDlGeJoaj%2FxwnA%2FG1F8nfb81UyhZaiISrKL%2F55Mxrg08Ck83JRCFMem621oDi%2B9ZpNNkI2YHDc9r%2BA%3D%3D&amp;c=AAD2s5C5Ud3HobvO3BMi%2F0GwcgdPqIjV64DydzIKHSfTXxBHqoneKxBN5W6dcahIBOkj%2BgqKTB79QcOL2cycOKwEm3%2FMMjDzbfFYUXBoWgA6Iae992SdICzJr7C5N%2B7WA5HmlD9IU3wnFdx3%2FKNwh%2BgsEpFiXDOXeFP8ypGhVDdCkYQghmlU5%2F8ur1gNv062FvMoMQdF0ec1MogMnHu%2BE%2FEK6btGsVeqBP1dNQ1vpMvahJvro4LouBHiAhvjWNNdjw%2BA7zAiO8%2BlKgmaIQ2OU056b5AmtCGRt0ZpU%2BHY1LibRci%2BUK71XTDVzuZKrRbOvgG9oEeNJuIvjCD9Dp8C%2FJWby1mhGIrsBB9lw5U3YSY3%2BMvAeHYTa%2F%2F2RkTO%2BOMPWMHHYeB%2BtV3PJY24qm0%2Fl9VgVWVHGPKpVZXKOGhyov0vtH5L4h%2BrzaKRQWv83mTMV%2FtTctgWrwmVSUv90ep14xGzN2vD%2FNvSJTkhBxNC2A%3D%3D">
&gt; Please note, it could take some time for the tracking information to show
&gt; on the above.
&gt; Please do not hesitate to send an email to
&gt; contact@mail.buddhastoneshop.com if you have any questions at all.
&gt;
&gt; Many thanks,
&gt;
&gt; The buddhastoneshop team
&gt;
&gt;
&gt; Order No. BS395355
&gt;
&gt; 28/04/2025
&gt; Shipping Address
&gt;
&gt; Susana Carreiro
&gt; Rua dr Mário charrua,17 , 1 Esq
&gt; Algés, PT-11 1495-235
&gt; Portugal
&gt; Tel. +351919098833
&gt; Shipping Method
&gt;
&gt; Yun Express
&gt; Premium Shipping（7-20 Business Days）
&gt; YT2510321272015248
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AACm2k%2FsmJO%2BJtWpk6Gf8VDBRWJaG24aulzlSNupTzOda4rnB6QfkbEwyPQE9MUd86HBnAkA6I1WuEV5OvH2PzS9d5yrjtlLj0uMkr%2B51mLbsOB8RNdz37hB3hgzpvUbJ9Yxm3F3afxQwpPiG1qIfbuXKHo7bXheUOp7jeABaVBIStGZ9b2zJJ6lKfCbyQeuPXbEWcgCXS3iAoEuaDe7Q4fku8w%2FcQ76Zvjs2SjAq5JGJH1bv6plmCGDncLR0g5o6YCFDHyjed8TjyOhKqyUFkhnt3Oowp%2BLItQ6vYjc74iNRgoc1pYzNJqDwUqo&amp;c=AACbg7UO5Mu17g9oESkbedMZu5rfHqbdgObiCQ5CnvVJTNUJ04O8FAIrbyW238ytxJtRKx1ma211EjRdUA1OI%2FO%2F5GSSBrxqxfmepkAyqR3u%2BddlekI%2F8lIlBQDNlDLgEoQ2d2qOE%2FuEIol5s0MGYukkVQFR4Ac72ecB73GJyk20IVKlDrMe9lWKydh67dvqKQGY0VVzq63PtvQH%2BI6BoFIy4xdq6U4QqQtkn64J17BudVQXdBvrjfNUuGEcV%2BPsmd%2Fvd4U%2BrBzZV2xYexZUYSFp0uAqmHTMHk9VN%2BzqRgY6fmJfhOkquUL6vP8fbKv9KdWODCR9RXd5EOxifCpnVQMwih8nk9UoDTO19GIJG5D4TAAkE2jVDVOm9teOF3%2BkNfVL08XfbDYY0TL1hOqCSDrowFbBf%2BMVEe7KLXD3g9sibIDYNn3ZXyPbxaUTaUowVil0sweM82%2BAlpjn%2Futkn%2Bw1a%2BKJisna28RFWNOnqg%3D%3D">
&gt;
&gt; Shipped Items
&gt;
&gt; 1/1
&gt; Items in Shipment
&gt; [image: Buddha Stones Gray Golden Koi Fish Lotus Print Gym Leggings
&gt; Women&#39;s Yoga Pants - CadetBlue / US6，UK/AU10，EU38 (M)]
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AACS%2FtSSfOZf4alwmN4zjBcbSNQG1YkU5UrApBP5srtxlIrY4ZQy%2BKFMwsNN7%2FuiVWGpswcw1hOu4S9bZqVJwFK%2B4TA2WxDflFseS0ZxakOpwKzuwZMQJIazU4Rdfa9dXz2c24p%2BnYYzkQbAXxasB7fms6HTvWmThOn7G77HBTbclZgkntNVvBREqrIQCjES9OqAtmWvHkTwsilAvzYTomz3Iw72NQ046JgEA2SvzqH2inXWU7PSxi9gFtE8OQ0NlA9YAwS4zaQhLobhKFmh1Ym5plvjOTyUYzgD12VE23RwX0QQWBdRulzhNvvLKSjf47zFoYMGl5qKUfK%2BXfEfSHWBmRUkfXVgRPdzuR%2Fw%2BMhbwWio0J0IiMDeYR2bxZ4nm92DpsyyS6eusIqGx%2BkcMatluuC7e%2FslA4ooTYmRpoxTm4HeE8EccoWuyGHC2CrIa4071ySB8Tx73W3hrtduewOekAr4ilxk9fkPP9oX8uNB4F8mycXN2vyhEH0%3D&amp;c=AAAcxkj0ZtnXAaR6Sa5Tgwas%2BQBS8cOqO33fjLILqQ7WYrxx5eep5KYM6z%2B5%2BVSJqZWFXrRiJrSn1YlG0yBMyu59DqEHo4QTk9LqQ4TVjCWRfH5rQjZV5jMBGqNsdRrUsvh1dh%2BnhyOVffiqzYrfnuuSsxASFaB%2FBbkEbmScHexUCX8loH3FEC89hDAWW%2B8a5fH1coSQiFVTqbnQiOM0hSAU3R5cBhIHIv2WJbBdrWH1otR27VkD%2B4GAPtVyQa5xOonOi6Z20%2BxZNlUZ5zjYHGf%2BcHHvyIHeZ88jr19brwc8rGreGfc%2FiGoNdN7Gs0j6ic%2FDGqz%2Bsn8aoGR9GJM9B2MrIHBW2keQdu%2FkgQEH4nkdGSwfxcdEMQagqCO%2BfKfrAh8Sct50faajty4Fk4ZsyI1DlKBTGBUMr1sghqOESCiuHPJTdpEFu1GAtgc2pZWDKDd1AsOxpzt7bP4VyTYqYgQydP6pqV2oYZebQA29Ew%3D%3D">
&gt;
&gt; Buddha Stones Gray Golden Koi Fish Lotus Print Gym Leggings Women&#39;s Yoga
&gt; Pants
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AAA78RSrHrDR6lvZw0Q5VXrlFGNFyq%2FBpS5e6mn49H0esQiyb6AMpVoBtHSk97BVIw8bxN6zermTuAGdt6xA4hDJ1aWXyBUIZ8fmoDPCjVnL6hJliiOZh6LeGVEFXUM7a2jPQRZMGjqFrmtIkGvRmuXVQ5vk6mAJ4REgTwriXEFRzWcLHdjVhYn%2FIr2F7ErwYjXnMsWeH%2Be7N%2Fic3XIzNenM1%2FzvlajUtVWe1S949UxUOiqzW7DzqAayiBivANRUYElfF3O73wDshc6vkPxyc3CCVZWVgrM%2FRhJ82OKwe92o%2BBjTAj8cdec8AK9OZ98siNTR8wLn7YDSh2qU1mgbJAuAAiytQZ%2Bz5TC%2Bs5Hh4LroYXOZjwTIrNCLtjCS1d1ZINEsU%2FqUPHoYI7AzstXOPy4hm9kfM7662Y%2BoO7G%2Fjxo2A6jJhHk5mTtSdG8JOZrB0n%2B9mJ9ys%2FVNqqY1QdEPr3DPWL7Ajd7axypwntJ9EWrqiYcCSCsQFGzGQec%3D&amp;c=AADZG25YvEbqpmCgH0w1zBUB%2FnO1hWOIVRDe%2Bm4gzfk5LPAuQkvP9bjOhaAE0lT9663Sg4UvhnrT46ytfkxHCBy5ksGg91ZZBmlD8FGG6GGQ9f3p5FrYYYMHtaoCVEwJMFoDu31G2ZKQNNTxULtWZ%2B6Omq5%2B%2Bbi44S9e1cEX%2BM4VSlN4YAmkBEhCcH6NsuixmYJv0WIVqbNKR9lT0VwOT%2B3AgQTvZVdXnFYEfc8LlODIUKy3wZ1K0yjt9qtU8wPj8LAys2RXPN8O37D%2BkPYNYTRMOtG8N60p%2BHDNUZGFW1fg%2BydmBkDGqt0CoADBC7Y%2B1LDYe%2BfmoKUtHBJEZcxtT1bQCY7a1ToQ8xA9ir%2F7JI86jML0RwAwgTLlUt4yU3MvlOIxayKlHaryEFvVodMIvEgczFvscBToqNnWH3c4PDMl4DPP5f3LGrpVB1SZE4IGR9gByXyHR1XkdT%2FufnmUo7iwmyWWj3QA4hfRR2z8Fg%3D%3D">
&gt; CadetBlue / US6，UK/AU10，EU38 (M)
&gt;
&gt; x 1
&gt; You May Also Like
&gt; [image: Buddha Stones Red String Jade Luck Fortune Knot Braided Couple
&gt; Bracelet]
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AAAqQnuZHabrRVyNF0dQAFAxTbWlImOct8vRhjXWQ58po88aMjE6tavdvMYsRsJwgWRvisWgPVZBEYp2QprOTAOHuQTkQcUFrgxSJ1qigwvWOU%2F%2Bh8sgzbiAMrmuCkUDmpm%2BuNMbR9pA98iBhwAMi7dS9nppYo3CBvKlCTeY4v5K9ITWmRTJcTCFbiWMM%2FSopxfxetMBRMZSlOig8rHCzgzla0sBHRkbnrT1Of1jodOmQ2raqGoDxulYHhw2z1IprjGCembWKibqSdu4i%2BaKutoMYbPIdGit16v4bm%2FAs%2F%2FD2WEv7b9iV0p4C11etHyFfGnHXbnIEVrebNPec3rUthj8rzjUqoKg6C0a8fXrEQN%2BMKxz6AWmSo5O8K%2Bi3twZjSo10UNfnq%2FpLBcIHc3asEmdNIhCYQkZWWHZmu%2FVWOIhrJg%2F6lgV0E180ZypywfLZQjrTJqyngTgY6AS%2F%2FXOS7p5uxcFQm0Oehynp%2BhCbn3TmcftKxMlGGyj0X4kn7We%2FMzZgeFJX3%2BkHZxYp1NDUhzcAh891Vvm5nl4yBp8QRI72ko74ZFf8WfbDlYaSqMn%2BX61LF6s6ZPUP0pV&amp;c=AADwen96Np%2BU8ud7D%2FkP3SK9j3F8l2UUgTa0px5tR7%2BmNnC9VA3OHiAkgTubDVmTtKlD6K5x5ZMdSKwPw5BamvEDKS7B1iYyGGhzQJEhHw%2FbwNEqOkbYyrZB7Xi6WfqkiEWWvuCjcvyeexTtnv6na6Zx9iioiyUOlYuL%2BAgnUDa6m955E55A%2FUt9h4qrh4kFSgwGCTRpW0rn5TUXzdl2qmjzwQAahU6yd20K79sMtQWmkWNOQA%2F2She6Vw04fEoUTSUXRiwKN5IKoIyE73eqBCJcovnbyma071T7HWSyCgRcmdnu8WP5Dh%2FYit8qDdsOpaEDoUH4wTxspwQgDXZ8yNReRarSTLUNXbXtHMvKDQcmR9eb35ccq7XZhTFI%2FykO3FKjLZQjjQlRBci7GcZSBT3SSzPLjuTQiCck8sIMBtV6yi%2Bt0M1LxMf%2BCTAfFtAX%2ByxEcMi0eyQWmUgOXonQihKP8nD65RyNSF14NWEimQ%3D%3D">
&gt;
&gt; Buddha Stones Red String Jade Luck Fortune Knot Braided Couple Bracelet
&gt;
&gt; SHOP NOW &gt;
&gt;
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AADIywq7fm9kpKqSdxyI%2B2mIHlM9rTtgBv9eH3hB9DH73KrG1mG%2BsJVjFPwG97XlChJs2yh5vgAt5GekABp35ZRN%2FulTYH2WFb0ZERRh20aH0l1ryD496XgrEAKVjd0Ha6nkXDpwCv%2BdVeaIXHOferxgW%2Bsvi4dqvtYJHbOt0pcm2kUhfm6bmIez8dUUrCzy6A1%2Bn9NWV0OiMEVugI9bTwuD2P5uGDzxcPM5AUfz6gvRt68RUjKVkoKQHYQ9LOgnH5w2KmTHPFBazegEYIKX9F8CuNJ%2FwgR2pSTnsI0cmlpy%2BtFkCtjTKBb5yTgZLRUUhE0W8yWlroCUu2PZSSvaukhwdAWRbdAEyKn5YDmODR3HpaPyRO3%2B%2BTd%2FDhgvkAgIjqgro78tCTSnqze%2B0z1u8lT8HXVgEZcR8bOmDhO%2BPVPBWja%2BSba%2FMictGN2ZkNbWILpEHIvepCDHPmfw5e%2BVA3t72VyQWbIB74f7NVABp8YfNqi5PyDnbLteM%2FX4pL3Gb98Hkb0%2FTGBxDeiMC2FAMkQRcJPWNoUHvjvc3miGLmVe%2B%2BLU2Q6RaKtzRQ%2BzAS4cm7WwgQuIuQQHYnYg&amp;c=AAA33A%2B50kbr8HvjgZV%2FoHhDryj6HqrRe%2FnViHeKWT0Mgm%2BZzVXjmwzNjmziuPBXeAAdnQMmCNyCvD7sY0W%2F82eo%2Bkf7EXz%2BD4qM8ePulbLOkQNgBiVlZyeoiG%2BBvKBvSSO7shldgqRYDVy64BN%2FxvR6L1IU7JW%2BOKMM7CFLiuDBTekpLCt7XlN8eFkEv0fmZItVwA7teI0m19seI3FUsNH8X0FJhn%2BAbmaNP8bxgXsMaZ8SBUaT1V%2Bm1se976ibyAfEuFzdUwiuidJTpfWKqUxpKnYAX4QffVHv99yOAl6%2By97X%2FqVsJLkMQSzgvY8adrRFwqJc8nT0TD97cXx0HWW4yePU0%2BtznEQ6hJCET7Fknq2tsVHj8T8pntQnem4%2BzwVYesSAvSQgzgHt%2Fb2vziCsEmkFR3L4npBi7jiQAz%2Bgae2VO49LCXw6qjzdEPg8buHP9PKa%2FpSrRHeG3tB7z7uEnxIyd9Bpt8lT5bqlVg%3D%3D">
&gt; [image: Buddha Stones Natural White Jade Koi Lucky Necklace]
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AADPxQx36gAS%2F77pOBguHGOnvNGK8JbUmgO8YHYJW42Sy88alvn346J4XP3%2FrAM%2FORhfaG9jmsQ%2BAXF%2F7J751h77XLOkYrNlM31qJ8PqV8vOh9TwRjlkC1LUJO8kWYPpOySMBGo8EkC%2BJJ1sCu83zFtXtTLm8%2FkqGdfIb4gM4CXabtmY0kiMhOQUW79lqUjyPMUWenORV7dIt6BTJLWSOYNKmEH%2BXapSW79rzcJitH1H2h4axxxHlrcpSvy9iqNQ%2FvjnCI%2B1V5jtAgQZ%2Bs4i0HqNk5igrXUN%2FwYCd0vmDm0dcCqM0WA5nr%2FSejyb0kLdZYhf8dlEKmElZvphJXU9iGr4Fc1n1VcKLYrS9v%2BkCID7Z%2BCPvu%2BIBzeFaGfP%2Bt%2FfChjrt2S%2FCCEKdGJKniHwOronBBo2Hlrq%2FVbEIFr0BVP1YWVwy4QGcrii7t7o387EmmL9s9G9f9u6CBGJP3Bdrt3lWbQi7Ps8Z2WEPmklYniTabb3%2FwP%2BrsGCrfEl%2BWrnK4qqfg2XGfJvDgOF%2BJDR9T%2BRdV4%3D&amp;c=AAC4wbmX6MH0lQ%2BkIWCCzoeCxBJcbI7jqXL8JBj2svoZnI2i8QuXZYn0JvXfoG4cwUhtr5LIg1YQ7SRVf6xWuyTcbOH3mDx%2BJEcRQ0kwM%2FOUqA2YSS47xHAsLi95Wi8fEzNmxvwfL856godjP3yfMlqvSBEUfh18jcFSSniE7f9vDV0M3HCGl0iDajbqXkegpAoSsU3GL2FeCmRFeP7a5DJOX1IAmovTNz1lhH3gukBm85EXuA2U9LzJ8gv0zDgtgT%2BQdzWwNqueMomNl34bA%2BKzT0l%2F2rfj%2BdDg6aSb6gL4wRDP9ZVnZEZs1Bf5j8f%2BmXo0SMa0PuGBNmZfrxOruBUAHpmeUYr9r%2FNZh6vEIitEUyVVjNd6Oo8jP7ucTASmyyXeRzV4iGLN3N4ci5FwLS%2FSB8f8nnUflC9mdJrzQZ%2FATsS%2FdnSaEodsgm7rAJT%2BBAofD6ouGK47C9acqsoxR8KWkV7BRPzcGyTopZkQ%2Fw%3D%3D">
&gt;
&gt; Buddha Stones Natural White Jade Koi Lucky Necklace
&gt;
&gt; SHOP NOW &gt;
&gt;
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AAAz82c%2BE9zVoDizTcKqxQhtn0vrcaYhFAFsEupsNvMDhMnWE16HbZbrlAKu95ZDpM7UU%2Bif91pmdHwqR%2Bi9CotVYoB05Y6LN7G4ezWAz7TuzHzXk%2BRkWzjVTGdYCWczcqv8IpTzNSIpSvvH4j%2Bys6XG95WQb38JaIP9CUNkebksQke9k3KD0GAWOEYrnQG59oD%2Bb%2BExLIl7PPxFYNhsN7NaS0Mg6ui10xgiyYYWB84FAxPbkJLJCx%2Fi55N%2BPpSFqLA44LiD3JxzwqEh8uRZUj4RrrZc%2B0sNfDMKV6IfJHc%2FGe6dDB25xNH%2FD%2FZbxa6wc%2FF3U%2Bxc%2BxXaffpV49QTEH0XzXRRmviKeKABvuxd9yoJt5K8YbNKHVjgvVKXbEQIHZHniPUnwPrRIr%2Bk%2B%2BQFbp6enJ8Y5buXCcZXLoCwpagVrk9g53E0T6VFWRYyRFlOeEtirqV%2B1iwf8fRiH90exfdmYkVEE3R0wc5eHvb5l3rG7Z2bAH9e%2Fft%2F56JOYvOGhu%2BFyhGClavTfZnmpy49%2FqxC6vc%3D&amp;c=AABU4NE0mh%2BfluyNkXJpTwSPmKb%2Bkuv04cbtxGMglzuMv0tMtScKH8sCNjQA4OhQ6Yo6Ctx7XRHGj4L226AAaPX36q5ZE%2FYhue94J%2BofY21bdz5sgHvog%2BhWOjaU4B7kj1FLzLeAyWQTaMwvKgnLcWSZ1mMqWYvQbanKWXVjr7nsqdtcXY9vjurl85oYKwmE2scFv5w5TfOibwaCAmghvcO3SDUFFevnZ62hTG29n9G736YDmChZ4DnNgm%2BL4MD9ZexwcD2MvmnlzOY7N6rwtrCLjaaXDMiXFDAHeaRKHZuZ0LZ4r5LZQeGfVjy2LXV3rAQW4hq9nN3Hd6GjTeTj7dnkCy6Vol2vYPPH%2BXg9GLVLt5aJYMdhLznsY%2BztbFXbNNW9kJ2v%2FFjion10KtEro7FfrfCWPSW4eiZ0dOFtP2z2C6lLYBJ%2BONSryZ9iw0WdoEf5YXt2%2FECy6%2FmKN4mhksOf5bmUQ%2BkZ2A3jtUWJnQ%3D%3D">
&gt; [image: Buddha Stones Natural White Jade Luck Bracelet]
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AADpVKItx2hM%2F%2FkRUgaq6wDQ7DrWTCgh38UakuVF%2FO96jzbcXVpD1wcSQ4qcVWtlp%2B1K7MiOdT7AFbPTpY34CyDUzgueOuYu9FenYc0IM5oIzDrP9Wb0SZnEbC%2BAEwsaGA0V0HpotzMZyM%2B9Wo68itcpx3d99ncQAY3sZxn0zsQpSXFuCHSuUdoJi7XbDQYcAFV%2B9zrHIBgtC3gWKvCqjYPjFdL9pYeyL7zXvllI3S5SHcYa1zS05S%2Fv4yEtSfRYdGr6ZOTeXo1oR0SdN%2FYpggjzWHnpqX5wQlBDR3gL4md7J%2B1Gk6jh3NegojwW28yoqp1NO8%2FudXwLDTBcydiTaCoP4TD3vcwptMrI3elEcNY%2Fs53DW0c5RtfIfgf3AMmhfJOco4oP3%2BGr3XJVccPzIld1hXYKMtwNt7ezxVsyM%2BXsaibcGQsrJuuGl%2BBt5MU07wcy%2BxcGayHUVpheponR89O9Cnw8dGqgAv3BSM0SxWzvll5XyJMt6OzUf31s%2BgkC9CKimdZOaUsQkqsPHQE%3D&amp;c=AAASgGZOMTUFnNkwsNekIrMt8u6llIsoIzmKDAbzna%2BaYfUEy6Bdhm9ysrrzoqcLFG04t%2B5MmcFPrN%2FHDCKOADqD9MCxeehyYxBSqxNZqdOYOTp8pAnje2tkpCe0UZaFZ7b8a0MyTG5van1vW0h8rx9N45rwDu0i55uHY2iRsZ9dxp%2BiD7bX1xRWKvGIS6srbA2OVNRBLzzGmkOMbV3ku6tr9Jfz0rSjHhfbHRilsNycMXY8ZaWAsUfdbMlc%2BXzQjA%2BzvJjcd4MTP9A8d8N9HpB9n16o9HAXZjRZGsMVNa6hyQ86ohmxTpP2F0WjMuOqkdVQ4wnac6sdOtcQfbfCn%2FYmLjc3mupgmYyUxY33Nb3F7jqsW4mZtiyofry6prRD27VkGFGCXXr1GDdcfBCq7e%2B3OFx0kYfOdeEInB2Dx7F0Du1us9I87mMEOOMQvk8S65iF%2FXWHlEcOwgJulBkFPPfV%2Fe%2B2ktKryLJdla0n0Q%3D%3D">
&gt;
&gt; Buddha Stones Natural White Jade Luck Bracelet
&gt;
&gt; SHOP NOW &gt;
&gt;
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AADhY5L58Zi998fAZXRRJgqgdguZY2PH3uAsU9ikTcT8gB1NGRLCVVIfDVlcWZZ6wngQtl5cCZ7wfA0i3WYKLwaBh9jSGVa8enhvFIuMLRsy0heb9S4HlGZVwUZ4mscGLRC3WbVF6NAiR50zjzQ8sIdVqN5Qwzp%2FiEvJzT6yl7Itsqp191SSMrq%2Fbo%2FV4qQBrIK98nkp0Y1IlWOFQuyrHXrFMN7ycNEUD4Ka1iksyuJLIqjIVyuUIBsI5glehZEELCcnEoPz286o38RVOT76PAyy5360H7B9TbiB6NvyNQBtsZPRynK%2BlVf612oUc5pEEpeImFz9Y4aTm6IPvWcLB6h4gh%2BNYVO5DPl4r8%2FaDNheExRXlkjs4uXKt2ZO8kIvVSr1rCXSahUxRJCbBLUxyHv48xBDIFpjY42UTvWIdG1lrNVU5cp3%2BMMFeg6vsSyVeQCbQ317op%2BauObSljvE%2F6ukMCWfmpsPnfzuga%2Bsj%2BRpjsWVKppXn%2FYxOK3btaWtUmnzkQIldQUvMeG%2B7E8%3D&amp;c=AABR36JvflkdnkNXoWPWSLATTYFfHYEV2Af3Mw2kWUcL0hIwEZpchjyYZres1a0HVJJcEOKa1Mjmb%2FuGKXPTMc6CNLYdYyOYEZMrvFTPGL75X8dl94CvxAfiH2PrxCKYSg85CecybXxxHzdWkDEcUT6vc0Ej%2F27uL8KetIVzTzysf9gDwjIr6Z1%2F3z%2FZNYxxEoXGl1xdG%2BKlOHd3qILY1zjKAliFWhMGlTVMleKWPcBblLRb1j49VheF%2F5JSCUhChyAsl4mB8EXPEBXR4dtRk1%2Fxy%2FIkho5%2BRWyWUJT6BBIdwZJVQlbds3I3P8eybpyoYtR%2F1d0fR98yPIN4S5P1SxSzk0dq9mcQfkCp%2BM9OCASM2WN3R5eUjMToX%2BltocU1kAl3VBHqgL67psIsm%2Bh%2FF1DF8uHHvqXSI14RvgkIpmI91%2BXAuoFUACxL0KZcibnSGroFAhaHWqvYOku7xKj6veYvd7TD6czrZGfGsf5ZSg%3D%3D">
&gt; [image: Buddha Stones White Jade Bodhi Lotus Mala Harmony Necklace
&gt; Bracelet]
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AABMzqcFh2Scx4FLisTjuZwKkJ45Ba04qm%2FwqEB3FymIN%2F2EO6O0gKxXBL30UEYH%2FOZaVq0J0TX63Uk5mRyT2Ceak71ruSX7F4v24flw7rr%2BrHMqcdVrmC3eNdPiVRAQnSp%2B%2FZ61l3MMooOSzlqJ4C6HFMUFcvUGphfh6CQjOdFX%2FvqHUTskSqsBFnzLX%2B%2By%2B527Rb9xmCrIngVTh2cm5dtvpbDrHRYI44Lgp3Jr5Ody4vBN5eU3p4uN0bpJ5FUXUDDLRyFS8cWvESn5NaQQ%2F5UZBj5kC5ex2vSftlpR2vdTV57ZOe6e5qgI4DUHTGin3bwJuIOOFRj0L1DjjDgQ2pjvYdQQfmDkt7lPXM%2BX9GyonDRxcZHnUUnYdaMUnwzvIzcEqVhmSMNlxxsWdgbNjtRcSCpd0%2Bzo8mlcyFADaAIKkqNmLT3Oo5JucqCjryxVHYzakT6KI2fsbRcTK%2FE653%2FYwDHT9VDQ58qLvuKh2A9tJb%2B7dJvEq9D2L9ZRSzyLBsUSnF5xkkUSQB%2BzXx3q6mRpepX1MzKMzpJLFPGYbKSZqLrcNVFM%2BKDycea5ifw7CBLxYg%3D%3D&amp;c=AACt5UP0WzTlFmaqk4obyhMu25qC9%2FhqrbjVrlp3ztK9ucbcl8Yi9LkCviuBovTJpDNzpVFu2cYsJGj0AAOZTr5UfZHfQEyxN9coi8g9Fx6RM2l%2FfQqczmmar88gU1h4C83cDPPy1BV8WftLKkZe9%2F0iyfX0EBdSIs2Z%2BirqzaqmQBsm74OQWMEARDTf2h%2BFuwgS3LZajMRfQZZCyBUyoN1Mm0tO3J5Mj9tA4mePiFSqxNAwxI3U5a7tJcy9bC%2FEeVZVuL5DhEuWkzYnBveVD45LwOQKg4gTgT6qFprxxN1xaacqNnvwUr1fZtt4UzGJ%2BV2fT5oTAgzQyx%2FbxvHMDzCPt%2BrDQdfJfBfxB%2BchjO7hlW2lLKVyZ1kk1pCdV0MQEzYR%2BvSDQWdXYxWQPPm%2FPgHyHRzQipYWV845Pcv9YVenoqKMZv5RMDnsk7AGnffko5kiIDhhupA73VM7zabdcxsnZxXOiDU1xj4V3JhDaA%3D%3D">
&gt;
&gt; Buddha Stones White Jade Bodhi Lotus Mala Harmony Necklace Bracelet
&gt;
&gt; SHOP NOW &gt;
&gt;
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AADuJdVesU41XWTyeoiWr1G9lKkQBie%2BzNrSsgM4uNUIxaD3ZTib%2BJtMTZjiHe4RybT4TjrVgjkapqmE7N5teUiJ9KzKKCmhn64J%2Fl%2BS6kHMR7DADy2F3P%2FrVnvhc83v3HPu%2BC9%2BXRLMqSByctDthZBfCXgocok%2BE%2BSxV5sD988MMRPYD%2BI8BkzgH0nBjS9%2B9XbvzaqL1yN%2FqpfF2xiRgBz13TRj%2Fg%2FuQa2IRe7F7J65yKwgyvN8DFlwscfj5fca9lfrLshzaI2nMdakTDdg8wP0UW6J4DS0%2BabYuoHRNAkJS3hCBX%2FEDC509sXpQj8r0G6ulNJoHrTwZqptbL8iWzvExyQQhXZgCks10vhPWVcRJqv5nWbmDO%2BhteH1uZxZ5%2B5G3HDCF3GGnp54mRraIMipxgKjUaXMTmyRTLFOKl7yjalLI4Kl%2Bek7FHcaxa7itM%2F7x9yXPAdrWT3ECvIBhNPwIlKaSa3Wm1r3oUO0j%2FQhAGsJJq0a7eZu8AyRFsqKmNz54CoNa67H%2FxDLnbuadHpVbp4EFB4kWkJxdRcdKR2yJVIHtJT9zwLlWgQ7jKKhVVfIwA%3D%3D&amp;c=AAA3ia31KvA%2BloZM2xJ2sz5snDoXBCG9HHbZXjcdj9T%2BIzJeMVa98l4VnKXSZhbcRzXpHV1PWQUwgLdlIuXr6q1LLUtR9qg3JBr%2B3X510xHPg3N8bil8vTHZZ6yV8O7jVABw2oMxPz5chF1%2FiJ%2FtP%2FYbC4yT1qtzFY8PJK49HORsB%2FGB3k8g6aT71dCbJz8GdBxKGWENS1ZEa4FHDfwN523YqMlBMH%2BywsxmMdZl%2FlHN9hiFh3wZcXrDkf86RTvkit7MJIeYIwzvZXq%2BbG2E34dC%2Fh1xtTEguN3ImEGjwdae%2FHxSKnADdrFeL63CEwYGkair8ovHhju3B9HfmynhDQSUyn3GzDT0hGpfpOUbtgWyyYjqiT%2F3HAQZdIndgDfBr0DBxlyOvk92Hkwi3k6OC65Dh3tAsPJJd5nFNO%2BQD%2FZG3EGpmjodDtQhtvrgAWTxhd7KU2IRmlc0VTwMH4HEqQcwU4WZmU1oijSLGE7vbA%3D%3D">
&gt; Shop
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AABx20wC5wm5KQJDtBJQ28bf6YclKx3CGfx8TKaWMCSnIXlAz2C%2B9AeWp8K%2BVvYKj3kHIceWO6fCTLYssG3rNj5AoKSRI%2FPgRJxi%2FAvVM0YAS%2Bfum8pTaetH4uKsrMmbWePZZQDhrewlCODYXb4wLTTVdwGFIcDMGlXyKNiMmYJPJhxsFMaS3%2FMXMFBhE%2B%2BZdabON6NtRjy7rJaJP8pnRgsVQh07w%2BctW4pu2iDBQfEtBMKNU%2FxfGkgeenOrb1n3GzI8rzJvlRu9itXF%2FNYSkv3PmN8PWo%2BPenZD5minuU0lOFhZVygX2augtrUaXGukd%2B6uTCu9z1vu5fID97YSZ3hq6QGIUChup%2BRwm11Y6GFSzfYhb1iHmWixwk9qdtjT77S7Pq3jypnthvknDroM8GROdRvVIK0YOTQRJGh4v2YFF4OVXpuPf9wTqSUu5%2B0Ywalt&amp;c=AAAtJuPl4iqyYQTh%2BOUCW3j4F5wr3DFKV36aqAQlAa6bA3p4e4LLHbY51hVIZ2CAuziWIRb1NzHSL1mpkmr9DkWUtD7ZwGnEGuYBvrFhp2%2FruS8D6Mf54%2FAwfBoaBXA3fz%2BZsLi3o1DkHPTLtRH2kADkdTCToB2ERsOxvs4IsKpDbZFYNwRzWHmOr%2Bh4XLkNxoKsyP8rB4lI1feeCN45zdKe9zJ%2Bj6mKjt8qXMhiItxjler%2B7aGPd2%2BcdH7ab9Pb%2FW6l4lpj77AOqYeU6JJspkGlGg5IyAJZEvFq6p5oM4ys4QOnDdUYpAo5Qp09nKxxXFisfxxVq182%2Bs5eqH81qtlW6fr2DQ062nOvdwzIqkcPzflm%2F3FppzV8xRvr06rR9v31aNKo%2F860NES100vuqWn4Lk2CF2PiJxDvdOS8RqJgx1ZJERC1%2B0OQuRu7PSGgE3vSltOCBsPBNkS37FEXW4AzWVjvQ7gkEfbG%2B0uoEA%3D%3D"> About
&gt; Buddha Stones
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AAAumYRtGU8S2OXFSfnRwnRjZG2tdC5QF%2FwuEzTuFD5XlbKrDB06o%2Bydl%2FN7pKd7x0bLZkdF8axmoSglZQ2tjUdKSSOM10H4tAcfUHC%2FnI4hmb0n35nnVF4oKznirDDkF%2FOoEWayTASHFNIbmNMVZhcXoX8wWMtLMltRsgJSOH7TwkWZ8Ql9l4G5uASRFLiZcHr5MJDZY97hJ8Hj%2FYgIfrYl42uycAwwDHaRYq9DQtqCqsmAKCX7XcdblQXm3jul2ccSf0byBrfFkBnSyeOxYq5ZXQ5uzSTESIsC7ZduFYkJrUDjYnkOANViFL%2BvlD8%2FfZdFAmuXLccv31%2BuVcCQv4De2DXNn0FUzurV1bNW31ZFPrDcaRsmdDhX7Zsbm4UiFFFnF%2BD7Cw06F3I53ab6Blq58Rp3FWfxtH5YV9JwXH6WfRnXoWBHELOaQIpjkGHZJ3mAso8%2BTbyyHf9KigiCXdJNTA%3D%3D&amp;c=AAAZ2SptjKxRBLdFUQ24eMsS9HRTIKHXIOywWcUF5QYU6w7Z0THj9GfODbE5A1hThdGrSx%2Fl%2F4J31MnqRprEy3L2yBa%2FcfA7tRYDpAEYrzbsXg1vpz28rmiyGs580G%2Fs7yO7yJHFIzJh4tHAsYRvYofjg%2Bsb7TZntTT06zJEu2cOMGHLks9UtXv2EhUxgIYBMBuPGLryrdNTg0oZfRdf6%2FQLAM0VmLI7jl2w3CPNBVcaNOOVrHLe3XCc6%2Flmy0CDXifQ3wOLdqyXY%2BkjYMBqfowQ8YVG2x8DrfXezYaEdX7kdRgevRkI0Z3Owm4KOpd7x%2F%2BaSPKyyBkM%2FG4SM1EoRX8j9fBkaeqJDNG9PnFXnAanWPAPrzADqbRJERTUgz8Bg7QsR4FJgW64BMx0BPhE6hGHLiKFhtof5%2B35OFZ7pGeq5FgF4jszvkxz20YYEVyT%2B7pti4mudkUzIOpDmPOhlZWBNvQua10XwYJEAK7h%2Fg%3D%3D"> Contact
&gt;
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AAA6phXamt1TuNAByCIawJuPMvmlpMuD%2FjJOQN8NR8xPQ3yNia0TmiVxkKeg1nPByoH1vDlAgVv9L52E3ivtkvCAKRsAlUoZpoaYIOKiG%2Bq%2BWJ95Jjpwpft7uxJlEv1NcDNAG4Em4C%2B9FDoyoX6FCS0AivuQ3o0L8IXm0ZLfYn8Jl%2BzzNIZKFYYIGECiyt8tEXFd5BjGDyQ0RXGKX8AlvlLyGJZNEX0yrcc3S9cVrTYwcsHSs0Q%2FPa7z1e9dBRUZ6dYdHKAgJtP1MBIY5fS7PJIo0DGceNVnSBIUx3lD9kwrCBRPDel17qvGObX7fPCIUz0GdLM7hc5Yz3ZcLdBZKry0XwOeqtTg%2FV53JmrxmKAuFuLohr2oJRFKtx54E5z6ZlIyy2%2BJmuGJvDJ4X2v%2Fs0MHm01eWeO69VPKObiMpbXz7TQWhqoiowgvMuPMXg%3D%3D&amp;c=AADEOMvy06EnWiDLx8eAXU8iC5BDqqygqN0jhBTw31d4yWrkCoA8iudjMvZc8983GcqDI5WOxPQH5gJYSeXYDiVURRLUwvJmMDzEzt34Eax3QaSqjuUbUydrV3mJ1XlkFLurXNeCDr54CsJ9vEbT57pQ2%2FIOH%2FfJrvZRVoC%2BjfdEyKEAypLqos%2FqNvdAUBVd2QW1JwlTpS26LEJP88JMYt8uvkA5azjBga92xZ2nn1Vfww1hUDnm%2BsXEdyAT29bMtTj2XqCFcOCMXIJ3T0vNF1CnKQZjq4iyoGKUl5Fcv46ubeG%2Fwqm4rFTWPuL7rurOQX0AlzsRW8bkTyPH6A7WHIrh8XuohpD79e9A2V5gpaMOzcAaIi53d%2FDA6IOxKLWKLtXNcKZbnUCb3%2BVXMLv9QlnCmW1X7MdB%2B4I752DJFGRF3Wr5k1oH%2B567AJFH1EUxQ25zB%2Fe9yNC14rmSLvYJXmxAdoAk%2FIxJsp63ioOTrg%3D%3D"> FAQ
&gt;
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AADuQHmPWS2S37oYpDGGLKi3T0osQOSvG86LjJ5ZvypB%2FN%2B%2F1nQZDPwrXJJEpsdLXFxrMvwXyRkDJOe7tlg20JkX7qninaEnlAoWUmvXLCBirmxugmbzjNIYs1K9TLyTGOW4qi0zucgUUpG%2Fwv9dzpdd%2Bacj1BDjf54MZJVte7e0OwJRSrPVA13dVdnRim6g4osPTQ8Jo9Sm3rh1SR8gzlqC%2BLtNjZRM87YoHQ5%2BS9H2dNr4B0n7SjaNJNjQsLiM3eXCwSVYI45YRApDMVbuuB82pnUBeGXOCgR4EzE2hoxQKWw52DvMmz7fmEuW7oHX%2FqVmMByFL6fKrYge3SIkTwRbqbRBBXwkLzczzxZ%2BRDPwvMiU8zvqrf3KgZiAi3lTcRLTU1jNZFF%2F%2FhigEx16Jg8Iicqt%2FAzYT159qxqirgMiQto%3D&amp;c=AAA0wv05SFoeMwO3FyDepnB3%2FWBkKaAcOW5qN1%2BuO9HjHJMFBOpShebVns1B%2B4OkLS6LN1wINoec52ue%2BEDNMYmI9ypacD%2BbsE6QTCnNFnG0gFAIbbTAhv4bLzj1BBAQZCfJklji%2FFuRl68hHAyYQjl9C6yy8mCPfToU2av7HWPv41to%2Bs4ClOeoRPgfGjFaQBpFLzSX3SGoZHDBVagc%2Fu33N5wUs2oroAeX5DwHq3diC0lqY38c8N8bWtSUHjIK47tjR%2BrBPZpNrK%2FA%2FNyxFFq9n1hO7miX0H3L%2FyRrQwPjp26qc%2BLSZIkybB0gX0K2L0bKVFzfaqPyNaTNTfsKzIaBBypNJ4QhYtoH%2BWcno%2FIyk7Nvubic2d3%2FeTS69HhMt5uuZKwecQ0uwsAxxH%2FVjMyYnbY2hq%2BuTFcjDJsvgzYH1bLdssmnVyLMzCclEZCu5OOo2OsMpHrGghCS5racpO8KaYD7rxwU3r9V02orMw%3D%3D">
&gt;
&gt; Connect
&gt; [image: Facebook]
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AACBv8NaaOIWGepRvB7VOLVCoezfDtHodud6Uj7Ia1vXLiebTquS6ScwiMGfuEJdhlzQit8w1aHZP1%2Fk3GrNVMWALXINcAwlPWq4bPBHqyyn9NB6ck4yJlZ7xtbAfv0RNqDeoel7CH4hdpEizEEz3WHKSMwR9oPNclRLg28ksIrwhkHCM4EF&amp;c=AABApm6AcD01LZ3ea6%2Bv82yhsch%2BATIHtvXHYjxqUgYzxUHNn%2BtCgo1LdkVSnpKLmbkqzT8chQhclc%2BqNA%2FGxc1DXzHP8JqoW5SJ4rffBkvlbkubHLrbV4IMZ648Jo9WrVOnFrcNfatjtzD6BrkUW7npNIY3%2Bg4GU84q0BYROUJb9JpQnWhm9kJcSailFLPek7w3lNstPk8La5OUWjQ0B6bDBZ8d2XjnUC8L%2F8Z4W4kyeXJCAqRPcXBUV297qMU%2Fnf1%2BoAt12DXk0Ow5%2Bbabkv9Ld8pvM3Oyv6RzE2kUKXoNzhl8tJ8udkLrsLnPdoWDQVg3tYEBLZu%2BYGAr3Mf4ljIFFW2xB4ACGkh7nw4aXsi99jST%2BdKO6OYU4iluQBB%2BWcV606Mrsb7znwTAV8bZ0yse9p5IxFDOPxvGDQxhIWp1GzhrL4o7WK%2FERnwPAgJpySCychnuxSyGaiolnINvO3jff3A47jxepMz7oTyuHw%3D%3D"> [image:
&gt; Instagram]
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AAB%2FWahEUhwkV4KwRtmAHFDwvmqSf0AEfcJzMAdyz3AgmWOrDSJ4%2FfbbAt%2Byhv33ETv2fj0BMsdRwZ7qnEBCmZWztIGJUBzu81jrYQmRtFWdS%2FiAedgLNM5M%2FjcexCYYRoB0Zjq0QLXk5S638KXNVO524LTkX4xuDjBBIDnWdBwZwUF1YglG&amp;c=AAB4tntn%2BZHScCu0czFte0FUCsO2a6n2lsnIhI%2FgAqOB8hEHz4oBXU9K%2BSOU%2FiE0JZt12o%2BAdKXw9WbLJ4jhvm6sbFfUwGjdoBMTWPM%2BsO7iZOqLjN%2FRyx5zLOhXxI46Wj%2BlKgoV4oQ9iC8leQJacg85Ocp3IQHGeGXKwbnEUfimkp77zEXb7ovIfmZaHEOpsUrbzJjOYF9ga5yD6EJP4epRVzGf32n%2Fo0oK5hJER7XF63DwIMBa9OtxOKoM3AJd9fqqDft44kHhevtiVr5GvybKNmff23W58wPzYBV6DCje5lawLCj5dGuDjakVItxjEoafI8KrAF5aOZc3H8OiykgCs7Hw11KfVu43Vlj%2F7cX%2B1%2F9b2UH4Sx14M2%2FroQ3W3rx5UcyRqzoenejR%2FlZCFWSeea3bXWaWjPbMzsTEnP6CeBStnVaM3AJEuj2GMPPMAagyGKOMMA4ZvTmbnwCMj%2F4LNk7s7kYkMtyQDiEfKw%3D%3D">
&gt; buddhastoneshop.com
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AAAEZZlmCjR6DDd%2BswbLYJfRkyWXwvjCRCRxhBXrOXddGUABpqZTqjNNXn1S1uI4wUXdDeQowISRkSko1OeYxAfxXDPaLNdkUDGto663epDIz4Dh%2BReyNZwGIs6F8q7Xm8akR%2FYyOkgHDs8dT4M%2FXOHiCwRECBOCQu5XWgGMWZyu32n03QIFXsMRwjfW%2Bje%2FmZIFkJvid82afPNOUR%2B3RAt%2BCljGHLVV76MpNWF4ArMgnF7e8Ina%2BxGGDbYTCEdSU2BzSLUXRkrbCRs7t0ngnpSwDPgJd%2BlQX0hLnSRobV7MXk3WE3ttmVTQVByrjZmOcaXrNz4Kxs%2B1Hkjvo6bom1sQ3X3jvs7rs5ZR%2BtPf%2FF0mC0t2xGPOct0bw4boOXR8dejLf3iTClKAeolGtEd2E%2B8%3D&amp;c=AAC%2FwN3iJ8ybBMjc0hlE1QJ3Fg7BhOAEAuySmwIIx2N8QQmBl3cV8vps3xnUDgkWb9%2B4ZaiQz8UUSMBpHCn3INsmnbz4jmHkA2MHitLU6y4wyjlXdITBfvFsrFACjTFlrH%2Fu93cs1OlSeRdnFLWdU1Khl7UDS2P1fH%2FOvBzLYrOhTqvwYwGZ4Sdtp35VCvtQ15HKc6%2BNFYs1aO7V8Rk35C2p1SerMK4Qoy6qLMM7dCyURoWV%2FGsKIONEvLRPwd%2BA%2FifF9uQXo90ICMjqsi3ceqFo%2F3LHhxJID%2FmOzm4awlZrxPk3SpE01Y83pq%2FR%2B5XiSLkkwF1jPdHOeJEGmEViLdC44A3tLaT3lGTiSyY%2BYxx1z%2BzKSF0pyWzGCEDKzN%2F8LsTdtCy5LUIyve7dJFe7ZkiD5FbV97DvDEZJXhiTiyusG847Kjx0oQ4agG%2BH4uIHFAJHZTGbWf9gmXgA9lVbJoGBnRoCh1ILyXWRhR2eBw%3D%3D">
&gt; buddhastoneshop
&gt; Nos. 721-725 Nathan Road, Mongkok
&gt; Mongkok, HK
&gt;
&gt; Copyright © 2025
&gt; [image: buddhastoneshop]
&gt; <https: buddhastoneshop.com="" _t="" c="" a1030004-183a83ca86934c92-774ad00e?l="AAAbjpEkD9FEcAMi5whPtfMgb%2B%2BazcVnQkPOJB3A0KJc1pc8PNTKh4qt7hmLxmooghzvyqRAwjj0y9gXrMvpUDiZG%2B26LVoY5RqbalY%2FThGXsr2fXUHkEtvaFk25gLF12FNNNd0xlWNDMBsfLiM3PpwIbIDkMEoU0qYuTDMu2Tg1YA%2Bd6bvlwsrIvrbsZvqMPaEqR3qg2Nb5g1pCIoFM9UghWCRUs3pxGUQJzQy5E%2FPCTCPrLKqS%2FoMu8zj1DngXBnwVSYsqPmQqTNVHpy7x49%2BVEYVbud5ImAHX9VyxgUNpvIvpMIgSGdGlKpp8vESHUctdD%2FjQ6b%2F1P9HWzb%2FOdeKgfVn565L4FXpnVF9CwU9g61UpYghYsp3vsrpfaPgrzT85%2B3GtE52M&amp;c=AACeNzXzzF472KSyFCVoNXicdKkYJ2B38Ir6nz1PcXOaMIy2XW1IV5ZNzHFfmGU93Mi%2FBUkIqQuMCLHsL6kj2Mz10KhVWsg16NwJwL2ZsqDxR2S8G3vENLd0Zje5xRYAG8H5Sr4tnAhBgv3ghIeBBIEOy0ON%2FAEk9qoZUfovhyfraSrUoLcevCSAK%2Fi%2BBM9XpwoxzhV0RWEWCQhVaDFfc8d3wqc%2FaSM1t%2BzJg1gmG4DbfPbgaB9%2Bkm%2F980zQPLpuXjuiKDLHv4K03%2B1DG5j4qXhxutC%2BsOVMgZW7ueCDvE1T65nuGd7c%2BKIGbcU7NvZpkyTzKHc%2FtPAsCz4MRhZJZuIDpeEUF478IOaabpxyCT33MR%2F%2FG90Lp%2FcThAh9ZN4iQBmxHnLdEruo%2BTf1zpTdUPIOAYNay%2F2GoBMKe%2FXgUekG%2BCHU3%2F5MDumj%2BC8doTwD1osb53LMzMYrn7dktKnsWfkC6cO98ABN1U1x4VhTEw%3D%3D">
&gt;

</https:></https:></https:></https:></https:></https:></https:></https:></https:></https:></https:></https:></https:></https:></https:></https:></https:></https:></https:></https:></https:>`
	fmt.Println(len(str))
}
