package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Blob interface {
	Root() string
	Save(kind string, id string, payload any) error
	Load(kind string, id string, out any) error
	Exists(kind string, id string) bool
	List(kind string) ([]string, error)
}

type FileStore struct {
	root string
}

func New(root string) *FileStore {
	return &FileStore{root: root}
}

func (s *FileStore) Root() string {
	return s.root
}

func (s *FileStore) Dir(kind string) string {
	return filepath.Join(s.root, kind)
}

func (s *FileStore) Path(kind string, id string) string {
	return filepath.Join(s.root, kind, id+".json")
}

func (s *FileStore) Save(kind string, id string, payload any) error {
	dir := s.Dir(kind)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, id+".*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return err
	}
	if err := os.Rename(name, s.Path(kind, id)); err != nil {
		os.Remove(name)
		return err
	}
	return nil
}

func (s *FileStore) Load(kind string, id string, out any) error {
	data, err := os.ReadFile(s.Path(kind, id))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

func (s *FileStore) Exists(kind string, id string) bool {
	_, err := os.Stat(s.Path(kind, id))
	return err == nil
}

func (s *FileStore) List(kind string) ([]string, error) {
	entries, err := os.ReadDir(s.Dir(kind))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	ids := []string{}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".json") {
			continue
		}
		ids = append(ids, strings.TrimSuffix(name, ".json"))
	}
	sort.Strings(ids)
	return ids, nil
}
