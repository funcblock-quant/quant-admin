package utils

import "testing"

func TestEncryptDecrypt(t *testing.T) {
	original := "12345" // 明文输入

	// 加密
	encrypt, err := Encrypt(original)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// 解密
	decrypt, err := Decrypt(encrypt)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	// 验证解密结果
	if decrypt != original {
		t.Errorf("Decryption mismatch: expected %s, got %s", original, decrypt)
	}

	t.Logf("Encrypt: %s, Decrypt: %s", encrypt, decrypt)
}

func TestEncryptDecryptWithSecretKey(t *testing.T) {
	original := "12345"                                     // 明文输入
	secretKey := []byte("your-32-byte-secret-key-for-aess") //输入你的32字节的密钥

	// 加密
	encrypt, err := EncryptWithSecretKey(secretKey, original)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// 解密
	decrypt, err := DecryptWithSecretKey(secretKey, encrypt)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	// 验证解密结果
	if decrypt != original {
		t.Errorf("Decryption mismatch: expected %s, got %s", original, decrypt)
	}

	t.Logf("Encrypt: %s, Decrypt: %s", encrypt, decrypt)
}
