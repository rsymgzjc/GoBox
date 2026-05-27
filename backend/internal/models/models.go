package models

import "time"

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	Name         string           `gorm:"size:120;not null" json:"name"`
	Email        string           `gorm:"size:160;uniqueIndex;not null" json:"email"`
	PasswordHash string           `gorm:"size:255;not null" json:"-"`
	Role         string           `gorm:"size:20;not null;default:user" json:"role"`
	LastLoginAt  *time.Time       `json:"lastLoginAt"`
	CreatedAt    time.Time        `json:"createdAt"`
	UpdatedAt    time.Time        `json:"updatedAt"`
	Preferences  []UserPreference `json:"preferences,omitempty"`
}

type Tool struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:120;not null" json:"name"`
	Slug        string    `gorm:"size:80;uniqueIndex;not null" json:"slug"`
	Category    string    `gorm:"size:80;not null" json:"category"`
	Description string    `gorm:"size:255;not null" json:"description"`
	InputHint   string    `gorm:"size:255" json:"inputHint"`
	OutputHint  string    `gorm:"size:255" json:"outputHint"`
	IsFeatured  bool      `gorm:"not null;default:false" json:"isFeatured"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ToolUsage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ToolSlug  string    `gorm:"size:80;index;not null" json:"toolSlug"`
	UserID    *uint     `gorm:"index" json:"userId"`
	Success   bool      `gorm:"not null" json:"success"`
	LatencyMS int64     `gorm:"not null" json:"latencyMs"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserPreference struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"userId"`
	Key       string    `gorm:"size:80;not null" json:"key"`
	Value     string    `gorm:"size:255;not null" json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type EmailVerificationCode struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Email      string     `gorm:"size:160;index;not null" json:"email"`
	Purpose    string     `gorm:"size:40;index;not null" json:"purpose"`
	CodeHash   string     `gorm:"size:64;not null" json:"-"`
	ExpiresAt  time.Time  `gorm:"index;not null" json:"expiresAt"`
	ConsumedAt *time.Time `json:"consumedAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

func DefaultTools() []Tool {
	return []Tool{
		{Name: "JSON 格式化", Slug: "json-format", Category: "编码转换", Description: "格式化与校验 JSON", InputHint: "粘贴 JSON 文本", OutputHint: "美化后的 JSON", IsFeatured: true},
		{Name: "JSON 压缩", Slug: "json-minify", Category: "编码转换", Description: "压缩 JSON 并移除空白", InputHint: "粘贴 JSON 文本", OutputHint: "压缩后的 JSON"},
		{Name: "Base64 编码", Slug: "base64-encode", Category: "编码转换", Description: "将文本转为 Base64", InputHint: "输入任意文本", OutputHint: "Base64 字符串", IsFeatured: true},
		{Name: "Base64 解码", Slug: "base64-decode", Category: "编码转换", Description: "将 Base64 转回文本", InputHint: "输入 Base64 字符串", OutputHint: "解码结果"},
		{Name: "URL 编码", Slug: "url-encode", Category: "网络开发", Description: "对 URL 参数进行编码", InputHint: "输入待编码文本", OutputHint: "编码结果"},
		{Name: "URL 解码", Slug: "url-decode", Category: "网络开发", Description: "解码 URL 参数", InputHint: "输入编码后的内容", OutputHint: "解码结果"},
		{Name: "MD5 计算", Slug: "hash-md5", Category: "安全加密", Description: "生成 MD5 摘要", InputHint: "输入任意文本", OutputHint: "MD5 值"},
		{Name: "SHA256 计算", Slug: "hash-sha256", Category: "安全加密", Description: "生成 SHA256 摘要", InputHint: "输入任意文本", OutputHint: "SHA256 值"},
		{Name: "Unix 时间戳转日期", Slug: "unix-to-time", Category: "日期时间", Description: "将 10 位或 13 位时间戳转为时间", InputHint: "输入时间戳", OutputHint: "本地时间"},
		{Name: "日期转 Unix 时间戳", Slug: "time-to-unix", Category: "日期时间", Description: "将日期转为 Unix 时间戳", InputHint: "输入 RFC3339 时间", OutputHint: "时间戳"},
		{Name: "UUID 生成", Slug: "uuid-generate", Category: "开发辅助", Description: "生成随机 UUID", InputHint: "无需输入", OutputHint: "UUID 列表", IsFeatured: true},
		{Name: "Slug 生成", Slug: "slugify", Category: "开发辅助", Description: "将标题转为 URL 友好 slug", InputHint: "输入标题", OutputHint: "slug"},
	}
}
