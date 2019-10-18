package proxy

import (
	//"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)


func GetHtml(url string)  string {
	client := http.Client{
		Timeout: time.Duration(3) * time.Second,
	}


	request, err := http.NewRequest("GET", url, nil)
        request.Header.Add("User-Agent",  "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36")
	if err != nil {
		return "Connect Error"
	}
	response, err := client.Do(request)
	if err != nil {
		return "Request Error"
	}
	defer response.Body.Close()

        body,_:=ioutil.ReadAll(response.Body)
        return  string(body)
}


func GetProxy(html string) []string{
	re := regexp.MustCompile(`<td>(.*)+</td>`)
	matched := re.FindAllStringSubmatch(html, -1)
	var result []string
	i := 0
	r := ""
	for _, match := range matched {

		i += 1
		if i % 5 == 1 {
			ip := strings.Replace(match[0],"<td>","",-1)
			ip = strings.Replace(ip,"</td>","",-1)
			//fmt.Println("ip:" + ip)
			r += "http://" + ip

		}
		if i% 5 == 2 {
			port := strings.Replace(match[0],"<td>","",-1)
			port = strings.Replace(port,"</td>","",-1)
			//fmt.Println("port:" + port)
			r += ":" + port + "/"
			result = append(result, r)
			r = ""
		}

	}


	return(result)

}


