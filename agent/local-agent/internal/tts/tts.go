package tts

import (
	"fmt"
	"os/exec"
	"sync"

	"go.uber.org/zap"
)

// TTSClient — клиент для голосовых оповещений.
//
//nolint:revive // exported: имя TTSClient намеренно (Client слишком общее в пакете tts).
type TTSClient struct {
	logger  *zap.Logger
	engine  string
	voice   string
	queue   []Announcement
	volume  int
	mu      sync.Mutex
	enabled bool
}

// Announcement — голосовое оповещение в очереди.
type Announcement struct {
	Text     string `json:"text"`
	Language string `json:"language"`
	Priority string `json:"priority"` // high, normal, low
}

// NewTTSClient создаёт клиент TTS.
func NewTTSClient(engine, voice string, volume int, enabled bool, logger *zap.Logger) *TTSClient {
	client := &TTSClient{
		engine:  engine,
		voice:   voice,
		volume:  volume,
		enabled: enabled,
		logger:  logger,
		queue:   make([]Announcement, 0),
	}

	// Запустить worker для обработки очереди
	go client.processQueue()

	return client
}

// Announce добавляет оповещение в очередь.
func (t *TTSClient) Announce(text, language, priority string) error {
	if !t.enabled {
		t.logger.Info("TTS disabled, skipping announcement")
		return nil
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	announcement := Announcement{
		Text:     text,
		Language: language,
		Priority: priority,
	}

	// Высокий приоритет — в начало очереди
	if priority == "high" {
		t.queue = append([]Announcement{announcement}, t.queue...)
	} else {
		t.queue = append(t.queue, announcement)
	}

	t.logger.Info("Announcement queued", zap.String("text", text), zap.String("priority", priority))
	return nil
}

func (t *TTSClient) processQueue() {
	for {
		t.mu.Lock()
		if len(t.queue) == 0 {
			t.mu.Unlock()
			continue
		}

		announcement := t.queue[0]
		t.queue = t.queue[1:]
		t.mu.Unlock()

		// Воспроизвести
		if err := t.speak(announcement.Text); err != nil {
			t.logger.Error("Failed to speak", zap.Error(err))
		}
	}
}

func (t *TTSClient) speak(text string) error {
	var cmd *exec.Cmd

	switch t.engine {
	case "rhvoice":
		// RHVoice для русского языка.
		//nolint:gosec // G204: text/voice из конфига и очереди, не из сетевого ввода.
		cmd = exec.Command("sh", "-c", "echo "+text+" | RHVoice-test -p "+t.voice)
	case "espeak":
		// eSpeak для английского
		cmd = exec.Command("espeak", "-v", "en", text)
	default:
		// Fallback: системная команда say (macOS)
		cmd = exec.Command("say", text)
	}

	t.logger.Info("Speaking text", zap.String("engine", t.engine), zap.String("text", text))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("TTS command failed: %w, output: %s", err, string(output))
	}

	return nil
}

// GetStatus возвращает статус TTS и очереди.
func (t *TTSClient) GetStatus() map[string]interface{} {
	t.mu.Lock()
	defer t.mu.Unlock()

	return map[string]interface{}{
		"enabled":      t.enabled,
		"engine":       t.engine,
		"voice":        t.voice,
		"queue_length": len(t.queue),
	}
}
