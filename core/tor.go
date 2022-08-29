package core

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/net/proxy"
	"log"
	"net"
	"strings"

	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func InitilizeTor(enableTor bool) bool {

	// The site that we want to scrape for data
	//var site string

	// Do we want to proxy through tor
	tor := true

	// Debugging
	debug := true

	// set the tor bits
	proxyAddress := "127.0.0.1"
	proxyPort := 9050

	if debug {

		//fmt.Println("site: ", site)
		fmt.Println("tor: ", tor)
		fmt.Println("debugging: ", debug)
		fmt.Println("Tor Address: ", proxyAddress)
		fmt.Println("Tor Port: ", proxyPort)
	}

	if tor {

		// Check to see if we running requests through tor
		torStatus, ipAddr := checkTorConnection(proxyAddress, proxyPort)

		//siteStatus := scrapeSite(site, proxyAddress, proxyPort)

		fmt.Println("Connection Status: ", torStatus)
		fmt.Println("IP Address: ", ipAddr)

		//fmt.Println(siteStatus)
		return true
	} else {
		return false
	}
}

// Check the proxy connection details
func checkTorConnection(address string, port int) (string, string) {

	torAddress := "https://check.torproject.org/"

	// create a socks5 dialer
	dialer := CreateDialer(address, port)

	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial
	// create a request
	req, err := http.NewRequest("GET", torAddress, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't create request:", err)
		os.Exit(2)
	}
	// use the http client to fetch the page
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't GET page:", err)
		os.Exit(3)
	}
	defer resp.Body.Close()

	// return statement for tor connection
	var stmt string

	// scrape the tor page to check if the connection is being sent over the proxy
	secure, ipAddr := checkTorResponse(resp)
	if secure {
		stmt = "Secure"
	} else {
		stmt = "Not Secure"
	}

	return stmt, ipAddr

}

func CreateDialer(ip string, port int) proxy.Dialer {
	address := ip + ":" + strconv.Itoa(port)

	// create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", address, nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}

	return dialer
}

// Parse the tor project site to ensure that the proxy is working. This will return a bool and the ip address
func checkTorResponse(resp *http.Response) (bool, string) {

	var secure = false
	var address string

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".content").Each(func(i int, s *goquery.Selection) {
		ans := s.Text()
		if strings.Contains(ans, "Congratulations.") {
			secure = true
		}

		ipAddr := net.ParseIP(s.Find("strong").Text())
		if ipAddr != nil {
			address = ipAddr.String()
		} else {
			address = ""
		}
	})

	return secure, address
}

// scrapeSite is a general purpose function to pull a site
//func scrapeSite(site string, address string, port int) string {
//	// TODO need to create the dialer once and then pass it around
//	torAddress := site
//
//	// create a socks5 dialer
//	dialer := CreateDialer(address, port)
//
//	// setup a http client
//	httpTransport := &http.Transport{}
//	httpClient := &http.Client{Transport: httpTransport}
//	// set our socks5 as the dialer
//	httpTransport.Dial = dialer.Dial
//	// create a request
//	req, err := http.NewRequest("GET", torAddress, nil)
//	if err != nil {
//		fmt.Fprintln(os.Stderr, "can't create request:", err)
//		os.Exit(2)
//	}
//	// use the http client to fetch the page
//	resp, err := httpClient.Do(req)
//	if err != nil {
//		fmt.Fprintln(os.Stderr, "can't GET page:", err)
//		os.Exit(3)
//	}
//	defer resp.Body.Close()
//
//	// return statment for tor connection
//	var stmt string
//
//	// scrape the site for the required data info
//	pullData(resp)
//
//	return stmt
//
//}

// pullData will find specific necessary site keywords
//func pullData(resp *http.Response) string {
//
//	buf := new(bytes.Buffer)
//	buf.ReadFrom(resp.Body)
//	s := buf.String() // Does a complete copy of the bytes in the buffer.
//
//	foo := findEmail(s)
//	bar := findBitcoin(s)
//
//	fmt.Println(foo)
//	fmt.Println(bar)
//
//	return "all done"
//
//	// // Load the HTML document
//	// doc, err := goquery.NewDocumentFromReader(resp.Body)
//	// if err != nil {
//	// 	log.Fatal(err)
//	// }
//
//	// Find the review items
//	// doc.Find(".content").Each(func(i int, s *goquery.Selection) {
//	// 	ans := s.Text()
//	// 	if strings.Contains(ans, "Congratulations.") {
//	// 		secure = true
//	// 	}
//
//	// 	ipAddr := net.ParseIP(s.Find("strong").Text())
//	// 	if ipAddr != nil {
//	// 		address = ipAddr.String()
//	// 	} else {
//	// 		address = ""
//	// 	}
//	// })
//
//}
