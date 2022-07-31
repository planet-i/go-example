package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type User struct {
	Name    string
	Website string
	Age     uint
	Male    bool
	Skill   []string
}

func main() {
	user := User{
		"搜索",
		"http://sousuo.com",
		20,
		true,
		[]string{"google", "baidu", "bing", "huohu"},
	}
	u, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Json 编码失败：%v\n", err)
		return
	}
	fmt.Printf("Json 格式数据：%s\n", u)

	var user2 User
	err = json.Unmarshal(u, &user2)
	if err != nil {
		fmt.Printf("JSON 解码失败：%v\n", err)
		return
	}
	fmt.Printf("JSON 解码结果: %#v\n", user2)

	u2 := []byte(`{"name": "搜索", "website": "https://google.com", "note": "请输入搜索词"}`)
	var user3 User
	err = json.Unmarshal(u2, &user3)
	if err != nil {
		fmt.Printf("JSON 解码失败：%v\n", err)
		return
	}
	fmt.Printf("JSON 解码结果: %#v\n", user3)

	u3 := []byte(`{"name": "搜索", "website": "https://souhu.com", "age": 18, "male": true, "skills": ["Google", "baidu"]}`)
	var user4 interface{}
	err = json.Unmarshal(u3, &user4)
	if err != nil {
		fmt.Printf("JSON 解码失败：%v\n", err)
		return
	}
	fmt.Printf("JSON 解码结果: %#v\n", user4)

	user5, ok := user4.(map[string]interface{})
	if ok {
		for k, v := range user5 {
			switch v2 := v.(type) {
			case string:
				fmt.Println(k, "is string", v2)
			case float64:
				fmt.Println(k, "is float64", v2)
			case bool:
				fmt.Println(k, "is bool", v2)
			case []interface{}:
				fmt.Println(k, "is an array:")
				for i, iv := range v2 {
					fmt.Println(i, iv)
				}
			default:
				fmt.Println(k, "is another type not handle yet")
			}
		}
	}
	Iojson()
}

func Iojson() {
	fmt.Println("请输入数据")
	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	for {
		var v map[string]interface{}
		if err := dec.Decode(&v); err != nil {
			log.Println(err)
			return
		}
		if err := enc.Encode(&v); err != nil {
			log.Println(err)
		}
	}
}
