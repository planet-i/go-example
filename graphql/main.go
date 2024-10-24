package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
)

func main() {
	TestGrahqlJson()
	TestGrahqlJson2()

	type SEO struct {
		Title string `json:"title"`
	}

	type Product struct {
		ID     string `json:"id"`
		Handle string `json:"handle"`
		SEO    SEO    `json:"seo"`
	}

	product := Product{
		ID:     "gid://shopify/Product/8552026079520",
		Handle: "abstract-pattern-cotton-tassel-sofa-protector",
		SEO:    SEO{Title: "Abstract Pattern Cotton Tassel Sofa Protector"},
	}

	// 序列化为 JSON
	jsonOutput, err := serializeToJSON(product)
	if err != nil {
		log.Fatalf("Error serializing to JSON: %v", err)
	}
	fmt.Println("JSON Output:")
	fmt.Println(jsonOutput)

	// 转换为 GraphQL Input 语法
	graphqlInput, err := toGraphQLInput(product)
	if err != nil {
		log.Fatalf("Error converting to GraphQL Input: %v", err)
	}
	fmt.Println("\nGraphQL Input:")
	fmt.Println(graphqlInput)
}

// Input struct
type Input struct {
	Id     string `json:"id"`
	Handle string `json:"handle,omitempty"`
	Seo    struct {
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
	} `json:"seo"`
}

func TestGrahqlJson() {
	params := Input{
		Id:     "gid://shopify/Product/8552026079520",
		Handle: "abstract-pattern-cotton-tassel-sofa-protector",
		Seo: struct {
			Title       string `json:"title,omitempty"`
			Description string `json:"description,omitempty"`
		}{
			Title:       "Abstract Pattern Cotton Tassel Sofa Protector",
			Description: "Abstract Pattern Cotton Tassel Sofa Protector",
		},
	}

	paramsJs, _ := MyJsGraphsql(params)
	fmt.Println(paramsJs)
}

func MyJsGraphsql(v interface{}) (string, error) {
	// 输入的 v 转换为JSON格式的字节切片
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	bjs := string(b)
	// var paramsT reflect.Type
	paramsT := reflect.TypeOf(v)
	// 通过 reflect.TypeOf 获取输入参数 v 的类型信息，并遍历所有字段
	for i := 0; i < paramsT.NumField(); i++ {
		// 对于每个字段，检查其 json 标签，并根据标签调整JSON字符串的格式。将字段名的""给去掉，比如 "id":"gid" -> id: "gid"
		f := paramsT.Field(i)
		// 判断f是不是对象类型，如果是的，就继续解析
		tag := strings.Split(f.Tag.Get("json"), ",")
		if len(tag) > 0 && tag[0] != "-" {
			bjs = strings.ReplaceAll(bjs, `"`+tag[0]+`":`, tag[0]+": ")
		}
		fieldType := f.Type
		fmt.Println("fieldType.Kind()", fieldType.Kind())
		// 如果字段是结构体类型，递归调用 MyJsGraphsql
		if fieldType.Kind() == reflect.Struct || fieldType.Kind() == reflect.Map {
			fieldValue := reflect.ValueOf(v).Field(i).Interface()
			nestedResult, err := MyJsGraphsql(fieldValue)
			if err != nil {
				return "", fmt.Errorf("处理嵌套字段时出错: %w", err)
			}
			bjs = strings.ReplaceAll(bjs, `"{"`+tag[0]+": ", nestedResult)
		}
	}

	return bjs, nil
}

// {"id":"gid://shopify/Product/8552026079520","handle":"abstract-pattern-cotton-tassel-sofa-protector","seo":{"title":"Abstract Pattern Cotton Tassel Sofa Protector"}}
// {id: "gid://shopify/Product/8552026079520",handle: "abstract-pattern-cotton-tassel-sofa-protector",seo: {"title":"Abstract Pattern Cotton Tassel Sofa Protector"}}
// 这样子处理，好像字段的字段没有得到处理呀
func TestGrahqlJson2() {
	params := Input{
		Id:     "gid://shopify/Product/8552026079520",
		Handle: "abstract-pattern-cotton-tassel-sofa-protector",
		Seo: struct {
			Title       string `json:"title,omitempty"`
			Description string `json:"description,omitempty"`
		}{
			Title: "Abstract Pattern Cotton Tassel Sofa Protector",
		},
	}

	paramsJs, _ := MyJsGraphsql2(params)
	fmt.Println(paramsJs)
}
func MyJsGraphsql2(v interface{}) (string, error) {
	// Convert input to JSON format
	// b, err := json.Marshal(v)
	// if err != nil {
	// 	return "", err
	// }
	// bjs := string(b)

	// Process the fields recursively
	result, err := processFields(reflect.ValueOf(v))
	if err != nil {
		return "", err
	}

	return result, nil
}

func processFields(value reflect.Value) (string, error) {
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return "", fmt.Errorf("expected a struct, got %s", value.Kind())
	}

	// Start building the output string
	var sb strings.Builder
	sb.WriteString("{")

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := value.Type().Field(i)

		// Get the JSON tag
		tag := fieldType.Tag.Get("json")
		tagParts := strings.Split(tag, ",")
		tagName := tagParts[0]

		// Handle omitted fields
		if len(tagParts) > 1 && tagParts[1] == "omitempty" && isEmptyValue(field) {
			continue
		}

		// Write the field name and value
		if sb.Len() > 1 {
			sb.WriteString(",")
		}
		sb.WriteString(tagName + ": ")

		// Handle nested structs
		if field.Kind() == reflect.Struct {
			nestedResult, err := processFields(field)
			if err != nil {
				return "", err
			}
			sb.WriteString(nestedResult)
		} else {
			// Convert field value to JSON string
			jsonValue, err := json.Marshal(field.Interface())
			if err != nil {
				return "", err
			}
			sb.WriteString(string(jsonValue))
		}
	}

	sb.WriteString("}")
	return sb.String(), nil
}

func isEmptyValue(v reflect.Value) bool {
	return v.Interface() == "" || v.IsZero()
}

// 将任意结构体序列化为 JSON
func serializeToJSON(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// 将任意结构体转换为 GraphQL Input 语法
func toGraphQLInput(data interface{}) (string, error) {
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	if val.Kind() != reflect.Struct {
		return "", fmt.Errorf("expected a struct, got %s", val.Kind())
	}

	input := "{"
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// 只处理导出字段
		if field.PkgPath == "" {
			input += fmt.Sprintf(`%s: "%v", `, field.Name, value.Interface())
		}
	}
	input = input[:len(input)-2] + "}" // 去掉最后的逗号和空格

	return input, nil
}
