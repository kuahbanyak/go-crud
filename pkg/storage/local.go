package storage

import (
    "io"
    "os"
    "path/filepath"
)

func SaveLocal(r io.Reader, filename string) (string, error) {
    base := os.Getenv("STORAGE_PATH")
    if base == "" {
        base = "./uploads"
    }
    if err := os.MkdirAll(base, 0755); err != nil {
        return "", err
    }
    path := filepath.Join(base, filename)
    f, err := os.Create(path)
    if err != nil { return "", err }
    defer f.Close()
    if _, err := io.Copy(f, r); err != nil { return "", err }
    return path, nil
}
