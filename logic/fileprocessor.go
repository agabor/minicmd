package logic

import (
	"strings"
)

func ParseCodeBlocks(response string) []CodeBlock {
	lines := strings.Split(response, "\n")
	var codeBlocks = make([]CodeBlock, 0)
	var lineBuffer = make([]string, 0)
	inBlock := false
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), BlockDelimiter) {
			if inBlock {
				if len(lineBuffer) > 0 {
					codeBlocks = append(codeBlocks, linesToCodeBlock(lineBuffer))
				}
				inBlock = false
				lineBuffer = make([]string, 0)
			} else {
				inBlock = true
			}
		} else if inBlock {
			lineBuffer = append(lineBuffer, line)
		}
	}

	if inBlock && len(lineBuffer) > 0 {
		codeBlocks = append(codeBlocks, linesToCodeBlock(lineBuffer))
	}

	return codeBlocks
}
