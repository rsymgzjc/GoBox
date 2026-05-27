package database

import (
	"context"
	"fmt"
	"time"

	"gobox/backend/internal/config"
	"gobox/backend/internal/models"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch cfg.Database.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.Database.DSN)
	default:
		dialector = sqlite.Open(cfg.Database.DSN)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := waitUntilReady(db, cfg, logger); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Tool{},
		&models.ToolUsage{},
		&models.UserPreference{},
		&models.EmailVerificationCode{},
	); err != nil {
		return nil, err
	}

	if err := seed(db, logger); err != nil {
		return nil, err
	}

	return db, nil
}

func waitUntilReady(db *gorm.DB, cfg *config.Config, logger *zap.Logger) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	attempts := 1
	if cfg.Database.Driver == "mysql" {
		attempts = 30
	}

	var lastErr error
	for i := 1; i <= attempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err = sqlDB.PingContext(ctx)
		cancel()
		if err == nil {
			return nil
		}

		lastErr = err
		if i < attempts {
			logger.Warn("database not ready yet",
				zap.String("driver", cfg.Database.Driver),
				zap.Int("attempt", i),
				zap.Error(err),
			)
			time.Sleep(2 * time.Second)
		}
	}

	return fmt.Errorf("database ping failed after %d attempts: %w", attempts, lastErr)
}

func seed(db *gorm.DB, logger *zap.Logger) error {
	tools := models.DefaultTools()
	for _, tool := range tools {
		var existing models.Tool
		if err := db.Where("slug = ?", tool.Slug).First(&existing).Error; err == nil {
			continue
		}
		if err := db.Create(&tool).Error; err != nil {
			return err
		}
	}

	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("admin123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.User{
		Name:         "GoBox 管理员",
		Email:        "admin@gobox.local",
		PasswordHash: string(passwordHash),
		Role:         models.RoleAdmin,
		LastLoginAt:  ptrTime(time.Now()),
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	logger.Info("seeded default admin", zap.String("email", admin.Email), zap.String("password", "admin123456"))
	return nil
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func DSNExample() string {
	return fmt.Sprintf("%s", "gobox:gobox123@tcp(127.0.0.1:3306)/gobox?charset=utf8mb4&parseTime=True&loc=Local")
}
