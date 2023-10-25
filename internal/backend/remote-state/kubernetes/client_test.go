// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubernetes

import (
	"context"
	"testing"

	"github.com/opentofu/opentofu/internal/backend"
	"github.com/opentofu/opentofu/internal/states/remote"
	"github.com/opentofu/opentofu/internal/states/statemgr"
)

func TestRemoteClient_impl(t *testing.T) {
	var _ remote.Client = new(RemoteClient)
	var _ remote.ClientLocker = new(RemoteClient)
}

func TestRemoteClient(t *testing.T) {
	testACC(t)
	defer cleanupK8sResources(t)

	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"secret_suffix": secretSuffix,
	}))

	ctx := context.Background()

	state, err := b.StateMgr(ctx, backend.DefaultStateName)
	if err != nil {
		t.Fatal(err)
	}

	remote.TestClient(t, state.(*remote.State).Client)
}

func TestRemoteClientLocks(t *testing.T) {
	testACC(t)
	defer cleanupK8sResources(t)

	b1 := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"secret_suffix": secretSuffix,
	}))

	b2 := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"secret_suffix": secretSuffix,
	}))

	ctx := context.Background()

	s1, err := b1.StateMgr(ctx, backend.DefaultStateName)
	if err != nil {
		t.Fatal(err)
	}

	s2, err := b2.StateMgr(ctx, backend.DefaultStateName)
	if err != nil {
		t.Fatal(err)
	}

	remote.TestRemoteLocks(t, s1.(*remote.State).Client, s2.(*remote.State).Client)
}

func TestForceUnlock(t *testing.T) {
	testACC(t)
	defer cleanupK8sResources(t)

	b1 := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"secret_suffix": secretSuffix,
	}))

	b2 := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"secret_suffix": secretSuffix,
	}))

	ctx := context.Background()

	// first test with default
	s1, err := b1.StateMgr(ctx, backend.DefaultStateName)
	if err != nil {
		t.Fatal(err)
	}

	info := statemgr.NewLockInfo()
	info.Operation = "test"
	info.Who = "clientA"

	lockID, err := s1.Lock(ctx, info)
	if err != nil {
		t.Fatal("unable to get initial lock:", err)
	}

	// s1 is now locked, get the same state through s2 and unlock it
	s2, err := b2.StateMgr(ctx, backend.DefaultStateName)
	if err != nil {
		t.Fatal("failed to get default state to force unlock:", err)
	}

	if err := s2.Unlock(ctx, lockID); err != nil {
		t.Fatal("failed to force-unlock default state")
	}

	// now try the same thing with a named state
	// first test with default
	s1, err = b1.StateMgr(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}

	info = statemgr.NewLockInfo()
	info.Operation = "test"
	info.Who = "clientA"

	lockID, err = s1.Lock(ctx, info)
	if err != nil {
		t.Fatal("unable to get initial lock:", err)
	}

	// s1 is now locked, get the same state through s2 and unlock it
	s2, err = b2.StateMgr(ctx, "test")
	if err != nil {
		t.Fatal("failed to get named state to force unlock:", err)
	}

	if err = s2.Unlock(ctx, lockID); err != nil {
		t.Fatal("failed to force-unlock named state")
	}
}
