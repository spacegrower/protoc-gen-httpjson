package plugin

import (
	"path/filepath"
	"testing"
)

func TestPathJoin(t *testing.T) {
	t.Log(filepath.Join("..", ".."))
}
