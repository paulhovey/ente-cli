package crypto

import (
	"bufio"
	"bytes"
	"cli-go/utils/encoding"
	"fmt"
	"github.com/jamesruan/sodium"
	"io"
	"log"
	"os"
)

// EncryptChaCha20poly1305 encrypts the given data using the ChaCha20-Poly1305 algorithm.
// Parameters:
//   - data: The plaintext data as a byte slice.
//   - key: The key for encryption as a byte slice.
//
// Returns:
//   - A byte slice representing the encrypted data.
//   - A byte slice representing the header of the encrypted data.
//   - An error object, which is nil if no error occurs.
func EncryptChaCha20poly1305(data []byte, key []byte) ([]byte, []byte, error) {
	var buf bytes.Buffer
	encoder := sodium.MakeSecretStreamXCPEncoder(sodium.SecretStreamXCPKey{Bytes: key}, &buf)
	_, err := encoder.WriteAndClose(data)
	if err != nil {
		log.Println("Failed to write to encoder", err)
		return nil, nil, err
	}
	return buf.Bytes(), encoder.Header().Bytes, nil
}

// decryptChaCha20poly1305 decrypts the given data using the ChaCha20-Poly1305 algorithm.
// Parameters:
//   - data: The encrypted data as a byte slice.
//   - key: The key for decryption as a byte slice.
//   - nonce: The nonce for decryption as a byte slice.
//
// Returns:
//   - A byte slice representing the decrypted data.
//   - An error object, which is nil if no error occurs.
func decryptChaCha20poly1305LibSodium(data []byte, key []byte, nonce []byte) ([]byte, error) {
	reader := bytes.NewReader(data)
	header := sodium.SecretStreamXCPHeader{Bytes: nonce}
	decoder, err := sodium.MakeSecretStreamXCPDecoder(
		sodium.SecretStreamXCPKey{Bytes: key},
		reader,
		header)
	if err != nil {
		log.Println("Failed to make secret stream decoder", err)
		return nil, err
	}
	// Buffer to store the decrypted data
	decryptedData := make([]byte, len(data))
	n, err := decoder.Read(decryptedData)
	if err != nil && err != io.EOF {
		log.Println("Failed to read from decoder", err)
		return nil, err
	}
	return decryptedData[:n], nil
}

func decryptChaCha20poly1305(data []byte, key []byte, nonce []byte) ([]byte, error) {
	decryptor, err := NewDecryptor(key, nonce)
	if err != nil {
		return nil, err
	}
	decoded, _, err := decryptor.Pull(data)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func SecretBoxOpen(c []byte, n []byte, k []byte) ([]byte, error) {
	var cp sodium.Bytes = c
	res, err := cp.SecretBoxOpen(sodium.SecretBoxNonce{Bytes: n}, sodium.SecretBoxKey{Bytes: k})
	return res, err
}

func SecretBoxOpenBase64(cipher string, nonce string, k []byte) ([]byte, error) {
	return SecretBoxOpen(encoding.DecodeBase64(cipher), encoding.DecodeBase64(nonce), k)
}

func SealedBoxOpen(cipherText []byte, publicKey, masterSecret []byte) ([]byte, error) {
	var cp sodium.Bytes = cipherText
	om, err := cp.SealedBoxOpen(sodium.BoxKP{
		PublicKey: sodium.BoxPublicKey{Bytes: publicKey},
		SecretKey: sodium.BoxSecretKey{Bytes: masterSecret},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open sealed box: %v", err)
	}
	return om, nil
}

func DecryptFile(encryptedFilePath string, decryptedFilePath string, key, nonce []byte) error {
	inputFile, err := os.Open(encryptedFilePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(decryptedFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	reader := bufio.NewReader(inputFile)
	writer := bufio.NewWriter(outputFile)

	header := sodium.SecretStreamXCPHeader{Bytes: nonce}
	decoder, err := sodium.MakeSecretStreamXCPDecoder(
		sodium.SecretStreamXCPKey{Bytes: key},
		reader,
		header)
	if err != nil {
		log.Println("Failed to make secret stream decoder", err)
		return err
	}

	buf := make([]byte, decryptionBufferSize)
	for {
		n, errErr := decoder.Read(buf)
		if errErr != nil && errErr != io.EOF {
			log.Println("Failed to read from decoder", errErr)
			return errErr
		}
		if n == 0 {
			break
		}
		if _, err := writer.Write(buf[:n]); err != nil {
			log.Println("Failed to write to output file", err)
			return err
		}
		if errErr == io.EOF {
			break
		}
	}
	if err := writer.Flush(); err != nil {
		log.Println("Failed to flush writer", err)
		return err
	}
	return nil
}