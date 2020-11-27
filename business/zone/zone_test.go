package zone_test

import (
	"log"
	"os"
	"testing"

	"github.com/sonemas/libereco/business/blockchain"
	"github.com/sonemas/libereco/business/zone"
	"github.com/sonemas/libereco/foundation/tests"
)

func TestZone(t *testing.T) {
	logger := log.New(os.Stdout, "ZONE : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	t.Log("Given the need to work with Zones")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen working with a blockchain Zone.", testID)
		{
			bc, err := blockchain.New()
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create blockchain: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create blockchain.", tests.Success, testID)

			_, err = zone.New(logger, "test", bc)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create zone: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create zone.", tests.Success, testID)
		}
	}
}
