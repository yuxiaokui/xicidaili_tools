package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	proxy "proxy_fuck/proxy"
	"strings"
	"flag"
	"sync"
	"time"
	"strconv"
)

var wg sync.WaitGroup

func check(proxy string) {
	proxyCheckUrl := "http://2000019.ip138.com"
	s := strings.Split(proxy,":")
	ip := strings.Replace(s[1],"//","",-1)
	port := strings.Replace(s[2],"/","",-1)
	if PortScan(ip,port) {
		cli := NewHttpClient(proxy)
    data,_ := HttpGET(cli, proxyCheckUrl)
    if strings.ContainsAny(string(data), ip) {
			fmt.Println(ip + ":" + port)
		}
	}
}

func PortScan(ip string,port string) bool {
    _, err  := net.DialTimeout("tcp", ip + ":" + port, time.Second*1)
    if err != nil{
        return false
    }
    return true
}

func NewHttpClient(proxyAddr string) *http.Client {
	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		return nil
	}

	netTransport := &http.Transport{
		//Proxy: http.ProxyFromEnvironment,
		Proxy: http.ProxyURL(proxy),
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(10))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * time.Duration(3),
	}

	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}

func HttpGET(client *http.Client, url string) (body []byte, err error) {
	rsp, err := client.Get(url)
	if err != nil {
		return
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK || err != nil{
		err = fmt.Errorf("HTTP GET Code=%v, URI=%v, err=%v", rsp.StatusCode, url, err)
		return
	}

	return ioutil.ReadAll(rsp.Body)
}

func work(proxys chan string){
	for len(proxys) >0 {
		proxy := <- proxys
		check(proxy)
	}
	wg.Done()
}

func help() {
	fmt.Println("[+] xi4okv proxy tools")
}

func main() {
	help()
	pageNumber := flag.Int("num",1,"num")
	flag.Parse()
	var proxyList []string
	for i := 1 ; i <= *pageNumber; i++ {
		url := "https://www.xicidaili.com/wt/" + strconv.Itoa(i)
		html := proxy.GetHtml(url)
		for _,proxy := range proxy.GetProxy(html) {
			proxyList = append(proxyList,proxy)
		}
	}
	proxys := make(chan string,len(proxyList))
	//把proxy列表放入channel
	for _, proxy := range proxyList {
		proxys <- proxy
	}
	// testproxy
	fmt.Println("Test Start")
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go work(proxys)
	}
	wg.Wait()
	fmt.Println("End!")


}
