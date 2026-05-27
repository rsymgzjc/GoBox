package tests

import (
	"testing"

	"gobox/backend/internal/services"
)

func TestSlugifyTool(t *testing.T) {
	result, _, err := services.ExecuteForTest("slugify", services.RunToolInput{Input: "Hello Go Box"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "hello-go-box" {
		t.Fatalf("unexpected slugify result: %s", result)
	}
}
