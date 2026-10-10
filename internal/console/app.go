package console

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"vyshkat/internal/domain"
	"vyshkat/internal/fleet"
)

// App handles user interaction. All fleet rules remain in fleet.Service.
type App struct {
	fleet *fleet.Service
	input *bufio.Scanner
	out   io.Writer
}

func NewApp(fleetService *fleet.Service, in io.Reader, out io.Writer) *App {
	return &App{
		fleet: fleetService,
		input: bufio.NewScanner(in),
		out:   out,
	}
}

// Run reads commands until the user chooses to exit or input ends.
func (a *App) Run() error {
	fmt.Fprintln(a.out, "=== ВышКат — управление парком ===")

	for {
		a.printMenu()
		command, err := a.readLine("Выберите действие: ")
		if err != nil {
			return a.finishOnEOF(err)
		}

		switch command {
		case "0":
			fmt.Fprintln(a.out, "Выход из программы.")
			return nil
		case "1":
			err = a.addTransport()
		case "2":
			err = a.addItem()
		case "3":
			fmt.Fprintf(a.out, "Транспортных единиц в активном парке: %d\n", a.fleet.TransportCount())
		case "4":
			fmt.Fprintf(a.out, "Суммарное суточное энергопотребление: %.2f кВт·ч\n", a.fleet.TotalDailyEnergyKWh())
		case "5":
			a.printBeginnerTransports()
		case "6":
			a.printInventory()
		default:
			fmt.Fprintln(a.out, "Неизвестная команда. Выберите пункт меню от 0 до 6.")
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				return a.finishOnEOF(err)
			}
			if scanErr := a.input.Err(); scanErr != nil {
				return scanErr
			}
			// Invalid user data must not crash the application.
			a.reportError(err)
		}
	}
}

func (a *App) reportError(err error) {
	switch {
	case errors.Is(err, fleet.ErrDuplicateInventory):
		fmt.Fprintln(a.out, "Ошибка: такой инвентарный номер уже существует.")
	case errors.Is(err, domain.ErrEmptyName):
		fmt.Fprintln(a.out, "Ошибка: наименование не может быть пустым.")
	case errors.Is(err, domain.ErrEmptyInventoryNumber):
		fmt.Fprintln(a.out, "Ошибка: инвентарный номер не может быть пустым.")
	default:
		fmt.Fprintf(a.out, "Не удалось выполнить действие: %v\n", err)
	}
}

func (a *App) printMenu() {
	fmt.Fprintln(a.out, "\n1. Добавить транспорт")
	fmt.Fprintln(a.out, "2. Добавить вещь на баланс")
	fmt.Fprintln(a.out, "3. Количество транспорта")
	fmt.Fprintln(a.out, "4. Суммарное энергопотребление")
	fmt.Fprintln(a.out, "5. Транспорт для новичков")
	fmt.Fprintln(a.out, "6. Весь инвентарь")
	fmt.Fprintln(a.out, "0. Выход")
}

