package main

import "testing"

func BenchmarkSORT(b *testing.B) {

	RedisInit()
	for n := 0; n < b.N; n++ {
		Search([]string{"tweet", "link"}, "n", 10)
	}
}

func BenchmarkZREVRANGEThenMGET(b *testing.B) {

	RedisInit()
	for n := 0; n < b.N; n++ {
		Search2([]string{"tweet", "link"}, "n", 10)
	}
}

func BenchmarkLUA(b *testing.B) {

	RedisInit()
	for n := 0; n < b.N; n++ {
		Search3([]string{"tweet", "link"}, "n", 10)
	}
}
