// VirtualBox vboxsf sendfile bug workaround
// use seccomp to make sendfile return ENOTSUP which makes go http fallback to io.Copy
//
// Licensed as public domain
// mattias.wadman@gmail.com

package disable_sendfile_vbox_linux

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

// https://github.com/torvalds/linux/blob/master/include/uapi/linux/filter.h
type sockFilter struct {
	code uint16
	jt   uint8
	jf   uint8
	k    uint32
}
type sockFprog struct {
	len  uint16
	filt []sockFilter
}

func init() {
	// only disable if running in virtualbox
	// TODO: some better way of detecting?
	mounts, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		return
	}
	if !strings.Contains(string(mounts), "vboxsf") {
		return
	}

	// https://github.com/torvalds/linux/blob/master/include/uapi/linux/seccomp.h
	const (
		seccompSetModeFilter          = 0x1        // set bpf filter mode
		seccompFilterFlagTsync        = 0x1        // same filter for all threads
		seccompRetErrno        uint32 = 0x00050000 // return errno value
		seccompRetAllow        uint32 = 0x7fff0000 // allow syscall
	)

	// syscall number for seccomp
	// https://github.com/torvalds/linux/tree/master/arch/x86/entry/syscalls
	var sysSeccomp uintptr
	if runtime.GOARCH == "amd64" {
		sysSeccomp = 317
	} else {
		sysSeccomp = 354
	}

	// for some reason not defined for some archs in syscall package
	// https://github.com/torvalds/linux/blob/master/include/uapi/linux/prctl.h
	const PR_SET_NO_NEW_PRIVS = 38

	// set no new privs so that we can install seccomp filter as non-root
	// https://github.com/torvalds/linux/blob/master/Documentation/userspace-api/no_new_privs.rst
	if r1, r2, errno :=
		syscall.Syscall(
			uintptr(syscall.SYS_PRCTL),
			uintptr(PR_SET_NO_NEW_PRIVS),
			uintptr(1),
			uintptr(0)); errno != 0 {
		fmt.Printf("WARNING: prctl PR_SET_NO_NEW_PRIVS FAILED r1=%d r2=%d %v\n", r1, r2, errno)
		return
	}

	/*
		BPF runs on this data:
		struct seccomp_data {
		  int nr;
		  __u32 arch;
		  __u64 instruction_pointer;
		  __u64 args[6];
		};
	*/
	filter := &sockFprog{
		len: 4,
		filt: []sockFilter{
			{syscall.BPF_LD + syscall.BPF_ABS + syscall.BPF_W, 0, 0, 0},                        // load syscall (nr at offset 0)
			{syscall.BPF_JMP + syscall.BPF_JEQ + syscall.BPF_K, 0, 1, syscall.SYS_SENDFILE},    // if sendfile
			{syscall.BPF_RET + syscall.BPF_K, 0, 0, seccompRetErrno | uint32(syscall.ENOTSUP)}, // true: ENOTSUP
			{syscall.BPF_RET + syscall.BPF_K, 0, 0, seccompRetAllow},                           // false: ALLOW
		},
	}

	if r1, r2, errno :=
		syscall.Syscall(
			sysSeccomp,
			uintptr(seccompSetModeFilter),
			uintptr(seccompFilterFlagTsync),
			uintptr(unsafe.Pointer(filter))); errno == 0 {
		fmt.Printf("WARNING: sendfile disabled\n")
	} else {
		fmt.Printf("WARNING: disable sendfile FAILED r1=%d r2=%d %v\n", r1, r2, errno)
	}
}
