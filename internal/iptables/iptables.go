package iptables

import (
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

type iptables struct {
	connbytesSupport bool
	markSupport      bool
	multiportSupport bool
	nfqwsPortsTcp    []string
	nfqwsPortsUdp    []string
}

func NewIptables() *iptables {
	return &iptables{
		connbytesSupport: viper.GetBool("iptables.connbytesSupport"),
		markSupport:      viper.GetBool("iptables.markSupport"),
		multiportSupport: viper.GetBool("iptables.multiportSupport"),
		nfqwsPortsTcp:    viper.GetStringSlice("nfqws.ports.tcp"),
		nfqwsPortsUdp:    viper.GetStringSlice("nfqws.ports.udp"),
	}
}

func (ip *iptables) SetupIptables() error {
	if ip.multiportSupport {
		if err := ip.setupIptablesMultiport(); err != nil {
			return err
		}
	} else {
		if err := ip.setupIptablesTcp(); err != nil {
			return err
		}
		if err := ip.setupIptablesUdp(); err != nil {
			return err
		}
	}

	if err := exec.Command("sysctl", "-w", "net.netfilter.nf_conntrack_tcp_be_liberal=1").Run(); err != nil {
		return err
	}

	return nil
}

func (ip *iptables) CleanIptables() error {
	if err := exec.Command("iptables", "-t", "mangle", "-F").Run(); err != nil {
		return err
	}
	return nil
}

func (ip *iptables) enableMarkAndConnbytes(iptablesArgs *[]string) {
	if ip.markSupport {
		*iptablesArgs = append(*iptablesArgs, []string{"-m", "mark", "!", "--mark", "0x40000000/0x40000000"}...)
	}
	if ip.connbytesSupport {
		*iptablesArgs = append(*iptablesArgs, []string{"-m", "connbytes", "--connbytes-dir=original",
			"--connbytes-mode=packets", "--connbytes", "1:6"}...)
	}
}

func (ip *iptables) setupIptablesMultiport() error {
	tcpPortsString := strings.Join(ip.nfqwsPortsTcp, ",")
	udpPortsString := strings.Join(ip.nfqwsPortsUdp, ",")

	cmdArgsTCP := []string{"-t", "mangle", "-I", "POSTROUTING", "-p", "tcp", "-m", "multiport",
		"--dports", tcpPortsString, "-j", "NFQUEUE", "--queue-num", "200", "--queue-bypass"}
	ip.enableMarkAndConnbytes(&cmdArgsTCP)
	if err := exec.Command("iptables", cmdArgsTCP...).Run(); err != nil {
		return &ErrSetupPortRule{
			Protocol:  "tcp",
			Multiport: true,
			Connbytes: ip.connbytesSupport,
			Mark:      ip.markSupport,
			Err:       err,
		}
	}

	cmdArgsUDP := []string{"-t", "mangle", "-I", "POSTROUTING", "-p", "udp", "-m", "multiport",
		"--dports", udpPortsString, "-j", "NFQUEUE", "--queue-num", "200", "--queue-bypass"}
	ip.enableMarkAndConnbytes(&cmdArgsUDP)
	if err := exec.Command("iptables", cmdArgsUDP...).Run(); err != nil {
		return &ErrSetupPortRule{
			Protocol:  "udp",
			Multiport: true,
			Connbytes: ip.connbytesSupport,
			Mark:      ip.markSupport,
			Err:       err,
		}
	}

	return nil
}

func (ip *iptables) setupIptablesTcp() error {
	for _, port := range ip.nfqwsPortsTcp {
		log.Printf("Добавление правила для TCP порта %s", port)
		cmdArgs := []string{"-t", "mangle", "-I", "POSTROUTING", "-p", "tcp", "--dport", port,
			"-j", "NFQUEUE", "--queue-num", "200", "--queue-bypass"}
		ip.enableMarkAndConnbytes(&cmdArgs)
		if err := exec.Command("iptables", cmdArgs...).Run(); err != nil {
			return &ErrSetupPortRule{
				Protocol:  "tcp",
				Multiport: false,
				Connbytes: ip.connbytesSupport,
				Mark:      ip.markSupport,
				Err:       err,
			}
		}
	}

	return nil
}

func (ip *iptables) setupIptablesUdp() error {
	for _, port := range ip.nfqwsPortsUdp {
		cmdArgs := []string{"-t", "mangle", "-I", "POSTROUTING", "-p", "udp", "--dport", port,
			"-j", "NFQUEUE", "--queue-num", "200", "--queue-bypass"}
		ip.enableMarkAndConnbytes(&cmdArgs)
		if err := exec.Command("iptables", cmdArgs...).Run(); err != nil {
			return &ErrSetupPortRule{
				Protocol:  "udp",
				Multiport: false,
				Connbytes: ip.connbytesSupport,
				Mark:      ip.markSupport,
				Err:       err,
			}
		}
	}

	return nil
}
