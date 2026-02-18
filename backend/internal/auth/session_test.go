package auth

import "testing"

func TestGenerateToken(t *testing.T) {
	tok1, err := GenerateToken()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(tok1) != 64 {
		t.Fatalf("token length = %d, want 64", len(tok1))
	}

	tok2, err := GenerateToken()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if tok1 == tok2 {
		t.Fatal("two calls produced the same token")
	}
}
