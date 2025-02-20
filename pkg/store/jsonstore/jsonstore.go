package jsonstore

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vitorfhc/webdiffer/pkg/helpers"
	"github.com/vitorfhc/webdiffer/pkg/store"
	"github.com/vitorfhc/webdiffer/pkg/types"
)

type JSONStore struct {
	filepath string
}

func NewJSONStore(filepath string) *JSONStore {
	return &JSONStore{filepath: filepath}
}

func (j *JSONStore) ListTargets() ([]types.Target, error) {
	data, err := helpers.ReadOrCreateJSON(j.filepath)
	if err != nil {
		return nil, err
	}

	var targets []types.Target
	err = json.Unmarshal(data, &targets)
	if err != nil {
		return nil, err
	}

	return targets, nil
}

func (j *JSONStore) UpdateTarget(target types.Target) error {
	data, err := helpers.ReadOrCreateJSON(j.filepath)
	if err != nil {
		return err
	}

	target.URL, err = helpers.NormalizeURL(target.URL)
	if err != nil {
		return err
	}

	var targets []types.Target
	err = json.Unmarshal(data, &targets)
	if err != nil {
		return err
	}

	updated := false
	for i, t := range targets {
		if t.URL == target.URL {
			targets[i] = target
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("target with URL %q not found", target.URL)
	}

	newData, err := json.Marshal(targets)
	if err != nil {
		return err
	}

	err = os.WriteFile(j.filepath, newData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (j *JSONStore) GetResult(target types.Target) (types.Result, error) {
	targets, err := j.ListTargets()
	if err != nil {
		return types.Result{}, err
	}

	for _, t := range targets {
		if t.URL == target.URL {
			return t.LastResult, nil
		}
	}

	return types.Result{}, fmt.Errorf("target with URL %q not found", target.URL)
}

func (j *JSONStore) InsertTarget(target types.Target) error {
	data, err := helpers.ReadOrCreateJSON(j.filepath)
	if err != nil {
		return err
	}

	target.URL, err = helpers.NormalizeURL(target.URL)
	if err != nil {
		return err
	}

	var targets []types.Target
	err = json.Unmarshal(data, &targets)
	if err != nil {
		return err
	}

	for _, t := range targets {
		if t.URL == target.URL {
			return fmt.Errorf("target with URL %q already exists", target.URL)
		}
	}

	targets = append(targets, target)

	newData, err := json.Marshal(targets)
	if err != nil {
		return err
	}

	err = os.WriteFile(j.filepath, newData, 0644)
	if err != nil {
		return err
	}

	return nil
}

var _ store.TargetStore = &JSONStore{}
