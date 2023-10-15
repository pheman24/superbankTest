package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/abhinav-TB/dantdb"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Account struct {
	Type         string `json:"type,omitempty"`
	Pid          string `json:"pid,omitempty"`
	Email        string `json:"email,omitempty"`
	Password     string `json:"password,omitempty"`
	Token        string `json:"Token,omitempty"`
	RefreshToken string `json:"RefreshToken,omitempty"`
}

// Scopes: OAuth 2.0 scopes provide a way to limit the amount of access that is granted to an access token.

var (
	googleOauthConfig *oauth2.Config
	db                = Database()
	fileName          = "pageOne.html"
)

func init() {
	// init ENV
	// initialize the variable using ENV values
	os.Setenv("GOOGLE_OAUTH_CLIENT_ID", "983032359383-fjusp2f0su6j71tbnajhjjb7hehbtjeo.apps.googleusercontent.com")
	os.Setenv("GOOGLE_OAUTH_CLIENT_SECRET", "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")
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8000/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {

	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)

	u := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, account, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Redirect or response with a token.
	// More code .....

	var f interface{}
	errv := json.Unmarshal(data, &f)
	if errv != nil {
		fmt.Println("Error parsing JSON: ", err)
	}

	// JSON object parses into a map with string keys
	itemsMap := f.(map[string]interface{})
	account.Email = itemsMap["email"].(string)

	db.Write("account", account.Email, account)
	movePageOne(w)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromGoogle(code string) ([]byte, *Account, error) {
	account := new(Account)
	// Use code to get token and get user info from Google.
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, account, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	account.Token = token.AccessToken
	account.RefreshToken = token.RefreshToken

	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, account, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, account, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, account, nil
}

func LoginByEmail(w http.ResponseWriter, r *http.Request) {
	pass := r.FormValue("password")
	email := r.FormValue("email")
	account := new(Account)
	db.Read("account", r.FormValue("email"), account)

	if strings.ToLower(email) != strings.ToLower(account.Email) {
		println("False " + pass)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	println("True " + pass)
	movePageOne(w)

}

func Database() dantdb.Driver {
	dir := "./database/"

	database, err := dantdb.New(dir) // creates a new database
	if err != nil {
		fmt.Println("Error", err)
	}
	return *database
}

func movePageOne(w http.ResponseWriter) (t *template.Template) {
	t, _ = template.ParseFiles(fileName)
	errs := t.ExecuteTemplate(w, fileName, nil)
	if errs != nil {
		fmt.Println(w, "UserInfo: %s\n", errs)
	}
	return
}
