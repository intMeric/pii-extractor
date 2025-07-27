package patterns

import (
	"regexp"
	"strings"
	"unicode"
)

// International/generic patterns
const (
	EmailPattern          = `(?i)\b([A-Za-z0-9!#$%&'*+\/=?^_{|.}~-]+@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)\b`
	IPv4Pattern           = `(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
	IPv6Pattern           = `(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){1}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(?:%.+)?\s*`
	IPPattern             = IPv4Pattern + `|` + IPv6Pattern
	CreditCardPattern     = `\b(?:(?:\d{4}[\s-]?){3}\d{4}|\d{15,16})\b`
	VISACreditCardPattern = `4\d{3}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}`
	MCCreditCardPattern   = `5[1-5]\d{2}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}`
	BtcAddressPattern     = `\b[13][a-km-zA-HJ-NP-Z1-9]{25,34}\b`
	IBANPattern           = `\b[A-Z]{2}\d{2}[A-Z0-9]{4,}\d{7,}[A-Z0-9]*\b`
)

// International/generic compiled patterns
var (
	EmailRegex          = regexp.MustCompile(EmailPattern)
	IPv4Regex           = regexp.MustCompile(IPv4Pattern)
	IPv6Regex           = regexp.MustCompile(IPv6Pattern)
	IPRegex             = regexp.MustCompile(IPPattern)
	CreditCardRegex     = regexp.MustCompile(CreditCardPattern)
	VISACreditCardRegex = regexp.MustCompile(VISACreditCardPattern)
	MCCreditCardRegex   = regexp.MustCompile(MCCreditCardPattern)
	BtcAddressRegex     = regexp.MustCompile(BtcAddressPattern)
	IBANRegex           = regexp.MustCompile(IBANPattern)
)

func Match(text string, regex *regexp.Regexp) []string {
	matches := regex.FindAllStringSubmatch(text, -1)
	if matches == nil {
		return []string{}
	}

	var results []string
	for _, match := range matches {
		if len(match) > 1 {
			// Use the first capture group if it exists
			results = append(results, match[1])
		} else {
			// Use the full match if no capture groups
			results = append(results, match[0])
		}
	}
	return results
}

// MatchAddresses is a specialized function for address patterns that handles
// splitting matches containing linking words like "to", "and", "y"
func MatchAddresses(text string, regex *regexp.Regexp) []string {
	matches := regex.FindAllString(text, -1)
	if matches == nil {
		return []string{}
	}

	var results []string
	linkingWords := []string{" to ", " and ", " y ", " et "}
	endLinkingWords := []string{" to", " and", " y", " et"}

	for _, match := range matches {
		// Check if the match contains linking words
		processed := false

		// First check for linking words in the middle
		for _, linkWord := range linkingWords {
			if strings.Contains(strings.ToLower(match), linkWord) {
				// Split the match at the linking word
				linkIdx := strings.Index(strings.ToLower(match), linkWord)
				if linkIdx >= 0 {
					firstPart := strings.TrimSpace(match[:linkIdx])
					secondStart := linkIdx + len(linkWord)

					if firstPart != "" {
						results = append(results, firstPart)
					}

					if secondStart < len(match) {
						secondPart := strings.TrimSpace(match[secondStart:])
						if secondPart != "" {
							results = append(results, secondPart)
						}
					}
					processed = true
					break
				}
			}
		}

		// If not processed, check for linking words at the end
		if !processed {
			for _, endWord := range endLinkingWords {
				if strings.HasSuffix(strings.ToLower(match), endWord) {
					trimmed := strings.TrimSpace(match[:len(match)-len(endWord)])
					if trimmed != "" {
						results = append(results, trimmed)
						processed = true
						break
					}
				}
			}
		}

		// If no linking words found, add the whole match
		if !processed {
			results = append(results, match)
		}
	}

	return results
}

// MatchWithIndices returns matches along with their start and end positions
func MatchWithIndices(text string, regex *regexp.Regexp) [][]int {
	return regex.FindAllStringIndex(text, -1)
}

// ExtractContext extracts the context around a match, prioritizing full sentences over word count
func ExtractContext(text string, start, end int) string {
	// First try to find a complete sentence
	sentence := extractSentence(text, start, end)
	if sentence != "" {
		return strings.TrimSpace(sentence)
	}

	// Fallback to 8 words before and after
	return extractWordContext(text, start, end)
}

// extractSentence tries to extract a complete sentence containing the match
func extractSentence(text string, start, end int) string {
	// Find sentence boundaries (., !, ?, or start/end of text)
	sentenceStart := start
	sentenceEnd := end

	// Look backwards for sentence start
	for i := start - 1; i >= 0; i-- {
		char := text[i]
		if char == '.' || char == '!' || char == '?' {
			sentenceStart = i + 1
			break
		}
		if i == 0 {
			sentenceStart = 0
		}
	}

	// Look forwards for sentence end
	for i := end; i < len(text); i++ {
		char := text[i]
		if char == '.' || char == '!' || char == '?' {
			sentenceEnd = i + 1
			break
		}
		if i == len(text)-1 {
			sentenceEnd = len(text)
		}
	}

	// Skip whitespace at the beginning
	for sentenceStart < len(text) && unicode.IsSpace(rune(text[sentenceStart])) {
		sentenceStart++
	}

	if sentenceStart < sentenceEnd && sentenceEnd <= len(text) {
		return text[sentenceStart:sentenceEnd]
	}

	return ""
}

// extractWordContext extracts 8 words before and after the match
func extractWordContext(text string, start, end int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	// Find the word indices that contain our match
	wordStart := -1
	wordEnd := -1
	currentPos := 0

	for i, word := range words {
		// Find the actual position of this word in the text
		wordStartPos := strings.Index(text[currentPos:], word)
		if wordStartPos == -1 {
			continue
		}
		wordStartPos += currentPos
		wordEndPos := wordStartPos + len(word)

		// Check if this word contains the start of our match
		if wordStartPos <= start && start < wordEndPos {
			wordStart = i
		}
		// Check if this word contains the end of our match  
		if wordStartPos < end && end <= wordEndPos {
			wordEnd = i
		}

		// Move to the position after this word for the next search
		currentPos = wordEndPos
		
		// Skip whitespace
		for currentPos < len(text) && unicode.IsSpace(rune(text[currentPos])) {
			currentPos++
		}

		if wordStart != -1 && wordEnd != -1 {
			break
		}
	}

	if wordStart == -1 || wordEnd == -1 {
		return ""
	}

	// Extract 8 words before and after
	contextStart := max(0, wordStart-8)
	contextEnd := min(len(words), wordEnd+8+1)

	return strings.Join(words[contextStart:contextEnd], " ")
}

// International/generic convenience functions
var Emails = func(text string) []string { return Match(text, EmailRegex) }
var IPv4s = func(text string) []string { return Match(text, IPv4Regex) }
var IPv6s = func(text string) []string { return Match(text, IPv6Regex) }
var IPs = func(text string) []string { return Match(text, IPRegex) }
var CreditCards = func(text string) []string { return Match(text, CreditCardRegex) }
var VISACreditCards = func(text string) []string { return Match(text, VISACreditCardRegex) }
var MCCreditCards = func(text string) []string { return Match(text, MCCreditCardRegex) }
var BtcAddresses = func(text string) []string { return Match(text, BtcAddressRegex) }
var IBANs = func(text string) []string { return Match(text, IBANRegex) }
