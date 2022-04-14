package main

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

var id string
var guild string

// https://stackoverflow.com/questions/35558166/when-to-randomize-auth-code-state-in-oauth2
func genState(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(data), nil
}

var cancelFunc func()
var conf *oauth2.Config

const stateLength = 10

var state string

type handler struct{}

const port = 3000

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check for valid request
	if r.FormValue("state") != state {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error logging in."))
		return
	} else {
		Write("login", "Logging in...")
	}

	// Exchange code for access token
	token, err := conf.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		Error("login", "Login error: %s", err.Error())
		os.Exit(1)
	}

	// Get user info
	res, err := conf.Client(context.Background(), token).Get("https://discord.com/api/users/@me")
	if err != nil || res.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error logging in."))

		if err != nil {
			Error("login", "Login error: %s", err.Error())
		} else {
			Error("login", "Login error: %s", res.Status)
		}
		os.Exit(1)
	}
	defer res.Body.Close()

	// Read
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error logging in."))

		Error("login", "Login error: %s", err.Error())
		os.Exit(1)
	}

	// Return
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged in successfully."))

	// Process
	var data map[string]interface{}
	Clear()
	Write("login", "Logged in!")
	json.Unmarshal(body, &data)
	id = data["id"].(string)

	// Close
	time.Sleep(time.Second / 5) // Wait for resp to be sent
	cancelFunc()
}

func Login() {
	conf = &oauth2.Config{
		RedirectURL: fmt.Sprintf("http://localhost:%d/", port),
		// This next 2 lines must be edited before running this.
		ClientID:     "964274065508556800",
		ClientSecret: "dRDvGpuHZAgH6u-F5_UHakTxZgLewhe4",
		Scopes:       []string{"identify"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://discord.com/api/oauth2/authorize",
			TokenURL:  "https://discord.com/api/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	var err error
	state, err = genState(stateLength)
	if err != nil {
		Error("login", "Login error: %s", err.Error())
		os.Exit(1)
	}
	url := conf.AuthCodeURL(state)
	Write("login", "Login at %s", url)

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: &handler{}}
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
	}
	cancelFunc = func() {
		listener.Close()
	}

	Write("login", "Waiting for login...")
	if err := server.Serve(listener); err != nil {
		if !strings.HasSuffix(err.Error(), "use of closed network connection") {
			Error("login", "Login error: %s", err.Error())
			os.Exit(1)
		}
	}
}

func GuildLogin(name string) {
	// TODO: Actually check if guild exists
	guild = name
}
