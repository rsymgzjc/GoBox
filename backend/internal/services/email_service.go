package services

import (
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"

	"gobox/backend/internal/config"

	"go.uber.org/zap"
)

type Mailer struct {
	cfg    *config.Config
	logger *zap.Logger
}

func NewMailer(cfg *config.Config, logger *zap.Logger) *Mailer {
	return &Mailer{cfg: cfg, logger: logger}
}

func (m *Mailer) SendRegisterCode(to, code string) error {
	if !m.cfg.SMTP.Enabled {
		if m.cfg.App.Env == "production" {
			return errors.New("邮件服务未配置，无法发送验证码")
		}
		m.logger.Info("email sending disabled, using preview code",
			zap.String("email", to),
			zap.String("code", code),
		)
		return nil
	}

	fromAddress, err := resolvedFromAddress(m.cfg.SMTP.From, m.cfg.SMTP.Username)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", m.cfg.SMTP.Username, m.cfg.SMTP.Password, m.cfg.SMTP.Host)
	subject := "GoBox 注册验证码"
	body := fmt.Sprintf("您正在注册 GoBox，验证码为：%s\n验证码 %d 分钟内有效。", code, m.cfg.Auth.VerificationCodeTTLMinutes)
	message := strings.Join([]string{
		fmt.Sprintf("From: %s", m.cfg.SMTP.From),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	return smtp.SendMail(m.cfg.SMTP.Address(), auth, fromAddress, []string{to}, []byte(message))
}

func resolvedFromAddress(from, fallback string) (string, error) {
	if from != "" {
		addr, err := mail.ParseAddress(from)
		if err == nil {
			return addr.Address, nil
		}
	}
	if fallback != "" {
		return fallback, nil
	}
	return "", errors.New("未配置有效的发件人地址")
}
