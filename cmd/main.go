package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/s3rj1k/go-fanotify/fanotify"
	"golang.org/x/sys/unix"
)

const (
	procFsFdInfo = "/proc/self/fd"
)

func main() {
	log.SetFlags(log.Lshortfile)

	nd, err := fanotify.Initialize(fanotify.FAN_CLASS_NOTIF|fanotify.FAN_UNLIMITED_QUEUE, os.O_RDONLY|unix.O_LARGEFILE)
	if err != nil {
		log.Fatalf("fanotify: %v\n", err)
	}

	var mountpoint string

	if val, ok := os.LookupEnv("MOUNT_POINT"); !ok {
		mountpoint = "/"
	} else {
		mountpoint = val
	}

	if err = nd.Mark(
		fanotify.FAN_MARK_ADD|
			fanotify.FAN_MARK_MOUNT,
		fanotify.FAN_MODIFY|
			fanotify.FAN_CLOSE_WRITE,
		-1,
		mountpoint,
	); err != nil {
		log.Fatalf("fanotify: %v\n", err)
	}

	f := func(nd *fanotify.NotifyFD) (string, error) {
		data, err := nd.GetEvent()
		if err != nil {
			return "", fmt.Errorf("fanotify: %w", err)
		}

		defer data.File.Close()

		if data.Version != fanotify.FANOTIFY_METADATA_VERSION {
			return "", fmt.Errorf("fanotify: wrong metadata version")
		}

		if int(data.PID) == os.Getpid() {
			return "", fmt.Errorf("fanotify: self PID")
		}

		path, err := os.Readlink(
			filepath.Join(
				procFsFdInfo,
				strconv.FormatUint(
					uint64(data.File.Fd()),
					10,
				),
			),
		)
		if err != nil {
			return "", fmt.Errorf("fanotify: %w", err)
		}

		if (data.Mask & fanotify.FAN_MODIFY) == fanotify.FAN_MODIFY {
			return fmt.Sprintf("PID:%d %s", data.PID, path), nil
		}

		if (data.Mask & fanotify.FAN_CLOSE_WRITE) == fanotify.FAN_CLOSE_WRITE {
			return fmt.Sprintf("PID:%d %s", data.PID, path), nil
		}

		return "", fmt.Errorf("fanotify: unknown event")
	}

	for {
		if str, err := f(nd); err == nil && len(str) > 0 {
			fmt.Printf("%s\n", str)
		}
	}
}
