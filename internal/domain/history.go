package domain

import (
	"time"

	"github.com/google/uuid"
)

// HistoryStatus — отдельный тип для статусов истории
type HistoryStatus string

const (
	StatusPending   HistoryStatus = "pending"    // В очереди на генерацию
	StatusInProcess HistoryStatus = "processing" // Прямо сейчас генерируется
	StatusCompleted HistoryStatus = "completed"  // Успешно завершено
	StatusFailed    HistoryStatus = "failed"     // Произошла ошибка
	StatusFrozen    HistoryStatus = "frozen"     // Заморожен
	StatusDelete    HistoryStatus = "delete"
)

type UserRequest struct {
	UserID      int    `json:"user_id"`     // Чья это история (Ownership)
	Description string `json:"description"` // Основное описание истории
	TokenSize   int    `json:"token_size"`  // Размер истории в токенах
	ImageSize   int    `json:"image_size"`  // Количество картинок в запросе
	Save        bool   `json:"save"`
}
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
	UUID        string      `json:"uuid"`       // Публичный ID для клиента
	BookTitle   string      `json:"book_title"` // Краткое название (для списка)
	UserRequest UserRequest `json:"user_request"`

	// Состояние (особенно важно для генерации)
	Status   HistoryStatus `json:"status"`          // "pending", "completed", "failed"
	ErrorMsg string        `json:"error,omitempty"` // Почему не создалось?

	// Временные метки (Timestamps)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"` // Меняется при изменении статуса

	Chapters []Chapter `json:"chapter"`
}

type HistoryRepository interface {
	// Create Создать запись
	Create(h *History) error
	// GetByID Получить историю по ID
	GetByID(id uuid.UUID) (*History, error)
	// GetByUserID получить все истории персонажа
	GetByUserID(id int) (*[]History, error)
	// Update Обновить данные (статус, текст или новую картинку)
	Update(h *History) error
	// Delete Удаляет все лишнее оставляя только UserRequest
	Delete(id uuid.UUID) error
	// Freeze Меняет статус на остановленно
	Freeze(id uuid.UUID, frozen bool) error
	// CountActiveTasks Посчитать, сколько задач сейчас "в работе" у юзера
	CountActiveTasks(userID string) (int, error)
}

type HistoryService interface {
	// Create Проверяет лимит и создает историю
	Create(req UserRequest) (hID uuid.UUID, err error)
	// Get Просто возвращает текущее состояние из БД
	Get(id uuid.UUID) (*History, error)
	// Freeze Поставить на паузу историю
	Freeze(id uuid.UUID, frozen bool) error
	// Delete Удаляет историю оставляя запрос
	Delete(id uuid.UUID) error
}
