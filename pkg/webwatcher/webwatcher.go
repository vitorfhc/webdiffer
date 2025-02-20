package webwatcher

import (
	"io"
	"net/http"

	"github.com/vitorfhc/webdiff/pkg/store"
	"github.com/vitorfhc/webdiff/pkg/types"
)

type WebDiffer struct {
	store store.TargetStore
}

func NewWebWatcher(store store.TargetStore) *WebDiffer {
	return &WebDiffer{store: store}
}

func (w *WebDiffer) Run() ([]types.Diff, error) {
	targets, err := w.store.ListTargets()
	if err != nil {
		return nil, err
	}

	var diffs []types.Diff
	for _, target := range targets {
		lastResult, err := w.store.GetResult(target)
		if err != nil {
			return nil, err
		}

		resp, err := http.Get(target.URL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		result := types.Result{
			StatusCode: resp.StatusCode,
			Body:       string(body),
		}

		if lastResult != result {
			w.store.UpdateTarget(types.Target{
				URL:        target.URL,
				LastResult: result,
			})

			diffs = append(diffs, types.Diff{
				Target: target,
				Old:    lastResult,
				New:    result,
			})
		}
	}

	return diffs, nil
}
