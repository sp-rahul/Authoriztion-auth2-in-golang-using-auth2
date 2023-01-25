// package main

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"

// 	_ "github.com/joho/godotenv/autoload"
// 	"golang.org/x/oauth2"
// 	"golang.org/x/oauth2/google"
// )

// var (
// 	googleOauthConfig = &oauth2.Config{
// 		RedirectURL:  "http://localhost:8080/callback",
// 		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
// 		ClientSecret: os.Getenv("GOOGLE_CLIENT_ID"),
// 		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
// 		Endpoint:     google.Endpoint,
// 	}
// 	//Todo Random
// 	randomState = "pseudo-random"
// )

// func main() {
// 	http.HandleFunc("/", handleHome)
// 	http.HandleFunc("/login", handleLogin)
// 	http.HandleFunc("/callbcak", handleCallback)
// 	http.ListenAndServe(":8080", nil)

// }
// func handleHome(w http.ResponseWriter, r *http.Request) {
// 	var html = `<html><body><a href="/login"> Google Log In</a></body></html>`
// 	fmt.Fprint(w, html)
// }
// func handleLogin(w http.ResponseWriter, r *http.Request) {
// 	url := googleOauthConfig.AuthCodeURL(randomState)
// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// }

// func handleCallback(w http.ResponseWriter, r *http.Request) {
// 	if r.FormValue("state") != randomState {
// 		fmt.Println("State is not valid")

// 		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
// 		return
// 	}
// 	//token, err := googleOauthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
// 	token, err := googleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
// 	if err != nil {
// 		fmt.Printf("Could not get token: %s \n", err.Error())
// 		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
// 		return
// 	}

// 	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token" + token.AccessToken)
// 	if err != nil {
// 		fmt.Printf("Could not create get token: %s \n", err.Error())
// 		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
// 		return
// 	}

// 	defer resp.Body.Close()
// 	content, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Printf("Could not parse token: %s \n", err.Error())
// 		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
// 		return
// 	}

// 	fmt.Fprintf(w, "Response %s", content)

// }

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	//"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
)

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL: "http://localhost:8080/callback",
		// ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		// ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),

		//Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},

		Scopes:   []string{"https://github.com/login/oauth/authorize"},
		Endpoint: github.Endpoint,
	}
}

// AuthURL:  "https://github.com/login/oauth/authorize",
//
//	TokenURL: "https://github.com/login/oauth/access_token",
//	TokenURL: "https://github.com/login/oauth/access_token",
func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/callback", handleGoogleCallback)
	http.ListenAndServe(":8080", nil)
}
func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html>
<body>
	<a href="/login"> Github Log In</a>
</body>
</html>`
	fmt.Fprintf(w, htmlIndex)
}

var (
	// TODO: randomize it
	oauthStateString = "pseudo-random"
)

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Fprintf(w, "Content: %s\n", content)
}
func getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	// response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	response, err := http.Get("https://github.com/login/oauth/access_token/userinfo?access_token=" + token.AccessToken)
	//https://github.com/login/oauth/access_token
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}
