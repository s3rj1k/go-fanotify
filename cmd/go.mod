module github.com/s3rj1k/go-fanotify/cmd

go 1.13

replace github.com/s3rj1k/go-fanotify/fanotify => ../fanotify

require (
	github.com/s3rj1k/go-fanotify/fanotify v0.0.0-20191119131115-7a5c7b812c5c
	golang.org/x/sys v0.0.0-20191120155948-bd437916bb0e
)
