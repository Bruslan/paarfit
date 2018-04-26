package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"../data"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const capcthaURL = "https://www.google.com/recaptcha/api/siteverify"

var isSmall = regexp.MustCompile(`^[a-z]+$`).MatchString
var isCap = regexp.MustCompile(`^[A-Z]+$`).MatchString
var isNumb = regexp.MustCompile(`^[0-9]+$`).MatchString

// captcha reponse struct
type ApiCaptchaResponse struct {
	success     bool
	challengeTs time.Time
	hostname    string
	errorCodes  []int
}

// configuration struct:
type Configuration struct {
	Address      string
	AddressSSL   string
	Redirect     string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

var logger *log.Logger
var Config Configuration

// triggers parameter loading(config.json), creates logfile ***
func init() {
	loadConfig()
	file, err := os.OpenFile("./services/server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file in services package", err)
	}
	logger = log.New(file, "Info ", log.Ldate|log.Ltime|log.Lshortfile)
}

// loads config.json parameters ***
func loadConfig() {
	absPath, _ := filepath.Abs("./services/config.json")
	file, err := os.Open(absPath)
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	Config = Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}

// verification of capctha:
func verifyCaptcha(remoteip, cap_resp string) (err error) {
	resp, err := http.PostForm(capcthaURL,
		url.Values{"secret": {"6LcBKkoUAAAAAF5UcvuWKV-7TqDXp9s1i_PAM3wn"},
			"remoteip": {remoteip}, "reponse": {cap_resp}})
	if err != nil {
		danger("HTTP post form captcha error:", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		danger("Read captcha body error:", err)
	}
	var res ApiCaptchaResponse
	err = json.Unmarshal(b, &res)
	if err != nil {
		danger("Parse json error: ", err)
	}
	if res.success {
		return nil
	}
	return err
}

// Checks if the user is logged in and has a session, if not err is not nil ***
func sessionCheck(writer http.ResponseWriter, request *http.Request) (sess data.Session, err error) {
	// check if cookie exists
	cookie, err := request.Cookie("_ianzncookie")
	if err == nil {
		// check if session valid
		device := request.Header["User-Agent"][0]
		sess = data.Session{Session_id: cookie.Value, Device: device}
		if ok := sess.SessValid(); !ok {
			err = errors.New("Session Invalid")
		}
	}
	return
}

// parse http message body to json:
func body_to_json(b io.Reader) (jsondata map[string]string) {

	// parses single json without inner inner jsons
	body, _ := ioutil.ReadAll(b)
	err := json.Unmarshal(body, &jsondata)
	if err != nil {
		warning(err, "Cannot parse javascript data signup", err)
	}
	return
}

// parse HTML templates
func parseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("html/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))
	return
}

// passses html to agent
func generateHTML(writer http.ResponseWriter, data interface{}, startfile string, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("html/%s.html", file))
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(writer, startfile, data)

}

// Convenience function to redirect to the error message page
func error_message(writer http.ResponseWriter, request *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(writer, request, strings.Join(url, ""), http.StatusFound)
}

// for logging
func info(args ...interface{}) {
	logger.SetPrefix("INFO ")
	logger.Println(args...)
}

func danger(args ...interface{}) {
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}

func warning(args ...interface{}) {
	logger.SetPrefix("WARNING ")
	logger.Println(args...)
}
