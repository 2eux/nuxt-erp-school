package rag

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

type FileType string

const (
	FileTypePDF   FileType = "pdf"
	FileTypeDOCX  FileType = "docx"
	FileTypeXLSX  FileType = "xlsx"
	FileTypePPTX  FileType = "pptx"
	FileTypeImage FileType = "image"
	FileTypeText  FileType = "text"
	FileTypeHTML  FileType = "html"
	FileTypeCSV   FileType = "csv"
	FileTypeUnknown FileType = "unknown"
)

type extMapping struct {
	exts     []string
	fileType FileType
}

var extensionMap = []extMapping{
	{[]string{".pdf"}, FileTypePDF},
	{[]string{".docx", ".doc"}, FileTypeDOCX},
	{[]string{".xlsx", ".xls", ".xlsm"}, FileTypeXLSX},
	{[]string{".pptx", ".ppt"}, FileTypePPTX},
	{[]string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp", ".tiff", ".svg"}, FileTypeImage},
	{[]string{".txt", ".md", ".markdown", ".rst"}, FileTypeText},
	{[]string{".html", ".htm"}, FileTypeHTML},
	{[]string{".csv"}, FileTypeCSV},
}

func DetectFileType(filename string) FileType {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, mapping := range extensionMap {
		for _, e := range mapping.exts {
			if e == ext {
				return mapping.fileType
			}
		}
	}
	return FileTypeUnknown
}

func ExtractText(data []byte, filename string) (string, error) {
	fileType := DetectFileType(filename)

	switch fileType {
	case FileTypePDF:
		return extractPDFText(data)
	case FileTypeDOCX:
		return extractDocxText(data)
	case FileTypeXLSX:
		return extractXlsxText(data)
	case FileTypePPTX:
		return extractPptxText(data)
	case FileTypeText, FileTypeHTML, FileTypeCSV:
		return extractPlainText(data)
	case FileTypeImage:
		return extractImageDescription(filename)
	default:
		attempt := extractPlainText(data)
		if attempt != "" && len(attempt) > 20 {
			return attempt, nil
		}
		return "", fmt.Errorf("unsupported file type: %s (.%s)", fileType, filepath.Ext(filename))
	}
}

func extractPDFText(data []byte) (string, error) {
	text := extractPDFSimpleText(data)
	if text != "" {
		return text, nil
	}

	return "[PDF document - text extraction requires pdfcpu or similar library]", nil
}

func extractPDFSimpleText(data []byte) string {
	content := string(data)
	var sb strings.Builder

	inStream := false
	inText := false
	inBT := false

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "stream") && !strings.Contains(trimmed, "endstream") {
			inStream = true
			continue
		}
		if trimmed == "endstream" {
			inStream = false
			inBT = false
			inText = false
			continue
		}
		if inStream && strings.HasPrefix(trimmed, "BT") {
			inBT = true
			continue
		}
		if inBT && trimmed == "ET" {
			inBT = false
			sb.WriteString("\n")
			continue
		}
		if inBT {
			textMatch := extractPDFTextOperators(trimmed)
			if textMatch != "" {
				sb.WriteString(textMatch)
			}
		}
	}

	result := sb.String()
	result = cleanExtractedText(result)

	return result
}

func extractPDFTextOperators(line string) string {
	if strings.Contains(line, "Tj") || strings.Contains(line, "TJ") || strings.Contains(line, "'") || strings.Contains(line, "\"") {
		if idx := strings.Index(line, "("); idx >= 0 {
			if endIdx := findClosingParen(line, idx); endIdx > idx {
				text := line[idx+1 : endIdx]
				text = decodePDFString(text)
				return text
			}
		}
	}
	return ""
}

func findClosingParen(s string, start int) int {
	depth := 0
	for i := start; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 && i > start {
				return i
			}
		case '\\':
			i++
		}
	}
	return -1
}

