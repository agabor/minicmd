package fileprocessor

import (
	"strconv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var unknownFileCounter = 0

type CodeBlock struct {
	blockHeader string
	filePath    string
	lines       []string
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

func (cb *CodeBlock) getFilePath(safe bool) string {
	filePath := ""
	lineIndex := 0

	for lineIndex < len(cb.lines) && strings.TrimSpace(cb.lines[lineIndex]) == "" {
		lineIndex++
	}

	if lineIndex < len(cb.lines) && strings.HasPrefix(strings.TrimSpace(cb.lines[lineIndex]), "#!") {
		lineIndex++
	}

	if lineIndex < len(cb.lines) {
		extractedPath := extractFilenameFromComment(cb.lines[lineIndex])
		if extractedPath != "" {
			filePath = extractedPath
			cb.lines = append(cb.lines[:lineIndex], cb.lines[lineIndex+1:]...)
		}
	}

	if filePath == "" {
		filePath = cb.blockHeader
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

	if safe {
		filePath += ".new"
	}
	return filePath

}

func (cb *CodeBlock) write(safe bool) error {

	cb.filePath = cb.getFilePath(safe)

	dir := filepath.Dir(cb.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}

	if err := os.WriteFile(cb.filePath, []byte(strings.Join(cb.lines, "\n")), 0644); err != nil {
		return fmt.Errorf("error writing file %s: %w", cb.filePath, err)
	}

	fmt.Printf("Written: %s\n", cb.filePath)
	return nil
}