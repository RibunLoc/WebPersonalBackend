package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RibunLoc/WebPersonalBackend/auth-service/model"
	"github.com/RibunLoc/WebPersonalBackend/auth-service/util"

	repository "github.com/RibunLoc/WebPersonalBackend/auth-service/repository"
)

type UserRegister struct {
	Repo *repository.RedisMongo
}

func (h *UserRegister) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"fullname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	passwordHash, err := util.HashPassword(body.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	now := util.CustomTime(time.Now())

	userNew := &model.User{
		Email:     body.Email,
		Password:  passwordHash,
		Fullname:  body.FullName,
		Role:      "user",
		IsActive:  true,
		CreatedAt: &now,
	}

	//check email
	existingUser, err := h.Repo.FindByEmail(r.Context(), userNew.Email)
	if err == nil && existingUser != nil {
		fmt.Printf("User exists: %v\n", existingUser.Email)
		w.WriteHeader(http.StatusConflict)
		return
	}

	if err := h.Repo.CreateUser(r.Context(), userNew); err != nil {
		fmt.Println("failed to create new user: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(userNew)
	if err != nil {
		fmt.Println("failed to convert user to Json: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)

}
