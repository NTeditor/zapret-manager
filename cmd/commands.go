package cmd

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Plaseholder Short",
	Long:  "Plaseholder Long",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Plaseholder Short",
	Long:  "Plaseholder Long",
	Run: func(cmd *cobra.Command, args []string) {

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
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Plaseholder Short",
	Long:  "Plaseholder Long",
	Run: func(cmd *cobra.Command, args []string) {
		err := exec.Command("pgrep", "nfqws").Run()
		if err, ok := err.(*exec.ExitError); ok {
			if err.ExitCode() == 1 {
				fmt.Println("Zapret не работает.")
				return
			}
			log.Fatal(err)
		} else {
			fmt.Println("Zapret работает.")
			return
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long:  `Version of the program and libraries`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`Zapret Magisk  v0.1     [GPL v3.0 License]

Libraries:
  cobra        v1.9.1   [Apache 2.0 License]
  viper        v1.20.1  [MIT License]`)
	},
}

var autostartCmd = &cobra.Command{
	Use:    "autostart",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("magisk.autostart") {

		} else {

		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd, stopCmd, enableCmd, disableCmd, restartCmd, statusCmd, versionCmd, autostartCmd)
}
