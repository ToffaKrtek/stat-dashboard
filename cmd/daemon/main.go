package main

import (
	"flag"
	"fmt"
	"github.com/ToffaKrtek/stat-dashboard/internal/dashboard"
	"github.com/ToffaKrtek/stat-dashboard/internal/stat"
	"github.com/ToffaKrtek/stat-dashboard/internal/tui"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

const statFileName = "metrics.csv"
const plotFileName = "metrics.png"

var title *string
var period *string
var repeat *string
var cliRun *bool

var interval time.Duration
var fulltime time.Duration

func getTitle() string {
	if title == nil {
		setFlags()
	}
	return *title + "-" + time.Now().Format("2006-01-02")
}

func setFlags() {
	title = flag.String("server", "server", "Название сервера")
	period = flag.String("period", "15", "Периодичность запуска сборщика метрик (в минутах)")
	repeat = flag.String("repeat", "day", "Периодичность полного цикла (day, week, hour, infinity, debug -- 1 минута и сбор метрик каждые 2 секунды")
	cliRun = flag.Bool("cli", false, "Запуск в командной строке")
	flag.Parse()
}

func getInterval() time.Duration {
	if interval > 0 {
		return interval
	}
	if period == nil {
		setFlags()
	}
	intVal, err := strconv.Atoi(*period)
	if err != nil || intVal < 1 {
		intVal = 15
	}
	interval = time.Duration(int64(intVal)) * time.Minute
	if *repeat == "debug" {
		interval = 2 * time.Second
	}
	if *repeat == "infinity" {
		interval = 2 * time.Second
	}
	return interval
}

func getFullTime() time.Duration {
	if fulltime > 0 {
		return fulltime
	}
	if repeat == nil {
		setFlags()
	}
	switch *repeat {
	case "day":
		fulltime = 24 * time.Hour
	case "week":
		fulltime = 7 * 24 * time.Hour
	case "hour":
		fulltime = 1 * time.Hour
	case "infinity":
		fulltime = 54 * 7 * 24 * time.Hour
	case "debug":
		fulltime = 1 * time.Minute
	default:
		fulltime = 24 * time.Hour
	}
	return fulltime
}

func main() {
	legends := []string{
		"Нагрузка на процессоры (1м)",
		"Нагрузка на процессоры (5м)",
		"Нагрузка на процессоры (15м)",
		"Оперативная память",
		"Диск",
	}
	stat.ExistCSV(statFileName)
	filen, err := dashboard.MakeDashBoardFromCSV(statFileName, plotFileName, getTitle(), legends)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(filen)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			now := time.Now()

			nextRun := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location())

			if now.After(nextRun) {
				nextRun = nextRun.Add(getFullTime())
			}
			time.Sleep(time.Until(nextRun))

			fmt.Println("Запуск генерации графика и сброса данных...")
			dashboard.MakeDashBoardFromCSV(statFileName, plotFileName, getTitle(), legends)
			stat.ClearCSV(statFileName)
		}
	}()

	if *cliRun {
		ticker := time.NewTicker(getInterval())
		defer ticker.Stop()
		go func() {
			for {
				select {
				case <-ticker.C:
					graph := tui.PrintDashBordFromCSV(statFileName, getTitle(), legends)
					clearScreen()
					fmt.Println(graph)
				case <-sigs:
					return
				}
			}
		}()
	}

	go func() {
		for {
			stat.WriteMetricsToCSV(statFileName)
			time.Sleep(getInterval())
		}
	}()
	<-sigs
	fmt.Println("Получен сигнал завершения. Выполняем финальные действия...")

	// Выполняем финальные действия перед завершением
	dashboard.MakeDashBoardFromCSV(statFileName, plotFileName, getTitle(), legends)
	fmt.Println("Дашборд создан перед завершением.")
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
