package logic

import (
	"fmt"
	"strings"
)

func ProcessCodeBlocks(response string, safe bool) error {
	lines := strings.Split(response, "\n")
	inCodeBlock := false
	var currentBlock *CodeBlock

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "````") {
			if inCodeBlock {
				if len(currentBlock.lines) > 0 {
					if err := currentBlock.write(safe); err != nil {
						return err
					}
				}
				inCodeBlock = false
				currentBlock = nil
			} else {
				inCodeBlock = true
				blockHeader := strings.TrimSpace(strings.Replace(line, "````", "", 1))
				currentBlock = &CodeBlock{blockHeader: blockHeader}
			}
		} else if inCodeBlock {
			currentBlock.lines = append(currentBlock.lines, line)
		}
	}

	if inCodeBlock && currentBlock != nil && len(currentBlock.lines) > 0 {
		if err := currentBlock.write(safe); err != nil {
			return err
		}
		return fmt.Errorf("incomplete code block: file %s was written but no closing backticks found", currentBlock.filePath)
	}

	return nil
}
