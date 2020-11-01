package config

import (
	"esTool/importer"
	"esTool/pkg/elastic"
	"esTool/pkg/logger"
	"github.com/ghodss/yaml"
	"gopkg.in/go-playground/validator.v9"
	//"gopkg.in/yaml.v2"
	"io/ioutil"
)

type EsImportConfig struct {
	ImportConfig importer.ConfigT `validate:"required,dive"`
	FilePath      string          `validate:"required"`
	MaxLineSize   int
	ElasticConfig elastic.ConfigT   `validate:"required"`
	LoggerConfig  logger.ConfigT    `validate:"required"`
}

func (ec *EsImportConfig) check () error{
	vd := validator.New()
	return  vd.Struct(ec)
}

var config *EsImportConfig

func GetConfig()*EsImportConfig{
	return config
}

func InitConfig(filePath string)error{
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil{
		return err
	}
	c := new(EsImportConfig)
	if err = yaml.Unmarshal(yamlFile, c); err != nil{
		return err
	}
	//fmt.Println(c)
	if err = c.check(); err != nil{
		return err
	}
	config = c
	return nil
}