//web
package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"web/session"
	_ "web/session/memory"
)

type UserInfo struct {
	name     string
	account  string
	password string
}

// validation
func (userInfo *UserInfo) signinInfoCheck() error {
	if len(userInfo.account) == 0 {
		return errors.New("Username shoudn't be empty!")
	}
	if len(userInfo.password) == 0 {
		return errors.New("Password shoudn't be empty!")
	}

	phoneVaild, _ := regexp.MatchString(`^1[3|4|5|7|8|9][0-9]{9}$`, userInfo.account)
	emailVaild, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, userInfo.password)
	if !phoneVaild && !emailVaild {
		return errors.New("Wrong account, must be phone of email!")
	}

	db, err := sql.Open("mysql", "root:lyk19921009@/Web")
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT * FROM user WHERE phone = ? AND password = ?", userInfo.account, userInfo.password)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return errors.New("Wrong account or password!")
	} else {
		return nil
	}
}

func (userInfo *UserInfo) signupInfoCheck() error {
	if len(userInfo.name) == 0 {
		return errors.New("Name shoudn't be empty!")
	}
	if len(userInfo.account) == 0 {
		return errors.New("Phone number shoudn't be empty!")
	}
	if len(userInfo.password) == 0 {
		return errors.New("Password shoudn't be empty!")
	}

	phoneVaild, _ := regexp.MatchString(`^1[3|4|5|7|8|9][0-9]{9}$`, userInfo.account)
	emailVaild, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, userInfo.password)
	if !phoneVaild && !emailVaild {
		return errors.New("Wrong account, must be phone of email!")
	}

	db, err := sql.Open("mysql", "root:lyk19921009@/Web")
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT user SET name=?,phone=?,password=?")
	if err != nil {
		return err
	}
	result, err := stmt.Exec(userInfo.name, userInfo.account, userInfo.password)
	if err != nil {
		return err
	}
	if r, _ := result.RowsAffected(); r == 0 {
		return errors.New("A user is already registered with this phone number!")
	} else {
		return nil
	}
}

var manager, _ = session.NewManager("memory", "login", 100)

// request handle
func signinHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signin process")
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		userInfo := UserInfo{"", r.Form["account"][0], r.Form["password"][0]}
		cookie, err := r.Cookie("login")
		var count = 0
		if err != nil || cookie.Value == "" {
			err := userInfo.signinInfoCheck()
			if err != nil {
				fmt.Fprintln(w, err.Error())
			} else {
				sess := manager.SessionStart(w, r)
				sess.Set("count", 1)
				count = 1
				fmt.Fprintln(w, "Signin success!")
			}
		} else {
			sess := manager.SessionStart(w, r)
			count = sess.Get("count").(int)
			count++
			sess.Set("count", count)
			fmt.Fprintln(w, "Signin success!")
		}
		fmt.Println("count:", count)
		fmt.Println("account:", r.Form["account"])
		fmt.Println("password:", r.Form["password"])
	}
}

func signupHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signup process")
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		//
		r.ParseForm()
		userInfo := UserInfo{r.Form["name"][0], r.Form["phone"][0], r.Form["password"][0]}
		err := userInfo.signinInfoCheck()
		if err != nil {
			fmt.Fprintln(w, err.Error())
		} else {
			fmt.Fprintln(w, "Signup success!")
		}
		fmt.Println("name:", r.Form["name"])
		fmt.Println("phone:", r.Form["phone"])
		fmt.Println("password:", r.Form["password"])
	}
}

func main() {
	http.HandleFunc("/signin", signinHandle)
	http.HandleFunc("/signup", signupHandle)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
