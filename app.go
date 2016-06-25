package main

import "net/http"
import "net/url"
import "io/ioutil"
import "fmt"

var ClientID = YOUR_CLIENT_ID
var ClientSecret = YOUR_CLIENT_SECRET
var RedirectURI = "http://localhost:8080/redirect"

func redirect(res http.ResponseWriter, req *http.Request) {

	code := req.FormValue("code")

	if len(code) != 0 {

		formResponse, err := http.PostForm("https://api.instagram.com/oauth/access_token", url.Values{"client_id": {ClientID}, "client_secret": {ClientSecret}, "grant_type": {"authorization_code"}, "redirect_uri": {RedirectURI}, "code": {code}})
		if err != nil {
			fmt.Println(err)
			http.NotFound(res, req)
			return
		}
		defer formResponse.Body.Close()

		if formResponse.StatusCode == 200 {

			body, _ := ioutil.ReadAll(formResponse.Body)

			res.Write(body)
			return
		}
		fmt.Println(formResponse.StatusCode)
		http.NotFound(res, req)
	}

}

func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "index.html")
}

func main() {
	http.HandleFunc("/redirect", redirect)
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8080", nil)
}
