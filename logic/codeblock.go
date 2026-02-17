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
	blockHeader string
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

func (cb *CodeBlock) getFilePath() string {
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

	return filePath

}

func (cb *CodeBlock) getContent() string {
	return BlockDelimiter + "\n" + strings.Join(cb.lines, "\n") + "\n" + BlockDelimiter
}

func (cb *CodeBlock) write(safe bool) error {
	filePath := cb.getFilePath()
	if safe {
		filePath += ".new"
	}
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}

	if err := os.WriteFile(filePath, []byte(strings.Join(cb.lines, "\n")), 0644); err != nil {
		return fmt.Errorf("error writing file %s: %w", filePath, err)
	}

	fmt.Printf("Written: %s\n", filePath)
	return nil
}
