// Package fanotify package provides a simple fanotify API
package fanotify

/*
	Headers information are taken from fanotify.h
*/

// Event types that user-space can register for.
const (
	FAN_ACCESS        = 0x00000001 // File was accessed
	FAN_MODIFY        = 0x00000002 // File was modified
	FAN_ATTRIB        = 0x00000004 // Metadata changed
	FAN_CLOSE_WRITE   = 0x00000008 // Writtable file closed
	FAN_CLOSE_NOWRITE = 0x00000010 // Unwrittable file closed
	FAN_OPEN          = 0x00000020 // File was opened
	FAN_MOVED_FROM    = 0x00000040 // File was moved from X
	FAN_MOVED_TO      = 0x00000080 // File was moved to Y
	FAN_CREATE        = 0x00000100 // Subfile was created
	FAN_DELETE        = 0x00000200 // Subfile was deleted
	FAN_DELETE_SELF   = 0x00000400 // Self was deleted
	FAN_MOVE_SELF     = 0x00000800 // Self was moved
	FAN_OPEN_EXEC     = 0x00001000 // File was opened for exec

	FAN_Q_OVERFLOW = 0x00004000 // Event queued overflowed

	FAN_OPEN_PERM      = 0x00010000 // File open in perm check
	FAN_ACCESS_PERM    = 0x00020000 // File accessed in perm check
	FAN_OPEN_EXEC_PERM = 0x00040000 // File open/exec in perm check

	FAN_ONDIR = 0x40000000 // Event occurred against dir

	FAN_EVENT_ON_CHILD = 0x08000000 // Interested in child events
)

// Helper event types that user-space can register for.
const (
	FAN_CLOSE = FAN_CLOSE_WRITE | FAN_CLOSE_NOWRITE // Close event
	FAN_MOVE  = FAN_MOVED_FROM | FAN_MOVED_TO       // Move event
)

// Flags used as first parameter to Initiliaze, 'fanotify_init()'.
const (
	FAN_CLOEXEC  = 0x00000001
	FAN_NONBLOCK = 0x00000002

	// These are NOT bitwise flags. Both bits are used together.
	FAN_CLASS_NOTIF       = 0x00000000
	FAN_CLASS_CONTENT     = 0x00000004
	FAN_CLASS_PRE_CONTENT = 0x00000008

	FAN_UNLIMITED_QUEUE = 0x00000010
	FAN_UNLIMITED_MARKS = 0x00000020
	FAN_ENABLE_AUDIT    = 0x00000040

	// Flags to determine fanotify event format.
	FAN_REPORT_TID = 0x00000100 // event->pid is thread ID
	FAN_REPORT_FID = 0x00000200 // Report unique file ID
)

// Flags used for the Mark Method 'fanotify_modify_mark()'.
const (
	FAN_MARK_ADD                 = 0x00000001
	FAN_MARK_REMOVE              = 0x00000002
	FAN_MARK_DONT_FOLLOW         = 0x00000004
	FAN_MARK_ONLYDIR             = 0x00000008
	FAN_MARK_IGNORED_MASK        = 0x00000020
	FAN_MARK_IGNORED_SURV_MODIFY = 0x00000040
	FAN_MARK_FLUSH               = 0x00000080

	// These are NOT bitwise flags. Both bits can be used togther.
	FAN_MARK_INODE      = 0x00000000
	FAN_MARK_MOUNT      = 0x00000010
	FAN_MARK_FILESYSTEM = 0x00000100
)

const (
	FANOTIFY_METADATA_VERSION = 3
)

const (
	FAN_EVENT_INFO_TYPE_FID = 1
)

// Legit userspace responses to a _PERM event.
const (
	FAN_ALLOW = 0x01
	FAN_DENY  = 0x02
	FAN_AUDIT = 0x10 // Bit mask to create audit record for result
)

// No fd set in event
const (
	FAN_NOFD = -1
)
