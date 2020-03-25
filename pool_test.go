package wpool

import (
	"fmt"
	"testing"
)

func TestPool(t *testing.T) {
	t.Error("test error")
}

func BenchmarkPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Println("hello bench")
	}
}
