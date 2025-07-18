package User

import (
	"REST-API-pet-proj/Internal/HttpServer/Api/User/Handlers"
	"REST-API-pet-proj/Internal/Storage"
	"REST-API-pet-proj/Internal/Storage/Sqlite"
	"REST-API-pet-proj/Models"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var IsLogin bool

func Registration(storage *Sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Models.UserRegistration

		if !Handlers.ParseAndValidateJSON(w, r, &req) {
			return
		}

		var hash string
		hash, err := Handlers.PasswordHash(req.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		err = Storage.SaveUser(
			req.Username,
			req.Email,
			hash)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte("User Registration Successful"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func Login(storage *Sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Models.UserRegistration

		if !Handlers.ParseAndValidateJSON(w, r, &req) {
			return
		}

		passwordHash, err := Storage.GetUserPassword(req.Username, req.Email)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = w.Write([]byte("Invalid username or email"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		_, err = w.Write([]byte("User Login Successful"))
		IsLogin = true
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func Data(storage *Sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var username = chi.URLParam(r, "username")

		userData, _ := Storage.GetUserData(username)
		if userData == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		err := json.NewEncoder(w).Encode(userData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
