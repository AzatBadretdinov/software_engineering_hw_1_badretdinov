package console

import (
	"errors"
	"fmt"
	"io"

	"vyshkat/internal/domain"
	"vyshkat/internal/fleet"
)

// App demonstrates every required scenario in a console application.
type App struct {
	fleet *fleet.Service
	out   io.Writer
}

func NewApp(fleetService *fleet.Service, out io.Writer) *App {
	return &App{fleet: fleetService, out: out}
}

func (a *App) Run() error {
	fmt.Fprintln(a.out, "=== ВышКат — учёт парка ===")
	fmt.Fprintln(a.out)

	scooter, err := domain.NewScooter("Электросамокат V1", "TR-001", 8, true, 1.8)
	if err != nil {
		return err
	}
	brokenBike, err := domain.NewElectricBike("Электровелосипед V1", "TR-002", 5, false, 2.4)
	if err != nil {
		return err
	}
	bicycle, err := domain.NewBicycle("Городской велосипед", "TR-003", 7, true)
	if err != nil {
		return err
	}

	fmt.Fprintln(a.out, "Приём техники после техосмотра:")
	for _, transport := range []domain.Transport{scooter, brokenBike, bicycle} {
		if err := a.addTransportWithReport(transport); err != nil {
			return err
		}
	}

	helmet, err := domain.NewHelmet("Шлем M", "IT-001")
	if err != nil {
		return err
	}
	station, err := domain.NewChargingStation("Зарядная станция", "IT-002", 4.2)
	if err != nil {
		return err
	}
	if err := a.fleet.AddItem(helmet); err != nil {
		return err
	}
	if err := a.fleet.AddItem(station); err != nil {
		return err
	}

	fmt.Fprintln(a.out)
	fmt.Fprintf(a.out, "Транспортных единиц в активном парке: %d\n", a.fleet.TransportCount())
	fmt.Fprintf(a.out, "Суммарное суточное энергопотребление: %.1f кВт·ч\n", a.fleet.TotalDailyEnergyKWh())

	fmt.Fprintln(a.out)
	fmt.Fprintf(a.out, "Техника для новичков (простота >= %d):\n", fleet.BeginnerEaseThreshold)
	for _, transport := range a.fleet.BeginnerFriendlyTransports() {
		fmt.Fprintf(a.out, "- %s [%s], простота: %d\n", transport.Name(), transport.InventoryNumber(), transport.BeginnerEase())
	}

	fmt.Fprintln(a.out)
	fmt.Fprintln(a.out, "Вся техника и вещи на балансе:")
	for _, item := range a.fleet.Inventory() {
		fmt.Fprintf(a.out, "- %s [%s]\n", item.Name(), item.InventoryNumber())
	}

	return nil
}

func (a *App) addTransportWithReport(transport domain.Transport) error {
	err := a.fleet.AddTransport(transport)
	switch {
	case err == nil:
		fmt.Fprintf(a.out, "- %s [%s]: ПРИНЯТ\n", transport.Name(), transport.InventoryNumber())
		return nil
	case errors.Is(err, fleet.ErrInspectionFailed):
		fmt.Fprintf(a.out, "- %s [%s]: ОТКЛОНЁН техосмотром\n", transport.Name(), transport.InventoryNumber())
		return nil
	default:
		return err
	}
}
