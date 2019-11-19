module github.com/s3rj1k/go-fanotify/cmd

go 1.13

replace github.com/s3rj1k/go-fanotify/fanotify => ../fanotify

require (
	github.com/s3rj1k/go-fanotify/fanotify v0.0.0-20191118194719-aa354f566745
	golang.org/x/sys v0.0.0-20191119060738-e882bf8e40c2
)
