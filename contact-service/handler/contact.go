package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/RibunLoc/WebPersonalBackend/contact-service/model"
	"github.com/RibunLoc/WebPersonalBackend/contact-service/repository"
	"github.com/RibunLoc/WebPersonalBackend/contact-service/util"
)

type ContactHandler struct {
	Repo     *repository.ContactMongo
	Verifier util.Verifier
	Emailer  util.EmailSender
}

type submitInput struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Token   string `json:"turnstile_token,omitempty"`
	CFToken string `json:"cf_turnstile_response,omitempty"`
}

func (h *ContactHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var in submitInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Email) == "" || len(strings.TrimSpace(in.Message)) < 5 {
		util.Error(w, http.StatusBadRequest, "name/email/message too short")
		return
	}
	// Lấy ip
	ip := util.ClientIP(r)
	token := strings.TrimSpace(in.Token)

	if token == "" {
		token = strings.TrimSpace(in.CFToken)
	}
	if h.Verifier != nil {
		if _, err := h.Verifier.Verify(r.Context(), token, ip); err != nil {
			util.Error(w, http.StatusBadRequest, "turnstile verification failed")
			return
		}
	}

	now := util.CustomTime(time.Now())
	c := model.Contact{
		Name:      strings.TrimSpace(in.Name),
		Email:     strings.TrimSpace(in.Email),
		Message:   strings.TrimSpace(in.Message),
		CreatedAt: &now,
	}
	if err := h.Repo.Create(r.Context(), &c); err != nil {
		util.Error(w, http.StatusInternalServerError, "failed to save contact")
		return
	}

	// Gửi email
	if h.Emailer != nil {
		if err := h.Emailer.Send("[Contact] "+c.Name,
			"New contact message:\r\n"+
				"Name: "+c.Name+"\r\n"+
				"Email: "+c.Email+"\r\n"+
				"Message:\r\n"+c.Message+"\r\n"+
				"Time: "+c.CreatedAt.String()+"\r\n",
		); err != nil {
			log.Printf("[warn] send mail failed: %v", err)
		} else {
			log.Printf("[email] sent notification")
		}
	} else {
		log.Println("[email] DISABLED (no emailer configured)")
	}

	util.JSON(w, http.StatusCreated, c)
}
