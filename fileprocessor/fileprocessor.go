package fileprocessor

import (
	"strconv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func createFile(filePath, content string) error {
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

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}

	if !strings.HasSuffix(content, "\n") && content != "" {
		content += "\n"
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error writing file %s: %w", filePath, err)
	}

	fmt.Printf("Written: %s\n", filePath)
	return nil
}

func processMarkdownBlocks(lines []string, safe bool) error {
	inCodeBlock := false
	var currentBlock *CodeBlock
	unknownFileCounter := 0

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			if inCodeBlock {
				if len(currentBlock.lines) > 0 {
					if currentBlock.filePath == "" {
						unknownFileCounter += 1
						currentBlock.filePath = "unknown" + strconv.Itoa(unknownFileCounter)
					}
					if err := currentBlock.write(safe); err != nil {
						return err
					}
				}
				inCodeBlock = false
				currentBlock = nil
			} else {
				inCodeBlock = true
				blockHeader := strings.TrimSpace(strings.Replace(line, "```", "", 1))
				currentBlock = &CodeBlock{blockHeader: blockHeader}
			}
		} else if inCodeBlock {
			currentBlock.lines = append(currentBlock.lines, line)
		}
	}

	if inCodeBlock && currentBlock != nil && len(currentBlock.lines) > 0 {
		if currentBlock.filePath == "" {
			return fmt.Errorf("content lines exist but no filepath was provided")
		}
		if err := currentBlock.write(safe); err != nil {
			return err
		}
		return fmt.Errorf("incomplete code block: file %s was written but no closing backticks found", currentBlock.filePath)
	}

	return nil
}

func ProcessCodeBlocks(response string, safe bool) error {
	lines := strings.Split(response, "\n")
	return processMarkdownBlocks(lines, safe)
}
