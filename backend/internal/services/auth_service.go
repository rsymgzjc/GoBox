package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"gobox/backend/internal/config"
	"gobox/backend/internal/models"
	jwtpkg "gobox/backend/pkg/jwt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const registerCodePurpose = "register"

type AuthService struct {
	db     *gorm.DB
	cfg    *config.Config
	logger *zap.Logger
	mailer *Mailer
}

type RegisterPayload struct {
	Name     string `json:"name" validate:"required,min=2,max=120"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
	Code     string `json:"code" validate:"required,len=6"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type SendRegisterCodePayload struct {
	Email string `json:"email" validate:"required,email"`
}

type SendRegisterCodeResult struct {
	CooldownSeconds int    `json:"cooldownSeconds"`
	ExpiresInMinute int    `json:"expiresInMinutes"`
	PreviewCode     string `json:"previewCode,omitempty"`
}

func NewAuthService(db *gorm.DB, cfg *config.Config, logger *zap.Logger) *AuthService {
	return &AuthService{
		db:     db,
		cfg:    cfg,
		logger: logger,
		mailer: NewMailer(cfg, logger),
	}
}

func (s *AuthService) SendRegisterCode(input SendRegisterCodePayload) (*SendRegisterCodeResult, error) {
	var existing models.User
	if err := s.db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		return nil, errors.New("该邮箱已注册")
	}

	if err := s.ensureCooldown(input.Email); err != nil {
		return nil, err
	}

	code, err := generateVerificationCode()
	if err != nil {
		return nil, err
	}

	record := models.EmailVerificationCode{
		Email:     input.Email,
		Purpose:   registerCodePurpose,
		CodeHash:  hashVerificationCode(s.cfg.Auth.JWTSecret, registerCodePurpose, input.Email, code),
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.Auth.VerificationCodeTTLMinutes) * time.Minute),
	}

	if err := s.db.Create(&record).Error; err != nil {
		return nil, err
	}

	if err := s.mailer.SendRegisterCode(input.Email, code); err != nil {
		_ = s.db.Delete(&record).Error
		return nil, err
	}

	result := &SendRegisterCodeResult{
		CooldownSeconds: s.cfg.Auth.VerificationCodeCooldownSecond,
		ExpiresInMinute: s.cfg.Auth.VerificationCodeTTLMinutes,
	}
	if !s.cfg.SMTP.Enabled && s.cfg.App.Env != "production" {
		result.PreviewCode = code
	}
	return result, nil
}

func (s *AuthService) Register(input RegisterPayload) (*models.User, string, error) {
	var user *models.User

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var existing models.User
		if err := tx.Where("email = ?", input.Email).First(&existing).Error; err == nil {
			return errors.New("该邮箱已注册")
		}

		if err := s.consumeRegisterCode(tx, input.Email, input.Code); err != nil {
			return err
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		newUser := models.User{
			Name:         input.Name,
			Email:        input.Email,
			PasswordHash: string(hash),
			Role:         models.RoleUser,
		}

		if err := tx.Create(&newUser).Error; err != nil {
			return err
		}

		user = &newUser
		return nil
	})
	if err != nil {
		return nil, "", err
	}

	token, err := s.issueToken(*user)
	return user, token, err
}

func (s *AuthService) Login(input LoginPayload) (*models.User, string, error) {
	var user models.User
	if err := s.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return nil, "", errors.New("邮箱或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, "", errors.New("邮箱或密码错误")
	}

	now := time.Now()
	user.LastLoginAt = &now
	_ = s.db.Save(&user).Error

	token, err := s.issueToken(user)
	return &user, token, err
}

func (s *AuthService) Profile(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Preferences").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) SavePreferences(userID uint, preferences map[string]string) error {
	for key, value := range preferences {
		pref := models.UserPreference{UserID: userID, Key: key}
		if err := s.db.Where(&pref).Assign(models.UserPreference{Value: value}).FirstOrCreate(&pref).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *AuthService) issueToken(user models.User) (string, error) {
	duration := time.Duration(s.cfg.Auth.TokenDuration) * time.Hour
	return jwtpkg.Generate(s.cfg.Auth.JWTSecret, duration, user.ID, user.Email, user.Role)
}

func (s *AuthService) ensureCooldown(email string) error {
	cooldown := s.cfg.Auth.VerificationCodeCooldownSecond
	if cooldown <= 0 {
		return nil
	}

	var record models.EmailVerificationCode
	err := s.db.Where(
		"email = ? AND purpose = ? AND created_at >= ?",
		email,
		registerCodePurpose,
		time.Now().Add(-time.Duration(cooldown)*time.Second),
	).Order("created_at desc").First(&record).Error
	if err != nil {
		return nil
	}

	remain := cooldown - int(time.Since(record.CreatedAt).Seconds())
	if remain < 1 {
		remain = 1
	}
	return fmt.Errorf("验证码发送过于频繁，请 %d 秒后再试", remain)
}

func (s *AuthService) consumeRegisterCode(tx *gorm.DB, email, code string) error {
	var record models.EmailVerificationCode
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("email = ? AND purpose = ? AND consumed_at IS NULL", email, registerCodePurpose).
		Order("created_at desc").
		First(&record).Error; err != nil {
		return errors.New("请先获取邮箱验证码")
	}

	if time.Now().After(record.ExpiresAt) {
		return errors.New("验证码已过期，请重新获取")
	}

	expected := hashVerificationCode(s.cfg.Auth.JWTSecret, registerCodePurpose, email, code)
	if record.CodeHash != expected {
		return errors.New("验证码错误")
	}

	now := time.Now()
	result := tx.Model(&models.EmailVerificationCode{}).
		Where("id = ? AND consumed_at IS NULL", record.ID).
		Update("consumed_at", &now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("验证码已失效，请重新获取")
	}

	return nil
}

func generateVerificationCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func hashVerificationCode(secret, purpose, email, code string) string {
	sum := sha256.Sum256([]byte(secret + "|" + purpose + "|" + email + "|" + code))
	return hex.EncodeToString(sum[:])
}
