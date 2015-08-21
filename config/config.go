package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type config struct {
	Server struct {
		Host   string `yaml:"host"`
		Port   int    `yaml:"port"`
		Secret string `yaml:"secret"`
	} `yaml:"server"`

	EmailActivation bool   `yaml:"email_activation"`
	UploadDir       string `yaml:"upload_dir"`
	AssetDir        string `yaml:"asset_dir"`
}

const (
	configDir   = "config"
	Development = "development"
	Production  = "production"
	Test        = "test"
	GoEnv       = "GO_ENV"
	GoCwd       = "GO_CWD"
)

var (
	extname = []string{"yml", "yaml"}
	BaseDir string
	Config  config
	Env     string
)

func init() {
	if Env = os.Getenv(GoEnv); Env == "" {
		Env = Development
	}

	switch Env {
	case Production:
		gin.SetMode(gin.ReleaseMode)
	case Test:
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	BaseDir = os.Getenv(GoCwd)

	if BaseDir == "" {
		BaseDir, _ = os.Getwd()
	}

	loaded := false

	for _, ext := range extname {
		path := path.Join(BaseDir, configDir, Env+"."+ext)

		if !exists(path) {
			continue
		}

		data, err := ioutil.ReadFile(path)

		if err != nil {
			panic(err)
		}

		if err = yaml.Unmarshal(data, &Config); err != nil {
			panic(err)
		}

		loaded = true
		break
	}

	if !loaded {
		panic(errors.New("Can't found config file for " + Env))
	}
}

func exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
