package console

import (
	"bytes"
	"strings"
	"testing"

	"vyshkat/internal/fleet"
	"vyshkat/internal/servicecenter"
)

func TestRunDemonstratesRequiredScenarios(t *testing.T) {
	var out bytes.Buffer
	app := NewApp(fleet.NewService(servicecenter.New()), &out)

	if err := app.Run(); err != nil {
		t.Fatal(err)
	}

	text := out.String()
	checks := []string{
		"ПРИНЯТ",
		"ОТКЛОНЁН техосмотром",
		"Транспортных единиц в активном парке: 2",
		"Суммарное суточное энергопотребление: 6.0 кВт·ч",
		"Техника для новичков",
		"Вся техника и вещи на балансе",
		"TR-001",
		"IT-001",
	}
	for _, check := range checks {
		if !strings.Contains(text, check) {
			t.Fatalf("output does not contain %q\n%s", check, text)
		}
	}
}
