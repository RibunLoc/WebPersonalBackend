package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/RibunLoc/WebPersonalBackend/api-gateway/internal/client"
	"github.com/RibunLoc/WebPersonalBackend/api-gateway/util"
)

type AuthProxy struct {
	AuthGRPC     *client.AuthGRPC
	AuthHTTPBase string
	HTTP         *http.Client
}

func NewAuthProxy(grpcCl *client.AuthGRPC, httpBase string) *AuthProxy {
	return &AuthProxy{
		AuthGRPC:     grpcCl,
		AuthHTTPBase: httpBase,
		HTTP: &http.Client{
			Timeout: 7 * time.Second,
		},
	}
}

// POST /auth/login  (gRPC -> auth-service)
func (h *AuthProxy) Login(w http.ResponseWriter, r *http.Request) {
	var in client.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	res, err := h.AuthGRPC.Login(r.Context(), in)
	if err != nil {
		util.Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	util.JSON(w, http.StatusOK, res)
}

// POST /auth/register  (HTTP -> auth-service REST)
func (h *AuthProxy) Register(w http.ResponseWriter, r *http.Request) {
	target, _ := url.Parse(h.AuthHTTPBase + "/auth/register")

	// đọc body và forward
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "cannot read body")
		return
	}
	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, target.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		util.Error(w, http.StatusInternalServerError, "cannot create forward request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.HTTP.Do(req)
	if err != nil {
		util.Error(w, http.StatusBadGateway, "auth-service unavailable")
		return
	}
	defer resp.Body.Close()

	// truyền nguyên status + body trả về cho client
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}
