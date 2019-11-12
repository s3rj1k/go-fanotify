// Copyright (c) 2012, Moritz Bitsch <moritzbitsch@googlemail.com>
//
// Permission to use, copy, modify, and/or distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

// Package fanotify package provides a simple fanotify api
package fanotify

import (
	"bufio"
	"encoding/binary"
	"os"

	"golang.org/x/sys/unix"
)

// Flags used as first parameter to Initiliaze
const (
	/* flags used for fanotify_init() */

	// FAN_CLOEXEC
	FanCloExec = 0x00000001
	// FAN_NONBLOCK
	FanNonBlock = 0x00000002

	/* These are NOT bitwise flags. Both bits are used togther. */

	// FAN_CLASS_NOTIF
	FanClassNotif = 0x00000000
	// FAN_CLASS_CONTENT
	FanClassContent = 0x00000004
	// FAN_CLASS_PRE_CONTENT
	FanClassPreContent = 0x00000008

	// FAN_ALL_CLASS_BITS
	FanAllClassBits = FanClassNotif |
		FanClassContent |
		FanClassPreContent

	// FAN_UNLIMITED_QUEUE
	FanUnlimitedQueue = 0x00000010
	// FAN_UNLIMITED_MARKS
	FanUnlimitedMarks = 0x00000020

	// FAN_ALL_INIT_FLAGS
	FanAllInitFlags = FanCloExec |
		FanNonBlock |
		FanAllClassBits |
		FanUnlimitedQueue |
		FanUnlimitedMarks
)

// Flags used for the Mark Method
const (
	/* flags used for fanotify_modify_mark() */

	// FAN_MARK_ADD
	FanMarkAdd = 0x00000001
	// FAN_MARK_REMOVE
	FanMarkRemove = 0x00000002
	// FAN_MARK_DONT_FOLLOW
	FanMarkDontFollow = 0x00000004
	// FAN_MARK_ONLYDIR
	FanMarkOnlyDir = 0x00000008
	// FAN_MARK_MOUNT
	FanMarkMount = 0x00000010
	// FAN_MARK_IGNORED_MASK
	FanMarkIgnoredMask = 0x00000020
	// FAN_MARK_IGNORED_SURV_MODIFY
	FanMarkIgnoredSurvModify = 0x00000040
	// FAN_MARK_FLUSH
	FanMarkFlush = 0x00000080

	// FAN_ALL_MARK_FLAGS
	FanAllMarkFlags = FanMarkAdd |
		FanMarkRemove |
		FanMarkDontFollow |
		FanMarkOnlyDir |
		FanMarkMount |
		FanMarkIgnoredMask |
		FanMarkIgnoredSurvModify |
		FanMarkFlush
)

// Event types
const (
	// FAN_ACCESS
	FanAccess = 0x00000001 /* File was accessed */
	// FAN_MODIFY
	FanModify = 0x00000002 /* File was modified */
	// FAN_CLOSE_WRITE
	FanCloseWrite = 0x00000008 /* Writtable file closed */
	// FAN_CLOSE_NOWRITE
	FanCloseNoWrite = 0x00000010 /* Unwrittable file closed */
	// FAN_OPEN
	FanOpen = 0x00000020 /* File was opened */

	// FAN_Q_OVERFLOW
	FanQOverflow = 0x00004000 /* Event queued overflowed */

	// FAN_OPEN_PERM
	FanOpenPerm = 0x00010000 /* File open in perm check */
	// FAN_ACCESS_PERM
	FanAccessPerm = 0x00020000 /* File accessed in perm check */

	// FAN_ONDIR
	FanOnDir = 0x40000000 /* event occurred against dir */

	// FAN_EVENT_ON_CHILD
	FanEventOnChild = 0x08000000 /* interested in child events */

	/* helper events */

	// FAN_CLOSE
	FanClose = FanCloseWrite | FanCloseNoWrite /* close */

	/*
		All of the events - we build the list by hand so that we can add flags in
		the future and not break backward compatibility.  Apps will get only the
		events that they originally wanted. Be sure to add new events here!
	*/

	// FAN_ALL_EVENTS
	FanAllEvents = FanAccess |
		FanModify |
		FanClose |
		FanOpen

	/*
		All events which require a permission response from userspace
	*/

	// FAN_ALL_PERM_EVENTS
	FanAllPermEvents = FanOpenPerm |
		FanAccessPerm

	// FAN_ALL_OUTGOING_EVENTS
	FanAllOutgoingEvents = FanAllEvents |
		FanAllPermEvents |
		FanQOverflow

	// FANOTIFY_METADATA_VERSION
	FanotifyMetadataVersion = 3

	// FAN_ALLOW
	FanAllow = 0x01
	// FAN_DENY
	FanDeny = 0x02
	// FAN_NOFD
	FanNoFD = -1
)

// Internal eventMetadata struct, used for fanotify comm
type eventMetadata struct {
	Len         uint32
	Version     uint8
	Reserved    uint8
	MetadataLen uint16
	Mask        uint64
	Fd          int32
	Pid         int32
}

// Internal response struct, used for fanotify comm
type response struct {
	Fd       int32
	Response uint32
}

// EventMetadata is a struct returned from NotifyFD.GetEvent
//
// The File member needs to be Closed after usage, to prevent
// an Fd leak
type EventMetadata struct {
	Len         uint32
	Version     uint8
	Reserved    uint8
	MetadataLen uint16
	Mask        uint64
	File        *os.File
	Pid         int32
}

// NotifyFD a notify handle, used by all notify functions
type NotifyFD struct {
	f *os.File
	r *bufio.Reader
}

// Initialize inits the notify support
func Initialize(faflags, openflags int) (*NotifyFD, error) {
	fd, _, errno := unix.Syscall(unix.SYS_FANOTIFY_INIT, uintptr(faflags), uintptr(openflags), uintptr(0))

	var err error
	if errno != 0 {
		err = errno
	}

	f := os.NewFile(fd, "")

	return &NotifyFD{f, bufio.NewReader(f)}, err
}

// GetEvent returns an event from the fanotify handle
func (nd *NotifyFD) GetEvent() (*EventMetadata, error) {
	ev := &eventMetadata{}

	err := binary.Read(nd.r, binary.LittleEndian, ev)
	if err != nil {
		return nil, err
	}

	res := &EventMetadata{ev.Len, ev.Version, ev.Reserved, ev.MetadataLen, ev.Mask, os.NewFile(uintptr(ev.Fd), ""), ev.Pid}

	return res, nil
}

// ResponseAllow sends an allow message back to fanotify, used for permission checks
func (nd *NotifyFD) ResponseAllow(ev *EventMetadata) error {
	resp := &response{Fd: int32(ev.File.Fd())}

	resp.Response = FanAllow

	return binary.Write(nd.f, binary.LittleEndian, resp)
}

// ResponseDeny sends a deny message back to fanotify, used for permission checks
func (nd *NotifyFD) ResponseDeny(ev *EventMetadata) error {
	resp := &response{Fd: int32(ev.File.Fd())}

	resp.Response = FanDeny

	return binary.Write(nd.f, binary.LittleEndian, resp)
}
