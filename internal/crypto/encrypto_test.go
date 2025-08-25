package crypto

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	setup()
	defer teardown()
	m.Run()
}
func setup() {
	clearContent()
}
func teardown() {
	clearContent()
}
func clearContent() {
	// testdata内のファイルの中身を空にする
	files, err := os.ReadDir("testdata")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		err := os.Truncate(filepath.Join("testdata", f.Name()), 0)
		if err != nil {
			panic(err)
		}
	}
}

func TestCryptFile_OK(t *testing.T) {
	type testcase struct {
		name          string
		plainText     string
		srcFile       string
		encryptedFile string
		decryptedFile string
		password      string
	}
	tests := []testcase{
		{"ok", "test", "testdata/plain.txt", "testdata/encrypted.bin", "testdata/decrypted.txt", "password"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// 準備
			err := os.WriteFile(test.srcFile, []byte(test.plainText), 0644)
			assert.NoError(t, err)

			// 暗号化
			err = EncryptFile(test.srcFile, test.encryptedFile, test.password)
			assert.NoError(t, err)

			// ファイルの内容を検証する
			encrypted, err := os.ReadFile(test.encryptedFile)
			require.NoError(t, err)
			assert.NotZero(t, len(encrypted))

			// 復号化
			err = DecryptFile(test.encryptedFile, test.decryptedFile, test.password)
			assert.NoError(t, err)

			// 元のファイルと一致するか確認する
			decrypted, err := os.ReadFile(test.decryptedFile)
			require.NoError(t, err)
			assert.NotZero(t, len(decrypted))
			assert.Equal(t, string(test.plainText), string(decrypted))
		})
	}
}
