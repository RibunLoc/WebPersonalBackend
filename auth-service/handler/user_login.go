package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	repository "github.com/RibunLoc/microservices-learn/auth-service/repository"
	"github.com/RibunLoc/microservices-learn/auth-service/util"
)

type UserLogin struct {
	Repo *repository.RedisMongo
}

func (h *UserLogin) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	// Lấy user từ DB
	user, err := h.Repo.FindByEmail(r.Context(), body.Email)
	if err != nil || user == nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// so sánh password
	if !util.CheckPasswordHash(body.Password, user.Password) {
		http.Error(w, "Invaled email or passowrd", http.StatusUnauthorized)
		return
	}

	token, err := util.GenerateJWT(user.ID.Hex(), h.Repo.JwtSecret)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	// Trả về thông tin user trừ password
	resBody := struct {
		ID       interface{} `json:"id"`
		Email    string      `json:"email"`
		Fullname string      `json:"fullname"`
		Role     string      `json:"role"`
		Token    string      `json:"token"`
	}{
		ID:       user.ID,
		Email:    user.Email,
		Fullname: user.Fullname,
		Role:     user.Role,
		Token:    token,
	}

	res, err := json.Marshal(resBody)
	if err != nil {
		fmt.Println("failed to connect login response to json: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
