package configuration

import "github.com/spf13/viper"

var v *viper.Viper

func init() {
	v = initConfiguration()
}

func initConfiguration() *viper.Viper {
	vNew := viper.New()
	vNew.AddConfigPath("$HOME/.db-backup")

	vNew.SetConfigType("yaml")
	vNew.SetConfigName("config")

	return vNew
}