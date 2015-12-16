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
	"time"
)

type UserInfo struct {
	name     string
	account  string
	password string
}

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

func signinHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signin process")
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		fmt.Println("cookies")
		if len(r.Cookies()) == 0 {
			expiration := time.Now()
			expiration = expiration.AddDate(1, 0, 0)
			cookie := http.Cookie{Name: "username", Value: "astaxie", Expires: expiration}
			http.SetCookie(w, &cookie)
		} else {
			for _, cookie := range r.Cookies() {
				fmt.Println(cookie.Name)
			}
		}
		fmt.Println("cookies over")

		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		userInfo := UserInfo{"", r.Form["account"][0], r.Form["password"][0]}
		err := userInfo.signinInfoCheck()
		if err != nil {
			fmt.Fprintln(w, err.Error())
		} else {
			fmt.Fprintln(w, "Signin success!")
		}

		fmt.Println("account:", r.Form["account"])
		fmt.Println("password:", r.Form["password"])
	}
}

func signupHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signup process")
	fmt.Println("method:", r.Method) //获取请求的方法
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
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
