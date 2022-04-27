# WFM TODO

## Interfaces
* WebDAV server
* FastCGI Interface
* Docker support

## Security
* userless/guest read-only mode, user rw
  requires custom login window
* two factor auth
  requires custom login window
* garbage collect old f2b entries
* f2b ddos prevention, sleep on too many bans?
* qps throttling
* download/upload throttling

## Layout / UI
* add flag to specify own favicon.ico
* top bar too long on mobile/small screen
* custom html login window - needed for two factor auth?
* editable and non editable documents by extension, also for git checkins
* thumbnail / icon view for pictures (cache thumbnails on server?)
* glob filter (*.*) in dir view
* errors in dialog boxes instead of plain text
* html as template

## File IO
* file search function
* path prefix, required for docker
* path prefix per user
* udf iso format https://github.com/mogaika/udf
* zip/unzip archives
* iso files recursive list
* zipped iso like .iso.gz, .iso.xz, .iso.lz
* auto unpack via mime type...
* add more formats like tgz/txz, etc
* du with xdev as a go routine
* git client https://github.com/go-git/go-git
* file locking https://github.com/gofrs/flock
* support for different filesystems, S3, SMB, archive files as io/fs
* archive files in main view / graphical/table form
