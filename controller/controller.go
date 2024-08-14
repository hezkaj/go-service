package controller

import (
	"ai-test/services"
	"ai-test/validation"
	"database/sql"
	"fmt"
	"net"
	"strings"

	"net/http"
)

func Router(db *sql.DB) {
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		access, refresh, err := services.GetTokens(id, ip, db)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintln(err)))
		}
		w.WriteHeader(201)
		write := fmt.Sprintf("access: %s, refresh: %s", access, refresh)
		w.Write([]byte(write))
	})
	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		bearer_token := r.Header.Get("Authorization")
		new_ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		_, ref_token, _ := strings.Cut(bearer_token, " ")
		err, mail, old_ip := validation.CheckValidAndGetMailAndIP(ref_token)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte("Authorization error"))
		}
		access, refresh, err := services.RefreshTokens(mail, new_ip, old_ip, db, ref_token)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintln(err)))
		}
		w.WriteHeader(201)
		write := fmt.Sprintf("access: %s, refresh: %s", access, refresh)
		w.Write([]byte(write))
	})
}
