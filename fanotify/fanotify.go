// Package fanotify package provides a simple fanotify API
package fanotify

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/sys/unix"
)

// Procfs constants
const (
	ProcFsFdInfo = "/proc/self/fd"
)

// EventMetadata is a struct returned from 'NotifyFD.GetEvent'.
type EventMetadata struct {
	unix.FanotifyEventMetadata
}

// GetPID return PID from event metadata.
func (metadata EventMetadata) GetPID() int {
	return int(metadata.Pid)
}

// GetPath returns path to file for FD inside event metadata.
func (metadata EventMetadata) GetPath() (string, error) {
	path, err := os.Readlink(
		filepath.Join(
			ProcFsFdInfo,
			strconv.FormatUint(
				uint64(metadata.Fd),
				10,
			),
		),
	)
	if err != nil {
		return "", fmt.Errorf("fanotify: %w", err)
	}

	return path, nil
}

// MatchMask returns 'true' when event metadata matches specified mask.
func (metadata EventMetadata) MatchMask(mask int) bool {
	return (metadata.Mask & uint64(mask)) == uint64(mask)
}

// File returns pointer to os.File created from event metadata supplied FD.
// File needs to be Closed after usage, to prevent an FD leak.
func (metadata EventMetadata) File() *os.File {
	return os.NewFile(uintptr(metadata.Fd), "")
}

// NotifyFD is a notify file handle, used by all fanotify functions.
type NotifyFD struct {
	fd   int
	file *os.File
}

// Initialize initializes the fanotify support.
func Initialize(fanotifyFlags uint, openFlags int) (*NotifyFD, error) {
	fd, err := unix.FanotifyInit(fanotifyFlags, uint(openFlags))
	if err != nil {
		return nil, fmt.Errorf("fanotify: %w", err)
	}

	return &NotifyFD{
		fd:   fd,
		file: os.NewFile(uintptr(fd), ""),
	}, err
}

// Mark implements Add/Delete/Modify for a fanotify mark.
func (handle *NotifyFD) Mark(flags uint, mask uint64, dirFd int, path string) error {
	if err := unix.FanotifyMark(handle.fd, flags, mask, dirFd, path); err != nil {
		return fmt.Errorf("fanotify: %w", err)
	}

	return nil
}

// GetEvent returns an event from the fanotify handle.
func (handle *NotifyFD) GetEvent(skipPIDs ...int) (*EventMetadata, error) {
	event := new(EventMetadata)

	err := binary.Read(bufio.NewReader(handle.file), binary.LittleEndian, event)
	if err != nil {
		return nil, fmt.Errorf("fanotify: %w", err)
	}

	if event.Vers != FANOTIFY_METADATA_VERSION {
		return nil, fmt.Errorf("fanotify: wrong metadata version")
	}

	for i := range skipPIDs {
		if int(event.Pid) == skipPIDs[i] {
			return nil, nil
		}
	}

	return event, nil
}

// ResponseAllow sends an allow message back to fanotify, used for permission checks.
func (handle *NotifyFD) ResponseAllow(ev *EventMetadata) error {
	if err := binary.Write(
		handle.file,
		binary.LittleEndian,
		&unix.FanotifyResponse{
			Fd:       ev.Fd,
			Response: FAN_ALLOW,
		},
	); err != nil {
		return fmt.Errorf("fanotify: %w", err)
	}

	return nil
}

// ResponseDeny sends a deny message back to fanotify, used for permission checks.
func (handle *NotifyFD) ResponseDeny(ev *EventMetadata) error {
	if err := binary.Write(
		handle.file,
		binary.LittleEndian,
		&unix.FanotifyResponse{
			Fd:       ev.Fd,
			Response: FAN_DENY,
		},
	); err != nil {
		return fmt.Errorf("fanotify: %w", err)
	}

	return nil
}
