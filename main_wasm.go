package main

import (
	_ "embed"
	"flag"
	"fmt"
	"syscall/js"

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

// This exports an add function.
// It takes in two 32-bit integer values
// And returns a 32-bit integer value.
// To make this function callable from JavaScript,
// we need to add the: "export add" comment above the function
//export configClash
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

	outb := configByte
	outb, err = convert.Patch(outb, s, include, exclude, nil)
	if err != nil {
		panic(err)
	}

	return string(outb)
}

func main() {
    fmt.Println("Hello, WebAssembly!")
    // Mount the function on the JavaScript global object.
    js.Global().Set("configClash", js.FuncOf(func(this js.Value, args []js.Value) any {
        if len(args) != 1 {
            //fmt.Println("invalid number of args")
            return "Invalid no of arguments passed"
        }
        input := args[0].String()
        fmt.Printf("input %s\n", input)
        return configClash(input)
    }))

    // Prevent the program from exiting.
    // Note: the exported func should be released if you don't need it any more,
    // and let the program exit after then. To simplify this demo, this is
    // omitted. See https://pkg.go.dev/syscall/js#Func.Release for more information.
    select {}
}