func decodePDFString(s string) string {
	s = strings.ReplaceAll(s, "\\(", "(")
	s = strings.ReplaceAll(s, "\\)", ")")
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\\r", "\r")
	s = strings.ReplaceAll(s, "\\t", "\t")

	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+3 < len(s) {
			octal := s[i+1 : i+4]
			if val := parseOctal(octal); val >= 0 {
				sb.WriteByte(byte(val))
				i += 3
				continue
			}
		}
		sb.WriteByte(s[i])
	}

	return sb.String()
}

func parseOctal(s string) int {
	if len(s) != 3 {
		return -1
	}
	for _, c := range s {
		if c < '0' || c > '7' {
			return -1
		}
	}
	return int(s[0]-'0')*64 + int(s[1]-'0')*8 + int(s[2]-'0')
}

func extractDocxText(data []byte) (string, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("open docx zip: %w", err)
	}

	for _, f := range reader.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("open document.xml: %w", err)
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				return "", fmt.Errorf("read document.xml: %w", err)
			}

			return extractDocxXMLText(content)
		}
	}

	return "", fmt.Errorf("document.xml not found in docx archive")
}

func extractDocxXMLText(data []byte) (string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var sb strings.Builder
	inParagraph := false

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "p":
				inParagraph = true
			case "t":
			case "tab":
				sb.WriteString("\t")
			case "br":
				sb.WriteString("\n")
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "p":
				if inParagraph {
					sb.WriteString("\n\n")
				}
				inParagraph = false
			case "r":
			}
		case xml.CharData:
			if inParagraph {
				text := string(t)
				sb.WriteString(text)
			}
		}
	}

	result := sb.String()
	result = cleanExtractedText(result)
	return result, nil
}

func extractXlsxText(data []byte) (string, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("open xlsx zip: %w", err)
	}

	var sb strings.Builder
	sharedStrings := make([]string, 0)

	for _, f := range reader.File {
		if f.Name == "xl/sharedStrings.xml" {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			content, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				continue
			}
			sharedStrings = parseSharedStrings(content)
			break
		}
	}

	for _, f := range reader.File {
		if strings.HasPrefix(f.Name, "xl/worksheets/sheet") && strings.HasSuffix(f.Name, ".xml") {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			content, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				continue
			}

			text := parseSheetXML(content, sharedStrings)
			sb.WriteString(text)
			sb.WriteString("\n")
		}
	}

	result := sb.String()
	result = cleanExtractedText(result)
	return result, nil
}

func parseSharedStrings(data []byte) []string {
	var strings []string
	decoder := xml.NewDecoder(bytes.NewReader(data))

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if se, ok := token.(xml.StartElement); ok && se.Name.Local == "t" {
			if charData, err := decoder.Token(); err == nil {
				if cd, ok := charData.(xml.CharData); ok {
					strings = append(strings, string(cd))
				}
			}
		}
	}

	return strings
}

func parseSheetXML(data []byte, sharedStrings []string) string {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var sb strings.Builder
	inRow := false
	cellCount := 0

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "row":
				inRow = true
				cellCount = 0
			case "c":
				if inRow {
					cellCount++
					cellText := getCellValue(t, decoder, sharedStrings)
					if cellText != "" {
						if cellCount > 1 {
							sb.WriteString("\t")
						}
						sb.WriteString(cellText)
					}
				}
			}
		case xml.EndElement:
			if t.Name.Local == "row" && inRow {
				sb.WriteString("\n")
				inRow = false
			}
		}
	}

	return sb.String()
}

func getCellValue(se xml.StartElement, decoder *xml.Decoder, sharedStrings []string) string {
	var cellType string
	for _, attr := range se.Attr {
		if attr.Name.Local == "t" {
			cellType = attr.Value
		}
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return ""
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "v" {
				if charData, err := decoder.Token(); err == nil {
					if cd, ok := charData.(xml.CharData); ok {
						val := string(cd)
						if cellType == "s" {
							idx := 0
							for _, c := range val {
								idx = idx*10 + int(c-'0')
							}
							if idx >= 0 && idx < len(sharedStrings) {
								return sharedStrings[idx]
							}
						}
						return val
					}
				}
			}
		case xml.EndElement:
			if t.Name.Local == "c" {
				return ""
			}
		}
	}
}

