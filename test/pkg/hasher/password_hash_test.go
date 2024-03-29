package hasher_test

import (
	"encoding/hex"
	"testing"

	"github.com/budgetin-app/user-service/app/pkg/hasher"
)

func TestGenerateRandomSalt(t *testing.T) {
	salt := hasher.GenerateRandomSalt()

	if len(salt) != hasher.DefaultSaltSize {
		t.Errorf("GenerateRandomSalt failed: expected salt size %d, got %d", hasher.DefaultSaltSize, len(salt))
	}
}

func TestGenerateRandomSalts(t *testing.T) {
	size := 32
	salt := hasher.GenerateRandomSalts(size)

	if len(salt) != size {
		t.Errorf("GenerateRandomSalts failed: expected salt size %d, got %d", size, len(salt))
	}
}

func TestGenerateHashPasswordBCrypt(t *testing.T) {
	password := []byte("secretpassword123")
	salt := hasher.GenerateRandomSalt()

	hasher := hasher.New(hasher.BCrypt)

	// Test GenerateHashPassword with BCrypt
	hashedPassword, err := hasher.GenerateHashPassword(password, salt)
	if err != nil {
		t.Errorf("GenerateHashPassword failed: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Error("GenerateHashPassword failed: hashed password is empty")
	}
}

func TestGenerateHashPasswordSHA256(t *testing.T) {
	password := []byte("secretpassword123")
	salt := hasher.GenerateRandomSalt()

	hasher := hasher.New(hasher.SHA256)

	// Test GenerateHashPassword with SHA256
	hashedPassword, err := hasher.GenerateHashPassword(password, salt)
	if err != nil {
		t.Errorf("GenerateHashPassword failed: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Error("GenerateHashPassword failed: hashed password is empty")
	}
}

func TestVerifyPasswordBCrypt(t *testing.T) {
	password := []byte("secretpassword123")
	salt := hasher.GenerateRandomSalt()

	hasher := hasher.New(hasher.BCrypt)

	// Generate hashed password with BCrypt
	hashedPassword, err := hasher.GenerateHashPassword(password, salt)
	if err != nil {
		t.Errorf("GenerateHashPassword failed: %v", err)
	}

	// Test VerifyPassword with BCrypt
	match, err := hasher.VerifyPassword(hashedPassword, password, salt)
	if err != nil {
		t.Errorf("VerifyPassword failed: %v", err)
	}

	if !match {
		t.Error("VerifyPassword failed: passwords do not match for BCrypt")
	}
}

func TestVerifyPasswordSHA256(t *testing.T) {
	password := []byte("secretpassword123")
	salt := hasher.GenerateRandomSalt()

	hasher := hasher.New(hasher.SHA256)

	// Generate hashed password with SHA256
	hashedPassword, err := hasher.GenerateHashPassword(password, salt)
	if err != nil {
		t.Errorf("GenerateHashPassword failed: %v", err)
	}

	// Test VerifyPassword with SHA256
	match, err := hasher.VerifyPassword(hashedPassword, password, salt)
	if err != nil {
		t.Errorf("VerifyPassword failed: %v", err)
	}

	if !match {
		t.Error("VerifyPassword failed: passwords do not match for SHA256")
	}
}

func TestGenerateHashPasswordInvalidAlgorithm(t *testing.T) {
	password := []byte("secretpassword123")
	salt := hasher.GenerateRandomSalt()

	hasher := hasher.New("invalid_algorithm")

	_, err := hasher.GenerateHashPassword(password, salt)
	if err == nil {
		t.Error("GenerateHashPassword did not return error for invalid algorithm")
	}
}

func TestVerifyPasswordInvalidAlgorithm(t *testing.T) {
	password := []byte("secretpassword123")
	salt := hasher.GenerateRandomSalt()

	hasher := hasher.New("invalid_algorithm")

	_, err := hasher.VerifyPassword([]byte("hashedPassword"), password, salt)
	if err == nil {
		t.Error("VerifyPassword did not return error for invalid algorithm")
	}
}

func TestGenerateHashPasswordEmptyPassword(t *testing.T) {
	emptyPassword := []byte("")
	salt := hasher.GenerateRandomSalt()

	hasher := hasher.New(hasher.BCrypt)

	// Test GenerateHashPassword with an empty password
	_, err := hasher.GenerateHashPassword(emptyPassword, salt)
	if err == nil {
		t.Error("GenerateHashPassword did not return error for empty password")
	}
}

func TestVerifyPasswordEmptyHashedPassword(t *testing.T) {
	password := []byte("secretpassword123")
	salt := hasher.GenerateRandomSalt()

	hasher := hasher.New(hasher.BCrypt)

	// Test VerifyPassword with an empty hashed password
	match, err := hasher.VerifyPassword([]byte(""), password, salt)
	if err == nil {
		t.Error("VerifyPassword did not return error for empty hashed password")
	}

	if match {
		t.Error("VerifyPassword unexpectedly succeeded with empty hashed password")
	}
}

func TestGenerateHashPasswordEmptySalt(t *testing.T) {
	password := []byte("secretpassword123")
	salt := []byte{}

	hasher := hasher.New(hasher.BCrypt)

	// Test GenerateHashPassword with an empty salt
	_, err := hasher.GenerateHashPassword(password, salt)
	if err == nil {
		t.Error("GenerateHashPassword did not return error for empty salt")
	}
}

func TestVerifyPasswordEmptySalt(t *testing.T) {
	password := []byte("secretpassword123")
	salt := []byte{}

	hasher := hasher.New(hasher.BCrypt)

	// Test VerifyPassword with an empty salt
	_, err := hasher.VerifyPassword([]byte("hashedPassword"), password, salt)
	if err == nil {
		t.Error("VerifyPassword did not return error for empty salt")
	}
}

func TestVerifyPasswordMismatchedAlgorithm(t *testing.T) {
	password := []byte("secretpassword123")
	salt := hasher.GenerateRandomSalt()

	hasherBCrypt := hasher.New(hasher.BCrypt)
	hasherSHA256 := hasher.New(hasher.SHA256)

	// Test GenerateHashPassword with BCrypt
	hashedPasswordBCrypt, err := hasherBCrypt.GenerateHashPassword(password, salt)
	if err != nil {
		t.Errorf("GenerateHashPassword failed: %v", err)
	}

	// Test VerifyPassword with SHA256 and mismatched algorithm
	match, err := hasherSHA256.VerifyPassword(hashedPasswordBCrypt, password, salt)
	if err == nil {
		t.Error("VerifyPassword did not return error for mismatched algorithm")
	}

	if match {
		t.Error("VerifyPassword unexpectedly succeeded with mismatched algorithm")
	}
}

func TestGenerateHashPasswordAndVerifyPasswordEmptySaltAndEmptyPassword(t *testing.T) {
	emptyPassword := []byte("")
	emptySalt := []byte("")

	hasherBCrypt := hasher.New(hasher.BCrypt)

	// Test GenerateHashPassword with an empty password and empty salt
	_, err := hasherBCrypt.GenerateHashPassword(emptyPassword, emptySalt)
	if err == nil {
		t.Error("GenerateHashPassword did not return error for empty password and empty salt")
	}

	// Test VerifyPassword with an empty hashed password, empty password, and empty salt
	match, err := hasherBCrypt.VerifyPassword([]byte(""), emptyPassword, emptySalt)
	if err == nil {
		t.Error("VerifyPassword did not return error for empty hashed password, empty password, and empty salt")
	}

	if match {
		t.Error("VerifyPassword unexpectedly succeeded with empty hashed password, empty password, and empty salt")
	}
}

func TestGenerateHashPasswordAndVerifyPasswordMismatchedSalt(t *testing.T) {
	password := []byte("secretpassword123")
	salt1 := hasher.GenerateRandomSalt()
	salt2 := hasher.GenerateRandomSalt()

	hasherBCrypt := hasher.New(hasher.BCrypt)

	// Test GenerateHashPassword with BCrypt
	hashedPassword, err := hasherBCrypt.GenerateHashPassword(password, salt1)
	if err != nil {
		t.Errorf("GenerateHashPassword failed: %v", err)
	}

	// Test VerifyPassword with BCrypt and mismatched salt
	match, _ := hasherBCrypt.VerifyPassword(hashedPassword, password, salt2)
	if match {
		t.Error("VerifyPassword unexpectedly succeeded with mismatched salt")
	}
}

func TestGenerateHashPasswordAndVerifyPasswordEmptyHashedPasswordAndMismatchedAlgorithm(t *testing.T) {
	password := []byte("secretpassword123")
	salt := hasher.GenerateRandomSalt()

	hasherBCrypt := hasher.New(hasher.BCrypt)
	hasherSHA256 := hasher.New(hasher.SHA256)

	// Test GenerateHashPassword with BCrypt
	_, err := hasherBCrypt.GenerateHashPassword(password, salt)
	if err != nil {
		t.Errorf("GenerateHashPassword failed: %v", err)
	}

	// Test VerifyPassword with an empty hashed password and mismatched algorithm
	match, err := hasherSHA256.VerifyPassword([]byte(""), password, salt)
	if err == nil {
		t.Error("VerifyPassword did not return error for empty hashed password and mismatched algorithm")
	}

	if match {
		t.Error("VerifyPassword unexpectedly succeeded with empty hashed password and mismatched algorithm")
	}
}

func TestVerifyPasswordStoredInDatabase(t *testing.T) {
	// Value in database found according to the specified user:
	// - hashedPassword: bcrypt$2a$10$AsYbsxQFNFvhsINCe.lF4.tbOaE2wuqhwfNJ9gICYFlcaFpq0ZrMS
	// - password: password@123
	// - salt: 933c48ba9868b52020bf3c7eddd91c25
	hashedPassword := []byte("bcrypt$2a$10$AsYbsxQFNFvhsINCe.lF4.tbOaE2wuqhwfNJ9gICYFlcaFpq0ZrMS")
	password := []byte("password@123")

	// Password stored in database should be decoded
	salt, err := hex.DecodeString("933c48ba9868b52020bf3c7eddd91c25")

	if err != nil {
		t.Errorf("VerifyPassword failed decode salt: %v", err)
	}

	hash := hasher.New(hasher.HashAlgorithm("bcrypt"))

	match, err := hash.VerifyPassword(hashedPassword, password, salt)
	if err != nil {
		t.Error("VerifyPassword unexpected error")
	}
	if !match {
		t.Error("VerifyPassword password verification failed")
	}
}
