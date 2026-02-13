package standards

import (
	"context"
	"os"
	"sync"
	"time"
)

type FileLoader struct {
	path     string
	cacheTTL time.Duration

	mu        sync.RWMutex
	content   string
	loadedAt  time.Time
}

func NewFileLoader(path string, cacheTTL time.Duration) *FileLoader {
	return &FileLoader{
		path:     path,
		cacheTTL: cacheTTL,
	}
}

func (f *FileLoader) Load(ctx context.Context) (string, error) {
	f.mu.RLock()
	if f.content != "" && time.Since(f.loadedAt) < f.cacheTTL {
		defer f.mu.RUnlock()
		return f.content, nil
	}
	f.mu.RUnlock()

	return f.Reload(ctx)
}

func (f *FileLoader) Reload(ctx context.Context) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	data, err := os.ReadFile(f.path)
	if err != nil {
		return "", err
	}

	f.content = string(data)
	f.loadedAt = time.Now()
	return f.content, nil
}
