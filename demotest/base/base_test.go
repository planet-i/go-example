// split/split_test.go

package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) { // 测试函数名必须以Test开头，必须接收一个*testing.T类型参数
	t.Parallel()
	got := Split("a:b:c", ":")         // 程序输出的结果
	want := []string{"a", "b", "c"}    // 期望的结果
	if !reflect.DeepEqual(want, got) { // 因为slice不能比较直接，借助反射包中的方法比较
		t.Errorf("expected:%v, got:%v", want, got) // 测试失败输出错误提示
	}
}

func TestSplitWithComplexSep(t *testing.T) {
	got := Split("abcd", "b")
	want := []string{"a", "cd"}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected:%v, got:%v", want, got)
	}
}

func TestTimeConsuming(t *testing.T) {
	fmt.Println("aaa")
	t.Skip("short模式下会跳过该测试用例")
	fmt.Println("bbb")
}
