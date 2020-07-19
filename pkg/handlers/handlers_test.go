package handlers

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_sendCode(t *testing.T) {
	type args struct {
		c AuthenticationContext
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sendCode(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("sendCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sendCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_verifyCode(t *testing.T) {
	type args struct {
		c AuthenticationContext
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := verifyCode(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("verifyCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthenticationRoute_Name(t *testing.T) {
	type fields struct {
		name   string
		title  string
		config map[string]interface{}
		route  func(*gin.RouterGroup)
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthenticationRoute{
				name:   tt.fields.name,
				title:  tt.fields.title,
				config: tt.fields.config,
				route:  tt.fields.route,
			}
			if got := a.Name(); got != tt.want {
				t.Errorf("AuthenticationRoute.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthenticationRoute_Title(t *testing.T) {
	type fields struct {
		name   string
		title  string
		config map[string]interface{}
		route  func(*gin.RouterGroup)
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthenticationRoute{
				name:   tt.fields.name,
				title:  tt.fields.title,
				config: tt.fields.config,
				route:  tt.fields.route,
			}
			if got := a.Title(); got != tt.want {
				t.Errorf("AuthenticationRoute.Title() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthenticationRoute_Config(t *testing.T) {
	type fields struct {
		name   string
		title  string
		config map[string]interface{}
		route  func(*gin.RouterGroup)
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthenticationRoute{
				name:   tt.fields.name,
				title:  tt.fields.title,
				config: tt.fields.config,
				route:  tt.fields.route,
			}
			if got := a.Config(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthenticationRoute.Config() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthenticationRoute_Router(t *testing.T) {
	type fields struct {
		name   string
		title  string
		config map[string]interface{}
		route  func(*gin.RouterGroup)
	}
	tests := []struct {
		name   string
		fields fields
		want   func(*gin.RouterGroup)
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthenticationRoute{
				name:   tt.fields.name,
				title:  tt.fields.title,
				config: tt.fields.config,
				route:  tt.fields.route,
			}
			if got := a.Router(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthenticationRoute.Router() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	type args struct {
		ar HandlerRegistrar
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Register(tt.args.ar)
		})
	}
}
