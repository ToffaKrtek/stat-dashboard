package tui

import (
	"fmt"
	"os"

	"github.com/ToffaKrtek/stat-dashboard/internal/dashboard"
	"golang.org/x/term"

	"github.com/guptarohit/asciigraph"
)

var colors = []asciigraph.AnsiColor{
	asciigraph.Red,
	asciigraph.Green,
	asciigraph.Blue,
	asciigraph.Orange,
	asciigraph.Purple,
	asciigraph.Cyan,
	asciigraph.Yellow,
}

func PrintDashBordFromCSV(filename string, title string, legends []string) string {

	// Запуск цикла обновления графика
	_, data, err := dashboard.ConvertCSVToLabelsAndData(filename)
	if err != nil {
		fmt.Println("Ошибка при обработки данных")
		return ""
	}
	if len(data) == 0 || len(data[0]) == 0 {
		return ""
	}
	width := GetTerminalWidth() - 15
	if len(data[0]) > width {
		for i := range data {
			data[i] = data[i][len(data[i])-width:]
		}
	}

	graph := asciigraph.PlotMany(
		data,
		asciigraph.Precision(3),
		asciigraph.SeriesColors(
			colors...,
		),
		asciigraph.SeriesLegends(legends...),
		asciigraph.Caption(title),
		asciigraph.LowerBound(0.0),
	)
	// Выводим график
	return graph
}

func GetTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Ошибка получения ширины терминала")
		return 50
	}
	return width
}
