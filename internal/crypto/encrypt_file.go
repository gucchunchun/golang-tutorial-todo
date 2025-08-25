package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"os"

	"golang.org/x/crypto/scrypt"
)

func EncryptFile(srcPath, dstPath, passphrase string) error {
	// ソースファイルの読み込み
	plain, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("read source: %w", err)
	}

	// ソルトの生成と鍵の導出
	salt := make([]byte, saltSize)
	if _, err = rand.Read(salt); err != nil {
		return fmt.Errorf("salt: %w", err)
	}
	key, err := scrypt.Key([]byte(passphrase), salt, scryptN, scryptR, scryptP, keySize)
	if err != nil {
		return fmt.Errorf("derive key: %w", err)
	}
	// ガーベッジコレクションに依存せず、鍵をメモリから確実に消去する
	defer zero(key)

	// AES-GCM での暗号化
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("gcm: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err = rand.Read(nonce); err != nil {
		return fmt.Errorf("nonce: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, plain, []byte(magic))

	// 暗号文を書き込むファイルの作成
	f, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("open dst: %w", err)
	}
	defer f.Close()

	// 復元に必要な情報を保存
	if _, err = f.Write([]byte(magic)); err != nil {
		return err
	}
	if _, err = f.Write(salt); err != nil {
		return err
	}
	if _, err = f.Write(nonce); err != nil {
		return err
	}
	if _, err = f.Write(ciphertext); err != nil {
		return err
	}
	return f.Sync()
}
