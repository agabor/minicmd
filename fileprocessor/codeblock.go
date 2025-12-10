package fileprocessor

import (
	"regexp"
	"strings"
)

type CodeBlock struct {
	blockHeader string
	filePath   string
	lines      []string
}

func extractFilenameFromComment(line string) string {
	patterns := []string{
		`^\s*//\s*(.+?)(?:\s*//.*)?$`,
		`^\s*#\s*(.+?)(?:\s*#.*)?$`,
		`^\s*#\s*//\s*(.+?)(?:\s*#.*)?$`,
		`^\s*/\*\s*(.+?)\s*\*/$`,
		`^\s*--\s*(.+?)(?:\s*--.*)?$`,
		`^\s*<!--\s*(.+?)\s*-->$`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			filename := strings.TrimSpace(matches[1])
			if strings.Contains(filename, "!") {
				continue
			}
			filename = regexp.MustCompile(`\s*\*+/$`).ReplaceAllString(filename, "")
			return filename
		}
	}
	return ""
}

func (cb *CodeBlock) resolveFilePath() {
	if cb.filePath == "" {
		cb.filePath = cb.blockHeader
	}
}

func (cb *CodeBlock) extractFilePathFromFirstLine() {
	if len(cb.lines) == 0 {
		return
	}

	extractedPath := extractFilenameFromComment(cb.lines[0])
	if extractedPath != "" {
		cb.filePath = extractedPath
		cb.lines = cb.lines[1:]
	}
}

func (cb *CodeBlock) content() string {
	return strings.Join(cb.lines, "\n")
}

func (cb *CodeBlock) write(safe bool) error {
	cb.extractFilePathFromFirstLine()
	cb.resolveFilePath()

	finalPath := cb.filePath
	if safe {
		finalPath += ".new"
	}

	return createFile(finalPath, cb.content())
}
