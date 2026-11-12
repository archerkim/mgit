package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrMissingParent = errors.New("missing parent")
	ErrNoPermission  = errors.New("no permission")
	ErrStaleLock     = errors.New("stale lock")
)

type lockfile struct {
	filePath string
	lock     *os.File
}

func (l *lockfile) hold_for_update() (bool, error) {
	if l.lock == nil {
		lock, err := os.OpenFile(lockPath(l.filePath), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				return false, nil
			} else if errors.Is(err, os.ErrNotExist) {
				return false, fmt.Errorf("%w: %v", ErrMissingParent, err)
			} else if errors.Is(err, os.ErrPermission) {
				return false, fmt.Errorf("%w: %v", ErrNoPermission, err)
			} else {
				return false, fmt.Errorf("failed to hold lock: %v", err)
			}
		}
		l.lock = lock
		return true, nil
	}
	return false, fmt.Errorf("Already locked by other process")
}

func (l *lockfile) write(str string) error {
	err := l.raise_on_stale_lock()
	if err != nil {
		return err
	}
	_, err = l.lock.WriteString(str)
	if err != nil {
		return err
	}
	return nil
}
func (l *lockfile) commit() error {
	err := l.raise_on_stale_lock()
	if err != nil {
		return err
	}
	defer l.lock.Close()

	err = os.Rename(lockPath(l.filePath), l.filePath)
	if err != nil {
		return err
	}
	l.lock = nil
	return nil
}

func (l *lockfile) raise_on_stale_lock() error {
	if l.lock == nil {
		return fmt.Errorf("%w: Not holding lock on %s", ErrStaleLock, lockPath(l.filePath))
	}
	return nil
}

func lockPath(filePath string) string {
	return filepath.Join(filePath, "lock")
}
