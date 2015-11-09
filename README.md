# gogengit
golang generate for git-based versioning

## Usage
Try `//go:generate gogengit --production=myver.go` in one of your files and then
run `go generate`. Don't like 'git' versioning and wanted beautiful hand-maintained
symbolic (or other) versioning? Try `gogengit --production=myver.go --version-file=VMINE`
and writing something symbolic to the VMINE file before running `go generate`.

