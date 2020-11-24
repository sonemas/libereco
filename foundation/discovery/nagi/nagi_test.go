package nagi

import (
	"log"
	"os"
	"testing"

	"github.com/sonemas/libereco/foundation/discovery"
	"github.com/sonemas/libereco/foundation/tests"
)

func TestNode(t *testing.T) {
	log := log.New(os.Stdout, "NODE : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	t.Log("Given the need to work with Nodes.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single Node.", testID)

		n, err := New(log, "0.0.0.0:50000/0")
		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to create a new node: %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to create a new node.", tests.Success, testID)

		peer1 := Peer{id: "1", addr: "0.0.0.0:50001"}
		peer2 := Peer{id: "2", addr: "0.0.0.0:50002"}

		// AddPeer
		n.AddPeer(&peer1)
		if len(n.peers) != 1 {
			t.Fatalf("\t%s\tTest %d:\tShould be able to add a new peer: %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to add a new peer.", tests.Success, testID)

		if n.joinedPeers[peer1.id].id != peer1.id {
			t.Fatalf("\t%s\tTest %d:\tShould be listed in new peer's table.", tests.Failed, testID)
		}
		t.Logf("\t%s\tTest %d:\tShould be listed in table of new peers.", tests.Success, testID)

		expected, got := 0, len(n.faultyPeers)
		if got != expected {
			t.Fatalf("\t%s\tTest %d:\tShould have %d inactive peers, but got %d.", tests.Failed, testID, expected, got)
		}
		t.Logf("\t%s\tTest %d:\tShould inactive peers to be %d.", tests.Success, testID, expected)

		// SetPeer
		if err := n.SetPeer(&peer2); err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to set a new peer: %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to set a new peer.", tests.Success, testID)

		err = n.SetPeer(&peer1)
		if err != discovery.ErrPeerExists {
			t.Fatalf("\t%s\tTest %d:\tShould get ErrPeerExists when trying to set an existing peer, but got %v.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould get ErrPeerExists when trying to set an existing peer.", tests.Success, testID)

		joinedPeers := n.JoinedPeers()
		if len(joinedPeers) != 1 && joinedPeers[0].id != peer1.id {
			t.Fatalf("\t%s\tTest %d:\tShould have only peer %s in joinedPeers: %v.", tests.Failed, testID, peer1.id, n.joinedPeers)
		}
		t.Logf("\t%s\tTest %d:\tShould have only peer %s in joinedPeers.", tests.Success, testID, peer1.id)

		// GetPeer
		p, err := n.GetPeer(peer1.id)
		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to get an existing peer: %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to get an existing peers.", tests.Success, testID)

		if p.id != peer1.id {
			t.Fatalf("\t%s\tTest %d:\tShould get peer id %s, but got: %s.", tests.Failed, testID, peer1.id, p.id)
		}
		t.Logf("\t%s\tTest %d:\tShould get peer id %s.", tests.Success, testID, peer1.id)

		if _, err = n.GetPeer("blahblah"); err != discovery.ErrPeerNotFound {
			t.Fatalf("\t%s\tTest %d:\tShould get ErrPeerNotFound when trying to get a non-existing peer, but got %v.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould get ErrPeerNotFound when trying to get a non-existing peer.", tests.Success, testID)

		// UpdatePeer
		updPeer := peer1
		updPeer.addr = "127.0.0.1:55551"

		if err := n.UpdatePeer(&updPeer, true); err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to update and mark a peer as inactive:  %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to update and mark a peer as inactive.", tests.Success, testID)

		if _, err := n.GetPeer(peer1.id); err != discovery.ErrPeerNotFound {
			t.Fatalf("\t%s\tTest %d:\tShould get ErrPeerNotFound when getting inactive client, but got:  %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould get ErrPeerNotFound when getting inactive client.", tests.Success, testID)

		if len(n.faultyPeers) > 0 && n.faultyPeers[peer1.id].id != peer1.id {
			t.Fatalf("\t%s\tTest %d:\tShould have updated peer in table of inactive peers: %+v.", tests.Failed, testID, n.faultyPeers)
		}
		t.Logf("\t%s\tTest %d:\tShould have updated peer in table of inactive peers.", tests.Success, testID)

		updPeer = peer2
		updPeer.addr = "127.0.0.1:55552"
		if err := n.UpdatePeer(&updPeer, false); err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to update without marking the peer as inactive:  %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to update without marking the peer as inactive.", tests.Success, testID)

		l := len(n.faultyPeers)
		if l != 1 {
			t.Fatalf("\t%s\tTest %d:\tShould have 1 peer marked as inactive, but got: %d.", tests.Failed, testID, l)
		}
		t.Logf("\t%s\tTest %d:\tShould have 1 peer marked as inactive.", tests.Success, testID)

		// InactivePeer
		if err := n.MarkPeerFaulty(peer2.id); err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to mark peer as inactive:  %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to mark peer as inactive.", tests.Success, testID)

		if l := len(n.faultyPeers); l != 2 {
			t.Fatalf("\t%s\tTest %d:\tShould have 2 peers marked as inactive, but got: %d.", tests.Failed, testID, l)
		}
		t.Logf("\t%s\tTest %d:\tShould have 2 peers marked as inactive.", tests.Success, testID)

		faultyPeers := n.FaultyPeers()
		if l := len(faultyPeers); l != 2 {
			t.Fatalf("\t%s\tTest %d:\tShould have 2 peers in faultyPeers, but got: %v.", tests.Failed, testID, l)
		}
		t.Logf("\t%s\tTest %d:\tShould have 2 peers in faultyPeers.", tests.Success, testID)

		// RemovePeer
		n.AddPeer(&peer2)

		if err := n.RemovePeer(peer2.id); err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to remove a peer: %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to remove a new peer.", tests.Success, testID)

		if _, err = n.GetPeer(peer2.id); err != discovery.ErrPeerNotFound {
			t.Fatalf("\t%s\tTest %d:\tShould get ErrPeerNotFound when trying to get a removed peer, but got %v.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould get ErrPeerNotFound when trying to get a removed peer.", tests.Success, testID)
	}
}
