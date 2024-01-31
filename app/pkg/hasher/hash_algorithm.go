package hasher

import "slices"

// Defined type for hashing algorithm
type HashAlgorithm string

const (
	// Hash algorithms
	BCrypt HashAlgorithm = "bcrypt"
	SHA256 HashAlgorithm = "sha256"
	// .. add more algorithm if needed

	// Default salt size
	DefaultSaltSize = 16
)

func IsAlgorithmAllowed(algorithm HashAlgorithm) bool {
	allowedAlgorithms := []HashAlgorithm{
		BCrypt,
		SHA256,
	}
	return slices.Contains(allowedAlgorithms, algorithm)
}
