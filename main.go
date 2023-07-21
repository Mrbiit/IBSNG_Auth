package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"time"
)

func Authenticate(usr, psw, ibs string) (bool, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return false, err
	}
	cl := &http.Client{
		Jar:       jar,
		Transport: &http.Transport{},
	}
	r1, err := cl.PostForm(ibs, url.Values{
		"normal_username": {usr},
		"normal_password": {psw},
		"lang":            {"en"},
		"x":               {"18"},
		"y":               {"10"},
	})
	if err != nil {
		return false, err
	}
	defer r1.Body.Close()

	io.ReadAll(r1.Body)

	resp, err := cl.Get(ibs + "/home.php")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	r := regexp.MustCompile("Form_Content_Row_Right_2col_dark\">.[^\\d]+(\\d+-\\d+-\\d+)")
	exp := r.FindStringSubmatch(string(b))[1]
	t, err := time.Parse("2006-01-02", exp)
	if err != nil {
		return false, err
	}
	return t.After(time.Now()), nil
}

func main() {
	ibs := fmt.Sprintf("http://%s/IBSng/user/", os.Args[1])
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		psw := r.URL.Query()["password"][0]
		usr := r.URL.Query()["username"][0]
		if good, err := Authenticate(usr, psw, ibs); err == nil {
			if good {
				fmt.Fprintf(w, "true")
			} else {
				fmt.Fprintf(w, "false")
			}
		} else {
			fmt.Fprintln(w, err.Error())
		}

	})
	log.Fatal(http.ListenAndServe(":"+os.Args[2], nil))
}
