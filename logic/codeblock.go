package logic

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var unknownFileCounter = 0

const BlockDelimiter = "``" + "``"

type CodeBlock struct {
	Path    string
	Content string
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

func linesToCodeBlock(lines []string) CodeBlock {
	filePath := ""
	lineIndex := 0

	for lineIndex < len(lines) && strings.TrimSpace(lines[lineIndex]) == "" {
		lineIndex++
	}

	if lineIndex < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[lineIndex]), "#!") {
		lineIndex++
	}

	if lineIndex < len(lines) {
		extractedPath := extractFilenameFromComment(lines[lineIndex])
		if extractedPath != "" {
			filePath = extractedPath
			lines = append(lines[:lineIndex], lines[lineIndex+1:]...)
		}
	}

	if filePath == "" {
		unknownFileCounter += 1
		filePath = "unknown" + strconv.Itoa(unknownFileCounter)
	}

	if strings.HasPrefix(filePath, "/") {
		relPath := filePath[1:]
		if _, err := os.Stat(relPath); err == nil {
			filePath = relPath
		} else {
			if _, err := os.Stat(filePath); err == nil {

			} else {
				filePath = relPath
			}
		}
	}

	return CodeBlock{Path: filePath, Content: joinLines(lines)}
}

func joinLines(lines []string) string {
	return strings.Join(lines, "\n")
}

func AsCodeBlock(path string, content string) string {
	return joinLines([]string{BlockDelimiter, "//" + path, content, BlockDelimiter})
}

func (cb *CodeBlock) Write(safe bool) error {
	filePath := cb.Path
	if safe {
		filePath += ".new"
	}
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}

	if err := os.WriteFile(filePath, []byte(cb.Content), 0644); err != nil {
		return fmt.Errorf("error writing file %s: %w", filePath, err)
	}

	fmt.Printf("Written: %s\n", filePath)
	return nil
}