func extractPptxText(data []byte) (string, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("open pptx zip: %w", err)
	}

	var sb strings.Builder
	slideCount := 0

	for _, f := range reader.File {
		if strings.HasPrefix(f.Name, "ppt/slides/slide") && strings.HasSuffix(f.Name, ".xml") {
			slideCount++
			rc, err := f.Open()
			if err != nil {
				continue
			}
			content, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				continue
			}

			sb.WriteString(fmt.Sprintf("--- Slide %d ---\n", slideCount))
			text := extractSlideXMLText(content)
			sb.WriteString(text)
			sb.WriteString("\n\n")
		}
	}

	result := sb.String()
	result = cleanExtractedText(result)

	if result == "" {
		return "[PowerPoint document - no text content extracted]", nil
	}

	return result, nil
}

func extractSlideXMLText(data []byte) string {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var sb strings.Builder

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "t" {
				if charData, err := decoder.Token(); err == nil {
					if cd, ok := charData.(xml.CharData); ok {
						text := string(cd)
						sb.WriteString(text)
					}
				}
			} else if t.Name.Local == "p" {
				// new paragraph
			}
		case xml.EndElement:
			if t.Name.Local == "p" {
				sb.WriteString("\n")
			}
		}
	}

	return sb.String()
}

func extractImageDescription(filename string) (string, error) {
	return fmt.Sprintf(
		"[Image file: %s - OCR text extraction is not available. Install an OCR engine (e.g., Tesseract) for image text extraction.]",
		filename,
	), nil
}

func extractPlainText(data []byte) (string, error) {
	if !utf8.Valid(data) {
		return "", fmt.Errorf("data is not valid UTF-8 text")
	}

	text := string(data)
	text = cleanExtractedText(text)
	return text, nil
}

func cleanExtractedText(text string) string {
	lines := strings.Split(text, "\n")
	var cleanLines []string
	prevEmpty := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if !prevEmpty {
				cleanLines = append(cleanLines, "")
				prevEmpty = true
			}
		} else {
			cleanLines = append(cleanLines, trimmed)
			prevEmpty = false
		}
	}

	result := strings.Join(cleanLines, "\n")
	result = strings.ReplaceAll(result, "\n\n\n", "\n\n")
	result = strings.TrimSpace(result)
	return result
}

func FixedSizeChunk(text string, size int, overlap int) []Chunk {
	if size <= 0 {
		size = 1500
	}
	if overlap < 0 {
		overlap = 200
	}
	if overlap >= size {
		overlap = size / 4
	}

	runes := []rune(text)
	var chunks []Chunk
	step := size - overlap
	if step <= 0 {
		step = size
	}

	for i := 0; i < len(runes); i += step {
		end := i + size
		if end > len(runes) {
			end = len(runes)
		}

		chunkText := string(runes[i:end])
		chunkText = strings.TrimSpace(chunkText)
		if chunkText == "" {
			continue
		}

		chunks = append(chunks, Chunk{
			Content: chunkText,
			Index:   len(chunks),
			Metadata: map[string]string{
				"chunk_type":  "fixed_size",
				"chunk_start": fmt.Sprintf("%d", i),
				"chunk_end":   fmt.Sprintf("%d", end),
			},
		})

		if end >= len(runes) {
			break
		}
	}

	return chunks
}

