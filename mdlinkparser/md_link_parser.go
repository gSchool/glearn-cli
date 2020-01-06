package mdlinkparser

import (
	"errors"
	"io"
	"strings"
)

// MDLinkParser performs our lexical analysis/scanning/parsing. Not a true lexer/parser
// because right now we don't need tokens, logic, ast, just to collect MD paths
type MDLinkParser struct {
	input        []rune
	char         rune     // current char under examination
	position     int      // current position in input (points to current char)
	readPosition int      // current reading position in input (after current char)
	Links        []string // collection of links paths
}

// New creates and returns a pointer to the MDLinkParser with it's input attached
func New(input string) *MDLinkParser {
	p := &MDLinkParser{input: []rune(input)}
	p.readChar()
	return p
}

// ParseLinks takes the input contents and parses it for our MD image links
func (p *MDLinkParser) ParseLinks() {
	for p.readPosition < len(p.input) {
		p.next()
	}
}

// readChar checks if were at EOF and if we are not, it sets the parser's char to
// the char at our readPosition and increments position and read position by one
func (p *MDLinkParser) readChar() {
	if p.readPosition >= len(p.input) {
		// End of input (haven't read anything yet or EOF)
		// 0 is ASCII code for "NUL" character
		p.char = 0
	} else {
		p.char = p.input[p.readPosition]
	}

	p.position = p.readPosition
	p.readPosition++
}

// skipWhitespace gets called on every iteration of "next" because we do
// not care about whitespace
func (p *MDLinkParser) skipWhitespace() {
	for p.char == ' ' || p.char == '\t' || p.char == '\n' || p.char == '\r' {
		p.readChar()
	}
}

// peek checks to see what the next char is by reading the input at our readPosition
// This does not consume any characters or increment our position
func (p *MDLinkParser) peek() rune {
	if p.readPosition >= len(p.input) {
		return 0
	}
	return p.input[p.readPosition]
}

func (p *MDLinkParser) extractLink() (string, error) {
	if p.readPosition >= len(p.input) {
		return "", io.EOF
	}

	for p.char != ']' && p.char != 0 {
		p.readChar()
	}

	if p.peek() != '(' {
		return "", errors.New("Not a valid MD image link")
	}

	// Consume ) and [ to get to start of the link
	p.readChar()
	p.readChar()

	linkPath, err := p.readImagePath()
	if err != nil {
		return "", errors.New("Please check your markdown syntax for your image links")
	}

	return linkPath, nil
}

func (p *MDLinkParser) readImagePath() (string, error) {
	var path []rune

	for p.char != ')' && p.char != 0 {
		path = append(path, p.char)
		p.readChar()
	}

	if p.char == 0 {
		return "", errors.New("Please check your markdown syntax for your image links")
	}

	return string(path), nil
}

// next switches through the lexer's current char and creates a new token.
// It then it calls readChar() to advance the lexer and it returns the token
func (p *MDLinkParser) next() {
	p.skipWhitespace()

	switch p.char {
	case '[':
		linkPath, err := p.extractLink()
		if err != nil {
			return
		}
		// We do not need to worry about hosted links, check for full url (http/https)
		if strings.HasPrefix(linkPath, "http") || strings.HasPrefix(linkPath, "https") || strings.HasSuffix(linkPath, ".md") {
			return
		}
		p.Links = append(p.Links, linkPath)
	case 0:
		p.readChar()
		return
	default:
		p.readChar()
		return
	}

	p.readChar()
}
