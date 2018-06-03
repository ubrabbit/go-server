package config

import (
	config "github.com/ubrabbit/go-config"
	. "github.com/ubrabbit/go-server/common"
	"log"
)

type ConfigFile struct {
	Path   string
	Config *config.Config
}

var (
	g_ConfigFile *ConfigFile = nil
)

func GetConfigFile() *ConfigFile {
	return g_ConfigFile
}

func InitConfig(filepath string) *ConfigFile {
	cfg, err := config.ReadDefault(filepath)
	if err != nil {
		log.Fatalf("Fail To Load Cfg: ", filepath, err)
	}
	obj := new(ConfigFile)
	obj.Path = filepath
	obj.Config = cfg
	g_ConfigFile = obj

	InitRunPath()
	return obj
}

func (self *ConfigFile) ReadConfig(section string) map[string]string {
	data := make(map[string]string, 0)
	if self.Config.HasSection(section) {
		section_list, err := self.Config.SectionOptions(section)
		if err != nil {
			log.Fatalf("Fail To Load Sections: ", section)
		}
		for _, v := range section_list {
			options, err := self.Config.String(section, v)
			if err == nil {
				data[v] = options
			}
		}
	}
	return data
}

func (self *ConfigFile) HasOption(section string, option string) bool {
	return self.Config.HasOption(section, option)
}

func (self *ConfigFile) ReadConfigOption(section string, option string) string {
	if !self.Config.HasOption(section, option) {
		return ""
	}
	value, err := self.Config.String(section, option)
	if err != nil {
		log.Fatalf("Fail To Load Options: ", section, option)
	}
	return value
}

func getSettingValue(setting map[string]string, key string, fatal int) string {
	value, ok := setting[key]
	if !ok && fatal > 0 {
		log.Fatalf("setting has no key %s", key)
	}
	return StripString(value)
}
