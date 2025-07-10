package cmd

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/nteditor/zapret-manager/internal/iptables"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func startZapret(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()
	if err := iptables.NewIptables().SetupIptables(ctx); err != nil {
		log.Fatalf("%v", err)
	}

	nfqwsArgs := append([]string{"/system/bin/nfqws", "--debug=android", "--daemon", "--qnum=200", "--uid=0:0"}, viper.GetStringSlice("nfqws.opt")...)
	output, err := exec.Command("/system/bin/nfqws", nfqwsArgs...).Output()
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			fmt.Printf("%s\n ExitCode: %d", err.Stderr, err.ExitCode())
		}
		log.Fatalf("%v", err)
	} else {
		fmt.Println(string(output))
	}
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Plaseholder Short",
	Long:  "Plaseholder Long",
	Run:   startZapret,
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Plaseholder Short",
	Long:  "Plaseholder Long",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := exec.Command("killall", "nfqws").Output()
		if err != nil {
			if err, ok := err.(*exec.ExitError); ok {
				fmt.Printf("%s\n ExitCode: %d", err.Stderr, err.ExitCode())
			}
			log.Fatalf("%v", err)
		} else {
			fmt.Println(string(output))
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
		defer cancel()
		if err := iptables.NewIptables().CleanIptables(ctx); err != nil {
			log.Fatalf("%v", err)
		}
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
			startZapret(cmd, args)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd, stopCmd, enableCmd, disableCmd, restartCmd, statusCmd, versionCmd, autostartCmd)
}
