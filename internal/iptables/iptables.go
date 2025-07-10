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
	nfqwsPortsTCP    []string
	nfqwsPortsUDP    []string
}

func NewIptables() *iptables {
	return &iptables{
		connbytesSupport: viper.GetBool("iptables.connbytesSupport"),
		markSupport:      viper.GetBool("iptables.markSupport"),
		multiportSupport: viper.GetBool("iptables.multiportSupport"),
		nfqwsPortsTCP:    viper.GetStringSlice("nfqws.ports.tcp"),
		nfqwsPortsUDP:    viper.GetStringSlice("nfqws.ports.udp"),
	}
}

func (ip *iptables) enableMarkAndConnbytes(iptablesArgs *[]string) {
	if ip.markSupport {
		*iptablesArgs = append(*iptablesArgs, []string{"-m", "mark", "!", "--mark", "0x40000000/0x40000000"}...)
	}
	if ip.connbytesSupport {
		*iptablesArgs = append(*iptablesArgs, []string{"-m", "connbytes", "--connbytes-dir=original", "--connbytes-mode=packets", "--connbytes", "1:6"}...)
	}
}

func (ip *iptables) setupIptablesMultiport() {
	tcpPortsString := strings.Join(ip.nfqwsPortsTCP, ",")
	udpPortsString := strings.Join(ip.nfqwsPortsUDP, ",")

	log.Printf("Настройка iptables с multiport для TCP: %s", tcpPortsString)
	cmdArgsTCP := []string{"-t", "mangle", "-I", "POSTROUTING", "-p", "tcp", "-m", "multiport", "--dports", tcpPortsString, "-j", "NFQUEUE", "--queue-num", "200", "--queue-bypass"}
	ip.enableMarkAndConnbytes(&cmdArgsTCP)

	if err := exec.Command("iptables", cmdArgsTCP...).Run(); err != nil {
		log.Fatalf("Ошибка настройки iptables (multiport TCP); %s", err)
	}

	log.Printf("Настройка iptables с multiport для UDP: %s", udpPortsString)
	cmdArgsUDP := []string{"-t", "mangle", "-I", "POSTROUTING", "-p", "udp", "-m", "multiport", "--dports", udpPortsString, "-j", "NFQUEUE", "--queue-num", "200", "--queue-bypass"}
	ip.enableMarkAndConnbytes(&cmdArgsUDP)

	if err := exec.Command("iptables", cmdArgsUDP...).Run(); err != nil {
		log.Fatalf("Ошибка настройки iptables (multiport UDP): %v", err)
	}
}

func (ip *iptables) setupIptablesTCP() {
	log.Println("Настройка iptables по отдельным портам.")
	for _, port := range ip.nfqwsPortsTCP {
		log.Printf("Добавление правила для TCP порта %s", port)
		cmdArgs := []string{"-t", "mangle", "-I", "POSTROUTING", "-p", "tcp", "--dport", port, "-j", "NFQUEUE", "--queue-num", "200", "--queue-bypass"}
		ip.enableMarkAndConnbytes(&cmdArgs)

		if err := exec.Command("iptables", cmdArgs...).Run(); err != nil {
			log.Printf("Ошибка добавления правила iptables для TCP порта %s; %s", port, err)
		}
	}
}

func (ip *iptables) setupIptablesUDP() {
	for _, port := range ip.nfqwsPortsUDP {
		log.Printf("Добавление правила для UDP порта %s", port)
		cmdArgs := []string{"-t", "mangle", "-I", "POSTROUTING", "-p", "udp", "--dport", port, "-j", "NFQUEUE", "--queue-num", "200", "--queue-bypass"}
		ip.enableMarkAndConnbytes(&cmdArgs)

		if err := exec.Command("iptables", cmdArgs...).Run(); err != nil {
			log.Printf("Ошибка добавления правила iptables для UDP порта %s: %v", port, err)
		}
	}
}

func (ip *iptables) SetupIptables() {
	if ip.multiportSupport {
		ip.setupIptablesMultiport()
	} else {
		ip.setupIptablesTCP()
		ip.setupIptablesUDP()
	}

	log.Println("Установка sysctl net.netfilter.nf_conntrack_tcp_be_liberal=1")
	if err := exec.Command("sysctl", "-w", "net.netfilter.nf_conntrack_tcp_be_liberal=1").Run(); err != nil {
		log.Printf("Предупреждение: Не удалось установить sysctl: %s", err)
	}

	log.Println("Настройка правил iptables завершена.")
}
