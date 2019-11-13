module github.com/s3rj1k/go-fanotify/cmd

go 1.13

replace github.com/s3rj1k/go-fanotify/fanotify => ../fanotify

require (
	github.com/s3rj1k/go-fanotify/fanotify v0.0.0-20191115105227-52c3f25bbad1
	golang.org/x/sys v0.0.0-20191113165036-4c7a9d0fe056
)
