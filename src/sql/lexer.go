package sql

type Token struct {
	Type  TokenType
	Value string
}

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

func IsAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func IsSpecial(c byte) bool {
	return c == '=' || c == '<' || c == '>' || c == '(' || c == ')' || c == '*' || c == ','
}

type TokenType int

const (
	TOKEN_TEXT TokenType = iota
	TOKEN_OPERATOR
	TOKEN_COMMA
	TOKEN_ASTERISK
	TOKEN_PARENTHESIS
)
