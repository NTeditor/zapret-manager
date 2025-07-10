package cmd

import (
	"fmt"
	"log"

	"github.com/nteditor/zapret-manager/internal/nfqws"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var nf = nfqws.NewNfqws()

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Plaseholder Short",
	Long:  "Plaseholder Long",
	Run: func(cmd *cobra.Command, args []string) {
		nf.Start()
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Plaseholder Short",
	Long:  "Plaseholder Long",
	Run: func(cmd *cobra.Command, args []string) {
		nf.Stop()
	},
}

var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Plaseholder Short",
	Long:  "Plaseholder Long",
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("magisk.autostart", true)
		viper.WriteConfig()
	},
}

var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Plaseholder Short",
	Long:  "Plaseholder Long",
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("magisk.autostart", false)
		viper.WriteConfig()
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Plaseholder Short",
	Long:  "Plaseholder Long",
	Run: func(cmd *cobra.Command, args []string) {
		if status, err := nf.Status(); err != nil {
			log.Fatalf("Не удалось проверить состояние nfqws; %v", err)
		} else {
			if status {
				nf.Stop()
			}
		}
		nf.Start()
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Plaseholder Short",
	Long:  "Plaseholder Long",
	Run: func(cmd *cobra.Command, args []string) {
		if status, err := nf.Status(); err != nil {
			log.Fatalf("Не удалось проверить состояние nfqws; %v", err)
		} else {
			if status {
				fmt.Println("Zapret работает.")
			} else {
				fmt.Println("Zapret не работает.")
			}
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Показать версию",
	Long:  `Отображает версию программы и используемых библиотек`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`--- Zapret Manager ---
"Версия: v0.1
Лицензия: GPL v3.0

Используемые библиотеки:
  • cobra: v1.9.1 (Лицензия: Apache 2.0)
  • viper: v1.20.1 (Лицензия: MIT)
----------------------`)
	},
}

var autostartCmd = &cobra.Command{
	Use:    "autostart",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("magisk.autostart") {
			nf.Start()
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd, stopCmd, enableCmd, disableCmd, restartCmd, statusCmd, versionCmd, autostartCmd)
}
