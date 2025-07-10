package valueobject

import (
	"testing"
)

func TestNewHash(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		expected  string
	}{
		{
			name:      "有效的哈希值（带0x前缀）",
			input:     "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			wantError: false,
			expected:  "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		},
		{
			name:      "有效的哈希值（不带0x前缀）",
			input:     "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			wantError: false,
			expected:  "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		},
		{
			name:      "有效的哈希值（大写字母）",
			input:     "0x1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF",
			wantError: false,
			expected:  "0x1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF",
		},
		{
			name:      "空字符串",
			input:     "",
			wantError: true,
		},
		{
			name:      "长度不正确（太短）",
			input:     "0x1234",
			wantError: true,
		},
		{
			name:      "长度不正确（太长）",
			input:     "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef12",
			wantError: true,
		},
		{
			name:      "包含无效字符",
			input:     "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdeg",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := NewHash(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("NewHash() 期望返回错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("NewHash() 返回了意外的错误: %v", err)
				return
			}

			if hash.String() != tt.expected {
				t.Errorf("NewHash() = %v, 期望 %v", hash.String(), tt.expected)
			}
		})
	}
}

func TestHashEquals(t *testing.T) {
	hash1, _ := NewHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	hash2, _ := NewHash("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	hash3, _ := NewHash("0x9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba")

	if !hash1.Equals(hash2) {
		t.Errorf("相同的哈希值应该相等")
	}

	if hash1.Equals(hash3) {
		t.Errorf("不同的哈希值不应该相等")
	}
}

func TestHashIsZero(t *testing.T) {
	zeroHash, _ := NewHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	nonZeroHash, _ := NewHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	if !zeroHash.IsZero() {
		t.Errorf("零哈希应该返回true")
	}

	if nonZeroHash.IsZero() {
		t.Errorf("非零哈希应该返回false")
	}
}

func TestNewHashFromBytes(t *testing.T) {
	bytes := make([]byte, 32)
	for i := 0; i < 32; i++ {
		bytes[i] = byte(i)
	}

	hash, err := NewHashFromBytes(bytes)
	if err != nil {
		t.Errorf("NewHashFromBytes() 返回了意外的错误: %v", err)
	}

	expected := "0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"
	if hash.String() != expected {
		t.Errorf("NewHashFromBytes() = %v, 期望 %v", hash.String(), expected)
	}

	// 测试错误的字节长度
	wrongBytes := make([]byte, 16)
	_, err = NewHashFromBytes(wrongBytes)
	if err == nil {
		t.Errorf("NewHashFromBytes() 应该对错误的字节长度返回错误")
	}
}

func TestHash_Value(t *testing.T) {
	hashStr := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	hash, err := NewHash(hashStr)
	if err != nil {
		t.Fatalf("NewHash() 失败: %v", err)
	}

	if hash.Value() != hashStr {
		t.Errorf("Value() = %v, 期望 %v", hash.Value(), hashStr)
	}
}

func TestHash_Bytes(t *testing.T) {
	hashStr := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	hash, err := NewHash(hashStr)
	if err != nil {
		t.Fatalf("NewHash() 失败: %v", err)
	}

	bytes := hash.Bytes()
	if len(bytes) != 32 {
		t.Errorf("Bytes() 长度 = %v, 期望 32", len(bytes))
	}

	// 验证第一个字节
	if bytes[0] != 0x12 {
		t.Errorf("Bytes()[0] = %v, 期望 0x12", bytes[0])
	}
}
