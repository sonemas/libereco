package node

import (
	"log"
	"os"
	"testing"

	"github.com/sonemas/libereco/business/tests"
)

func TestNode(t *testing.T) {
	log := log.New(os.Stdout, "API : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	t.Log("Given the need to work with Nodes.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single Node.", testID)

		n, err := New(log, "0.0.0.0:50000")
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

		if n.newPeers[peer1.id].id != peer1.id {
			t.Fatalf("\t%s\tTest %d:\tShould be listed in new peer's table.", tests.Failed, testID)
		}
		t.Logf("\t%s\tTest %d:\tShould be listed in table of new peers.", tests.Success, testID)

		expected, got := 0, len(n.inactivePeers)
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
		if err != ErrPeerExists {
			t.Fatalf("\t%s\tTest %d:\tShould get ErrPeerExists when trying to set an existing peer, but got %v.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould get ErrPeerExists when trying to set an existing peer.", tests.Success, testID)

		newPeers := n.NewPeers()
		if len(newPeers) != 1 && newPeers[0].id != peer1.id {
			t.Fatalf("\t%s\tTest %d:\tShould have only peer %s in newPeers: %v.", tests.Failed, testID, peer1.id, n.newPeers)
		}
		t.Logf("\t%s\tTest %d:\tShould have only peer %s in newPeers.", tests.Success, testID, peer1.id)

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

		if _, err = n.GetPeer("blahblah"); err != ErrPeerNotFound {
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

		if _, err := n.GetPeer(peer1.id); err != ErrPeerNotFound {
			t.Fatalf("\t%s\tTest %d:\tShould get ErrPeerNotFound when getting inactive client, but got:  %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould get ErrPeerNotFound when getting inactive client.", tests.Success, testID)

		if len(n.inactivePeers) > 0 && n.inactivePeers[peer1.id].id != peer1.id {
			t.Fatalf("\t%s\tTest %d:\tShould have updated peer in table of inactive peers: %+v.", tests.Failed, testID, n.inactivePeers)
		}
		t.Logf("\t%s\tTest %d:\tShould have updated peer in table of inactive peers.", tests.Success, testID)

		updPeer = peer2
		updPeer.addr = "127.0.0.1:55552"
		if err := n.UpdatePeer(&updPeer, false); err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to update without marking the peer as inactive:  %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to update without marking the peer as inactive.", tests.Success, testID)

		l := len(n.inactivePeers)
		if l != 1 {
			t.Fatalf("\t%s\tTest %d:\tShould have 1 peer marked as inactive, but got: %d.", tests.Failed, testID, l)
		}
		t.Logf("\t%s\tTest %d:\tShould have 1 peer marked as inactive.", tests.Success, testID)

		// InactivePeer
		if err := n.MarkPeerInactive(peer2.id); err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to mark peer as inactive:  %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to mark peer as inactive.", tests.Success, testID)

		if l := len(n.inactivePeers); l != 2 {
			t.Fatalf("\t%s\tTest %d:\tShould have 2 peers marked as inactive, but got: %d.", tests.Failed, testID, l)
		}
		t.Logf("\t%s\tTest %d:\tShould have 2 peers marked as inactive.", tests.Success, testID)

		inactivePeers := n.InactivePeers()
		if l := len(inactivePeers); l != 2 {
			t.Fatalf("\t%s\tTest %d:\tShould have 2 peers in inactivePeers, but got: %v.", tests.Failed, testID, l)
		}
		t.Logf("\t%s\tTest %d:\tShould have 2 peers in inactivePeers.", tests.Success, testID)

		// RemovePeer
		n.AddPeer(&peer2)

		if err := n.RemovePeer(peer2.id); err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to remove a peer: %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to remove a new peer.", tests.Success, testID)

		if _, err = n.GetPeer(peer2.id); err != ErrPeerNotFound {
			t.Fatalf("\t%s\tTest %d:\tShould get ErrPeerNotFound when trying to get a removed peer, but got %v.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould get ErrPeerNotFound when trying to get a removed peer.", tests.Success, testID)
	}
}
