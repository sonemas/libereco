package blockchain

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sonemas/libereco/business/hash"

	"github.com/pkg/errors"
)

// Block is a record in the chain.
type Block struct {
	Version    uint
	Difficulty uint
	PrevHash   string
	Hash       string
	Data       []byte
	Timestamp  int64
	Nonce      uint
}

func (b Block) CalcHash() string {
	return hash.Hash(fmt.Sprintf("%d%d%s%v%d", b.Version, b.Difficulty, b.PrevHash, b.Data, b.Timestamp))
}

func (b Block) CalcNonce(ctx context.Context, difficulty uint) (uint, error) {
	expected := strings.Repeat("0", int(b.Difficulty))

	nonce := uint(0)
	check := ""
	for {
		select {
		case <-ctx.Done():
			return 0, errors.Wrap(ctx.Err(), "calculating nonce")
		default:
		}

		check = hash.Hash(fmt.Sprintf("%s%d", b.Hash, nonce))
		if len(check) > 0 && check[0:difficulty] == expected {
			break
		}
		nonce++
	}

	return nonce, nil
}

type Blockchain struct {
	mu           sync.RWMutex
	blocks       []Block
	difficulty   uint
	NonceTimeOut time.Duration
}

type BlockchainOption func(*Blockchain) error

func WithNonceTimeout(v time.Duration) BlockchainOption {
	return func(bc *Blockchain) error {
		bc.NonceTimeOut = v
		return nil
	}
}

func WithDifficulty(v uint) BlockchainOption {
	return func(bc *Blockchain) error {
		bc.difficulty = v
		return nil
	}
}

func New(opts ...BlockchainOption) (*Blockchain, error) {
	bc := Blockchain{
		difficulty:   3,
		NonceTimeOut: 5 * time.Minute,
	}

	for _, opt := range opts {
		if err := opt(&bc); err != nil {
			return nil, errors.Wrapf(err, "failed to run option %T", opt)
		}
	}

	return &bc, nil
}

func (bc *Blockchain) Height() int {
	return len(bc.blocks)
}
func (bc *Blockchain) LastBlock() Block {
	height := bc.Height()
	if height == 0 {
		return Block{Hash: "0000000000000000000000000000000000000000000000000000000000000000"}
	}

	return bc.blocks[bc.Height()-1]
}

func (bc *Blockchain) Validate(b Block) error {
	lastBlock := bc.LastBlock()

	if lastBlock.Hash != b.PrevHash {
		return fmt.Errorf("previous hash should be %s, but got %s (%d)", lastBlock.Hash, b.PrevHash, bc.Height())
	}

	if (b.PrevHash == "0000000000000000000000000000000000000000000000000000000000000000" && b.Version != 0) || (b.PrevHash != "0000000000000000000000000000000000000000000000000000000000000000" && bc.LastBlock().Version+1 != b.Version) {
		return fmt.Errorf("version, should be %d, but got %d", bc.LastBlock().Version+1, b.Version)
	}

	if b.CalcHash() != b.Hash {
		return fmt.Errorf("invalid hash value")
	}

	check := hash.Hash(fmt.Sprintf("%s%d", b.Hash, b.Nonce))
	expected := strings.Repeat("0", int(b.Difficulty))

	if check[0:b.Difficulty] != expected {
		return fmt.Errorf("invalid nonce value")
	}

	return nil
}

func (bc *Blockchain) NewBlock(data []byte) error {
	if bc.Height() == 0 {
		return fmt.Errorf("blockchain has not been initialized")
	}

	bc.mu.Lock()
	defer bc.mu.Unlock()

	if data == nil {
		data = []byte{}
	}

	lastBlock := bc.LastBlock()

	b := Block{
		Version:    lastBlock.Version + 1,
		Difficulty: bc.difficulty,
		PrevHash:   lastBlock.Hash,
		Data:       data,
		Timestamp:  time.Now().UTC().Unix(),
	}
	b.Hash = b.CalcHash()

	ctx, cancel := context.WithTimeout(context.Background(), bc.NonceTimeOut)
	defer cancel()
	n, err := b.CalcNonce(ctx, bc.difficulty)
	if err != nil {
		return errors.Wrap(err, "calculating nonce")
	}
	b.Nonce = n

	if err := bc.Validate(b); err != nil {
		return err
	}

	bc.blocks = append(bc.blocks, b)

	return nil
}

func (bc *Blockchain) Init(v interface{}, args map[string]interface{}) error {
	data, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("v should be []byte")
	}

	bc.mu.Lock()
	defer bc.mu.Unlock()

	if len(bc.blocks) != 0 {
		return nil
	}

	b := Block{
		Version:    0,
		Difficulty: bc.difficulty,
		PrevHash:   "0000000000000000000000000000000000000000000000000000000000000000",
		Data:       data,
		Timestamp:  time.Now().UTC().Unix(),
	}
	b.Hash = b.CalcHash()

	ctx, cancel := context.WithTimeout(context.Background(), bc.NonceTimeOut)
	defer cancel()
	n, err := b.CalcNonce(ctx, bc.difficulty)
	if err != nil {
		return errors.Wrap(err, "calculating nonce")
	}
	b.Nonce = n

	if err := bc.Validate(b); err != nil {
		return err
	}

	bc.blocks = append(bc.blocks, b)

	return nil
}
