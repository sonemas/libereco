package kv

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/sonemas/libereco/business/hash"
)

// KVOption is an option to modify the settings
// of a KVStorage.
type KVOption func(*KVStorage) error

// WithFetching is an option to enable fetching.
func WithFetching() KVOption {
	return func(s *KVStorage) error {
		s.allowFetching = true
		return nil
	}
}

// WithDeleting is an option to enable deleting.
func WithDeleting() KVOption {
	return func(s *KVStorage) error {
		s.allowDeleting = true
		return nil
	}
}

// KVStorage provides basic key/value storage and implements
// the BasicStorage, Fetcher and Deleter interface.
type KVStorage struct {
	mu            sync.RWMutex
	filename      string
	logger        *log.Logger
	storage       map[string][]byte
	version       uint
	allowFetching bool
	allowDeleting bool
}

// New is a factory function to create a new initialized KVStorage
func New(logger *log.Logger, filename string, opts ...KVOption) (*KVStorage, error) {
	s := KVStorage{
		filename: filename,
		logger:   logger,
	}

	err := s.load()
	if err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "loading state")
	}

	if os.IsNotExist(err) {
		for _, opt := range opts {
			if err := opt(&s); err != nil {
				return nil, errors.Wrapf(err, "executing option %T", opt)
			}
		}
	}

	return &s, nil
}

func (s *KVStorage) Hash() string {
	return hash.Hash(fmt.Sprintf("%v%v%v%d", s.storage, s.allowFetching, s.allowDeleting, s.version))
}

func (s *KVStorage) Version() uint {
	return s.version
}

func (s *KVStorage) String() string {
	return fmt.Sprintf("%s-%d", s.Hash(), s.version)
}

// state is the state of a KVStorage as it will be stored/retrieved from disk.
type state struct {
	Version       uint
	Storage       map[uint]map[string][]byte
	Hashes        map[uint]string
	AllowFetching bool
	AllowDeleting bool
}

func load(filename string) (*state, error) {
	snapshot := state{}

	f, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if os.IsNotExist(err) {
		return &snapshot, err
	}
	if err != nil {
		return nil, errors.Wrap(err, "opening file")
	}
	defer f.Close()

	if err := gob.NewDecoder(f).Decode(&snapshot); err != nil {
		return nil, errors.Wrap(err, "decoding data")
	}

	return &snapshot, nil
}

var (
	ErrStorageHashesLengthMismatch = errors.New("storage and hashes length mismatch")

	// ErrVersionExists is an error indicating that the operation
	// tries to save a version that already exists.
	ErrVersionExists = errors.New("version exists")

	ErrHashMismatch = errors.New("hash mismatch")
)

func (s *KVStorage) save() error {
	snapshot, err := load(s.filename)
	if !os.IsNotExist(err) {
		return err
	}

	if snapshot.Hashes == nil {
		snapshot.Hashes = make(map[uint]string)
	}

	if snapshot.Storage == nil {
		snapshot.Storage = make(map[uint]map[string][]byte)
	}

	if _, exists := snapshot.Storage[s.version]; exists {
		return ErrVersionExists
	}

	snapshot.Storage[s.version] = s.storage
	snapshot.Hashes[s.version] = s.Hash()

	if len(snapshot.Storage) != len(snapshot.Hashes) {
		return ErrStorageHashesLengthMismatch
	}

	snapshot.Version = s.version
	snapshot.AllowDeleting = s.allowDeleting
	snapshot.AllowFetching = s.allowFetching

	f, err := os.OpenFile(s.filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrap(err, "opening file")
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(snapshot); err != nil {
		return errors.Wrap(err, "encoding data")
	}

	return nil
}

func (s *KVStorage) load() error {
	snapshot, err := load(s.filename)
	if err != nil {
		return err
	}

	tmp := KVStorage{
		storage:       snapshot.Storage[snapshot.Version],
		allowDeleting: snapshot.AllowDeleting,
		allowFetching: snapshot.AllowFetching,
		version:       snapshot.Version,
	}

	if tmp.Hash() != snapshot.Hashes[snapshot.Version] {
		return ErrHashMismatch
	}

	s.storage = snapshot.Storage[snapshot.Version]
	s.allowDeleting = snapshot.AllowDeleting
	s.allowFetching = snapshot.AllowFetching
	s.version = snapshot.Version

	return nil
}

func init() {
	gob.Register(state{})
}
