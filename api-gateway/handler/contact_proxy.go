package handler

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/RibunLoc/WebPersonalBackend/api-gateway/internal/client"
	"github.com/RibunLoc/WebPersonalBackend/api-gateway/util"
)

type ContactProxy struct {
	ContactGRPC *client.ContactGRPC
}

func NewContactProxy(cg *client.ContactGRPC) *ContactProxy {
	return &ContactProxy{ContactGRPC: cg}
}

func clientIP(r *http.Request) string {
	if v := r.Header.Get("CF-Connecting-IP"); v != "" {
		return strings.TrimSpace(v)
	}
	if v := r.Header.Get("X-Real-IP"); v != "" {
		return strings.TrimSpace(v)
	}
	if v := r.Header.Get("X-Forwarded-For"); v != "" {
		return strings.TrimSpace(v)
	}
	h, _, _ := net.SplitHostPort(r.RemoteAddr)
	return h
}

func (h *ContactProxy) Submit(w http.ResponseWriter, r *http.Request) {
	var in client.ContactSubmitInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	// validate cơ bản
	if len(strings.TrimSpace(in.Name)) < 2 || len(strings.TrimSpace(in.Email)) < 5 || len(strings.TrimSpace(in.Message)) < 5 {
		util.Error(w, http.StatusBadRequest, "name/email/message too short")
		return
	}
	in.RemoteIP = clientIP(r)
	out, err := h.ContactGRPC.Submit(r.Context(), in)
	if err != nil {
		util.Error(w, http.StatusBadGateway, "contact-service unavailable: "+err.Error())
		return
	}
	util.JSON(w, http.StatusCreated, out)
}
