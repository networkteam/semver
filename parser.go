package semver

import (
	"fmt"
	"strconv"
	"strings"
)

/*

Backus–Naur Form Grammar for Valid SemVer Versions

<valid semver> ::= <version core>
                 | <version core> "-" <pre-release>
                 | <version core> "+" <build>
                 | <version core> "-" <pre-release> "+" <build>

<version core> ::= <major> "." <minor> "." <patch>

<major> ::= <numeric identifier>

<minor> ::= <numeric identifier>

<patch> ::= <numeric identifier>

<pre-release> ::= <dot-separated pre-release identifiers>

<dot-separated pre-release identifiers> ::= <pre-release identifier>
                                          | <pre-release identifier> "." <dot-separated pre-release identifiers>

<build> ::= <dot-separated build identifiers>

<dot-separated build identifiers> ::= <build identifier>
                                    | <build identifier> "." <dot-separated build identifiers>

<pre-release identifier> ::= <alphanumeric identifier>
                           | <numeric identifier>

<build identifier> ::= <alphanumeric identifier>
                     | <digits>

<alphanumeric identifier> ::= <non-digit>
                            | <non-digit> <identifier characters>
                            | <identifier characters> <non-digit>
                            | <identifier characters> <non-digit> <identifier characters>

<numeric identifier> ::= "0"
                       | <positive digit>
                       | <positive digit> <digits>

<identifier characters> ::= <identifier character>
                          | <identifier character> <identifier characters>

<identifier character> ::= <digit>
                         | <non-digit>

<non-digit> ::= <letter>
              | "-"

<digits> ::= <digit>
           | <digit> <digits>

<digit> ::= "0"
          | <positive digit>

<positive digit> ::= "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"

<letter> ::= "A" | "B" | "C" | "D" | "E" | "F" | "G" | "H" | "I" | "J"
           | "K" | "L" | "M" | "N" | "O" | "P" | "Q" | "R" | "S" | "T"
           | "U" | "V" | "W" | "X" | "Y" | "Z" | "a" | "b" | "c" | "d"
           | "e" | "f" | "g" | "h" | "i" | "j" | "k" | "l" | "m" | "n"
           | "o" | "p" | "q" | "r" | "s" | "t" | "u" | "v" | "w" | "x"
           | "y" | "z"


*/

type Parser struct {
	input string
	pos   int
}

func NewParser(input string) *Parser {
	return &Parser{input: input, pos: 0}
}

// ParseVersion parses a valid semantic version (<valid semver>)
func (p *Parser) ParseVersion() (*Version, error) {
	major, minor, patch, err := p.parseVersionCore()
	if err != nil {
		return nil, fmt.Errorf("invalid version core: %w", err)
	}

	var preRelease, build string
	if p.match('-') {
		p.pos++
		preRelease, err = p.parsePreRelease()
		if err != nil {
			return nil, fmt.Errorf("invalid pre-release: %w", err)
		}
	}
	if p.match('+') {
		p.pos++
		build, err = p.parseBuild()
		if err != nil {
			return nil, fmt.Errorf("invalid build: %w", err)
		}
	}
	if p.pos < len(p.input) {
		return nil, &ParseError{Position: p.pos, Message: fmt.Sprintf("unexpected trailing characters: %q", p.input[p.pos:])}
	}

	return &Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		PreRelease: preRelease,
		Build:      build,
	}, nil
}

func (p *Parser) parseVersionCore() (major int, minor int, patch int, err error) {
	major, err = p.parseNumericIdentifier()
	if err != nil {
		return 0, 0, 0, fmt.Errorf("major: %w", err)
	}

	if !p.consume('.') {
		return 0, 0, 0, &ParseError{Position: p.pos, Message: "missing dot separator"}
	}

	minor, err = p.parseNumericIdentifier()
	if err != nil {
		return 0, 0, 0, fmt.Errorf("minor: %w", err)
	}

	if !p.consume('.') {
		return 0, 0, 0, &ParseError{Position: p.pos, Message: "missing dot separator"}
	}

	patch, err = p.parseNumericIdentifier()
	if err != nil {
		return 0, 0, 0, fmt.Errorf("patch: %w", err)
	}

	return major, minor, patch, nil
}

func (p *Parser) parseNumericIdentifier() (int, error) {
	if p.match('0') {
		p.pos++

		// Check if next character is a digit for better error messages
		if p.matchDigit() {
			return 0, &ParseError{
				Position: p.pos - 1,
				Message:  "leading zero is not allowed",
			}
		}

		return 0, nil
	}

	var sb strings.Builder

	err := p.appendPositiveDigit(&sb)
	if err != nil {
		return 0, err
	}

	p.appendDigits(&sb)

	num, err := strconv.Atoi(sb.String())
	if err != nil {
		return 0, err
	}

	return num, nil
}

func (p *Parser) consume(ch byte) bool {
	if p.pos < len(p.input) && p.input[p.pos] == ch {
		p.pos++
		return true
	}
	return false
}

func (p *Parser) match(ch byte) bool {
	if p.pos < len(p.input) && p.input[p.pos] == ch {
		return true
	}
	return false
}

func (p *Parser) appendPositiveDigit(sb *strings.Builder) error {
	if p.pos >= len(p.input) {
		return &ParseError{
			Position: p.pos,
			Message:  "unexpected end of input",
		}
	}

	if p.input[p.pos] < '1' || p.input[p.pos] > '9' {
		return &ParseError{
			Position: p.pos,
			Message:  fmt.Sprintf("expected positive digit, got %c", p.input[p.pos]),
		}
	}

	sb.WriteByte(p.input[p.pos])
	p.pos++
	return nil
}

// appendDigits appends as many digits as possible to the string builder
func (p *Parser) appendDigits(sb *strings.Builder) {
	for p.pos < len(p.input) && p.input[p.pos] >= '0' && p.input[p.pos] <= '9' {
		sb.WriteByte(p.input[p.pos])
		p.pos++
	}
}

func (p *Parser) matchDigit() bool {
	if p.pos < len(p.input) && p.input[p.pos] >= '0' && p.input[p.pos] <= '9' {
		return true
	}
	return false
}

func (p *Parser) matchLetter() bool {
	if p.pos < len(p.input) && (p.input[p.pos] >= 'A' && p.input[p.pos] <= 'Z' || p.input[p.pos] >= 'a' && p.input[p.pos] <= 'z') {
		return true
	}
	return false
}

func (p *Parser) parsePreRelease() (string, error) {
	var sb strings.Builder

	err := p.appendAlphanumericIdentifier(&sb)
	if err != nil {
		return "", err
	}

	for p.consume('.') {
		sb.WriteByte('.')
		err := p.appendAlphanumericIdentifier(&sb)
		if err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}

func (p *Parser) appendAlphanumericIdentifier(sb *strings.Builder) error {
	if p.pos >= len(p.input) {
		return &ParseError{
			Position: p.pos,
			Message:  "unexpected end of input",
		}
	}

	for p.matchLetter() || p.matchDigit() || p.match('-') {
		sb.WriteByte(p.input[p.pos])
		p.pos++
	}

	if sb.Len() == 0 {
		return &ParseError{
			Position: p.pos,
			Message:  fmt.Sprintf("expected alphanumeric identifier, got %c", p.input[p.pos]),
		}
	}

	return nil
}

func (p *Parser) parseBuild() (string, error) {
	var sb strings.Builder

	err := p.appendAlphanumericIdentifier(&sb)
	if err != nil {
		return "", err
	}

	for p.consume('.') {
		sb.WriteByte('.')
		err := p.appendAlphanumericIdentifier(&sb)
		if err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}
