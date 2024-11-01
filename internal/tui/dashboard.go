package tui

import (
	"fmt"
	"github.com/ToffaKrtek/stat-dashboard/internal/dashboard"

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
