package stat

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

var fullTime *time.Duration
var statFileName *string
var plotFileName *string
var title *string
var legends *[]string

func Run(
	mchan chan bool,
) {
	go func() {
		for {
			now := time.Now()

			nextRun := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location())

			if now.After(nextRun) {
				nextRun = nextRun.Add(*fullTime)
			}
			time.Sleep(time.Until(nextRun))

			fmt.Println("Запуск генерации графика и сброса данных...")
			//TODO:: send by channel
			//cus didn't wait before clear
			//MakeDashBoardFromCSV(*statFileName, *plotFileName, *title, *legends)
			mchan <- true
			ClearCSV(*statFileName)
		}
	}()
}

type Metrics struct {
	Load1         float64
	Load5         float64
	Load15        float64
	MemoryPercent float64
	DiskPercent   float64
}

func getMetrics() *Metrics {
	var metrics Metrics
	countCpu := 1
	logicalCores, err := cpu.Counts(true)
	if err != nil {
		log.Printf("Ошибка получения CpuCounts: %v", err)
	} else {
		countCpu = logicalCores
	}
	loadAvg, err := load.Avg()
	if err != nil {
		log.Printf("Ошибка получения LoadAverage: %v", err)
	} else {
		metrics.Load1 = normalization(loadAvg.Load1, float64(countCpu))
		metrics.Load5 = normalization(loadAvg.Load5, float64(countCpu))
		metrics.Load15 = normalization(loadAvg.Load15, float64(countCpu))
	}
	memStat, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Ошибка получения VirtualMemory: %v", err)
	} else {
		metrics.MemoryPercent = normalization(memStat.UsedPercent, 100)
	}
	diskStat, err := disk.Usage("/")
	if err != nil {
		log.Printf("Ошибка получения DiskUsage: %v", err)
	} else {
		metrics.DiskPercent = normalization(diskStat.UsedPercent, 100)
	}
	return &metrics
}

func roundToOneDecimal(value float64) float64 {
	return float64(int(value*100+0.5)) / 100
}

func normalization(value float64, max float64) float64 {
	return roundToOneDecimal(value / max)
}

func WriteMetricsToCSV(filename string) error {
	metrics := getMetrics()
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	timestamp := time.Now().Format(time.RFC3339)
	record := []string{
		timestamp,
		fmt.Sprintf("%.2f", metrics.Load1),
		fmt.Sprintf("%.2f", metrics.Load5),
		fmt.Sprintf("%.2f", metrics.Load15),
		fmt.Sprintf("%.2f", metrics.MemoryPercent),
		fmt.Sprintf("%.2f", metrics.DiskPercent),
	}
	writer.Write(record)
	writer.Flush()

	return nil
}

func ClearCSV(filename string) {
	file, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Ошибка очистки CSV файла:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	//writer.Write([]string{"Timestamp", "Load1", "Load5", "Load15", "MemoryPercent", "DiskPercent"})
	writer.Flush()
}
func ExistCSV(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Файл не существует, создаем его
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("не удалось создать файл: %v", err)
		}
		defer file.Close() // Закрываем файл после создания
		fmt.Println("Файл создан:", filename)
	} else if err != nil {
		fmt.Printf("ошибка при проверке файла: %v", err)
	} else {
		fmt.Println("Файл уже существует:", filename)
	}
}
