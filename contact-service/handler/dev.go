package handler

import (
	"net/http"

	"github.com/RibunLoc/WebPersonalBackend/contact-service/util"
)

// handler/dev.go
func (h *ContactHandler) DevSendMail(w http.ResponseWriter, r *http.Request) {
	if h.Emailer == nil {
		util.Error(w, http.StatusBadRequest, "emailer disabled")
		return
	}
	if err := h.Emailer.Send("[Contact][DEV] test", "This is a test email."); err != nil {
		util.Error(w, http.StatusBadGateway, "send failed: "+err.Error())
		return
	}
	util.JSON(w, http.StatusOK, map[string]string{"status": "sent"})
}
