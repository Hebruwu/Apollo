package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"apollo.io/objects/request"
	"apollo.io/objects/response"
	"apollo.io/objects/servershared"
	"apollo.io/serverutils"
	"apollo.io/services"
)

type userRoute struct {
	services.UserService
	serverutils.ServerConfig
}

func NewUsersBase(serverConfig serverutils.ServerConfig) http.Handler {
	users := http.NewServeMux()
	ur := userRoute{
		services.NewUserService(serverConfig.PostgresConnection, serverConfig.Logger),
		serverConfig,
	}

	users.HandleFunc("POST /login", ur.authenticate)
	users.HandleFunc("POST /register", ur.createUser)
	users.HandleFunc("POST /logout", ur.logoutUser)
	users.HandleFunc("POST /refresh", ur.refreshToken)
	users.HandleFunc("PUT /{id}", ur.updateUser)
	users.HandleFunc("DELETE /{id}", ur.deleteUser)

	return users
}

func (ur userRoute) authenticate(_ http.ResponseWriter, _ *http.Request) {
	// TODO: Implement getUser
}

// Handles user creation.
func (ur userRoute) createUser(w http.ResponseWriter, r *http.Request) {
	var user request.NewUser
	err := json.NewDecoder(r.Body).Decode(&user)

	// Malformed payload.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = ur.CreateUser(r.Context(), user.Username, user.Email, user.Password)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		var encoderError error
		if errors.Is(err, servershared.ErrUsernameAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			encoderError = json.NewEncoder(w).Encode(response.StatusResponse{Error: response.UsernameAlreadyExists, Success: ""})
		} else {
			ur.Logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			encoderError = json.NewEncoder(w).Encode(response.StatusResponse{Error: response.UnexpectedError, Success: ""})
		}
		if encoderError != nil {
			ur.Logger.Error(encoderError.Error())
			http.Error(w, response.UnexpectedError, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response.StatusResponse{Success: response.UserCreated, Error: ""})
	if err != nil {
		ur.Logger.Error(err.Error())
		http.Error(w, response.UnexpectedError, http.StatusInternalServerError)
	}
}

func (ur userRoute) updateUser(_ http.ResponseWriter, _ *http.Request) {
	// TODO: Implement updateUser
}

func (ur userRoute) deleteUser(_ http.ResponseWriter, _ *http.Request) {
	// TODO: Implement deleteUser
}

func (ur userRoute) logoutUser(_ http.ResponseWriter, _ *http.Request) {
	// TODO: Implement logoutUser
}

func (ur userRoute) refreshToken(_ http.ResponseWriter, _ *http.Request) {
	// TODO: Implement refreshToken
}
