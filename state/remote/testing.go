package remote

import (
	"bytes"
	"testing"

	"github.com/hashicorp/terraform/state"
	"github.com/hashicorp/terraform/terraform"
)

// TestClient is a generic function to test any client.
func TestClient(t *testing.T, c Client) {
	var buf bytes.Buffer
	s := state.TestStateInitial()
	if err := terraform.WriteState(s, &buf); err != nil {
		t.Fatalf("err: %s", err)
	}
	data := buf.Bytes()

	if err := c.Put(data); err != nil {
		t.Fatalf("put: %s", err)
	}

	p, err := c.Get()
	if err != nil {
		t.Fatalf("get: %s", err)
	}
	if !bytes.Equal(p.Data, data) {
		t.Fatalf("bad: %#v", p)
	}

	if err := c.Delete(); err != nil {
		t.Fatalf("delete: %s", err)
	}

	p, err = c.Get()
	if err != nil {
		t.Fatalf("get: %s", err)
	}
	if p != nil {
		t.Fatalf("bad: %#v", p)
	}
}

// Test the lock implementation for a remote.Client.
// This test requires 2 client instances, in oder to have multiple remote
// clients since some implementations may tie the client to the lock, or may
// have reentrant locks.
func TestRemoteLocks(t *testing.T, a, b Client) {
	lockerA, ok := a.(state.Locker)
	if !ok {
		t.Fatal("client A not a state.Locker")
	}

	lockerB, ok := b.(state.Locker)
	if !ok {
		t.Fatal("client B not a state.Locker")
	}

	infoA := state.NewLockInfo()
	infoA.Operation = "test"
	infoA.Who = "clientA"

	infoB := state.NewLockInfo()
	infoB.Operation = "test"
	infoB.Who = "clientB"

	if _, err := lockerA.Lock(infoA); err != nil {
		t.Fatal("unable to get initial lock:", err)
	}

	if _, err := lockerB.Lock(infoB); err == nil {
		lockerA.Unlock("")
		t.Fatal("client B obtained lock while held by client A")
	} else {
		t.Log("lock info error:", err)
	}

	if err := lockerA.Unlock(""); err != nil {
		t.Fatal("error unlocking client A", err)
	}

	if _, err := lockerB.Lock(infoB); err != nil {
		t.Fatal("unable to obtain lock from client B")
	}

	if err := lockerB.Unlock(""); err != nil {
		t.Fatal("error unlocking client B:", err)
	}

	// unlock should be repeatable
	if err := lockerA.Unlock(""); err != nil {
		t.Fatal("Unlock error from client A when state was not locked:", err)
	}
}
