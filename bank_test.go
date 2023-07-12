package main

import (
	"fmt"
	"testing"
)

func Test_bankService_HandlePinganBankTransactionDetail(t *testing.T) {
	// 创建一个空的 map
	m := make(map[string]int)

	m["b"] = 3
	m["c"] = 4

	// 获取键值对
	value, err := m["a"]
	if err {
		fmt.Println("Value:", value)
	} else {
		fmt.Println("Key 'a' not found")
	}

	// 遍历所有键值对
	for key, value := range m {
		fmt.Println("Key:", key, "Value:", value)
	}

}
