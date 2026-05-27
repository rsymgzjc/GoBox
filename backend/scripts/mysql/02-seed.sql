USE `gobox`;

INSERT INTO `tools` (`name`, `slug`, `category`, `description`, `input_hint`, `output_hint`, `is_featured`)
VALUES
  ('JSON 格式化', 'json-format', '编码转换', '格式化与校验 JSON', '粘贴 JSON 文本', '美化后的 JSON', 1),
  ('JSON 压缩', 'json-minify', '编码转换', '压缩 JSON 并移除空白', '粘贴 JSON 文本', '压缩后的 JSON', 0),
  ('Base64 编码', 'base64-encode', '编码转换', '将文本转为 Base64', '输入任意文本', 'Base64 字符串', 1),
  ('Base64 解码', 'base64-decode', '编码转换', '将 Base64 转回文本', '输入 Base64 字符串', '解码结果', 0),
  ('URL 编码', 'url-encode', '网络开发', '对 URL 参数进行编码', '输入待编码文本', '编码结果', 0),
  ('URL 解码', 'url-decode', '网络开发', '解码 URL 参数', '输入编码后的内容', '解码结果', 0),
  ('MD5 计算', 'hash-md5', '安全加密', '生成 MD5 摘要', '输入任意文本', 'MD5 值', 0),
  ('SHA256 计算', 'hash-sha256', '安全加密', '生成 SHA256 摘要', '输入任意文本', 'SHA256 值', 0),
  ('Unix 时间戳转日期', 'unix-to-time', '日期时间', '将 10 位或 13 位时间戳转为时间', '输入时间戳', '本地时间', 0),
  ('日期转 Unix 时间戳', 'time-to-unix', '日期时间', '将日期转为 Unix 时间戳', '输入 RFC3339 时间', '时间戳', 0),
  ('UUID 生成', 'uuid-generate', '开发辅助', '生成随机 UUID', '无需输入', 'UUID 列表', 1),
  ('Slug 生成', 'slugify', '开发辅助', '将标题转为 URL 友好 slug', '输入标题', 'slug', 0)
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `category` = VALUES(`category`),
  `description` = VALUES(`description`),
  `input_hint` = VALUES(`input_hint`),
  `output_hint` = VALUES(`output_hint`),
  `is_featured` = VALUES(`is_featured`);

INSERT INTO `users` (`name`, `email`, `password_hash`, `role`, `last_login_at`)
VALUES
  ('GoBox 管理员', 'admin@gobox.local', '$2a$10$cLhiYLL//SAoXO06tJbB2uG.AQW1l4uXu00kyHnII1rwvuYYGUtiG', 'admin', NOW())
ON DUPLICATE KEY UPDATE
  `name` = VALUES(`name`),
  `role` = VALUES(`role`);
