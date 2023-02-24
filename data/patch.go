package data

import (
	"fmt"
	"reflect"
	"sort"
)

// Patch represents a series of modifications to Data.
type Patch struct {
	actions []*PatchAction
}

// PatchAction presents a patch action.
type PatchAction struct {
	Deletes []string        `data:"deletes"`
	Updates map[string]Data `data:"updates"`
}

// NewPatch creates a new Patch.
func NewPatch() *Patch {
	return &Patch{}
}

// Add adds a new patch action.
// Every patch action is composed of deletes and updates.
// When applying, deletes will be deleted first, then updates will be merged into Data.
// The key of updates is a query of Data. When applying, updates will be merged into Data
//
// Note that `Patch#Apply`/`Patch#ApplyTo` uses `Merge`/`MergeTo` to update Data.
// The merge series functions will traverse map/slice deeply,
// which means that the new value cannot simply overwrite the old value.
// If you want the new value to overwrite the old value instead of merging,
// use deletes to delete the old value first and then merge.
func (patch *Patch) Add(deletes []string, updates map[string]Data) {
	patch.actions = append(patch.actions, &PatchAction{
		Deletes: deletes,
		Updates: updates,
	})
}

// Actions returns all actions.
func (patch *Patch) Actions() []*PatchAction {
	return patch.actions
}

// Apply copies d and applies all changes to the copy of d.
//
// Apply returns error when:
//
//   - If a query in updates cannot find the corresponding element.
//   - If a query in updates finds a result that is not a RawData.
func (patch *Patch) Apply(d Data) (applied Data, err error) {
	d = d.Clone()

	if err = patch.ApplyTo(&d); err != nil {
		return
	}

	applied = d
	return
}

// ApplyTo applies all changes to target.
//
// ApplyTo returns error in the same conditions as Apply.
func (patch *Patch) ApplyTo(target *Data) error {
	if target == nil {
		return nil
	}

	for _, action := range patch.actions {
		if err := action.ApplyTo(target); err != nil {
			return err
		}
	}

	return nil
}

// ApplyTo applies an action to target.
func (action *PatchAction) ApplyTo(target *Data) error {
	data := target.data

	// Delete first.
	data.Delete(action.Deletes...)
	target.data = data // Delete may reset data.

	if len(action.Updates) == 0 {
		return nil
	}

	// Thenn update.
	queries := make([]string, 0, len(action.Updates))

	for query := range action.Updates {
		queries = append(queries, query)
	}

	// The queries should be sorted in ascending order to make sure that
	// the upper level data will be updated first.
	// For example, if both "a" and "a.b" should be updated, it updates "a" first.
	sort.Strings(queries)

	for _, query := range queries {
		v := data.Query(query)

		if v == nil {
			return fmt.Errorf("fail to apply patch due to invalid query `%v` when updating", query)
		}

		d, ok := v.(RawData)

		if !ok {
			return fmt.Errorf("fail to apply patch due to query `%v` pointing to a value in unsupported type", query)
		}

		merge(reflect.ValueOf(d), action.Updates[query].data)
	}

	return nil
}
