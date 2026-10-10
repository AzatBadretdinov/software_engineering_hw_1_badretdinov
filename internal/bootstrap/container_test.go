package bootstrap

import (
	"testing"

	"vyshkat/internal/console"
)

func TestContainerResolvesConsoleApp(t *testing.T) {
	container, err := BuildContainer()
	if err != nil {
		t.Fatal(err)
	}

	if err := container.Invoke(func(app *console.App) {
		if app == nil {
			t.Fatal("container resolved nil app")
		}
	}); err != nil {
		t.Fatal(err)
	}
}
