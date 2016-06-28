package main

import "net/http"
import "io/ioutil"
import "fmt"
import "html/template"
import "golang.org/x/oauth2"

var ClientID = YOUR_CLIENT_ID

var ClientSecret = YOUR_CLIENT_SECRET

var RedirectURI = "http://localhost:8080/redirect"

var authURL = "https://api.instagram.com/oauth/authorize"

var tokenURL = "https://api.instagram.com/oauth/access_token"

var templ = template.Must(template.New("index2.html").ParseFiles("index2.html"))

var igConf *oauth2.Config

func redirect(res http.ResponseWriter, req *http.Request) {

	code := req.FormValue("code")

	if len(code) != 0 {
		tok, err := igConf.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Println(err)
			http.NotFound(res, req)
			return
		}

		if tok.Valid() {
			client := igConf.Client(oauth2.NoContext, tok)

			request, err := http.NewRequest("GET", "https://api.instagram.com/v1/users/self/?access_token="+tok.AccessToken, nil)
			if err != nil {
				fmt.Println(err)
				http.NotFound(res, req)
				return
			}

			resp, err := client.Do(request)
			if err != nil {
				fmt.Println(err)
				http.NotFound(res, req)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				http.NotFound(res, req)
				return
			}

			res.Write(body)
		}

		http.NotFound(res, req)
	}

}

func homePage(res http.ResponseWriter, req *http.Request) {
	url := igConf.AuthCodeURL("", oauth2.AccessTypeOffline)
	fmt.Println(url)
	err := templ.Execute(res, url)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	igConf = &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		RedirectURL: RedirectURI,
		Scopes:      []string{"public_content", "comments"},
	}

	http.HandleFunc("/redirect", redirect)
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8080", nil)
}
