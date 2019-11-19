package main

import (
	"fmt"
	"log"
	"os"

	"github.com/s3rj1k/go-fanotify/fanotify"
	"golang.org/x/sys/unix"
)

func main() {
	log.SetFlags(log.Lshortfile)

	notify, err := fanotify.Initialize(fanotify.FAN_CLASS_NOTIF|fanotify.FAN_UNLIMITED_QUEUE, os.O_RDONLY|unix.O_LARGEFILE)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	var mountpoint string

	if val, ok := os.LookupEnv("MOUNT_POINT"); !ok {
		mountpoint = "/"
	} else {
		mountpoint = val
	}

	if err = notify.Mark(
		fanotify.FAN_MARK_ADD|
			fanotify.FAN_MARK_MOUNT,
		fanotify.FAN_MODIFY|
			fanotify.FAN_CLOSE_WRITE,
		-1,
		mountpoint,
	); err != nil {
		log.Fatalf("%v\n", err)
	}

	f := func(notify *fanotify.NotifyFD) (string, error) {
		data, err := notify.GetEvent(os.Getpid())
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		defer data.File().Close()

		if data == nil {
			return "", nil
		}

		path, err := data.GetPath()
		if err != nil {
			return "", err
		}

		if data.MatchMask(fanotify.FAN_MODIFY) {
			return fmt.Sprintf("PID:%d %s", data.GetPID(), path), nil
		}

		if data.MatchMask(fanotify.FAN_CLOSE_WRITE) {
			return fmt.Sprintf("PID:%d %s", data.GetPID(), path), nil
		}

		return "", fmt.Errorf("fanotify: unknown event")
	}

	for {
		if str, err := f(notify); err == nil && len(str) > 0 {
			fmt.Printf("%s\n", str)
		}
	}
}
