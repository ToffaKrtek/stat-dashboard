package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/ToffaKrtek/stat-dashboard/internal/stat"
	"github.com/ToffaKrtek/stat-dashboard/internal/tui"
)

const statFileName = "metrics.csv"
const plotFileName = "metrics.png"

func main() {
	legends := []string{
		"Нагрузка на процессоры (1м)",
		"Нагрузка на процессоры (5м)",
		"Нагрузка на процессоры (15м)",
		"Оперативная память",
		"Диск",
	}
	stat.ExistCSV(statFileName)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for {
			select {
			case <-ticker.C:
				graph := tui.PrintDashBordFromCSV(statFileName, "", legends)
				clearScreen()
				fmt.Println(graph)
			case <-sigs:
				wg.Done()
				return
			}
		}
	}()

	wg.Wait()
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
