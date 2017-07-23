package main

import "testing"

// use values that you know are right
var mulTests = []struct {
	a, b     int
	expected int
}{
	{1, 1, 1},
	{2, 2, 4},
	{3, 3, 9},
	{4, 4, 16},
	{5, 5, 25},
}

func TestMul(t *testing.T) {
	for _, mt := range mulTests {
		if v := Mul(mt.a, mt.b); v != mt.expected {
			t.Errorf("Mul(%d, %d) returned %d, expected %d", mt.a, mt.b, v, mt.expected)
		}
	}
}

// use values that you know are right
var fibTests = []uint64{1, 1, 2, 3, 5, 8, 13, 21, 34, 55}

func TestFibFunc(t *testing.T) {
	fn := FibFunc()
	for i, v := range fibTests {
		if val := fn(); val != v {
			t.Fatalf("at index %d, expected %d, got %d.", i, v, val)
		}
	}
}

func BenchmarkFibFunc(b *testing.B) {
	fn := FibFunc()
	for i := 0; i < b.N; i++ {
		_ = fn()
	}
}
