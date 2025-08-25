package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/scrypt"
)

func DecryptFile(srcPath, dstPath, passphrase string) error {
	// 暗号化ファイルの読み込み
	f, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open src: %w", err)
	}
	defer f.Close()

	// マジックヘッダー、ソルト、ノンス、暗号文の読み込み
	header := make([]byte, len(magic))
	if _, err = io.ReadFull(f, header); err != nil {
		return fmt.Errorf("read header: %w", err)
	}
	if string(header) != magic {
		return ErrInvalidFormat
	}

	salt := make([]byte, saltSize)
	if _, err = io.ReadFull(f, salt); err != nil {
		return fmt.Errorf("read salt: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err = io.ReadFull(f, nonce); err != nil {
		return fmt.Errorf("read nonce: %w", err)
	}

	ciphertext, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read ciphertext: %w", err)
	}

	// 鍵の導出
	key, err := scrypt.Key([]byte(passphrase), salt, scryptN, scryptR, scryptP, keySize)
	if err != nil {
		return fmt.Errorf("derive key: %w", err)
	}
	defer zero(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("gcm: %w", err)
	}

	// 復号化
	plain, err := gcm.Open(nil, nonce, ciphertext, []byte(magic))
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}

	// 書き込み
	if err = os.WriteFile(dstPath, plain, 0o600); err != nil {
		return fmt.Errorf("write dst: %w", err)
	}
	return nil
}