func (a *App) readLine(prompt string) (string, error) {
	fmt.Fprint(a.out, prompt)
	if !a.input.Scan() {
		if err := a.input.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return strings.TrimSpace(a.input.Text()), nil
}

func (a *App) readEase() (int, error) {
	for {
		value, err := a.readLine("Простота для новичка (1–10): ")
		if err != nil {
			return 0, err
		}
		ease, err := strconv.Atoi(value)
		if err == nil && ease >= 1 && ease <= 10 {
			return ease, nil
		}
		fmt.Fprintln(a.out, "Введите целое число от 1 до 10.")
	}
}

func (a *App) readOperational() (bool, error) {
	for {
		value, err := a.readLine("Транспорт исправен? (да/нет): ")
		if err != nil {
			return false, err
		}
		switch strings.ToLower(value) {
		case "да", "д", "yes", "y", "true", "1":
			return true, nil
		case "нет", "н", "no", "n", "false", "0":
			return false, nil
		default:
			fmt.Fprintln(a.out, "Ответьте 'да' или 'нет'.")
		}
	}
}

func (a *App) readEnergy() (float64, error) {
	for {
		value, err := a.readLine("Суточное энергопотребление (кВт·ч, например 1.8): ")
		if err != nil {
			return 0, err
		}
		energy, err := strconv.ParseFloat(strings.ReplaceAll(value, ",", "."), 64)
		if err == nil && !math.IsNaN(energy) && !math.IsInf(energy, 0) && energy >= 0 {
			return energy, nil
		}
		fmt.Fprintln(a.out, "Введите неотрицательное число (например, 1.8).")
	}
}

func (a *App) addTransport() error {
	fmt.Fprintln(a.out, "\nВид транспорта: 1 — электросамокат; 2 — велосипед; 3 — электровелосипед")
	kind, err := a.readLine("Тип: ")
	if err != nil {
		return err
	}
	if kind != "1" && kind != "2" && kind != "3" {
		return fmt.Errorf("неизвестный вид транспорта %q", kind)
	}

	name, err := a.readLine("Наименование: ")
	if err != nil {
		return err
	}
	id, err := a.readLine("Инвентарный номер: ")
	if err != nil {
		return err
	}
	ease, err := a.readEase()
	if err != nil {
		return err
	}
	operational, err := a.readOperational()
	if err != nil {
		return err
	}

	var transport domain.Transport
	switch kind {
	case "1":
		energy, readErr := a.readEnergy()
		if readErr != nil {
			return readErr
		}
		transport, err = domain.NewScooter(name, id, ease, operational, energy)
	case "2":
		transport, err = domain.NewBicycle(name, id, ease, operational)
	case "3":
		energy, readErr := a.readEnergy()
		if readErr != nil {
			return readErr
		}
		transport, err = domain.NewElectricBike(name, id, ease, operational, energy)
	}
	if err != nil {
		return err
	}

	// The inspection decision is delegated to the injected fleet service.
	err = a.fleet.AddTransport(transport)
	if errors.Is(err, fleet.ErrInspectionFailed) {
		fmt.Fprintf(a.out, "%s [%s]: ОТКЛОНЁН техосмотром\n", transport.Name(), transport.InventoryNumber())
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Fprintf(a.out, "%s [%s]: ПРИНЯТ в парк\n", transport.Name(), transport.InventoryNumber())
	return nil
}

func (a *App) addItem() error {
	fmt.Fprintln(a.out, "\nВид вещи: 1 — шлем; 2 — зарядная станция")
	kind, err := a.readLine("Тип: ")
	if err != nil {
		return err
	}
	if kind != "1" && kind != "2" {
		return fmt.Errorf("неизвестный вид вещи %q", kind)
	}

	name, err := a.readLine("Наименование: ")
	if err != nil {
		return err
	}
	id, err := a.readLine("Инвентарный номер: ")
	if err != nil {
		return err
	}

	var item domain.InventoryItem
	switch kind {
	case "1":
		item, err = domain.NewHelmet(name, id)
	case "2":
		energy, readErr := a.readEnergy()
		if readErr != nil {
			return readErr
		}
		item, err = domain.NewChargingStation(name, id, energy)
	}
	if err != nil {
		return err
	}
	if err := a.fleet.AddItem(item); err != nil {
		return err
	}
	fmt.Fprintf(a.out, "%s [%s]: добавлено на баланс\n", item.Name(), item.InventoryNumber())
	return nil
}

func (a *App) printBeginnerTransports() {
	fmt.Fprintf(a.out, "Техника для новичков (простота >= %d):\n", fleet.BeginnerEaseThreshold)
	transports := a.fleet.BeginnerFriendlyTransports()
	if len(transports) == 0 {
		fmt.Fprintln(a.out, "Пока нет подходящей техники.")
	}
	for _, transport := range transports {
		fmt.Fprintf(a.out, "- %s [%s], простота: %d\n", transport.Name(), transport.InventoryNumber(), transport.BeginnerEase())
	}
}

func (a *App) printInventory() {
	fmt.Fprintln(a.out, "Вся техника и вещи на балансе:")
	items := a.fleet.Inventory()
	if len(items) == 0 {
		fmt.Fprintln(a.out, "Баланс пока пуст.")
	}
	for _, item := range items {
		fmt.Fprintf(a.out, "- %s [%s]\n", item.Name(), item.InventoryNumber())
	}
}

func (a *App) finishOnEOF(err error) error {
	if !errors.Is(err, io.EOF) {
		return err
	}
	fmt.Fprintln(a.out, "\nВвод завершён.")
	return nil
}
