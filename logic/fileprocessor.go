package logic

import (
	"strings"
)

func ParseCodeBlocks(response string) []CodeBlock {
	lines := strings.Split(response, "\n")
	var currentBlock *CodeBlock
	var codeBlocks = make([]CodeBlock, 0)

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), BlockDelimiter) {
			if currentBlock != nil {
				if len(currentBlock.lines) > 0 {
					codeBlocks = append(codeBlocks, *currentBlock)
				}
				currentBlock = nil
			} else {
				blockHeader := strings.TrimSpace(strings.Replace(line, BlockDelimiter, "", 1))
				currentBlock = &CodeBlock{blockHeader: blockHeader}
			}
		} else if currentBlock != nil {
			currentBlock.lines = append(currentBlock.lines, line)
		}
	}

	if currentBlock != nil && len(currentBlock.lines) > 0 {
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
