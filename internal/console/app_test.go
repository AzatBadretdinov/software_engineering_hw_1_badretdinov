package console

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"vyshkat/internal/fleet"
	"vyshkat/internal/servicecenter"
)

func runSession(t *testing.T, input string) (string, *fleet.Service) {
	t.Helper()
	var out bytes.Buffer
	service := fleet.NewService(servicecenter.New())
	app := NewApp(service, strings.NewReader(input), &out)
	if err := app.Run(); err != nil {
		t.Fatalf("Run() returned error: %v\nOutput:\n%s", err, out.String())
	}
	return out.String(), service
}

func TestRunDemonstratesRequiredScenariosFromUserInput(t *testing.T) {
	// Arrange: commands are typed by the user rather than hardcoded in App.
	input := strings.Join([]string{
		"1", "1", "Электросамокат V1", "TR-001", "8", "да", "1.8", // accepted
		"1", "3", "Электровелосипед V1", "TR-002", "5", "нет", "2.4", // rejected
		"1", "2", "Городской велосипед", "TR-003", "7", "да", // accepted
		"2", "1", "Шлем M", "IT-001", // item without energy
		"2", "2", "Зарядная станция", "IT-002", "4.2", // item with energy
		"3", "4", "5", "6", "0",
	}, "\n") + "\n"

	// Act
	output, service := runSession(t, input)

	// Assert
	for _, expected := range []string{
		"TR-001]: ПРИНЯТ",
		"TR-002]: ОТКЛОНЁН техосмотром",
		"TR-003]: ПРИНЯТ",
		"IT-001]: добавлено на баланс",
		"IT-002]: добавлено на баланс",
		"Транспортных единиц в активном парке: 2",
		"Суммарное суточное энергопотребление: 6.00 кВт·ч",
		"Техника для новичков",
		"Вся техника и вещи на балансе",
		"Выход из программы.",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("output does not contain %q\n%s", expected, output)
		}
	}
	if service.TransportCount() != 2 || len(service.Inventory()) != 4 {
		t.Errorf("unexpected state: transports=%d, inventory=%d", service.TransportCount(), len(service.Inventory()))
	}
	if strings.Contains(output, "- Электровелосипед V1 [TR-002]") {
		t.Error("rejected transport must not appear in inventory report")
	}
}

func TestRunRecoversFromInvalidCommandsAndDuplicateIDs(t *testing.T) {
	input := strings.Join([]string{
		"bad",        // unknown menu command
		"1", "wrong", // unknown transport type
		"1", "2", "Велосипед", "TR-1", "wrong", "11", "6", "maybe", "да", // retry invalid ease and boolean
		"2", "1", "Шлем", "TR-1", // duplicate inventory number
		"3", "6", "0",
	}, "\n") + "\n"

	output, service := runSession(t, input)
	for _, expected := range []string{
		"Неизвестная команда",
		"неизвестный вид транспорта",
		"Введите целое число от 1 до 10",
		"Ответьте 'да' или 'нет'",
		"такой инвентарный номер уже существует",
		"Транспортных единиц в активном парке: 1",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("output does not contain %q\n%s", expected, output)
		}
	}
	if service.TransportCount() != 1 || len(service.Inventory()) != 1 {
		t.Fatal("duplicate item or invalid transport changed fleet state")
	}
}

func TestRunEmptyInventoryAndEndOfInput(t *testing.T) {
	// EOF is a valid way to finish input, e.g. when commands come from a file.
	output, service := runSession(t, "3\n4\n5\n6\n")
	for _, expected := range []string{
		"Транспортных единиц в активном парке: 0",
		"Суммарное суточное энергопотребление: 0.00 кВт·ч",
		"Пока нет подходящей техники.",
		"Баланс пока пуст.",
		"Ввод завершён.",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("missing %q", expected)
		}
	}
	if len(service.Inventory()) != 0 {
		t.Fatal("empty session must not create inventory")
	}
}

func TestRunRejectsInvalidItemAndRetriesEnergyInput(t *testing.T) {
	input := strings.Join([]string{
		"2", "1", "  ", "IT-empty", // validation error for empty item name
		"2", "2", "Станция", "IT-2", "abc", "-1", "NaN", "2,5", // invalid energy, comma accepted
		"4", "6", "0",
	}, "\n") + "\n"
	output, service := runSession(t, input)
	if !strings.Contains(output, "наименование не может быть пустым") {
		t.Errorf("missing validation error: %s", output)
	}
	if strings.Count(output, "Введите неотрицательное число") != 3 {
		t.Errorf("expected three invalid energy prompts: %s", output)
	}
	if !strings.Contains(output, "Суммарное суточное энергопотребление: 2.50 кВт·ч") {
		t.Errorf("incorrect consumption: %s", output)
	}
	if len(service.Inventory()) != 1 {
		t.Fatalf("expected 1 valid inventory item, got %d", len(service.Inventory()))
	}
}

func TestRunEndsGracefullyDuringIncompleteInput(t *testing.T) {
	output, service := runSession(t, "1\n1\nСамокат\n")
	if !strings.Contains(output, "Ввод завершён.") {
		t.Errorf("unexpected output: %s", output)
	}
	if service.TransportCount() != 0 {
		t.Fatal("incomplete input must not create a transport")
	}
}

func TestExampleInputFileRunsAllReports(t *testing.T) {
	data, err := os.ReadFile("../../examples/demo-input.txt")
	if err != nil {
		t.Fatal(err)
	}
	output, service := runSession(t, string(data))
	if service.TransportCount() != 2 || len(service.Inventory()) != 4 {
		t.Fatalf("unexpected example state: transports=%d, inventory=%d", service.TransportCount(), len(service.Inventory()))
	}
	if !strings.Contains(output, "Суммарное суточное энергопотребление: 6.00 кВт·ч") ||
		!strings.Contains(output, "Выход из программы.") {
		t.Errorf("example scenario did not finish correctly: %s", output)
	}
}
