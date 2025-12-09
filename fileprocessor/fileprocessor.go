package fileprocessor

import (
        "strconv"
        "fmt"
        "os"
        "path/filepath"
        "regexp"
        "strings"
)

func extractFilenameFromComment(line string) string {
        patterns := []string{
                `^\s*//\s*(.+?)(?:\s*//.*)?$`,     // // filename
                `^\s*#\s*(.+?)(?:\s*#.*)?$`,       // # filename
                `^\s*#\s*//\s*(.+?)(?:\s*#.*)?$`,  // # // filename
                `^\s*/\*\s*(.+?)\s*\*/$`,          // /* filename */
                `^\s*--\s*(.+?)(?:\s*--.*)?$`,     // -- filename (SQL)
                `^\s*<!--\s*(.+?)\s*-->$`,         // <!-- filename --> (HTML)
        }

        for _, pattern := range patterns {
                re := regexp.MustCompile(pattern)
                matches := re.FindStringSubmatch(line)
                if len(matches) > 1 {
                        filename := strings.TrimSpace(matches[1])
                        if strings.Contains(filename, "!") {
                                continue
                        }
                        // Remove any trailing comment markers
                        filename = regexp.MustCompile(`\s*\*+/$`).ReplaceAllString(filename, "")
                        return filename
                }
        }
        return ""
}

func createFile(filePath, content string) error {
        // Handle absolute paths
        if strings.HasPrefix(filePath, "/") {
                // First check if file exists as absolute path within working directory
                relPath := filePath[1:] // Remove leading slash
                if _, err := os.Stat(relPath); err == nil {
                        filePath = relPath
                } else {
                        // Try using absolute path if file exists
                        if _, err := os.Stat(filePath); err == nil {
                                // Use absolute path
                        } else {
                                // Default to relative path if file doesn't exist
                                filePath = relPath
                        }
                }
        }

        // Create directory if it doesn't exist
        dir := filepath.Dir(filePath)
        if err := os.MkdirAll(dir, 0755); err != nil {
                return fmt.Errorf("error creating directory %s: %w", dir, err)
        }

        // Write content to file
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
        filePath := ""
        blockHeader := ""
        unknownFileCounter := 0
        var contentLines []string

        for _, line := range lines {
                if strings.Contains(line, "```") {
                        if inCodeBlock {
                                // End of code block - create file
                                if len(contentLines) > 0 {
                                        content := strings.Join(contentLines, "\n")
                                        if filePath == "" {
                                                unknownFileCounter += 1
                                                filePath = "unknown" + strconv.Itoa(unknownFileCounter)
                                        }
                                        finalPath := filePath
                                        if safe {
                                                finalPath += ".new"
                                        }
                                        if err := createFile(finalPath, content); err != nil {
                                                return err
                                        }
                                }
                                inCodeBlock = false
                                filePath = ""
                                blockHeader = ""
                                contentLines = nil
                        } else {
                                // Start of code block
                                inCodeBlock = true
                                blockHeader = strings.TrimSpace(strings.Replace(line, "```", "", 1))
                        }
                } else if inCodeBlock {
                        if filePath == "" {
                                // Check if this line contains the filename
                                extractedPath := extractFilenameFromComment(line)
                                if extractedPath != "" {
                                        filePath = extractedPath
                                        continue
                                } else {
                                        filePath = blockHeader
                                }
                        }
                        contentLines = append(contentLines, line)
                } else {
                        inCodeBlock = true
                }
        }

        // Check for leftover content after processing all lines
        if len(contentLines) > 0 {
                if filePath == "" {
                        return fmt.Errorf("content lines exist but no filepath was provided")
                }
                content := strings.Join(contentLines, "\n")
                finalPath := filePath
                if safe {
                        finalPath += ".new"
                }
                if err := createFile(finalPath, content); err != nil {
                        return err
                }
                return fmt.Errorf("incomplete code block: file %s was written but no closing backticks found", filePath)
        }

        return nil
}

func ProcessCodeBlocks(response string, safe bool) error {
        lines := strings.Split(response, "\n")
        return processMarkdownBlocks(lines, safe)
}