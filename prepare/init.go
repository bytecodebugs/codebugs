package prepare

import (
	"codebugs/config"
	"codebugs/log"
	"github.com/spf13/pflag"
)

var configPath = pflag.StringP("config", "c", "data/config.temp.yaml", "config file path.")

func Init() error {
	pflag.Parse()
	// init config
	if err := config.Init(*configPath); err != nil {
		return err
	}
	if err := log.Init(); err != nil {
		return err
	}
	watchSignal()
	singleton()
	return nil
}
