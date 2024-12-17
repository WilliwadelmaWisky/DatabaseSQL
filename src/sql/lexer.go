package sql

// A single token
type Token struct {
	Type  TokenType
	Value string
}

// Enum to represent a type of a token, values are named with a TOKEN prefix
type TokenType int

const (
	// Token represents text
	TOKEN_TEXT TokenType = iota
	// Token represents a single operator `< > =`
	TOKEN_OPERATOR
	// Token represents a single comma `,`
	TOKEN_COMMA
	// Token represents a single asterisk `*`
	TOKEN_ASTERISK
	// Token represents a single parenthesis `( )`
	TOKEN_PARENTHESIS
)

// Get a TokenType enum value based of the input string
func GetTokenType(value string) TokenType {
	switch value {
	case "*":
		return TOKEN_ASTERISK
	case ",":
		return TOKEN_COMMA
	case "=", "<", ">":
		return TOKEN_OPERATOR
	case "(", ")":
		return TOKEN_PARENTHESIS
	default:
		return TOKEN_TEXT
	}
}

// Read all tokens from an input byte array
func Tokenize(b []byte) []*Token {
	start := 0
	tokens := []*Token{}
	isQuoted := false

	for i := 0; i < len(b); i++ {
		if isQuoted {
			if b[i] != '\'' {
				continue
			}

			value := string(b[start+1 : i])
			token := &Token{
				Type:  GetTokenType(value),
				Value: value,
			}
			tokens = append(tokens, token)
			start = i + 1
			isQuoted = false
			continue
		}

		if b[i] == '\'' {
			if start <= i-1 {
				value := string(b[start:i])
				token := &Token{
					Type:  GetTokenType(value),
					Value: value,
				}
				tokens = append(tokens, token)
				start = i
			}

			isQuoted = true
			continue
		}

		if IsSpecial(b[i]) {
			if start <= i-1 {
				value := string(b[start:i])
				token := &Token{
					Type:  GetTokenType(value),
					Value: value,
				}
				tokens = append(tokens, token)
				start = i
			}

			value := string(b[start : i+1])
			token := &Token{
				Type:  GetTokenType(value),
				Value: value,
			}
			tokens = append(tokens, token)
			start = i + 1
			continue
		}

		if IsAlphaNumeric(b[i]) {
			if i == len(b)-1 {
				value := string(b[start : i+1])
				token := &Token{
					Type:  GetTokenType(value),
					Value: value,
				}
				tokens = append(tokens, token)
				break
			}

			continue
		}

		// End of token
		if start != i {
			value := string(b[start:i])
			token := &Token{
				Type:  GetTokenType(value),
				Value: value,
			}
			tokens = append(tokens, token)
		}

		start = i + 1
	}

	return tokens
}

// Check if a character is in the alphabet or a number ([a-z] or [0-9])
func IsAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

// Check if a character is a special character in the sql syntax.
// `= < > ( ) * ,`
func IsSpecial(c byte) bool {
	return c == '=' || c == '<' || c == '>' || c == '(' || c == ')' || c == '*' || c == ','
}
