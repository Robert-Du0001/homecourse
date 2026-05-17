package utils

import (
	"path/filepath"
	"strings"
)

// 判断文件是否可在浏览器中直接预览
//
//	@param filename 文件名
//	@return bool 是否可直接预览
func IsInlinePreview(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf",
		".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp", ".ico",
		".txt", ".html", ".htm", ".css", ".js", ".json", ".xml", ".md",
		".mp4", ".webm", ".ogg", ".mp3", ".wav":
		return true
	}
	return false
}
