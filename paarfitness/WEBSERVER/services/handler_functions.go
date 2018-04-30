package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"../data"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Bruse Api:
func Bruse(writer http.ResponseWriter, request *http.Request) {

	var jsonStr = []byte(`{"title":"Bruse"}`)
	u, _ := url.ParseRequestURI("http://129.187.229.141:8080/fb")

	client := &http.Client{}
	r, _ := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonStr))
	r.Header.Add("Content-Type", "application/json")

	resp, _ := client.Do(r)
	fmt.Println(resp.Status)

	// resp, err := http.Get("http://129.187.229.141:8080/fb")
	// if err != nil {
	// 	fmt.Println("Error calling Brus")
	// }
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	writer.Write(body)

	//fmt.Fprintf(writer, body)
}

// redirect to HTTPS:
func Redirect(writer http.ResponseWriter, request *http.Request) {
	host := strings.Split(request.Host, ":")[0]
	http.Redirect(writer, request, "https://"+host+ Config.AddressSSL, http.StatusMovedPermanently)
}

// Terms and Conditions, Privacy, Third Party, GET /terms
func About(writer http.ResponseWriter, request *http.Request) {
	generateHTML(writer, nil, "layout", "about", "terms", "privacy", "thirdparty")
}

// GET /err?msg=
func Err(writer http.ResponseWriter, request *http.Request) {

	vals := request.URL.Query()
	//fmt.Println("Printing out request.URL.Query() values", vals)
	_, err := sessionCheck(writer, request)
	if err != nil {
		generateHTML(writer, vals.Get("msg"), "layout", "publayout","public.navbar")
	} else {
		generateHTML(writer, vals.Get("msg"), "layout", "privlayout", "private.navbar", "error")
	}
}

// GET /
func Index(writer http.ResponseWriter, request *http.Request) {

	// check if user has valid session
	_, err := sessionCheck(writer, request)
	if err != nil {
		generateHTML(writer, nil, "layout", "publayout","login")
	} else {
		generateHTML(writer, nil, "layout", "privlayout", "news", "dashboard", "settings", "startpage")
	}
}

// POST /signup
func SignupAccount(writer http.ResponseWriter, request *http.Request) {

	// parse signup data
	jsdata := body_to_json(request.Body)

	// verify age > 16 years
	tm, _ := time.Parse("01/02/2006", jsdata["age_month"]+"/15/1983")
	day, _ := strconv.Atoi(jsdata["age_day"])
	year, _ := strconv.Atoi(jsdata["age_year"])
	birthday := time.Date(year, tm.Month(), day, 0, 0, 0, 0, time.UTC)

	// verify captcha:
	/*err := request.ParseForm()
	if err != nil {
		warning(err, "Cannot parse form /signup", err)
	}*/
	// cap_resp := request.PostFormValue("g-recaptcha-response")
	// remoteip := strings.Split(request.RemoteAddr, ":")[0]
	// fmt.Println(cap_resp)
	// err = verifyCaptcha(remoteip, cap_resp)
	// if err != nil {
	// 	fmt.Println("capctha error: ", err)
	// 	generateHTML(writer, "Captcha Error", "layout", "signup.layout", "signup.err")
	//	return
	// }

	// create user in database and check if already exists:
	user := data.User{
		Company:    jsdata["company"],
		First_name: jsdata["first_name"],
		Last_name:  jsdata["last_name"],
		Email:      jsdata["email"],
		Country:    jsdata["country"],
		Pass:       jsdata["passw1"],
		Birthday:   birthday,
	}
	if err, stmt := user.Create(); err != nil || stmt != "" {
		writer.Header().Set("Content-Type", "application/json")
		json_content, _ := json.Marshal(map[string]string{"success": "false", "msg": stmt})
		fmt.Fprintf(writer, string(json_content))

	} else {
		writer.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(writer, `{"success": "true", "msg": "You registered successfully."}`)
	}
}

// POST /authenticate
func Authenticate(writer http.ResponseWriter, request *http.Request) {

	err := request.ParseForm()
	if err != nil {
		warning(err, "Cannot parse form /authenticate", err)
	}

	user, err := data.GetUserByEmail(request.PostFormValue("loginmail"))
	if err != nil {
		info(err, "Cannot find email")
	}
	if user.Pass == data.Encrypt(request.PostFormValue("loginpw")) {

		// create session
		device := request.Header["User-Agent"][0]
		sess := data.Session{User_id: user.User_id, Device: device}
		// check if inactive session of user id and device already in session table

		err := sess.InactiveExists()
		if err != nil {
			info(err, "Could not delete inactive existing session")
		}
		// create new session table entry
		sess, err = user.CreateSession(device)
		if err != nil {
			info(err, "Cannot create session")
		}
		cookie := http.Cookie{
			Name:     "_ianzncookie",
			Value:    sess.Session_id,
			HttpOnly: true,
		}
		http.SetCookie(writer, &cookie)
	}
	http.Redirect(writer, request, "/", 302)
}

// Get /delete_account
func DelAccount(writer http.ResponseWriter, request *http.Request) {

	// check cookie uuid
	cookie, err := request.Cookie("_ianzncookie")
	if err != http.ErrNoCookie {

		// get user Id in session
		sess := data.Session{Session_id: cookie.Value}
		user, err := sess.User()
		if err != nil {
			warning("Could not find user to session uuid", err)
		}

		// delete all user sessions and user
		if err = user.DeleteSessions(); err != nil {
			warning("Could not delete all session from user", err)
		}
		if err = user.Delete(); err != nil {
			warning("Could not delete User", err)
		}
	}
	http.Redirect(writer, request, "/", 302)
}

// GET /logout
func Logout(writer http.ResponseWriter, request *http.Request) {

	// check cookie uuid
	cookie, err := request.Cookie("_ianzncookie")
	if err != http.ErrNoCookie {
		info("Failed to get cookie", err)
		sess := data.Session{Session_id: cookie.Value}
		if err = sess.SetInactive(); err != nil {
			warning("Could not set Session to inactive in lougout function", err)
		}
	}
	http.Redirect(writer, request, "/", 302)
}
