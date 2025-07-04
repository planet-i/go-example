package main

import (
	"testing"
)

func TestGetFilenameFromUrl(t *testing.T) {
	tests := []struct {
		name         string // 测试用例名称
		url          string // 输入URL
		expectedFull string // 期望的完整文件名
		expectedName string // 期望的文件名（不含扩展名）
	}{
		{
			name:         "多级路径带扩展名",
			url:          "https://cdn.example.net/a/b/c/document.pdf",
			expectedFull: "document.pdf",
			expectedName: "document",
		},
		{
			name:         "无扩展名文件",
			url:          "https://example.com/data/README",
			expectedFull: "README",
			expectedName: "README",
		},
		{
			name:         "文件名含多点",
			url:          "https://example.com/files/archive.tar.gz",
			expectedFull: "archive.tar.gz",
			expectedName: "archive", // 注意：函数设计只取第一个点前的部分
		},
		{
			name:         "URL以斜杠结尾",
			url:          "https://example.com/directory/",
			expectedFull: "",
			expectedName: "",
		},
		{
			name:         "空URL",
			url:          "",
			expectedFull: "",
			expectedName: "",
		},
		{
			name:         "带查询参数",
			url:          "https://example.com/download?file=report.docx",
			expectedFull: "download?file=report.docx",
			expectedName: "download?file=report", // 注意：这里没有处理查询参数
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			full, name := GetFilenameFromUrl(tt.url)
			if full != tt.expectedFull || name != tt.expectedName {
				t.Errorf(
					"GetFilenameFromUrl(%q) = (%q, %q); want (%q, %q)",
					tt.url,
					full,
					name,
					tt.expectedFull,
					tt.expectedName,
				)
			}
		})
	}
}
