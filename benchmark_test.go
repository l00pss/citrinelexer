package citrinelexer

import (
	"testing"
)

func BenchmarkLexer(b *testing.B) {
	input := `SELECT users.name, users.email, profiles.bio 
	         FROM users 
	         INNER JOIN profiles ON users.id = profiles.user_id 
	         WHERE users.age >= 18 AND users.status = 'active' 
	         ORDER BY users.created_at DESC 
	         LIMIT 100 OFFSET 0;`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(input)
		for {
			tok := lexer.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

func BenchmarkSingleCharTokens(b *testing.B) {
	input := `;;;,,,())***//%`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(input)
		for {
			tok := lexer.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

func BenchmarkKeywordLookup(b *testing.B) {
	input := `SELECT FROM WHERE INSERT UPDATE DELETE CREATE TABLE TRUNCATE DROP ALTER`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(input)
		for {
			tok := lexer.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}
