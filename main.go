package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var dst = flag.String("t", "https://api.openai.com", "targetURL,-t https://api.openai.com")
var bind = flag.String("bind", ":8080", "bind host and port,like 127.0.0.1:8080")

var keyFile = flag.String("keyfile", "", "")
var certFile = flag.String("certfile", "", "")

func main() {
	flag.Parse()
	// 创建目标 URL
	targetURL, err := url.Parse(*dst)
	if err != nil {
		log.Fatal(err)
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 处理请求的处理函数
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handle %s from %s\n", r.RequestURI, r.RemoteAddr)
		// 修改请求的 Host 字段为目标 URL 的主机名
		r.Host = targetURL.Host
		// 设置请求的 URL.Scheme 和 URL.Host 字段为目标 URL 的相应字段
		r.URL.Scheme = targetURL.Scheme
		r.URL.Host = targetURL.Host
		// 转发请求到目标 URL
		proxy.ServeHTTP(w, r)
	})

	// 启动代理服务器
	log.Println("bind", *bind, "proxy to ", (*targetURL).String())
	if *keyFile != "" && *certFile != "" {
		log.Fatal(http.ListenAndServeTLS(*bind, *certFile, *keyFile, nil))
	} else {
		log.Fatal(http.ListenAndServe(*bind, nil))
	}

}
