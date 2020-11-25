package kv

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/sonemas/libereco/foundation/storage"
	"github.com/sonemas/libereco/foundation/tests"
)

func TestKVStorage(t *testing.T) {

	logger := log.New(os.Stdout, "KV : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	t.Log("Given the need to work with KVStorage")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a default KVStorage.", testID)
		{
			file := fmt.Sprintf("test%d.lio", testID)
			defer os.Remove(file)
			s, err := New(logger, file)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create KVStorage: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create KVStorage.", tests.Success, testID)

			if err := s.Init(map[string][]byte{}, nil); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to init KVStorage: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to init KVStorage.", tests.Success, testID)

			if err := s.Set("0", []byte{100}, nil); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to set key/value: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to set key/value.", tests.Success, testID)

			got, expected := s.String(), "0156b02d1162ce75726795015bb0b96b1fef2e850f92a4f1d4ee2ade25aa9602-1"
			if got != expected {
				t.Fatalf("\t%s\tTest %d:\tShould have hash-version pair %q, but got %q.", tests.Failed, testID, expected, got)
			}
			t.Logf("\t%s\tTest %d:\tShould have hash-version pair %q.", tests.Success, testID, expected)

			if err := s.Set("0", []byte{200}, nil); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to set same key/value again: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to set key/value.", tests.Success, testID)

			got, expected = s.String(), "0cfa24a8460990a176f485e6b6af19cda0574e7cc9337c6ca1612339c74f61ae-2"
			if got != expected {
				t.Fatalf("\t%s\tTest %d:\tShould have hash-version pair %q, but got %q.", tests.Failed, testID, expected, got)
			}
			t.Logf("\t%s\tTest %d:\tShould have hash-version pair %q.", tests.Success, testID, expected)

			b, err := s.Get("0", nil)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to get key/value: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to get key/value.", tests.Success, testID)

			{
				got, expected := int(b[0]), 200
				if got != expected {
					t.Fatalf("\t%s\tTest %d:\tShould have value %d, but got %d.", tests.Failed, testID, expected, got)
				}
				t.Logf("\t%s\tTest %d:\tShould have value %d.", tests.Success, testID, expected)
			}

			if err := s.Delete("0", nil); err != storage.ErrOperationNotAllowed {
				t.Fatalf("\t%s\tTest %d:\tShould get ErrOperationNotAllowed when trying to delete, but got: %v.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould get ErrOperationNotAllowed when trying to delete.", tests.Success, testID)

			if _, err := s.Fetch(nil); err != storage.ErrOperationNotAllowed {
				t.Fatalf("\t%s\tTest %d:\tShould get ErrOperationNotAllowed when trying to fetch, but got: %v.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould get ErrOperationNotAllowed when trying to fetch.", tests.Success, testID)
		}
	}

}
