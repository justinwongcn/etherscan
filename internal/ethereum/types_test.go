package ethereum

import (
	"testing"
	"time"
)

func TestDefaultClientOptions(t *testing.T) {
	opts := DefaultClientOptions()

	// 验证返回的选项不为nil
	if opts == nil {
		t.Errorf("DefaultClientOptions() 返回了nil")
	}

	// 验证默认值
	if opts.MaxConns <= 0 {
		t.Errorf("MaxConns 应该大于0，实际值: %d", opts.MaxConns)
	}

	if opts.IdleTimeout <= 0 {
		t.Errorf("IdleTimeout 应该大于0，实际值: %v", opts.IdleTimeout)
	}

	if opts.MaxIdleConns < 0 {
		t.Errorf("MaxIdleConns 不应该为负数，实际值: %d", opts.MaxIdleConns)
	}

	// 验证合理的默认值
	expectedMaxConns := 30
	if opts.MaxConns != expectedMaxConns {
		t.Errorf("MaxConns 期望 %d，实际 %d", expectedMaxConns, opts.MaxConns)
	}

	expectedIdleTimeout := 3 * time.Minute
	if opts.IdleTimeout != expectedIdleTimeout {
		t.Errorf("IdleTimeout 期望 %v，实际 %v", expectedIdleTimeout, opts.IdleTimeout)
	}

	expectedMaxIdleConns := 5
	if opts.MaxIdleConns != expectedMaxIdleConns {
		t.Errorf("MaxIdleConns 期望 %d，实际 %d", expectedMaxIdleConns, opts.MaxIdleConns)
	}
}

func TestClientOptions_Basic(t *testing.T) {
	// 测试基本的ClientOptions创建
	opts := ClientOptions{
		MaxConns:     5,
		IdleTimeout:  30 * time.Second,
		HealthCheck:  true,
		MaxIdleConns: 3,
	}

	if opts.MaxConns != 5 {
		t.Errorf("MaxConns 期望 5，实际 %d", opts.MaxConns)
	}

	if opts.IdleTimeout != 30*time.Second {
		t.Errorf("IdleTimeout 期望 30s，实际 %v", opts.IdleTimeout)
	}

	if !opts.HealthCheck {
		t.Errorf("HealthCheck 期望 true，实际 %v", opts.HealthCheck)
	}

	if opts.MaxIdleConns != 3 {
		t.Errorf("MaxIdleConns 期望 3，实际 %d", opts.MaxIdleConns)
	}
}

// 基准测试
func BenchmarkDefaultClientOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = DefaultClientOptions()
	}
}

func TestClientOptions_Copy(t *testing.T) {
	// 测试选项的复制和修改
	original := DefaultClientOptions()

	// 创建副本并修改
	copy := *original
	copy.MaxConns = 20
	copy.IdleTimeout = 60 * time.Second

	// 验证原始选项没有被修改
	if original.MaxConns == copy.MaxConns {
		t.Errorf("原始选项被意外修改")
	}

	if original.IdleTimeout == copy.IdleTimeout {
		t.Errorf("原始选项被意外修改")
	}

	// 验证副本的值
	if copy.MaxConns != 20 {
		t.Errorf("副本的MaxConns期望20，实际%d", copy.MaxConns)
	}

	if copy.IdleTimeout != 60*time.Second {
		t.Errorf("副本的IdleTimeout期望60s，实际%v", copy.IdleTimeout)
	}
}
