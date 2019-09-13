package conman

// CryptoProvider is abstraction over crypto operations
type CryptoProvider interface {
	// Encrypt takes plain data and returns encrypted data or failure
	Encrypt([]byte) ([]byte, error)
	// Decrypt take encrypted data and returns plain data or failure
	Decrypt([]byte) ([]byte, error)
}
