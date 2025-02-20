package store

import (
	"github.com/vitorfhc/webdiffer/pkg/types"
)

type TargetStore interface {
	InsertTarget(target types.Target) error
	ListTargets() ([]types.Target, error)
	UpdateTarget(target types.Target) error
	GetResult(target types.Target) (types.Result, error)
}
