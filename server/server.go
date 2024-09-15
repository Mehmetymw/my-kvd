package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

const (
	MAX_PASSWORD_LENGTH = 16
)

type Store struct {
	data map[string]string
	mu   sync.RWMutex
	file string
}

func NewStore(filename string) *Store {
	store := &Store{
		data: make(map[string]string),
		file: filename,
	}
	store.load()
	return store
}

func (store *Store) Get(key string) (string, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	value, exists := store.data[key]
	return value, exists
}

func (store *Store) Set(key, value string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.data[key] = value
	store.save()
}

func (store *Store) SetPassword(key string, length int) error {
	if length < MAX_PASSWORD_LENGTH {
		return fmt.Errorf("password length should be at least %v characters", MAX_PASSWORD_LENGTH)
	}

	value, err := GenerateStrongPassword(length)
	if err != nil {
		return fmt.Errorf("password cannot be generated: %w", err)
	}
	store.Set(key, value)
	return nil
}

func generateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func GenerateStrongPassword(length int) (string, error) {
	if length < MAX_PASSWORD_LENGTH {
		return "", fmt.Errorf("password length should be at least %v characters", MAX_PASSWORD_LENGTH)
	}

	bytes, err := generateRandomBytes(length)
	if err != nil {
		return "", err
	}
	password := hex.EncodeToString(bytes)

	if len(password) > length {
		return password[:length], nil
	}
	return password, nil
}

func (store *Store) save() {
	store.mu.RLock()
	defer store.mu.RUnlock()
	data, err := json.Marshal(store.data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}
	err = os.WriteFile(store.file, data, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func (store *Store) load() {
	file, err := os.Open(store.file)
	if err != nil {
		if os.IsNotExist(err) {
			store.data = make(map[string]string)
			return
		}
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	if err := json.Unmarshal(data, &store.data); err != nil {
		fmt.Println("Error unmarshalling data:", err)
	}
}
