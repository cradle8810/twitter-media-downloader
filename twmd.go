package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	URL "net/url"
	"os"
	"strings"
	"time"

	twitterscraper "github.com/Toeplitz/twitter-scraper-auth-fix"
	"github.com/mmpx12/optionparser"
)

var (
	usr     string
	proxy   string
	vidz    bool
	imgs    bool
	version = "1.14.2_mod"
	scraper *twitterscraper.Scraper
	client  *http.Client
	size    = "orig"
)

func videoSingle(tweet *twitterscraper.Tweet, output string) {
	if tweet == nil {
		return
	}
	if len(tweet.Videos) > 0 {
		for _, i := range tweet.Videos {
			url := strings.Split(i.URL, "?")[0]
			if usr != "" {
				fmt.Println(url)
			} else {
				fmt.Println(url)
			}
		}
	}
}

func photoSingle(tweet *twitterscraper.Tweet, output string) {
	if tweet == nil {
		return
	}
	if len(tweet.Photos) > 0 {
		for _, i := range tweet.Photos {
			var url string
			if !strings.Contains(i.URL, "video_thumb/") {
				if size == "orig" || size == "small" {
					url = i.URL + "?name=" + size
				} else {
					url = i.URL
				}
				if usr != "" {
					fmt.Println(url)
				} else {
					fmt.Println(url)
				}
			}
		}
	}
}

func processCookieString(cookieStr string) []*http.Cookie {
	cookiePairs := strings.Split(cookieStr, "; ")
	cookies := make([]*http.Cookie, 0)
	expiresTime := time.Now().AddDate(1, 0, 0)

	for _, pair := range cookiePairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}

		name := parts[0]
		value := parts[1]
		value = strings.Trim(value, "\"")

		cookie := &http.Cookie{
			Name:     name,
			Value:    value,
			Path:     "/",
			Domain:   ".x.com",
			Expires:  expiresTime,
			HttpOnly: true,
			Secure:   true,
		}

		cookies = append(cookies, cookie)
	}
	return cookies
}

// DON'T LOOP THIS!!!!! DIE WITH RETURN CODE 2 IMMIDIATELY WHEN FAILED!!!!!
func askPass() {
	var auth_token, ct0 string
	fmt.Printf("auth_token cookie: ")
	fmt.Scanln(&auth_token)
	fmt.Printf("ct0 cookie: ")
	fmt.Scanln(&ct0)
	scraper.SetAuthToken(twitterscraper.AuthToken{Token: auth_token, CSRFToken: ct0})
	if !scraper.IsLoggedIn() {
		fmt.Fprintln(os.Stderr,"Bad Cookies.")
		os.Exit(2) // die()
	}
	cookies := scraper.GetCookies()
	js, _ := json.Marshal(cookies)
	f, _ := os.OpenFile("twmd_cookies.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	defer f.Close()
	f.Write(js)
}

func Login(useCookies bool) {
	if useCookies {
		if _, err := os.Stat("twmd_cookies.json"); errors.Is(err, fs.ErrNotExist) {
			fmt.Print("Enter cookies string: ")
			var cookieStr string
			cookieStr, _ = bufio.NewReader(os.Stdin).ReadString('\n')
			cookieStr = strings.TrimSpace(cookieStr)

			cookies := processCookieString(cookieStr)
			scraper.SetCookies(cookies)

			// Save cookies to file
			js, _ := json.MarshalIndent(cookies, "", "  ")
			f, _ := os.OpenFile("twmd_cookies.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
			defer f.Close()
			f.Write(js)
		} else {
			f, _ := os.Open("twmd_cookies.json")
			var cookies []*http.Cookie
			json.NewDecoder(f).Decode(&cookies)
			scraper.SetCookies(cookies)
			fmt.Fprintln(os.Stderr,scraper.IsLoggedIn())
		}
	} else {
		if _, err := os.Stat("twmd_cookies.json"); errors.Is(err, fs.ErrNotExist) {
			askPass()
		} else {
			f, _ := os.Open("twmd_cookies.json")
			var cookies []*http.Cookie
			json.NewDecoder(f).Decode(&cookies)
			scraper.SetCookies(cookies)
		}
	}

	if !scraper.IsLoggedIn() {
		if useCookies {
			fmt.Fprintln(os.Stderr,"Invalid cookies. Please try again.")
			os.Remove("twmd_cookies.json")
		} else {
			os.Exit(2)
		}
	} else {
		fmt.Fprintln(os.Stderr,"Logged in.")
	}
}

func singleTweet(output string, id string) {
	tweet, err := scraper.GetTweet(id)
	if err != nil {
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}
	if tweet == nil {
		fmt.Fprintln(os.Stderr,"Error retrieve tweet")
		return
	}
	videoSingle(tweet, output)
	photoSingle(tweet, output)
}

func main() {
	var (
		single       = ""
		output       = ""
		printversion bool
		login        = true
		useCookies   bool
	)
	op := optionparser.NewOptionParser()
	op.Banner = "twmd: Apiless twitter media downloader\n\nUsage:"
	op.On("-t", "--tweet TWEET_ID", "Single tweet to download", &single)
	op.On("-V", "--version", "Print version and exit", &printversion)
	op.Parse()

	if printversion {
		fmt.Fprintln(os.Stderr,"version:", version)
		os.Exit(1)
	}
	if usr == "" && single == "" {
		fmt.Fprintln(os.Stderr,"You must specify an user (-u --user) or a tweet (-t --tweet)")
		os.Exit(1)
	}

	client = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: time.Duration(5) * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   time.Duration(5) * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			DisableKeepAlives:     true,
		},
	}
	if proxy != "" {
		proxyURL, _ := URL.Parse(proxy)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
	}

	scraper = twitterscraper.New()
	scraper.WithReplies(true)
	scraper.SetProxy(proxy)

	// Modified login handling
	if login || useCookies {
		Login(useCookies)
	}

	singleTweet(output, single)
	os.Exit(0)
}
