package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/xmdhs/clash2singbox/convert"
	"github.com/xmdhs/clash2singbox/model/clash"
	"gopkg.in/yaml.v3"
)

var (
	url      string
	path     string
	outPath  string
	include  string
	exclude  string
	insecure bool
	ignore   bool
)

//go:embed config.json.template
var configByte []byte

func init() {
	flag.StringVar(&url, "url", "", "订阅地址，多个链接使用 | 分割")
	flag.StringVar(&path, "i", "", "本地 clash 文件")
	flag.StringVar(&outPath, "o", "config.json", "输出文件")
	flag.StringVar(&include, "include", "", "urltest 选择的节点")
	flag.StringVar(&exclude, "exclude", "", "urltest 排除的节点")
	flag.BoolVar(&insecure, "insecure", false, "所有节点不验证证书")
	flag.BoolVar(&ignore, "ignore", true, "忽略无法转换的节点")
	flag.Parse()
}

func configClash(config string) string {
	c := clash.Clash{}
	var err error
	var b []byte
	b = []byte(config)

	err = yaml.Unmarshal(b, &c)
	if err != nil {
		panic(err)
	}

	s, err := convert.Clash2sing(c)
	if err != nil {
		fmt.Println(err)
	}

	outb, err := os.ReadFile(outPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			outb = configByte
		} else {
			panic(err)
		}
	}

	outb, err = convert.Patch(outb, s, include, exclude, nil)
	if err != nil {
		panic(err)
	}

	return string(outb)
}

func main(){
  var myConfig string = `port: 7890
socks-port: 7891
redir-port: 7892
mixed-port: 7893
tproxy-port: 7895
ipv6: false
mode: rule
log-level: silent
allow-lan: true
external-controller: 0.0.0.0:9090
secret: ""
bind-address: "*"
unified-delay: true
profile:
  store-selected: true
dns:
  enable: true
  ipv6: false
  enhanced-mode: redir-host
  listen: 0.0.0.0:7874
  nameserver:
    - 8.8.8.8
    - 1.0.0.1
    - https://dns.google/dns-query
  fallback:
    - 1.1.1.1
    - 8.8.4.4
    - https://cloudflare-dns.com/dns-query
    - 112.215.203.254
  default-nameserver:
    - 8.8.8.8
    - 1.1.1.1
    - 112.215.203.254
proxies:
  - name: RELAY-104.17.253.39-0991
    server: cfcdn2.sanfencdn.net
    port: 2052
    type: vmess
    uuid: 78d6fa88-0d27-414f-8f51-1e16186d588f
    alterId: 0
    cipher: auto
    tls: false
    skip-cert-verify: true
    servername: us8.sanfencdn2.com
    network: ws
    ws-opts:
      path: /
      headers:
        Host: us8.sanfencdn2.com
    udp: true
proxy-groups:
  - name: FASTSSH-SSHKIT-HOWDY
    type: select
    proxies:
      - RELAY-104.17.253.39-0991
      - DIRECT
rules:
  - MATCH,FASTSSH-SSHKIT-HOWDY`
  fmt.Println(configClash(myConfig))
}
