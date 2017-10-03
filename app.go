package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"github.com/urfave/negroni"
	"os"
	"github.com/satori/go.uuid"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

const (
	githubOauthApiBaseUrl = "https://github.com/login/oauth"
)

var (
	validStates map[string]bool
)

func main() {
	validStates = make(map[string]bool)

	clientId := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	redirectUrl := os.Getenv("CALLBACK_URL")

	r  := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "The callback service is online.")
	})

	r.HandleFunc("/begin", func(w http.ResponseWriter, r *http.Request) {
		// Generate a random state val
		state := newState()

		destination, _ := url.Parse(githubOauthApiBaseUrl+"/authorize")
		q := destination.Query()
		q.Set("client_id", clientId)
		q.Set("redirect_url", redirectUrl)
		q.Set("state", state)
		q.Set("allow_signup", "true")
		destination.RawQuery = q.Encode()

		http.Redirect(w, r, destination.String(), http.StatusTemporaryRedirect)
	})

	r.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		code := q.Get("code")
		state := q.Get("state")

		if !validateState(state) {
			http.Error(w, "Invalid state", http.StatusUnauthorized)
			return
		}

		requestUrl, _ := url.Parse(githubOauthApiBaseUrl + "/access_token")
		resp, err := http.PostForm(requestUrl.String(), url.Values{
			"client_id": []string{clientId},
			"client_secret": []string{clientSecret},
			"code": []string{code},
			"redirect_uri": []string{redirectUrl},
			"state": []string{state},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		v, err := url.ParseQuery(string(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		accessToken := v.Get("access_token")
		tokenType := v.Get("token_type")

		output := struct {
			AccessToken string `json:"access_token"`
			TokenType string `json:"token_type"`
		}{
			AccessToken: accessToken,
			TokenType: tokenType,
		}

		err = json.NewEncoder(w).Encode(&output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	n := negroni.New()
	n.UseHandler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	n.Run(":" + port)
}

func newState() string {
	state := uuid.NewV4().String()
	validStates[state] = true

	return state
}

func validateState(state string) bool {
	valid, exists := validStates[state]

	return exists && valid
}
