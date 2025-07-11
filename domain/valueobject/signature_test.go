package valueobject

import (
	"testing"
)

// TestNewSignature 测试创建新的签名值对象
func TestNewSignature(t *testing.T) {
	tests := []struct {
		name    string
		v       string
		r       string
		s       string
		wantErr bool
	}{
		{
			name:    "有效的签名",
			v:       "27",
			r:       "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			s:       "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321",
			wantErr: false,
		},
		{
			name:    "EIP-155签名",
			v:       "37",
			r:       "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			s:       "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321",
			wantErr: false,
		},
		{
			name:    "空的v值",
			v:       "",
			r:       "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			s:       "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321",
			wantErr: true,
		},
		{
			name:    "空的r值",
			v:       "27",
			r:       "",
			s:       "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321",
			wantErr: true,
		},
		{
			name:    "空的s值",
			v:       "27",
			r:       "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			s:       "",
			wantErr: true,
		},
		{
			name:    "零值r",
			v:       "27",
			r:       "0x0",
			s:       "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321",
			wantErr: true,
		},
		{
			name:    "零值s",
			v:       "27",
			r:       "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			s:       "0x0",
			wantErr: true,
		},
		{
			name:    "无效的十六进制r",
			v:       "27",
			r:       "0xgg34567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			s:       "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sig, err := NewSignature(tt.v, tt.r, tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if sig.V() != tt.v {
					t.Errorf("V() = %v, want %v", sig.V(), tt.v)
				}
				if sig.R() != tt.r {
					t.Errorf("R() = %v, want %v", sig.R(), tt.r)
				}
				if sig.S() != tt.s {
					t.Errorf("S() = %v, want %v", sig.S(), tt.s)
				}
			}
		})
	}
}

// TestSignatureEquals 测试签名相等性比较
func TestSignatureEquals(t *testing.T) {
	sig1, _ := NewSignature("27", "0x1234", "0x5678")
	sig2, _ := NewSignature("27", "0x1234", "0x5678")
	sig3, _ := NewSignature("28", "0x1234", "0x5678")

	tests := []struct {
		name     string
		sig1     *Signature
		sig2     *Signature
		expected bool
	}{
		{
			name:     "相同的签名",
			sig1:     sig1,
			sig2:     sig2,
			expected: true,
		},
		{
			name:     "不同的v值",
			sig1:     sig1,
			sig2:     sig3,
			expected: false,
		},
		{
			name:     "与nil比较",
			sig1:     sig1,
			sig2:     nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sig1.Equals(tt.sig2)
			if result != tt.expected {
				t.Errorf("Equals() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestSignatureIsValid 测试签名有效性检查
func TestSignatureIsValid(t *testing.T) {
	tests := []struct {
		name     string
		sig      *Signature
		expected bool
	}{
		{
			name: "有效的签名",
			sig: &Signature{
				v: "27",
				r: "0x1234567890abcdef",
				s: "0xfedcba0987654321",
			},
			expected: true,
		},
		{
			name: "空的v值",
			sig: &Signature{
				v: "",
				r: "0x1234567890abcdef",
				s: "0xfedcba0987654321",
			},
			expected: false,
		},
		{
			name: "零值r",
			sig: &Signature{
				v: "27",
				r: "0x0",
				s: "0xfedcba0987654321",
			},
			expected: false,
		},
		{
			name: "零值s",
			sig: &Signature{
				v: "27",
				r: "0x1234567890abcdef",
				s: "0x0",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sig.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestSignatureIsZero 测试零值签名检查
func TestSignatureIsZero(t *testing.T) {
	tests := []struct {
		name     string
		sig      *Signature
		expected bool
	}{
		{
			name:     "零值签名",
			sig:      ZeroSignature(),
			expected: true,
		},
		{
			name: "非零值签名",
			sig: &Signature{
				v: "27",
				r: "0x1234",
				s: "0x5678",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sig.IsZero()
			if result != tt.expected {
				t.Errorf("IsZero() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestSignatureIsEIP155 测试EIP-155签名检查
func TestSignatureIsEIP155(t *testing.T) {
	tests := []struct {
		name     string
		sig      *Signature
		expected bool
	}{
		{
			name: "传统签名v=27",
			sig: &Signature{
				v: "27",
				r: "0x1234",
				s: "0x5678",
			},
			expected: false,
		},
		{
			name: "传统签名v=28",
			sig: &Signature{
				v: "28",
				r: "0x1234",
				s: "0x5678",
			},
			expected: false,
		},
		{
			name: "EIP-155签名",
			sig: &Signature{
				v: "37",
				r: "0x1234",
				s: "0x5678",
			},
			expected: true,
		},
		{
			name: "十六进制传统签名v=0x1b",
			sig: &Signature{
				v: "0x1b",
				r: "0x1234",
				s: "0x5678",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sig.IsEIP155Signature()
			if result != tt.expected {
				t.Errorf("IsEIP155Signature() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestSignatureString 测试字符串表示
func TestSignatureString(t *testing.T) {
	sig := &Signature{
		v: "27",
		r: "0x1234",
		s: "0x5678",
	}

	expected := "Signature{v: 27, r: 0x1234, s: 0x5678}"
	result := sig.String()

	if result != expected {
		t.Errorf("String() = %v, want %v", result, expected)
	}
}
