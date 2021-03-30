package common

import (
	"time"
)

const ConfigReload = time.Minute * 5

type PlatformInfo struct {
	Key         string                 `yaml:"key" required:"true"`
	Description string                 `yaml:"description" default:""`
	Disabled    bool                   `yaml:"disabled" default:"false"`
	AppID       string                 `yaml:"app_id" required:"true"`
	SecretKey   string                 `yaml:"secret_key" default:"" json:"-"`
	PublicKey   string                 `yaml:"public_key" default:"" json:"-"`
	Comment     string                 `yaml:"comment"`
	Args        map[string]interface{} `yaml:"args"`
}

type App struct {
	AppName        string `yaml:"app_name" required:"true"`
	AppDescription string `yaml:"app_description" default:""`
	DatabaseID     string `yaml:"database_id" required:"true"`
	Servers        []struct {
		ID              string `yaml:"id" default:""`
		Description     string `yaml:"description" default:""`
		Scheme          string `yaml:"scheme" default:"http"`
		Address         string `yaml:"address" default:"localhost:8080"`
		BasePath        string `yaml:"base_path" default:""`
		OnetimeTokenUrl string `yaml:"onetime_token_url" default:"" json:"-"`
		RootToken       string `yaml:"root_token" default:"" json:"-"`
	} `yaml:"servers"`
	DateOffset time.Duration  `yaml:"date_offset" default:"4h"`
	Platforms  []PlatformInfo `yaml:"platforms"`
}
type RealtimeServerResource struct {
	Region         string `yaml:"region" required:"true"`
	Specs          string `yaml:"specs" required:"true"`
	Network        string `yaml:"network" required:"true"`
	Protocol       string `yaml:"protocol" required:"true"`
	Host           string `yaml:"host" required:"true"`
	Port           uint16 `yaml:"port" required:"true"`
	Mtu            int    `yaml:"mtu" json:"mtu" required:"true"`
	SndWnd         int    `yaml:"snd_wnd" json:"sndWnd" required:"true"`
	RcvWnd         int    `yaml:"rcv_wnd" json:"rcvWnd" required:"true"`
	NoDelay        int    `yaml:"no_delay" json:"noDelay" required:"true"`
	Interval       int    `yaml:"interval" json:"interval" required:"true"`
	Resend         int    `yaml:"resend" json:"resend" required:"true"`
	Nc             int    `yaml:"nc" json:"nc" required:"true"`
	RootToken      string `yaml:"root_token" required:"true"`
	TokenAPIScheme string `yaml:"token_api_scheme" required:"true"`
	TokenAPIHost   string `yaml:"token_api_host" required:"true"`
	TokenAPIPort   uint16 `yaml:"token_api_port" required:"true"`
	TokenAPIPath   string `yaml:"token_api_path" required:"true"`
}

var Config = struct {
	RootToken               string                   `yaml:"root_token" default:""`
	Apps                    []App                    `yaml:"apps"`
	RealtimeServerResources []RealtimeServerResource `yaml:"realtime_server_resources"`
}{
	Apps:                    make([]App, 0),
	RealtimeServerResources: make([]RealtimeServerResource, 0),
}
