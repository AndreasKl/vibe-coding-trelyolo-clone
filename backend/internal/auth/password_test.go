package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("correctpassword")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if !CheckPassword(hash, "correctpassword") {
		t.Fatal("expected CheckPassword to return true for correct password")
	}
}

func TestCheckPasswordWrong(t *testing.T) {
	hash, err := HashPassword("correctpassword")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if CheckPassword(hash, "wrongpassword") {
		t.Fatal("expected CheckPassword to return false for wrong password")
	}
}
