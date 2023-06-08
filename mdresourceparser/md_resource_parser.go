package mdresourceparser

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// MDResourceParser performs our lexical analysis/scanning/parsing. Not a true lexer/parser
// because right now we don't need tokens, logic, ast, just to collect MD paths
type MDResourceParser struct {
	input                []rune
	char                 rune        // current char under examination
	position             int         // current position in input (points to current char)
	readPosition         int         // current reading position in input (after current char)
	newline              bool        // sets to true when char was preceded by a new line character \n
	Links                []string    // collection of links paths
	DockerDirectoryPaths []string    // collection of docker_directory_paths
	dockerDirMatch       *startMatch // keeps track of parsing a newline matcher for docker directories
	TestFilePaths        []string    // collection of test file paths
	testFileMatch        *startMatch // keeps track of parsing a newline matcher for test file matches
	SetupFilePaths       []string    // collection of setup file paths
	setupFileMatch       *startMatch // keeps track of parsing a newline matcher for setup file matches
}

// startMatch keeps track of a new line to detect if it starts with the given match string
type startMatch struct {
	match    string // the string for startMatch to find
	value    string // after the match is found, the value which comes after the start string
	position int    // position within the match
}

// New creates and returns a pointer to the MDResourceParser with it's input attached
func New(input []rune) *MDResourceParser {
	p := &MDResourceParser{
		input:          input,
		dockerDirMatch: &startMatch{match: "docker_directory_path:"},
		setupFileMatch: &startMatch{match: "setup_file:"},
		testFileMatch:  &startMatch{match: "test_file:"},
	}
	p.readChar()
	return p
}

// ParseResources takes the input contents and parses it for our MD image links
func (p *MDResourceParser) ParseResources() {
	for p.readPosition < len(p.input) {
		p.next()
	}
}

// readChar checks if were at EOF and if we are not, it sets the parser's char to
// the char at our readPosition and increments position and read position by one
func (p *MDResourceParser) readChar() {
	if p.readPosition >= len(p.input) {
		// End of input (haven't read anything yet or EOF)
		// 0 is ASCII code for "NUL" character
		p.char = 0
	} else {
		p.char = p.input[p.readPosition]
	}

	p.newline = false
	if p.char == '\n' {
		p.newline = true
	}

	p.position = p.readPosition
	p.readPosition++
}

// skipWhitespace gets called on every iteration of "next" because we do
// not care about whitespace
func (p *MDResourceParser) skipWhitespace() {
	for p.char == ' ' || p.char == '\t' || p.char == '\r' {
		p.readChar()
	}
}

// peek checks to see what the next char is by reading the input at our readPosition
// This does not consume any characters or increment our position
func (p *MDResourceParser) peek() rune {
	if p.readPosition >= len(p.input) {
		return 0
	}
	return p.input[p.readPosition]
}

// hasPathBullet reports if the next two characters are '* ', and advances the readPosition past them
// otherwise it reports false and does not advance the read position
func (p *MDResourceParser) hasPathBullet() bool {
	if p.position >= len(p.input) || p.position+1 >= len(p.input) {
		return false
	}
	if (p.char == '*' || p.char == '-') && p.peek() == ' ' {
		p.readChar()
		p.readChar()
		return true
	}
	return false
}

func (p *MDResourceParser) extractPath() error {
	if p.readPosition >= len(p.input) {
		return io.EOF
	}

	// any desired path ispresent if the line starts with '* '
	// if it has this prefix, check if one of the three desired matches are present
	if !p.hasPathBullet() {
		return fmt.Errorf("no match")
	}

	switch p.char {
	case 'd':
		// if the startMatch matches, continue readChar while it matches and increase the startChar

		for _, matchChar := range p.dockerDirMatch.match {
			if matchChar != p.char {
				return fmt.Errorf("no match")
			}

			if p.readPosition >= len(p.input) {
				return io.EOF
			}
			p.readChar()
		}

		path, err := p.readUntilChar('\n')
		if err != nil {
			return nil
		}

		p.DockerDirectoryPaths = append(p.DockerDirectoryPaths, strings.TrimSpace(path))
	case 't':
		// if the startMatch matches, continue readChar while it matches and increase the startChar
		for _, matchChar := range p.testFileMatch.match {
			if matchChar != p.char {
				return fmt.Errorf("no match")
			}

			if p.readPosition >= len(p.input) {
				return io.EOF
			}
			p.readChar()
		}

		path, err := p.readUntilChar('\n')
		if err != nil {
			return err
		}
		p.TestFilePaths = append(p.TestFilePaths, strings.TrimSpace(path))
	case 's':
		for _, matchChar := range p.setupFileMatch.match {
			if matchChar != p.char {
				return fmt.Errorf("no match")
			}

			if p.readPosition >= len(p.input) {
				return io.EOF
			}
			p.readChar()
		}

		path, err := p.readUntilChar('\n')
		if err != nil {
			return nil
		}
		p.SetupFilePaths = append(p.SetupFilePaths, strings.TrimSpace(path))
	default:
		return fmt.Errorf("no match")
	}

	//// if the startMatch matches, continue readChar while it matches and increase the startChar
	//for _, matchChar := range p.dockerDirMatch.match {
	//	if matchChar != p.char {
	//		return "", fmt.Errorf("no match")
	//	}

	//	if p.readPosition >= len(p.input) {
	//		return "", io.EOF
	//	}
	//	p.readChar()
	//}

	// if a full match is found, extract the remaining characters until a newline. If the p.char doesn't match, error
	//path, err := p.readUntilChar('\n')
	//if err != nil {
	//	return "", nil
	//}
	//return path, nil
	return nil
}

func (p *MDResourceParser) matchError(matchChar rune) error {
	if matchChar != p.char {
		return fmt.Errorf("no match")
	}

	if p.readPosition >= len(p.input) {
		return io.EOF
	}
	p.readChar()
	return nil
}

func (p *MDResourceParser) extractLink() (string, error) {
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

func (p *MDResourceParser) readImagePath() (string, error) {
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

func (p *MDResourceParser) readUntilChar(c rune) (string, error) {
	var path []rune

	for p.char != c && p.char != 0 {
		path = append(path, p.char)
		p.readChar()
	}

	if p.char == 0 {
		return "", fmt.Errorf("End of file reading until character %v", p.char)
	}

	return string(path), nil
}

// next switches through the lexer's current char and creates a new token.
// It then it calls readChar() to advance the lexer and it returns the token
func (p *MDResourceParser) next() {
	p.skipWhitespace()

	switch p.char {
	case '\n':
		// move to the newline
		p.readChar()
		// attempt to extract path
		p.extractPath()
		return
	case '[':
		linkPath, err := p.extractLink()
		if err != nil {
			return
		}
		// We do not need to worry about hosted links, check for full url (http/https)
		if linkPath == "" || strings.HasPrefix(linkPath, "http") || strings.HasPrefix(linkPath, "https") {
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
