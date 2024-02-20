package hasher

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher is the interface for password hasing
type PasswordHasher interface {
	GenerateHashPassword(password []byte, salt []byte) ([]byte, error)
	VerifyPassword(hashedPassword []byte, password []byte, salt []byte) (bool, error)
}

// PasswordHasher is the implementation of the PasswordHasher interface
type PasswordHasherImpl struct {
	algorithm HashAlgorithm
}

// New create an object of the PasswordHasherImpl
func New(algorithm HashAlgorithm) *PasswordHasherImpl {
	return &PasswordHasherImpl{algorithm: algorithm}
}

// GenerateHashPassword generate hashed password with the given salt
func (h PasswordHasherImpl) GenerateHashPassword(password []byte, salt []byte) ([]byte, error) {
	// Check for allowed algorithm
	if !IsAlgorithmAllowed(h.algorithm) {
		return nil, fmt.Errorf("algorithm '%s' is not supported", h.algorithm)
	}

	// Verify salt and password
	if len(password) <= 0 {
		return nil, errors.New("can't hash an empty password")
	}
	if len(salt) <= 0 {
		return nil, errors.New("can't hash an empty salt")
	}

	// Create a salted password
	saltedPassword := append([]byte(password), salt...)

	var hashedPassword []byte
	var err error
	switch h.algorithm {
	case BCrypt:
		hashedPassword, err = bcrypt.GenerateFromPassword(saltedPassword, bcrypt.DefaultCost)
	case SHA256:
		hasher := sha256.New()
		hasher.Write(saltedPassword)
		hashedPassword = hasher.Sum(nil)
	default:
		return nil, fmt.Errorf("hash algorithm '%s' not available", h.algorithm)
	}

	// Check for errors
	if err != nil {
		return nil, fmt.Errorf("failed hash password using algorithm (%s)", h.algorithm)
	}

	return append([]byte(h.algorithm), hashedPassword...), nil
}

// VerifyPassword compare hashed password with the clear password and given salt
func (h PasswordHasherImpl) VerifyPassword(hashedPassword []byte, password []byte, salt []byte) (bool, error) {
	// Check for allowed algorithm
	if !IsAlgorithmAllowed(h.algorithm) {
		return false, fmt.Errorf("algorithm '%s' is not supported", h.algorithm)
	}

	// Verify salt and password
	if len(hashedPassword) <= 0 {
		return false, errors.New("can't hash an empty hashedPassword")
	}
	if len(password) <= 0 {
		return false, errors.New("can't hash an empty password")
	}
	if len(salt) <= 0 {
		return false, errors.New("can't hash an empty salt")
	}

	// Split the algorithm used to hash password and the actual generated hash password
	storedAlgorithm := HashAlgorithm(hashedPassword[0:len(h.algorithm)])
	hashedPassword = hashedPassword[len(h.algorithm):]

	// Ensure the stored algorithm matches the expected algorithm
	if h.algorithm != storedAlgorithm {
		return false, fmt.Errorf("mismatched algorithm stored '%s', found for verifying '%s'", storedAlgorithm, h.algorithm)
	}

	// Create a salted password
	saltedPassword := append([]byte(password), salt...)

	switch h.algorithm {
	case BCrypt:
		return bcrypt.CompareHashAndPassword(hashedPassword, saltedPassword) == nil, nil
	case SHA256:
		hasher := sha256.New()
		hasher.Write(saltedPassword)
		return bytes.Equal(hashedPassword, hasher.Sum(nil)), nil
	default:
		return false, fmt.Errorf("hash algorithm '%s' not available", h.algorithm)
	}
}

// GenerateRandomSalt helper function to generate a random salt
func GenerateRandomSalt() []byte {
	return GenerateRandomSalts(DefaultSaltSize)
}

// GenerateRandomSalt helper function to generate a random salt with the specified salt size
func GenerateRandomSalts(size int) []byte {
	if size <= 0 {
		size = DefaultSaltSize
	}

	salt := make([]byte, size)
	if _, err := rand.Read(salt[:]); err != nil {
		log.Errorf("failed to generate random salt with size %d", size)
		return nil
	}

	return salt
}
