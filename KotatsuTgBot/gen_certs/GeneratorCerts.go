package gen_certs

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/fogleman/gg"
)

func GenCerts(last_name string, first_name string, paternity string, input_image_path string, user_tg_id_str string) (output_image_path string) {
	// Загрузка вашей картинки
	file, err := os.Open(input_image_path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	dc := gg.NewContextForImage(img)

	dc.SetColor(color.Black) // Установка цвета текста (черный в данном случае)
	if err := dc.LoadFontFace("./fonts/JetBrainsMono-Regular.ttf", 72); err != nil {
		panic(err)
	}

	startY1 := 0.0
	startY2 := 0.0
	startY3 := 0.0

	// Определите начальные координаты Y
	if last_name == "" && first_name == "" && paternity == "" {
		fmt.Println("Откат")
	} else if last_name != "" && first_name == "" && paternity == "" {
		startY1 = 1640.0
	} else if last_name != "" && first_name != "" && paternity == "" {
		startY1 = 1600.0
		startY2 = 1700.0
	} else if last_name != "" && first_name == "" && paternity != "" {
		startY1 = 1600.0
		startY3 = 1700.0
	} else if last_name == "" && first_name != "" && paternity == "" {
		startY2 = 1640.0
	} else if last_name == "" && first_name != "" && paternity != "" {
		startY2 = 1600.0
		startY3 = 1700.0
	} else if last_name == "" && first_name == "" && paternity != "" {
		startY3 = 1640.0
	} else {
		startY1 = 1540.0
		startY2 = 1640.0
		startY3 = 1740.0
	}

	// Вычислите ширину изображения
	imgWidth := float64(img.Bounds().Max.X)

	// Функция для рассчета координаты X для выравнивания по центру
	centerX := func(text string) float64 {
		spec, _ := dc.MeasureString(text)
		return (imgWidth - spec) / 2
	}

	textX1 := centerX(last_name)
	textX2 := centerX(first_name)
	textX3 := centerX(paternity)

	dc.DrawStringAnchored(last_name, textX1, startY1, 0, 0.5)
	dc.DrawStringAnchored(first_name, textX2, startY2, 0, 0.5)
	dc.DrawStringAnchored(paternity, textX3, startY3, 0, 0.5)

	dc.Stroke()

	pic_path := "./img/temp/" + user_tg_id_str + "_output.png"

	// Сохранение результата в файл
	if err := dc.SavePNG("./img/temp/" + user_tg_id_str + "_output.png"); err != nil {
		fmt.Println("Ошибка при сохранении изображения:", err)
	}

	return pic_path

}
