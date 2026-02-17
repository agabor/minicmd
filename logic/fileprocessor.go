package logic

import (
	"strings"
)

func ParseCodeBlocks(response string) []CodeBlock {
	lines := strings.Split(response, "\n")
	inCodeBlock := false
	var currentBlock *CodeBlock
	var codeBlocks = make([]CodeBlock, 0)

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "````") {
			if inCodeBlock {
				if len(currentBlock.lines) > 0 {
					codeBlocks = append(codeBlocks, *currentBlock)
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
		codeBlocks = append(codeBlocks, *currentBlock)
	}

	return codeBlocks
}

func ProcessCodeBlocks(response string, safe bool) error {
	var err error
	for _, codeBlock := range ParseCodeBlocks(response) {
		err = codeBlock.write(safe)
	}
	return err
}
