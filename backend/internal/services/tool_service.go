package services

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gobox/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ToolService struct {
	db *gorm.DB
}

type RunToolInput struct {
	Input   string            `json:"input"`
	Options map[string]string `json:"options"`
}

type ToolResult struct {
	ToolSlug  string                 `json:"toolSlug"`
	Input     string                 `json:"input"`
	Output    string                 `json:"output"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
	Success   bool                   `json:"success"`
	LatencyMS int64                  `json:"latencyMs"`
}

func NewToolService(db *gorm.DB) *ToolService {
	return &ToolService{db: db}
}

func (s *ToolService) List(category string) ([]models.Tool, error) {
	var tools []models.Tool
	query := s.db.Order("is_featured desc, category asc, name asc")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	return tools, query.Find(&tools).Error
}

func (s *ToolService) Run(slug string, input RunToolInput, userID *uint) (*ToolResult, error) {
	start := time.Now()
	output, meta, err := executeTool(slug, input)
	result := &ToolResult{
		ToolSlug: slug,
		Input:    input.Input,
		Output:   output,
		Meta:     meta,
		Success:  err == nil,
	}
	result.LatencyMS = time.Since(start).Milliseconds()

	usage := models.ToolUsage{
		ToolSlug:  slug,
		UserID:    userID,
		Success:   result.Success,
		LatencyMS: result.LatencyMS,
	}
	_ = s.db.Create(&usage).Error

	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ToolService) Summary() (map[string]interface{}, error) {
	var toolCount int64
	var usageCount int64
	var userCount int64

	if err := s.db.Model(&models.Tool{}).Count(&toolCount).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.ToolUsage{}).Count(&usageCount).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return nil, err
	}

	type topTool struct {
		ToolSlug string `json:"toolSlug"`
		Count    int64  `json:"count"`
	}
	var top []topTool
	if err := s.db.Model(&models.ToolUsage{}).
		Select("tool_slug, count(*) as count").
		Group("tool_slug").
		Order("count desc").
		Limit(5).
		Scan(&top).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"toolCount":  toolCount,
		"usageCount": usageCount,
		"userCount":  userCount,
		"topTools":   top,
	}, nil
}

func executeTool(slug string, input RunToolInput) (string, map[string]interface{}, error) {
	switch slug {
	case "json-format":
		var any interface{}
		if err := json.Unmarshal([]byte(input.Input), &any); err != nil {
			return "", nil, err
		}
		buf, err := json.MarshalIndent(any, "", "  ")
		return string(buf), map[string]interface{}{"length": len(buf)}, err
	case "json-minify":
		var any interface{}
		if err := json.Unmarshal([]byte(input.Input), &any); err != nil {
			return "", nil, err
		}
		buf, err := json.Marshal(any)
		return string(buf), map[string]interface{}{"length": len(buf)}, err
	case "base64-encode":
		out := base64.StdEncoding.EncodeToString([]byte(input.Input))
		return out, map[string]interface{}{"length": len(out)}, nil
	case "base64-decode":
		buf, err := base64.StdEncoding.DecodeString(strings.TrimSpace(input.Input))
		return string(buf), map[string]interface{}{"length": len(buf)}, err
	case "url-encode":
		return url.QueryEscape(input.Input), nil, nil
	case "url-decode":
		out, err := url.QueryUnescape(input.Input)
		return out, nil, err
	case "hash-md5":
		sum := md5.Sum([]byte(input.Input))
		return hex.EncodeToString(sum[:]), nil, nil
	case "hash-sha256":
		sum := sha256.Sum256([]byte(input.Input))
		return hex.EncodeToString(sum[:]), nil, nil
	case "unix-to-time":
		text := strings.TrimSpace(input.Input)
		v, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return "", nil, err
		}
		if len(text) == 13 {
			v = v / 1000
		}
		t := time.Unix(v, 0)
		return t.Format(time.RFC3339), map[string]interface{}{"timezone": t.Location().String()}, nil
	case "time-to-unix":
		t, err := time.Parse(time.RFC3339, strings.TrimSpace(input.Input))
		if err != nil {
			return "", nil, err
		}
		return strconv.FormatInt(t.Unix(), 10), nil, nil
	case "uuid-generate":
		count := 1
		if raw := input.Options["count"]; raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err == nil && parsed > 0 && parsed <= 20 {
				count = parsed
			}
		}
		values := make([]string, 0, count)
		for i := 0; i < count; i++ {
			values = append(values, uuid.NewString())
		}
		return strings.Join(values, "\n"), map[string]interface{}{"count": count}, nil
	case "slugify":
		return slugify(input.Input), nil, nil
	default:
		return "", nil, errors.New("tool not implemented")
	}
}

func slugify(input string) string {
	lower := strings.ToLower(strings.TrimSpace(input))
	space := strings.ReplaceAll(lower, "_", "-")
	space = strings.ReplaceAll(space, " ", "-")
	reg := regexp.MustCompile(`[^a-z0-9\-]+`)
	space = reg.ReplaceAllString(space, "")
	space = strings.Trim(space, "-")
	return strings.ReplaceAll(space, "--", "-")
}

func ToolLabels() string {
	return fmt.Sprintf("%d", len(models.DefaultTools()))
}

func ExecuteForTest(slug string, input RunToolInput) (string, map[string]interface{}, error) {
	return executeTool(slug, input)
}
