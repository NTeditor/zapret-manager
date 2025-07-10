package nfqws

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/nteditor/zapret-manager/internal/iptables"
	"github.com/spf13/viper"
)

type nfqws struct {
	opt         []string
	markSupport bool
}

func NewNfqws() *nfqws {
	return &nfqws{
		opt:         viper.GetStringSlice("nfqws.opt"),
		markSupport: viper.GetBool("iptables.markSupport"),
	}
}

func (nf *nfqws) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	iptables := iptables.NewIptables()
	if err := iptables.SetupIptables(ctx); err != nil {
		if err := iptables.CleanIptables(ctx); err != nil {
			log.Printf("%v", err)
		}
		log.Fatalf("%v", err)
	}

	nfqwsArgs := []string{"--debug=android", "--daemon", "--qnum=200", "--uid=0:0"}
	nf.enableMark(&nfqwsArgs)

	output, err := exec.CommandContext(ctx, "/system/bin/nfqws", append(nfqwsArgs, nf.opt...)...).CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			log.Fatalf("%v", err)
		}
	}
}

func (nf *nfqws) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	output, err := exec.CommandContext(ctx, "killall", "nfqws").CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			log.Fatalf("%v", err)
		}
	}

	if err := iptables.NewIptables().CleanIptables(ctx); err != nil {
		log.Fatalf("%v", err)
	}
}

func (nf *nfqws) Status() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	if err := exec.CommandContext(ctx, "pgrep", "nfqws").Run(); err != nil {
		if err, ok := err.(*exec.ExitError); ok && err.ExitCode() != 0 {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

func (nf *nfqws) enableMark(nfqwsArgs *[]string) {
	if nf.markSupport {
		*nfqwsArgs = append(*nfqwsArgs, "--dpi-desync-fwmark=0x40000000")
	}
}
