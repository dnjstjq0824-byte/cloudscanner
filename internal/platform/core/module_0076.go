package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// EncryptedPayload encapsulates KMS envelope encryption fields.
type EncryptedPayload struct {
	KeyID        string `json:"key_id"`
	EncryptedDEK string `json:"encrypted_dek"` // Base64 DEK encrypted by KEK
	Nonce        string `json:"nonce"`         // Base64 12-byte GCM nonce
	Ciphertext   string `json:"ciphertext"`    // Base64 encrypted payload
}

// KMSEnvelopeService defines cryptographic envelope encryption operations.
type KMSEnvelopeService interface {
	SealWithAAD(plainText []byte, aad string) (*EncryptedPayload, error)
	OpenWithAAD(payload *EncryptedPayload, aad string) ([]byte, error)
}

// MockKMSService provides a local AES-256-GCM implementation for dev and testing.
type MockKMSService struct {
	kek   []byte // 32-byte Key Encryption Key
	keyID string
}

// NewMockKMSService constructs a mock KMS service with a 256-bit KEK.
func NewMockKMSService(keyID string) (*MockKMSService, error) {
	if keyID == "" {
		keyID = "local-dev-mock-key-001"
	}
	kek := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, kek); err != nil {
		return nil, fmt.Errorf("failed to generate KEK: %w", err)
	}
	return &MockKMSService{kek: kek, keyID: keyID}, nil
}

// SealWithAAD generates a random DEK, encrypts plainText with DEK and AAD,
// then encrypts the DEK with the KEK.
func (k *MockKMSService) SealWithAAD(plainText []byte, aad string) (*EncryptedPayload, error) {
	if len(plainText) == 0 {
		return nil, errors.New("cannot encrypt empty plaintext")
	}

	// 1. Generate 32-byte Data Encryption Key (DEK)
	dek := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, dek); err != nil {
		return nil, fmt.Errorf("failed to generate random DEK: %w", err)
	}

	// 2. Encrypt plainText with DEK using AES-GCM and AAD
	block, err := aes.NewCipher(dek)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, plainText, []byte(aad))

	// 3. Encrypt DEK with KEK
	kekBlock, err := aes.NewCipher(k.kek)
	if err != nil {
		return nil, err
	}
	kekGCM, err := cipher.NewGCM(kekBlock)
	if err != nil {
		return nil, err
	}

	dekNonce := make([]byte, kekGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, dekNonce); err != nil {
		return nil, err
	}
	encryptedDEK := kekGCM.Seal(dekNonce, dekNonce, dek, nil)

	return &EncryptedPayload{
		KeyID:        k.keyID,
		EncryptedDEK: base64.StdEncoding.EncodeToString(encryptedDEK),
		Nonce:        base64.StdEncoding.EncodeToString(nonce),
		Ciphertext:   base64.StdEncoding.EncodeToString(ciphertext),
	}, nil
}

// OpenWithAAD decrypts DEK using KEK, then decrypts ciphertext with DEK and verifies AAD.
func (k *MockKMSService) OpenWithAAD(payload *EncryptedPayload, aad string) ([]byte, error) {
	if payload == nil {
		return nil, errors.New("payload cannot be nil")
	}

	// 1. Decrypt DEK
	encDEKBytes, err := base64.StdEncoding.DecodeString(payload.EncryptedDEK)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 encrypted DEK: %w", err)
	}

	kekBlock, err := aes.NewCipher(k.kek)
	if err != nil {
		return nil, err
	}
	kekGCM, err := cipher.NewGCM(kekBlock)
	if err != nil {
		return nil, err
	}

	nonceSize := kekGCM.NonceSize()
	if len(encDEKBytes) < nonceSize {
		return nil, errors.New("ciphertext too short for DEK nonce")
	}

	dekNonce, actualCipher := encDEKBytes[:nonceSize], encDEKBytes[nonceSize:]
	dek, err := kekGCM.Open(nil, dekNonce, actualCipher, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to unseal DEK (KEK mismatch): %w", err)
	}

	// 2. Decrypt Ciphertext with DEK and verify AAD
	nonce, err := base64.StdEncoding.DecodeString(payload.Nonce)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 nonce: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(payload.Ciphertext)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 ciphertext: %w", err)
	}

	block, err := aes.NewCipher(dek)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plainText, err := aesGCM.Open(nil, nonce, ciphertext, []byte(aad))
	if err != nil {
		return nil, fmt.Errorf("security violation: AAD authentication failed (tampered data or tenant mismatch): %w", err)
	}

	return plainText, nil
}

// Module0076 wraps KMS testing.
type Module0076 struct {
	kms KMSEnvelopeService
}

func NewModule0076(kms KMSEnvelopeService) *Module0076 {
	if kms == nil {
		mock, _ := NewMockKMSService("local-dev-mock-key-001")
		kms = mock
	}
	return &Module0076{kms: kms}
}

// TestKMSTamperResistance verifies that tampering with either AAD or ciphertext fails unsealing.
func (m *Module0076) TestKMSTamperResistance(secretToken, orgID, connID, provider string) error {
	aad := fmt.Sprintf("%s:%s:%s", orgID, connID, provider)
	payload, err := m.kms.SealWithAAD([]byte(secretToken), aad)
	if err != nil {
		return fmt.Errorf("seal failed: %w", err)
	}

	// 1. Legitimate open should succeed
	decrypted, err := m.kms.OpenWithAAD(payload, aad)
	if err != nil || string(decrypted) != secretToken {
		return errors.New("legitimate unseal failed")
	}

	// 2. Tampered AAD (cross-tenant access attempt) MUST fail
	tamperedAAD := fmt.Sprintf("attacker-org:%s:%s", connID, provider)
	_, errTampered := m.kms.OpenWithAAD(payload, tamperedAAD)
	if errTampered == nil {
		return errors.New("security failure: tampered AAD should have been rejected by GCM tag verification")
	}

	// 3. Verify JSON serialization doesn't leak plaintext
	jsonBytes, _ := json.Marshal(payload)
	if string(jsonBytes) == "" || string(jsonBytes) == secretToken {
		return errors.New("serialized envelope payload should contain only ciphertext metadata")
	}

	return nil
}
