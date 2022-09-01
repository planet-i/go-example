package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"
	"testing"
)

func TestErrorWrap(t *testing.T) {
	type args struct {
		h handler
	}
	tests := []struct {
		name string
		args args
		want func(c *gin.Context)
	}{
		{
			name: "normal",
			args: struct{ h handler }{
				h: func(c *gin.Context) error {
					err := errors.New("this is a normal error")
					return err
				},
			},
			want: func(c *gin.Context) {
				return
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrorWrap(tt.args.h); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("ErrorWrap() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
		})
	}
}
