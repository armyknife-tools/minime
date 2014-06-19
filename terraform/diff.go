package terraform

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
)

// Diff tracks the differences between resources to apply.
type Diff struct {
	Resources map[string]*ResourceDiff
	once      sync.Once
}

// ReadDiff reads a diff structure out of a reader in the format that
// was written by WriteDiff.
func ReadDiff(src io.Reader) (*Diff, error) {
	var result *Diff

	dec := gob.NewDecoder(src)
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// WriteDiff writes a diff somewhere in a binary format.
func WriteDiff(d *Diff, dst io.Writer) error {
	return gob.NewEncoder(dst).Encode(d)
}

func (d *Diff) init() {
	d.once.Do(func() {
		if d.Resources == nil {
			d.Resources = make(map[string]*ResourceDiff)
		}
	})
}

// String outputs the diff in a long but command-line friendly output
// format that users can read to quickly inspect a diff.
func (d *Diff) String() string {
	var buf bytes.Buffer

	names := make([]string, 0, len(d.Resources))
	for name, _ := range d.Resources {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		rdiff := d.Resources[name]

		crud := "UPDATE"
		if rdiff.RequiresNew() {
			crud = "CREATE"
		}

		buf.WriteString(fmt.Sprintf(
			"%s: %s\n",
			crud,
			name))

		keyLen := 0
		keys := make([]string, 0, len(rdiff.Attributes))
		for key, _ := range rdiff.Attributes {
			keys = append(keys, key)
			if len(key) > keyLen {
				keyLen = len(key)
			}
		}
		sort.Strings(keys)

		for _, attrK := range keys {
			attrDiff := rdiff.Attributes[attrK]

			v := attrDiff.New
			if attrDiff.NewComputed {
				v = "<computed>"
			}

			newResource := ""
			if attrDiff.RequiresNew {
				newResource = " (forces new resource)"
			}

			buf.WriteString(fmt.Sprintf(
				"  %s:%s %#v => %#v%s\n",
				attrK,
				strings.Repeat(" ", keyLen-len(attrK)),
				attrDiff.Old,
				v,
				newResource))
		}
	}

	return buf.String()
}

// ResourceDiff is the diff of a resource from some state to another.
type ResourceDiff struct {
	Attributes map[string]*ResourceAttrDiff
}

// ResourceAttrDiff is the diff of a single attribute of a resource.
type ResourceAttrDiff struct {
	Old         string // Old Value
	New         string // New Value
	NewComputed bool   // True if new value is computed (unknown currently)
	RequiresNew bool   // True if change requires new resource
}

// RequiresNew returns true if the diff requires the creation of a new
// resource (implying the destruction of the old).
func (d *ResourceDiff) RequiresNew() bool {
	if d == nil {
		return false
	}

	for _, rd := range d.Attributes {
		if rd.RequiresNew {
			return true
		}
	}

	return false
}
