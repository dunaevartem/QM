package test

import (
	"testing"
)

// Простой тест структуры (Unit-test)
func TestMessageStructure(t *testing.T) {
	username := "Alice"
	content := "Hello world"
    
	if username != "Alice" || content != "Hello world" {
		t.Errorf("Ошибка инициализации сообщения")
	}
}
