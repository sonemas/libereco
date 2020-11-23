package node

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/sonemas/libereco/business/protobuf/networking"
	"github.com/sonemas/libereco/business/tests"
	"google.golang.org/grpc"
)

var nodes []*Node

func tearDown() {
	for _, node := range nodes {
		node.Shutdown(context.Background())
	}
}

func newNode() (int, error) {
	i := len(nodes)

	logger := log.New(os.Stdout, fmt.Sprintf("NODE %d : ", i), log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	n, err := New(logger, fmt.Sprintf("0.0.0.0:%d/%d", 50000+i, i))
	if err != nil {
		return i, err
	}
	n.DialOptions(grpc.WithInsecure())
	n.PingInterval = 5 * time.Hour // Disable automated pings while testing.

	nodes = append(nodes, n)

	go func() {
		if err := n.Serve(); err != nil {
			panic(err)
		}
	}()

	return i, err
}

func TestGRPC(t *testing.T) {
	defer tearDown()

	t.Log("Given the need to work with Nodes via GRPC")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a bootstrap Node.", testID)

		i, err := newNode()
		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to create node %d: %s.", tests.Failed, testID, i, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to create node %d.", tests.Success, testID, i)

		// err = nodes[0].Register(&networking.RegisterRequest{Id: "9999", Addr: "0.0.0.0:9999"}, nil)
		// if err != nil {
		// 	t.Fatalf("\t%s\tTest %d:\tShould be able to register with node %d: %s.", tests.Failed, testID, i, err)
		// }
		// t.Logf("\t%s\tTest %d:\tShould be able to register with node %d.", tests.Success, testID, i)

		// Register
		conn, err := grpc.Dial(nodes[0].addr, grpc.WithInsecure())
		client := networking.NewNetworkingServiceClient(conn)
		defer conn.Close()

		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to dial node %d: %s.", tests.Failed, testID, i, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to dial node %d.", tests.Success, testID, i)

		stream, err := client.Register(context.Background(), &networking.RegisterRequest{Id: "9999", Addr: "0.0.0.0:9999"})
		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to register with node %d: %s.", tests.Failed, testID, i, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to register with node %d.", tests.Success, testID, i)

		peers := []*networking.Node{}
		for {
			peer, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to get strean data from node %d: %s.", tests.Failed, testID, i, err)
			}
			peers = append(peers, peer)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to get strean data from node %d.", tests.Success, testID, i)

		if l := len(peers); l != 1 {
			t.Fatalf("\t%s\tTest %d:\tShould have received 1 peer via stream, but got: %d.", tests.Failed, testID, l)
		}
		t.Logf("\t%s\tTest %d:\tShould have received 1 peer via stream.", tests.Success, testID)

		{
			got, expected := peers[0].Id, "9999"
			if got != expected {
				t.Fatalf("\t%s\tTest %d:\tShould have received peer with ID %q via stream, but got: %q.", tests.Failed, testID, expected, got)
			}
			t.Logf("\t%s\tTest %d:\tShould have received peer with ID %q via stream.", tests.Success, testID, expected)
		}

		{
			got, expected := peers[0].Status, networking.Node_NODE_STATUS_JOINED
			if got != expected {
				t.Fatalf("\t%s\tTest %d:\tShould have received peer with status %v via stream, but got: %v.", tests.Failed, testID, expected, got)
			}
			t.Logf("\t%s\tTest %d:\tShould have received peer with status %v via stream.", tests.Success, testID, expected)
		}

		if l := len(nodes[0].peers); l != 1 {
			t.Fatalf("\t%s\tTest %d:\tShould have 1 peer registered with node %d, but got: %d :\n%v.", tests.Failed, testID, i, l, nodes[0])
		}
		t.Logf("\t%s\tTest %d:\tShould have 1 peer registered with node %d.", tests.Success, testID, i)

		if l := len(nodes[0].newPeers); l != 0 {
			t.Fatalf("\t%s\tTest %d:\tShould have 0 new peers registered with node %d, but got: %d.", tests.Failed, testID, i, l)
		}
		t.Logf("\t%s\tTest %d:\tShould have 0 new peers registered with node %d.", tests.Success, testID, i)

		// Ping
		err = nodes[0].MarkPeerInactive("9999")
		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to mark peer as inactive:  %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to mark peer as inactive.", tests.Success, testID)

		stream, err = client.Ping(context.Background(), &networking.EmptyRequest{})
		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to register with node %d: %s.", tests.Failed, testID, i, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to register with node %d.", tests.Success, testID, i)

		peers2 := []*networking.Node{}
		for {
			peer, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to get strean data from node %d: %s.", tests.Failed, testID, i, err)
			}
			peers2 = append(peers2, peer)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to get strean data from node %d.", tests.Success, testID, i)

		if l := len(peers2); l != 1 {
			t.Fatalf("\t%s\tTest %d:\tShould have received 1 peer via stream, but got: %d.", tests.Failed, testID, l)
		}
		t.Logf("\t%s\tTest %d:\tShould have received 1 peer via stream.", tests.Success, testID)

		{
			got, expected := peers2[0].Id, "9999"
			if got != expected {
				t.Fatalf("\t%s\tTest %d:\tShould have received peer with ID %q via stream, but got: %q.", tests.Failed, testID, expected, got)
			}
			t.Logf("\t%s\tTest %d:\tShould have received peer with ID %q via stream.", tests.Success, testID, expected)
		}

		{
			got, expected := peers2[0].Status, networking.Node_NODE_STATUS_FAILED
			if got != expected {
				t.Fatalf("\t%s\tTest %d:\tShould have received peer with status %v via stream, but got: %v.", tests.Failed, testID, expected, got)
			}
			t.Logf("\t%s\tTest %d:\tShould have received peer with status %v via stream.", tests.Success, testID, expected)
		}

		if l := len(nodes[0].peers); l != 0 {
			t.Fatalf("\t%s\tTest %d:\tShould have 0 peers registered with node %d, but got: %d :\n%v.", tests.Failed, testID, i, l, nodes[0])
		}
		t.Logf("\t%s\tTest %d:\tShould have 0 peer registered with node %d.", tests.Success, testID, i)

		if l := len(nodes[0].newPeers); l != 0 {
			t.Fatalf("\t%s\tTest %d:\tShould have 0 faulty peers registered with node %d, but got: %d.", tests.Failed, testID, i, l)
		}
		t.Logf("\t%s\tTest %d:\tShould have 0 faulty peers registered with node %d.", tests.Success, testID, i)

		// PingReq
		i2, err := newNode()
		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to create node %d: %s.", tests.Failed, testID, i2, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to create node %d.", tests.Success, testID, i2)

		res, err := client.PingReq(context.Background(), &networking.PingRequest{Node: &networking.Node{Id: "0", Addr: nodes[1].addr}})
		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to ping node %d: %s.", tests.Failed, testID, i, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to ping node %d.", tests.Success, testID, i)

		if !res.Success {
			t.Fatalf("\t%s\tTest %d:\tShould get success response to ping request.", tests.Failed, testID)
		}
		t.Logf("\t%s\tTest %d:\tShould get success response to ping request.", tests.Success, testID)
	}

}
