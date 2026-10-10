package bootstrap

import (
	"io"
	"os"

	"go.uber.org/dig"

	"vyshkat/internal/console"
	"vyshkat/internal/fleet"
	"vyshkat/internal/servicecenter"
)

// BuildContainer creates the application's composition root.
// Business packages never resolve dependencies from the container themselves.
func BuildContainer() (*dig.Container, error) {
	container := dig.New()

	providers := []any{
		func() fleet.VehicleInspector { return servicecenter.New() },
		fleet.NewService,
		func() io.Reader { return os.Stdin },
		func() io.Writer { return os.Stdout },
		console.NewApp,
	}

	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			return nil, err
		}
	}
	return container, nil
}

// Run resolves the top-level console application and executes it.
func Run() error {
	container, err := BuildContainer()
	if err != nil {
		return err
	}

	return container.Invoke(func(app *console.App) error {
		return app.Run()
	})
}
