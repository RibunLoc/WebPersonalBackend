package util

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/smtp"
	"strings"
	"time"
)

// EmailSender là interface để handler gọi  gửi email
type EmailSender interface {
	Send(subject, body string) error
}

// Cấu hình SMTP
type SMTPConfig struct {
	Host     string   // smtp.gmail.com
	Port     int      // 645 (TLS) | 587 (STARTTLS) | 1025 (MailHog dev )
	Username string   // user/email SMTP
	Password string   // app password/API SMTP
	From     string   //địa chỉ gửi
	To       []string // danh sách người nhận
	Timeout  time.Duration
}

// Triển khai EmailSender với SMTP
type SMTPSender struct {
	cfg SMTPConfig
}

// Helper khởi tạo
func NewSMTPSender(cfg SMTPConfig) *SMTPSender {
	if cfg.Timeout == 0 {
		cfg.Timeout = 7 * time.Second
	}
	return &SMTPSender{cfg: cfg}
}

// Gửi email dạng text/plain. Tự chọn TLS theo port
// - 465: implicit TLS
// - 587: STARTTLS (yêu cầu server hỗ trợ )
// - khacsL thử STARTTLLS, nếu không hỗ trợ thì gửi plain (dùng cho MailHog dev)
// helper: chuẩn hoá CRLF + quoted-printable
func qpEncode(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\n", "\r\n")
	var b bytes.Buffer
	w := quotedprintable.NewWriter(&b)
	_, _ = w.Write([]byte(s))
	_ = w.Close()
	out := b.String()
	if !strings.HasSuffix(out, "\r\n") {
		out += "\r\n"
	}
	return out
}

// helper: escape HTML tối thiểu
func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func (s *SMTPSender) Send(subject, body string) error {
	if len(s.cfg.To) == 0 {
		return fmt.Errorf("no recipients")
	}
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	// ----- Build MIME multipart/alternative (giống Mailtrap) -----
	boundary := fmt.Sprintf("boundary-%d", time.Now().UnixNano())
	now := time.Now().UTC().Format(time.RFC1123Z)
	subj := mime.QEncoding.Encode("utf-8", subject)

	// plain part
	plain := qpEncode(body)

	// html part (đơn giản: bọc vào <pre> + escape)
	html := "<!doctype html>\r\n<html><body>" +
		"<pre style=\"font-family:sans-serif;white-space:pre-wrap\">" +
		htmlEscape(body) +
		"</pre></body></html>"
	htmlQP := qpEncode(html)

	var msg bytes.Buffer
	// headers
	fmt.Fprintf(&msg, "From: %s\r\n", s.cfg.From)
	fmt.Fprintf(&msg, "To: %s\r\n", strings.Join(s.cfg.To, ", "))
	fmt.Fprintf(&msg, "Subject: %s\r\n", subj)
	fmt.Fprintf(&msg, "Date: %s\r\n", now)
	fmt.Fprintf(&msg, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&msg, "Content-Type: multipart/alternative; boundary=\"%s\"\r\n", boundary)
	fmt.Fprintf(&msg, "\r\n") // header/body separator

	// -- text/plain
	fmt.Fprintf(&msg, "--%s\r\n", boundary)
	fmt.Fprintf(&msg, "Content-Type: text/plain; charset=\"utf-8\"\r\n")
	fmt.Fprintf(&msg, "Content-Transfer-Encoding: quoted-printable\r\n")
	fmt.Fprintf(&msg, "Content-Disposition: inline\r\n\r\n")
	msg.WriteString(plain)

	// -- text/html
	fmt.Fprintf(&msg, "--%s\r\n", boundary)
	fmt.Fprintf(&msg, "Content-Type: text/html; charset=\"utf-8\"\r\n")
	fmt.Fprintf(&msg, "Content-Transfer-Encoding: quoted-printable\r\n")
	fmt.Fprintf(&msg, "Content-Disposition: inline\r\n\r\n")
	msg.WriteString(htmlQP)

	// end
	fmt.Fprintf(&msg, "--%s--\r\n", boundary)

	// ----- SMTP: 465 TLS | 587 STARTTLS -----
	if s.cfg.Port == 465 {
		return s.sendImplicitTLS(addr, msg.String())
	}

	conn, err := (&net.Dialer{Timeout: s.cfg.Timeout}).Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		return err
	}
	defer c.Quit()

	// EHLO 1 lần TRƯỚC khi kiểm tra STARTTLS
	if err := c.Hello("localhost"); err != nil {
		return fmt.Errorf("ehlo failed: %w", err)
	}

	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsCfg := &tls.Config{ServerName: s.cfg.Host, MinVersion: tls.VersionTLS12}
		if err := c.StartTLS(tlsCfg); err != nil {
			return fmt.Errorf("starttls failed: %w", err)
		}
		// KHÔNG Hello lại
	} else if s.cfg.Port == 587 {
		return fmt.Errorf("server does not support STARTTLS on port 587")
	}

	if s.cfg.Username != "" {
		if err := c.Auth(s.smtpAuth()); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}
	}

	if err := c.Mail(s.cfg.From); err != nil {
		return err
	}
	for _, rcpt := range s.cfg.To {
		if err := c.Rcpt(rcpt); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte(msg.String())); err != nil {
		return err
	}
	return w.Close()
}

func (s *SMTPSender) sendImplicitTLS(addr string, msg string) error {
	tlsCfg := &tls.Config{ServerName: s.cfg.Host, MinVersion: tls.VersionTLS12}
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: s.cfg.Timeout}, "tcp", addr, tlsCfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		return err
	}
	defer c.Quit()

	// EHLO 1 lần (sau TLS handshake)
	if err := c.Hello("localhost"); err != nil {
		return fmt.Errorf("ehlo failed: %w", err)
	}

	// AUTH
	if s.cfg.Username != "" {
		if err := c.Auth(s.smtpAuth()); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}
	}

	// MAIL FROM / RCPT TO / DATA
	if err := c.Mail(s.cfg.From); err != nil {
		return err
	}
	for _, rcpt := range s.cfg.To {
		if err := c.Rcpt(rcpt); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte(msg)); err != nil {
		return err
	}
	return w.Close()
}

func (s *SMTPSender) smtpAuth() smtp.Auth {
	// PLAIN auth: phổ biến với Gmail/App Password, Mailtrap...
	return smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
}
