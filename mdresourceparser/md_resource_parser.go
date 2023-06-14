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
	input          []rune
	char           rune       // current char under examination
	position       int        // current position in input (points to current char)
	readPosition   int        // current reading position in input (after current char)
	newline        bool       // sets to true when char was preceded by a new line character \n
	Links          []string   // collection of links paths
	dockerDirMatch *pathMatch // keeps track of parsing a newline matcher for docker directories
	testFileMatch  *pathMatch // keeps track of parsing a newline matcher for test file matches
	setupFileMatch *pathMatch // keeps track of parsing a newline matcher for setup file matches
	dataPathMatch  *pathMatch // keeps track of parsing a newline matcher for data path file matches
}

// pathMatch keeps track of a new line to detect if it starts with the given match string
type pathMatch struct {
	match string   // the string for startMatch to find
	paths []string // stored paths
}

// appendPath takes a given parser and reads from it to determine a match with the pathMatch.
// If one is found it appends the path to the pathMatch. The parser advances position either until
// a mismatched character is found, or in the case of a match the parser advances until the end of the found path.
func (pm *pathMatch) appendPath(parser *MDResourceParser) error {
	for _, matchChar := range pm.match {
		err := parser.matchError(matchChar)
		if err != nil {
			return err
		}
	}

	path, err := parser.readUntilChar('\n')
	if err != nil {
		return err
	}

	pm.paths = append(pm.paths, strings.TrimSpace(path))
	return nil
}

// New creates and returns a pointer to the MDResourceParser with its input attached
func New(input []rune) *MDResourceParser {
	p := &MDResourceParser{
		input:          input,
		dockerDirMatch: &pathMatch{match: "docker_directory_path:", paths: []string{}},
		setupFileMatch: &pathMatch{match: "setup_file:", paths: []string{}},
		testFileMatch:  &pathMatch{match: "test_file:", paths: []string{}},
		dataPathMatch:  &pathMatch{match: "data_path:", paths: []string{}},
	}
	p.readChar()
	return p
}

// ParseResources takes the input contents and parses it for our links and other
// curriculum content that represents files in the repository. It returns the resources found
func (p *MDResourceParser) ParseResources() (dataPaths, dockerDirPaths, testFilePaths, setupFilePaths []string) {
	for p.readPosition < len(p.input) {
		p.next()
	}
	return p.dataPathMatch.paths,
		p.dockerDirMatch.paths,
		p.testFileMatch.paths,
		p.setupFileMatch.paths
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

// hasPathBullet reports if the two characters at current position and read position are '* ' or '- ',
// and advances the readPosition past them. otherwise it reports false and does not advance the read position
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

// extractPath finds paths defined on challenges to other files and directories
// it appends the found paths to the parser
func (p *MDResourceParser) extractPath() error {
	if p.readPosition >= len(p.input) {
		return io.EOF
	}

	// any desired path is present if the line starts with '* '
	if !p.hasPathBullet() {
		return fmt.Errorf("no match")
	}

	switch p.char {
	case 'd':
		// 'd' can be 'data_path' or 'docker_directory_path'
		if p.peek() == 'a' {
			return p.dataPathMatch.appendPath(p)
		} else {
			return p.dockerDirMatch.appendPath(p)
		}
	case 't':
		return p.testFileMatch.appendPath(p)
	case 's':
		return p.setupFileMatch.appendPath(p)
	default:
		return fmt.Errorf("no match")
	}

	return nil
}

// matchError is used by the parser to determine if a rune equas the current character, and protects against EOF
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
