package domain

import "time"

type PartChapter struct {
	Number int `json:"number"` // Номер для последовательности, начиная с 1, чисто для последовательности
	// Подзаголовок если надо, если нет то пустой и не видно разделения
	Subtitle string `json:"subtitle"`
	// Или текст или изображение или оба (тогда сначало идет текст после изображение)
	Text     string `json:"text"`      // Абзац(ы) части
	ImageUrl string `json:"image_url"` // Изображение части
}
type Chapter struct {
	Title string        `json:"title"` // Заголовок главы
	Parts []PartChapter `json:"part"`  // Части главы
}
type History struct {
	UUID   string `json:"uuid"`    // Публичный ID для клиента
	UserID string `json:"user_id"` // Чья это история (Ownership)
	Title  string `json:"title"`   // Краткое название (для списка)

	// Состояние (особенно важно для генерации)
	Status   string `json:"status"`          // "pending", "completed", "failed"
	ErrorMsg string `json:"error,omitempty"` // Почему не создалось?

	// Временные метки (Timestamps)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"` // Меняется при изменении статуса

	Chapters []Chapter `json:"chapter"`
}

type HistoryRepository interface {
}

type HistoryService interface {
}
