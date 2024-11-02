package main

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"
	"unsafe"

	goruntime "runtime"

	"github.com/beevik/ntp"
	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App представляет собой структуру вашего приложения.
type App struct {
	ctx context.Context
}

// NewApp создает новое приложение.
func NewApp() *App {
	return &App{}
}

// Startup вызывается при запуске приложения.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

var (
	selectedServer    string
	predefinedServers = []string{
		"pool.ntp.org",
		"time.google.com",
		"time.windows.com",
		"time.apple.com",
		"europe.pool.ntp.org",
		"north-america.pool.ntp.org",
		"asia.pool.ntp.org",
		"192.168.0.100",
	}
	autoSyncTicker *time.Ticker
	autoSyncQuit   chan struct{}
)

// GetPredefinedServers возвращает список предустановленных NTP-серверов.
func (a *App) GetPredefinedServers() []string {
	return predefinedServers
}

// SetSelectedServer устанавливает выбранный NTP-сервер.
func (a *App) SetSelectedServer(server string) {
	selectedServer = server
}

// SyncTime синхронизирует системное время с указанным сервером.
func (a *App) SyncTime(server string) (string, error) {
	ntpTime, err := getNTPTime(server)
	if err != nil {
		return "", fmt.Errorf("Ошибка получения времени с сервера %s: %v", server, err)
	}

	err = setSystemTime(ntpTime)
	if err != nil {
		return "", fmt.Errorf("Ошибка установки системного времени: %v", err)
	}

	formattedTime := ntpTime.Format("15:04:05 MST 2006-01-02")
	return formattedTime, nil
}

// StartAutoSync запускает автоматическую синхронизацию времени.
func (a *App) StartAutoSync(interval int) error {
	if autoSyncTicker != nil {
		autoSyncTicker.Stop()
		close(autoSyncQuit)
	}

	autoSyncTicker = time.NewTicker(time.Duration(interval) * time.Second)
	autoSyncQuit = make(chan struct{})

	go func() {
		for {
			select {
			case <-autoSyncQuit:
				return
			case <-autoSyncTicker.C:
				if selectedServer == "" {
					wruntime.EventsEmit(a.ctx, "log", "NTP-сервер не выбран для синхронизации")
					continue
				}
				ntpTime, err := a.SyncTime(selectedServer)
				if err != nil {
					wruntime.EventsEmit(a.ctx, "log", err.Error())
					continue
				}
				wruntime.EventsEmit(a.ctx, "log", fmt.Sprintf("Системное время успешно синхронизировано: %s", ntpTime))
			}
		}
	}()

	return nil
}

// getNTPTime получает текущее время с указанного NTP-сервера.
func getNTPTime(server string) (time.Time, error) {
	return ntp.Time(server)
}

// setSystemTime устанавливает системное время.
func setSystemTime(t time.Time) error {
	if goruntime.GOOS == "windows" {
		t = t.UTC()
		st := syscall.Systemtime{
			Year:         uint16(t.Year()),
			Month:        uint16(t.Month()),
			DayOfWeek:    uint16(t.Weekday()),
			Day:          uint16(t.Day()),
			Hour:         uint16(t.Hour()),
			Minute:       uint16(t.Minute()),
			Second:       uint16(t.Second()),
			Milliseconds: uint16(t.Nanosecond() / 1e6),
		}

		ret, _, err := syscall.NewLazyDLL("kernel32.dll").NewProc("SetSystemTime").Call(uintptr(unsafe.Pointer(&st)))
		if ret == 0 {
			return fmt.Errorf("не удалось установить системное время: %v", err)
		}
		return nil
	} else {
		cmd := exec.Command("sudo", "date", "-s", t.Format("2006-01-02 15:04:05"))
		return cmd.Run()
	}
}
