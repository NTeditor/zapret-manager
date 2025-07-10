package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nteditor/zapret-manager/cmd"
	"github.com/spf13/viper"
)

func main() {
	cmd.Execute()
}

func init() {
	const (
		CONFIG_PATH = "/data/adb/zapret"
		CONFIG_NAME = "config"
		CONFIG_TYPE = "json"
	)
	var configFile = fmt.Sprintf("%s/%s.%s", CONFIG_PATH, CONFIG_NAME, CONFIG_TYPE)

	viper.SetConfigName(CONFIG_NAME)
	viper.SetConfigType(CONFIG_TYPE)
	viper.AddConfigPath(CONFIG_PATH)

	viper.SetDefault("iptables.multiportSupport", false)
	viper.SetDefault("iptables.markSupport", true)
	viper.SetDefault("iptables.connbytesSupport", false)

	// viper.SetDefault("nfqws.opt")
	viper.SetDefault("nfqws.ports.tcp", []string{"80", "443"})
	viper.SetDefault("nfqws.ports.udp", []string{"443", "50000-50099"})

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := os.MkdirAll(CONFIG_PATH, 0755); err != nil {
				log.Fatalf("Не удалось создать директорию: %s; %s", configFile, err)
			}
			if err := viper.WriteConfigAs(configFile); err != nil {
				log.Fatalf("Не удалось создать файл: %s; %s", configFile, err)
			}
			if err := os.Chown(configFile, 0, 0); err != nil {
				log.Fatalf("Не удалось установить владельца для файла: %s; %s", configFile, err)
			}
			if err := os.Chmod(configFile, 0644); err != nil {
				log.Fatalf("Не удалось установить права для файла: %s; %s", configFile, err)
			}
			log.Print("Файл конфигурации успешно создан")
		} else {
			log.Fatalf("Ошибка чтения файла конфигурации; %s", err)
		}
	}
}
