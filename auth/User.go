package auth

import "net/http"

type User struct {
	IsLogin bool
}

func (l *User) CheckLogin(r *http.Request) {
	token, _ := r.Cookie("token")
	if token == nil {
		l.IsLogin = false
		return
	}

	l.IsLogin = true
}
