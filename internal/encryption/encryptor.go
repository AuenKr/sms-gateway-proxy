package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"

	"go.uber.org/fx"

	"sms-gateway/internal/config"
)

const algorithm = "aes-256-cbc/pbkdf2-sha1"

type Encryptor struct {
	passphrase string
	iterations int
}

type NewParams struct {
	fx.In

	Config config.Config
}

type NewResult struct {
	fx.Out

	Encryptor *Encryptor
}

func New(in NewParams) NewResult {
	return NewResult{Encryptor: &Encryptor{passphrase: in.Config.Passphrase, iterations: in.Config.Iterations}}
}

func (e *Encryptor) Encrypt(plainText string) (string, error) {
	salt := make([]byte, aes.BlockSize)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	key := pbkdf2SHA1([]byte(e.passphrase), salt, e.iterations, 32)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}

	padded := pkcs7Pad([]byte(plainText), aes.BlockSize)
	cipherText := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, salt).CryptBlocks(cipherText, padded)

	return fmt.Sprintf(
		"$%s$i=%d$%s$%s",
		algorithm,
		e.iterations,
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(cipherText),
	), nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	if padding == 0 {
		padding = blockSize
	}
	return append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

func pbkdf2SHA1(password, salt []byte, iterations, keyLen int) []byte {
	hLen := sha1.Size
	numBlocks := (keyLen + hLen - 1) / hLen
	derived := make([]byte, 0, numBlocks*hLen)

	for block := 1; block <= numBlocks; block++ {
		derived = append(derived, pbkdf2Block(password, salt, iterations, block)...)
	}

	return derived[:keyLen]
}

func pbkdf2Block(password, salt []byte, iterations, blockNum int) []byte {
	mac := hmac.New(sha1.New, password)
	mac.Write(salt)
	mac.Write([]byte{byte(blockNum >> 24), byte(blockNum >> 16), byte(blockNum >> 8), byte(blockNum)})
	u := mac.Sum(nil)
	out := append([]byte(nil), u...)

	for i := 1; i < iterations; i++ {
		mac = hmac.New(sha1.New, password)
		mac.Write(u)
		u = mac.Sum(nil)
		for j := range out {
			out[j] ^= u[j]
		}
	}

	return out
}
