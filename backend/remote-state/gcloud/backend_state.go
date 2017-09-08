package gcloud

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/hashicorp/terraform/backend"
	"github.com/hashicorp/terraform/state"
	"github.com/hashicorp/terraform/state/remote"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/iterator"
)

const (
	stateFileSuffix = ".tfstate"
	lockFileSuffix  = ".tflock"
)

// States returns a list of names for the states found on GCS. The default
// state is always returned as the first element in the slice.
func (b *Backend) States() ([]string, error) {
	states := []string{backend.DefaultStateName}

	bucket := b.storageClient.Bucket(b.bucketName)
	objs := bucket.Objects(b.storageContext, &storage.Query{
		Delimiter: "/",
		Prefix:    b.stateDir,
	})
	for {
		attrs, err := objs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("querying Cloud Storage failed: %v", err)
		}

		name := path.Base(attrs.Name)
		if !strings.HasSuffix(name, stateFileSuffix) {
			continue
		}
		st := strings.TrimSuffix(name, stateFileSuffix)

		if st != backend.DefaultStateName {
			states = append(states, st)
		}
	}

	sort.Strings(states[1:])
	return states, nil
}

// DeleteState deletes the named state. The "default" state cannot be deleted.
func (b *Backend) DeleteState(name string) error {
	if name == backend.DefaultStateName {
		return fmt.Errorf("cowardly refusing to delete the %q state", name)
	}

	client, err := b.remoteClient(name)
	if err != nil {
		return err
	}

	return client.Delete()
}

// remoteClient returns a RemoteClient for the named state.
func (b *Backend) remoteClient(name string) (*RemoteClient, error) {
	if name == "" {
		return nil, fmt.Errorf("%q is not a valid state name", name)
	}

	return &RemoteClient{
		storageContext: b.storageContext,
		storageClient:  b.storageClient,
		bucketName:     b.bucketName,
		stateFilePath:  path.Join(b.stateDir, name+stateFileSuffix),
		lockFilePath:   path.Join(b.stateDir, name+lockFileSuffix),
	}, nil
}

// State reads and returns the named state from GCS. If the named state does
// not yet exist, a new state file is created.
func (b *Backend) State(name string) (state.State, error) {
	client, err := b.remoteClient(name)
	if err != nil {
		return nil, err
	}

	st := &remote.State{Client: client}
	lockInfo := state.NewLockInfo()
	lockInfo.Operation = "init"
	lockID, err := st.Lock(lockInfo)
	if err != nil {
		return nil, err
	}

	// Local helper function so we can call it multiple places
	unlock := func(baseErr error) error {
		if err := st.Unlock(lockID); err != nil {
			const unlockErrMsg = `%v
Additionally, unlocking the state file on Google Cloud Storage failed:

  Error message: %q
  Lock ID (gen): %v
  Lock file URL: %v

You may have to force-unlock this state in order to use it again.
The GCloud backend acquires a lock during initialization to ensure
the initial state file is created.`
			return fmt.Errorf(unlockErrMsg, baseErr, err.Error(), lockID, client.lockFileURL())
		}

		return baseErr
	}

	// Grab the value
	if err := st.RefreshState(); err != nil {
		return nil, unlock(err)
	}

	// If we have no state, we have to create an empty state
	if v := st.State(); v == nil {
		if err := st.WriteState(terraform.NewState()); err != nil {
			return nil, unlock(err)
		}
		if err := st.PersistState(); err != nil {
			return nil, unlock(err)
		}
	}

	// Unlock, the state should now be initialized
	if err := unlock(nil); err != nil {
		return nil, err
	}

	return st, nil
}
