package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"gopkg.in/cas.v2"
)

var V *viper.Viper

func initialConfig() {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	V = v
}

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(url), nil
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		username := cas.Username(r)
		// attributes := cas.Attributes(r)
		r.Header.Add(V.GetString("grafana.username_key"), username)
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	flag.Parse()
	defer glog.Flush()

	initialConfig()
	casurl, _ := url.Parse(V.GetString("cas.url"))
	client := cas.NewClient(&cas.Options{URL: casurl})

	root := chi.NewRouter()
	root.Use(client.Handler)

	proxy, err := NewProxy(V.GetString("grafana.url"))
	if err != nil {
		panic(err)
	}
	root.HandleFunc("/*", ProxyRequestHandler(proxy))

	port := V.GetInt("server.port")
	glog.Infof("start server on 0.0.0.0:%d\n", port)
	server := &http.Server{
		Addr:    ":" + fmt.Sprintf("%d", port),
		Handler: client.Handle(root),
	}
	if err := server.ListenAndServe(); err != nil {
		glog.Fatal(err)
	}
}