func SemanticChunk(text string, maxSize int) []Chunk {
	if maxSize <= 0 {
		maxSize = 1500
	}

	paragraphs := splitParagraphs(text)
	var chunks []Chunk
	var currentText strings.Builder
	currentSize := 0

	flushChunk := func() {
		if currentText.Len() > 0 {
			content := strings.TrimSpace(currentText.String())
			if content != "" {
				chunks = append(chunks, Chunk{
					Content: content,
					Index:   len(chunks),
					Metadata: map[string]string{
						"chunk_type": "semantic",
					},
				})
			}
			currentText.Reset()
			currentSize = 0
		}
	}

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			flushChunk()
			continue
		}

		paraRunes := len([]rune(para))

		if currentSize+paraRunes > maxSize && currentSize > 0 {
			flushChunk()
		}

		if currentText.Len() > 0 {
			currentText.WriteString("\n\n")
			currentSize += 2
		}
		currentText.WriteString(para)
		currentSize += paraRunes
	}

	flushChunk()

	return chunks
}

func RecursiveChunk(text string, maxSize int, overlap int) []Chunk {
	if maxSize <= 0 {
		maxSize = 1500
	}
	if overlap < 0 {
		overlap = 200
	}

	textLen := len([]rune(text))
	if textLen <= maxSize {
		trimmed := strings.TrimSpace(text)
		if trimmed == "" {
			return nil
		}
		return []Chunk{{
			Content: trimmed,
			Index:   0,
			Metadata: map[string]string{
				"chunk_type": "recursive",
			},
		}}
	}

	var chunks []Chunk
	separators := []string{"\n\n", "\n", ". ", "? ", "! ", ";", ", ", " "}

	var splitFunc func(string, int, int) []Chunk
	splitFunc = func(s string, depth int, parentSize int) []Chunk {
		if depth >= len(separators) {
			return FixedSizeChunk(s, maxSize, overlap)
		}

		sLen := len([]rune(s))
		if sLen <= maxSize {
			trimmed := strings.TrimSpace(s)
			if trimmed == "" {
				return nil
			}
			return []Chunk{{
				Content: trimmed,
				Index:   0,
				Metadata: map[string]string{
					"chunk_type": "recursive",
				},
			}}
		}

		sep := separators[depth]
		parts := strings.Split(s, sep)

		if len(parts) == 1 {
			return splitFunc(s, depth+1, parentSize)
		}

		var result []Chunk
		for _, part := range parts {
			subChunks := splitFunc(strings.TrimSpace(part), depth, parentSize)
			result = append(result, subChunks...)
		}

		return result
	}

	chunks = splitFunc(text, 0, maxSize)

	finalChunks := mergeChunks(chunks, maxSize, overlap)
	for i := range finalChunks {
		finalChunks[i].Index = i
	}

	return finalChunks
}

func mergeChunks(chunks []Chunk, maxSize int, overlap int) []Chunk {
	if len(chunks) <= 1 {
		return chunks
	}

	var merged []Chunk
	var currentContent strings.Builder
	currentSize := 0

	flushChunk := func() {
		if currentContent.Len() > 0 {
			content := strings.TrimSpace(currentContent.String())
			if content != "" {
				merged = append(merged, Chunk{
					Content: content,
					Index:   len(merged),
					Metadata: map[string]string{
						"chunk_type": "recursive_merged",
					},
				})
			}
			currentContent.Reset()
			currentSize = 0
		}
	}

	for _, chunk := range chunks {
		chunkRunes := len([]rune(chunk.Content))

		if currentSize+chunkRunes > maxSize && currentSize > 0 {
			flushChunk()
		}

		if currentContent.Len() > 0 {
			currentContent.WriteString(" ")
			currentSize++
		}
		currentContent.WriteString(chunk.Content)
		currentSize += chunkRunes
	}

	flushChunk()

	if overlap > 0 && len(merged) > 1 {
		for i := 1; i < len(merged); i++ {
			prev := []rune(merged[i-1].Content)
			curr := []rune(merged[i].Content)

			overlapLen := overlap
			if overlapLen > len(prev) {
				overlapLen = len(prev)
			}

			overlapText := string(prev[len(prev)-overlapLen:])
			merged[i].Content = overlapText + " " + string(curr)
		}
	}

	return merged
}

