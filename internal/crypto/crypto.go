package crypto

import (
	"errors"
)

const (
	// magic ：暗号化されたデータを識別するための4バイトマジックヘッダー
	magic = "ENC1"
	// saltSize ：鍵が毎度ユニークになるようにパスワードにプラスして使用するソルトのサイズ（16byteが標準）
	saltSize = 16
	// nonceSize ："number used once", AES-GCM で利用するノンス (初期化ベクタ) のサイズ （GCMでは12byteが標準）
	nonceSize = 12
	// keySize ：導出される対称鍵のサイズ (32 バイト, 256 ビット)
	keySize = 32
	// scryptN ；scrypt のコストパラメータ N (2^15) で、計算量とメモリ使用量を増加させ攻撃を困難にする
	scryptN = 1 << 15
	// scryptR ：scrypt のブロックサイズパラメータ r (8) で、メモリ使用パターンに影響する
	scryptR = 8
	// scryptP ：scrypt の並列化パラメータ p (1) で、同時並列計算スレッド数を制御する
	scryptP = 1
)

var (
	ErrInvalidFormat = errors.New("crypto: invalid encrypted file format")
)

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
