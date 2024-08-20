package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type OAuth struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	
	// Expiry is the optional expiration time of the access token.
	//
	// If zero, TokenSource implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	//Expiry time.Time `json:"expiry,omitempty"`


	Error        string
}


// Login to Spotify using the OAuth method
func LoginOAuth(clientId, clientSecret, callbackURL string) (string, error) {
	if callbackURL == "" {
		callbackURL = "http://localhost:8888/callback"
	}
	token, err := getOAuthToken(clientId, clientSecret, callbackURL)
	if err != nil {
		return "", err
	}
	tok := token.AccessToken
	fmt.Println("Got oauth token:\n", tok)
	return tok, nil
}

func GetOauthAccessToken(code string, redirectUri string, clientId string, clientSecret string) (*OAuth, error) {
	val := url.Values{}
	val.Set("grant_type", "authorization_code")
	val.Set("code", code)
	val.Set("redirect_uri", redirectUri)
	val.Set("client_id", clientId)
	val.Set("client_secret", clientSecret)

	resp, err := http.PostForm("https://accounts.spotify.com/api/token", val)
	if err != nil {
		// Retry since there is an nginx bug that causes http2 streams to get
		// an initial REFUSED_STREAM response
		// https://github.com/curl/curl/issues/804
		resp, err = http.PostForm("https://accounts.spotify.com/api/token", val)
		if err != nil {
			return nil, err
		}
	}
	defer resp.Body.Close()
	auth := OAuth{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &auth)
	fmt.Printf("auth:\n%+v\n", auth)
	if err != nil {
		return nil, err
	}
	if auth.Error != "" {
		return nil, fmt.Errorf("error getting token %v", auth.Error)
	}
	return &auth, nil
}

/*
func StartLocalOAuthServer(clientId string, clientSecret string, callback string) (string, chan OAuth) {
	ch := make(chan OAuth)

	urlPath := "https://accounts.spotify.com/authorize?" +
		"client_id=" + clientId +
		"&response_type=code" +
		"&redirect_uri=" + callback +
		"&scope=streaming"

	router := http.NewServeMux()
	server := &http.Server{
		// TODO pull port from callback
		Addr:    ":5000",
		Handler: router,
	}

	router.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		auth, err := GetOauthAccessToken(params.Get("code"), callback, clientId, clientSecret)
		if err != nil {
			fmt.Fprintf(w, "Error getting token: %q", err)
			return
		}
		fmt.Fprintf(w, "Got token, logging in.")
		ch <- *auth
		close(ch)

		time.Sleep(time.Second * 5)
		_ = server.Shutdown(context.Background())
	})

	go func() {
		_ = server.ListenAndServe()
	}()

	return urlPath, ch
}
*/


func getOAuthToken(clientId string, clientSecret string, callback string) (*OAuth, error) {
	ch := make(chan *OAuth)

	fmt.Println("go to this url")
	urlPath := "https://accounts.spotify.com/authorize?" +
		"client_id=" + clientId +
		"&response_type=code" +
		"&redirect_uri=" + callback +
		"&scope=streaming"
	fmt.Println(urlPath)
	
	// router := http.NewServeMux()
	// server := &http.Server{
	// 	// TODO pull port from callback
	// 	Addr:    ":5000",
	// 	Handler: router,
	// }
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		auth, err := GetOauthAccessToken(params.Get("code"), callback, clientId, clientSecret)
		if err != nil {
			fmt.Fprintf(w, "Error getting token %q", err)
			return 
		}
		fmt.Fprintf(w, "Got token, loggin in")
		ch <- auth
		
		// time.Sleep(time.Second * 1)
		// _ = server.Shutdown(context.Background())
	})

	go func() {
		log.Fatal(http.ListenAndServe(":5000", nil))
	}()
	
		// Wait then bail
	select {
	case <-time.After(time.Second * 60):
		return nil, errors.New("timed out waiting for auth")
	case validAuth := <-ch:
			return validAuth, nil
	}


}