func splitParagraphs(text string) []string {
	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")

	raw := strings.Split(normalized, "\n\n")
	var paragraphs []string
	for _, p := range raw {
		p = strings.TrimSpace(p)
		if p != "" {
			paragraphs = append(paragraphs, p)
		}
	}

	return paragraphs
}

func ChunkByTokens(text string, maxTokens int) []Chunk {
	if maxTokens <= 0 {
		maxTokens = 500
	}

	words := strings.Fields(text)
	tokensPerWord := 1.3

	var chunks []Chunk
	var currentWords []string
	currentTokens := 0.0

	flushChunk := func() {
		if len(currentWords) > 0 {
			content := strings.Join(currentWords, " ")
			chunks = append(chunks, Chunk{
				Content: strings.TrimSpace(content),
				Index:   len(chunks),
				Metadata: map[string]string{
					"chunk_type":  "token_based",
					"token_count": fmt.Sprintf("%.0f", currentTokens),
				},
			})
			currentWords = nil
			currentTokens = 0
		}
	}

	for _, word := range words {
		wordTokens := float64(len([]rune(word))) / 4.0
		if wordTokens < 1 {
			wordTokens = 1
		} else {
			wordTokens *= tokensPerWord
		}

		if currentTokens+wordTokens > float64(maxTokens) && len(currentWords) > 0 {
			flushChunk()
		}

		currentWords = append(currentWords, word)
		currentTokens += wordTokens
	}

	flushChunk()

	return chunks
}

func ChunkByMarkdownHeaders(text string, maxSize int) []Chunk {
	if maxSize <= 0 {
		maxSize = 1500
	}

	lines := strings.Split(text, "\n")
	var chunks []Chunk
	var currentSection string
	var currentContent strings.Builder
	currentSize := 0

	flushChunk := func(section string) {
		if currentContent.Len() > 0 {
			content := strings.TrimSpace(currentContent.String())
			if content != "" {
				chunks = append(chunks, Chunk{
					Content: content,
					Index:   len(chunks),
					Metadata: map[string]string{
						"chunk_type": "markdown_section",
						"section":    section,
					},
				})
			}
			currentContent.Reset()
			currentSize = 0
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "#") {
			flushChunk(currentSection)

			currentSection = trimmed
			currentContent.WriteString(trimmed)
			currentContent.WriteString("\n")
			currentSize = len([]rune(trimmed)) + 1
		} else {
			lineRunes := len([]rune(line))
			if currentSize+lineRunes > maxSize && currentSize > 0 {
				flushChunk(currentSection)
				if trimmed != "" {
					currentContent.WriteString(trimmed)
					currentSize = lineRunes
				}
			} else {
				if currentContent.Len() > 0 {
					currentContent.WriteString("\n")
					currentSize++
				}
				currentContent.WriteString(line)
				currentSize += lineRunes
			}
		}
	}

	flushChunk(currentSection)

	if len(chunks) == 0 {
		return FixedSizeChunk(text, maxSize, maxSize/4)
	}

	return chunks
}

func NormalizeText(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = strings.ReplaceAll(text, "\u00a0", " ")

	var sb strings.Builder
	for _, r := range text {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			continue
		}
		if unicode.Is(unicode.Cf, r) {
			continue
		}
		sb.WriteRune(r)
	}

	lines := strings.Split(sb.String(), "\n")
	var cleanLines []string
	prevEmpty := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if !prevEmpty {
				cleanLines = append(cleanLines, "")
				prevEmpty = true
			}
		} else {
			cleanLines = append(cleanLines, trimmed)
			prevEmpty = false
		}
	}

	result := strings.Join(cleanLines, "\n")
	result = strings.TrimSpace(result)
	return result
}
