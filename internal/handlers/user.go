package handlers

import (
	"context"
	"encoding/json"
	"github.com/altiby/son/internal/domain"
	"github.com/altiby/son/internal/protocol"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
	"gopkg.in/validator.v2"
	"net/http"
)

type UserService interface {
	RegisterUser(ctx context.Context, user domain.User, password string) (domain.User, error)
	AuthorizeUser(ctx context.Context, id string, password string) (domain.User, error)
	GetUserByID(ctx context.Context, id string) (domain.User, error)
}

type UserHandler struct {
	userService UserService
	router      chi.Router
}

func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u.router.ServeHTTP(w, r)
}

type UserLoginRequest struct {
	Id       string `json:"id" validate:"nonzero,nonnil"`
	Password string `json:"password" validate:"nonzero,nonnil"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

func (u *UserHandler) login(w http.ResponseWriter, r *http.Request) {
	var request UserLoginRequest
	requestID := middleware.GetReqID(r.Context())

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		msg := "decode request failed"
		log.Err(err).Msg(msg)
		protocol.WriteError(w, requestID, http.StatusBadRequest, msg)
		return
	}

	requestValidator := validator.NewValidator()
	if err := requestValidator.Validate(request); err != nil {
		msg := "validate request failed"
		log.Err(err).Msg(msg)
		protocol.WriteError(w, requestID, http.StatusBadRequest, msg)
		return
	}

	user, err := u.userService.AuthorizeUser(r.Context(), request.Id, request.Password)
	if err == domain.ErrUserNotFound {
		msg := "user not found"
		log.Err(err).Msg(msg)
		protocol.WriteError(w, requestID, http.StatusNotFound, msg)
		return
	}
	if err != nil {
		msg := "authorize user failed"
		log.Err(err).Msg(msg)
		protocol.WriteError(w, requestID, http.StatusInternalServerError, msg)
		return
	}

	protocol.WriteOk(w, UserLoginResponse{
		Token: user.ID,
	})
}

type UserRegisterRequest struct {
	FirstName  string `json:"first_name" validate:"nonzero,nonnil"`
	SecondName string `json:"second_name" validate:"nonzero,nonnil"`
	Birthdate  string `json:"birthdate" validate:"nonzero,nonnil"`
	Biography  string `json:"biography" validate:"nonzero,nonnil"`
	City       string `json:"city" validate:"nonzero,nonnil"`
	Password   string `json:"password" validate:"nonzero,nonnil"`
}

func (r UserRegisterRequest) ToDomainUser() domain.User {
	return domain.User{
		FirstName:  r.FirstName,
		SecondName: r.SecondName,
		Birthdate:  r.Birthdate,
		Biography:  r.Biography,
		City:       r.City,
	}
}

type UserRegisterResponse struct {
	UserID string `json:"user_id"`
}

func (u *UserHandler) register(w http.ResponseWriter, r *http.Request) {
	var request UserRegisterRequest
	requestID := middleware.GetReqID(r.Context())

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		msg := "decode request failed"
		log.Err(err).Msg(msg)
		protocol.WriteError(w, requestID, http.StatusBadRequest, msg)
		return
	}

	requestValidator := validator.NewValidator()
	if err := requestValidator.Validate(request); err != nil {
		msg := "validate request failed"
		log.Err(err).Msg(msg)
		protocol.WriteError(w, requestID, http.StatusBadRequest, msg)
		return
	}

	user, err := u.userService.RegisterUser(r.Context(), request.ToDomainUser(), request.Password)
	if err != nil {
		msg := "create user failed"
		log.Err(err).Msg(msg)
		protocol.WriteError(w, requestID, http.StatusInternalServerError, msg)
		return
	}

	protocol.WriteOk(w, UserRegisterResponse{
		UserID: user.ID,
	})

}

type GetUserResponse struct {
	FirstName  string `json:"first_name" validate:"nonzero,nonnil"`
	SecondName string `json:"second_name" validate:"nonzero,nonnil"`
	Birthdate  string `json:"birthdate" validate:"nonzero,nonnil"`
	Biography  string `json:"biography" validate:"nonzero,nonnil"`
	City       string `json:"city" validate:"nonzero,nonnil"`
}

func (r *GetUserResponse) FromDomainUser(user domain.User) {
	r.Biography = user.Biography
	r.City = user.City
	r.Birthdate = user.Birthdate
	r.FirstName = user.FirstName
	r.SecondName = user.SecondName
}

func (u *UserHandler) get(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	requestID := middleware.GetReqID(r.Context())

	if len(userID) == 0 {
		msg := "user id is empty error"
		log.Error().Msg(msg)
		protocol.WriteError(w, requestID, http.StatusBadRequest, msg)
		return
	}

	user, err := u.userService.GetUserByID(r.Context(), userID)

	if err == domain.ErrUserNotFound {
		msg := "user not found"
		log.Err(err).Msg(msg)
		protocol.WriteError(w, requestID, http.StatusNotFound, msg)
		return
	}
	if err != nil {
		msg := "get user failed"
		log.Err(err).Msg(msg)
		protocol.WriteError(w, requestID, http.StatusInternalServerError, msg)
		return
	}

	userResponse := GetUserResponse{}
	userResponse.FromDomainUser(user)

	protocol.WriteOk(w, userResponse)
}

func NewUserHandler(service UserService) *UserHandler {
	handler := &UserHandler{userService: service}
	r := chi.NewRouter()
	r.Post("/login", handler.login)
	r.Post("/register", handler.register)
	r.Get("/get/{id}", handler.get)
	handler.router = r
	return handler
}
