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

// Package fanotify package provides a simple fanotify API
package fanotify

import (
	"bufio"
	"encoding/binary"
	"os"

	"golang.org/x/sys/unix"
)

// Internal eventMetadata struct, used for fanotify communication,
// 'fanotify_event_metadata' in fanotify.h
type eventMetadata struct {
	EventLength    uint32
	Version        uint8
	Reserved       uint8
	MetadataLength uint16
	Mask           uint64
	Fd             int32
	PID            int32
}

// Internal response struct, used for fanotify communication,
// 'fanotify_response' in fanotify.h
type response struct {
	Fd       int32
	Response uint32
}

// EventMetadata is a struct returned from NotifyFD.GetEvent
//
// The File member needs to be Closed after usage,
// to prevent an FD leak
type EventMetadata struct {
	EventLength    uint32
	Version        uint8
	Reserved       uint8
	MetadataLength uint16
	Mask           uint64
	File           *os.File
	PID            int32
}

// NotifyFD is a notify handle, used by all fanotify functions
type NotifyFD struct {
	f *os.File
	r *bufio.Reader
}

// Initialize initializes the fanotify support
func Initialize(faFlags, openFlags int) (*NotifyFD, error) {
	fd, _, errno := unix.Syscall(
		unix.SYS_FANOTIFY_INIT,
		uintptr(faFlags),
		uintptr(openFlags),
		uintptr(0),
	)

	var err error
	if errno != 0 {
		err = errno
	}

	f := os.NewFile(fd, "")

	return &NotifyFD{
		f,
		bufio.NewReader(f),
	}, err
}

// GetEvent returns an event from the fanotify handle
func (nd *NotifyFD) GetEvent() (*EventMetadata, error) {
	ev := new(eventMetadata)

	err := binary.Read(nd.r, binary.LittleEndian, ev)
	if err != nil {
		return nil, err
	}

	return &EventMetadata{
		ev.EventLength,
		ev.Version,
		ev.Reserved,
		ev.MetadataLength,
		ev.Mask,
		os.NewFile(uintptr(ev.Fd), ""),
		ev.PID,
	}, nil
}

// ResponseAllow sends an allow message back to fanotify, used for permission checks
func (nd *NotifyFD) ResponseAllow(ev *EventMetadata) error {
	return binary.Write(
		nd.f,
		binary.LittleEndian,
		&response{
			Fd:       int32(ev.File.Fd()),
			Response: FAN_ALLOW,
		},
	)
}

// ResponseDeny sends a deny message back to fanotify, used for permission checks
func (nd *NotifyFD) ResponseDeny(ev *EventMetadata) error {
	return binary.Write(
		nd.f,
		binary.LittleEndian,
		&response{
			Fd:       int32(ev.File.Fd()),
			Response: FAN_DENY,
		},
	)
}
