package request
import (
	//"log"
	"fmt"
	"net"
	"net/http"
	"io/ioutil"
)

type Tor struct {
	user_agent string
	referer string
}

func (tor *Tor) NewIP() {
	conn, err := net.Dial("tcp", "127.0.0.1:9051")
	if err != nil {
		fmt.Printf("Connecting to torcr failed: %s", err.Error())
	}
	//authenticate with torrc
	_, err = conn.Write([]byte("authenticate \"god\"\n"))
	if err != nil {
		fmt.Printf("Writing auth to torcr failed: %s", err.Error())
	}
	
	_, err = conn.Write([]byte("signal newnym \"\"\n"))
	if err != nil {
		fmt.Printf("Writing newnym to torcr failed: %s", err.Error())
	}
}

func (tor *Tor) Connect() *http.Client {
	tor.user_agent = GetUserAgent()
	dial := DialSocksProxy(SOCKS5, "127.0.0.1:9050")
	transport := &http.Transport{Dial: dial}
	return &http.Client{Transport: transport}
}

func (tor *Tor) Browse(url string) (response *http.Response, err error) {
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", tor.user_agent)
	request.Header.Set("Referer", tor.referer)
	response, err = tor.Connect().Do(request)
	return
}

func (tor *Tor) Open(url string, referer string) (string, error) {
	tor.referer = referer
	response, err := tor.Browse(url)
	defer response.Body.Close()

	if err == nil {
		body_str, err := ioutil.ReadAll(response.Body)
		if err == nil {
			body := string(body_str)
			return body, nil
		}
	}

	return "", err
}