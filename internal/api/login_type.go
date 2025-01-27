package api

import (
	"fmt"

	"github.com/google/uuid"
)

type SRPAttributes struct {
	SRPUserID         uuid.UUID `json:"srpUserID" binding:"required"`
	SRPSalt           string    `json:"srpSalt"  binding:"required"`
	MemLimit          int       `json:"memLimit" binding:"required"`
	OpsLimit          int       `json:"opsLimit" binding:"required"`
	KekSalt           string    `json:"kekSalt" binding:"required"`
	IsEmailMFAEnabled bool      `json:"isEmailMFAEnabled" binding:"required"`
}

type ApiError struct {
	Message    string
	StatusCode int
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("status %d with err: %s", e.StatusCode, e.Message)
}

type CreateSRPSessionResponse struct {
	SessionID uuid.UUID `json:"sessionID" binding:"required"`
	SRPB      string    `json:"srpB" binding:"required"`
}

// KeyAttributes stores the key related attributes for a user
type KeyAttributes struct {
	KEKSalt                  string `json:"kekSalt" binding:"required"`
	KEKHash                  string `json:"kekHash"`
	EncryptedKey             string `json:"encryptedKey" binding:"required"`
	KeyDecryptionNonce       string `json:"keyDecryptionNonce" binding:"required"`
	PublicKey                string `json:"publicKey" binding:"required"`
	EncryptedSecretKey       string `json:"encryptedSecretKey" binding:"required"`
	SecretKeyDecryptionNonce string `json:"secretKeyDecryptionNonce" binding:"required"`
	MemLimit                 int    `json:"memLimit" binding:"required"`
	OpsLimit                 int    `json:"opsLimit" binding:"required"`
}

type AuthorizationResponse struct {
	ID                 int64          `json:"id"`
	KeyAttributes      *KeyAttributes `json:"keyAttributes,omitempty"`
	EncryptedToken     string         `json:"encryptedToken,omitempty"`
	Token              string         `json:"token,omitempty"`
	TwoFactorSessionID string         `json:"twoFactorSessionID"`
	// SrpM2 is sent only if the user is logging via SRP
	// SrpM2 is the SRP M2 value aka the proof that the server has the verifier
	SrpM2 *string `json:"srpM2,omitempty"`
}

func (a *AuthorizationResponse) IsMFARequired() bool {
	return a.TwoFactorSessionID != ""
}
