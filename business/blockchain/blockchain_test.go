package blockchain_test

import (
	"testing"
	"time"

	"github.com/sonemas/libereco/business/blockchain"
)

func TestBlockchain(t *testing.T) {
	bc, err := blockchain.New(blockchain.WithDifficulty(2), blockchain.WithNonceTimeout(30*time.Second))
	if err != nil {
		t.Fatal(err)
	}

	if err := bc.Init([]byte("Geri vyrai geroj girioj gerą girą gerai gėrė. Geriems vyrams geroj girioj gerą girą gera gerti. Geriems vyrams geroj girioj gerą girą gera gert."), nil); err != nil {
		t.Fatal(err)
	}

	l := bc.LastBlock()
	t.Logf("[CREATED GENESIS BLOCK]\nVersion:\t%d\nDifficulty:\t%d\nPrevious hash:\t%s\nHash:\t\t%s\nData:\t\t%s\nTimestamp:\t%d\nNonce:\t\t%d", l.Version, l.Difficulty, l.PrevHash, l.Hash, l.Data, l.Timestamp, l.Nonce)

	if err := bc.NewBlock([]byte("Labas rytas!")); err != nil {
		t.Fatal(err)
	}
	l = bc.LastBlock()
	t.Logf("[CREATED NEW BLOCK]\nVersion:\t%d\nDifficulty:\t%d\nPrevious hash:\t%s\nHash:\t\t%s\nData:\t\t%s\nTimestamp:\t%d\nNonce:\t\t%d", l.Version, l.Difficulty, l.PrevHash, l.Hash, l.Data, l.Timestamp, l.Nonce)
}
