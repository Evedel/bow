##### Go VirtualBox vboxsf sendfile bug workaround

If you serve static content from a shared folder you might have run into a
vboxsf file corruption bug. This hack disables the sendfile syscall for
the go process which will force the standard library to fallback to userland
buffered IO.

References:  
[Ticket #9069 shared folder doesn't seem to update](https://www.virtualbox.org/ticket/9069)  
[net: Add ability to disable sendfile](https://github.com/golang/go/issues/9694)

##### Usage

Save [disable_sendfile_vbox_linux.go](disable_sendfile_vbox_linux.go) to somewhere in your go project.
Or do
```go
import (
	_ "github.com/wader/disable_sendfile_vbox_linux"
}
```
in a source file.

##### License

Public domain. Your free to do whatever you want.
