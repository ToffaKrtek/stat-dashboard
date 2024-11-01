package dashboard

import (
	"encoding/csv"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func MakeDashBoardFromCSV(csvFileName string, pngFileName string, title string, legends []string) (string, error) {
	labels, data, err := ConvertCSVToLabelsAndData(csvFileName)
	if err != nil {
		return "", err
	}
	return makeDashboard(pngFileName, title, data, labels, legends)
}

func ConvertCSVToLabelsAndData(fileName string) ([]float64, [][]float64, error) {
	// Открываем CSV файл
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// Читаем содержимое CSV файла
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	// Инициализируем срезы для меток и данных
	var labels []float64
	var data [][]float64

	for i, record := range records {
		if i == 0 {
			// Инициализируем data с нужным количеством колонок
			numColumns := len(record) - 1 // Количество колонок данных (без меток)
			data = make([][]float64, numColumns)
			continue
		}

		// Форматируем время и добавляем в labels
		layout := time.RFC3339 // Используем стандартный формат RFC3339
		if len(record) > 0 {
			dateTime, err := time.Parse(layout, record[0]) // Замените формат на ваш
			if err != nil {
				return nil, nil, err
			}
			labels = append(labels, float64(dateTime.Unix())) // Оставляем только время
		}

		// Обрабатываем остальные колонки
		for j, value := range record[1:] {
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, nil, err
			}
			data[j] = append(data[j], floatValue) // Добавляем значение в соответствующую колонку
		}
	}

	return labels, data, nil
}

func makeDashboard(filename string, title string, data [][]float64, labels []float64, legendLabels []string) (string, error) {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = "Время"
	p.Y.Label.Text = "Значение"

	colors := []color.Color{
		color.RGBA{255, 0, 0, 255},   // Красный
		color.RGBA{0, 255, 0, 255},   // Зеленый
		color.RGBA{0, 0, 255, 255},   // Синий
		color.RGBA{255, 165, 0, 255}, // Оранжевый
		color.RGBA{128, 0, 128, 255}, // Пурпурный
		color.RGBA{0, 255, 255, 255}, // Циан
		color.RGBA{255, 255, 0, 255}, // Желтый
	}

	for i, metric := range data {
		line, err := plotter.NewLine(plotter.XYs{})
		if err != nil {
			return "", err
		}

		for j, value := range metric {
			//line.XYs = append(line.XYs, plotter.XY{X: float64(j), Y: value})
			line.XYs = append(line.XYs, plotter.XY{X: labels[j], Y: value})
		}

		line.Color = colors[i%len(colors)]
		p.Add(line)
		p.Legend.Add(legendLabels[i], line)

	}
	p.Legend.Top = true
	p.X.Tick.Marker = plot.TimeTicks{
		Format: "15:04:05",                                                  // Формат отображения времени
		Time:   func(t float64) time.Time { return time.Unix(int64(t), 0) }, // Преобразование float64 в time.Time
	}
	if err := p.Save(8*vg.Inch, 4*vg.Inch, filename); err != nil {
		return "", err
	}
	return filepath.Abs(filename)
}
