package security

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
)

func ConvertPKCSToOpenSSH(pkcsKey string) (string, bool) {
	block, _ := pem.Decode([]byte(pkcsKey))
	if block == nil {
		return "", false
	}

	// PEMヘッダーに応じてPKCS#1またはPKCS#8としてパース
	var parsedKey any
	var err error
	switch block.Type {
	case "RSA PUBLIC KEY": // PKCS#1
		parsedKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return "", false
		}
	case "PUBLIC KEY": // PKCS#8
		parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return "", false
		}
	default:
		return "", false
	}

	// 型アサーションでRSA公開鍵であることを確認
	rsaKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return "", false
	}

	// GoのRSA公開鍵からSSH公開鍵オブジェクトを生成
	sshKey, err := ssh.NewPublicKey(rsaKey)
	if err != nil {
		return "", false
	}

	// authorized_keys形式のバイト配列にマーシャリング
	// この処理で "ssh-rsa" のプレフィックスとBase64エンコードが行われます。
	authorizedKeyBytes := ssh.MarshalAuthorizedKey(sshKey)

	return string(authorizedKeyBytes), true
}

func ParseSSHKey(key string) (algorithm, keyBody, comment string, err error) {
	// 文字列の先頭と末尾の空白を削除し、中の連続した空白を区切り文字として分割
	parts := strings.Fields(key)

	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("invalid ssh key format: expected at least 2 parts, got %d", len(parts))
	}

	algorithm = parts[0]
	keyBody = parts[1]

	// 3つ目の要素（コメント）があれば設定し、なければデフォルトの空文字列のままにする
	if len(parts) > 2 {
		comment = strings.Join(parts[2:], " ") // コメントにスペースが含まれる場合を考慮
	}

	return
}
