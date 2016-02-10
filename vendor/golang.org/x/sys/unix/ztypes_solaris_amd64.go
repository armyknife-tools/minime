// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs types_solaris.go

package unix

const (
	sizeofPtr      = 0x8
	sizeofShort    = 0x2
	sizeofInt      = 0x4
	sizeofLong     = 0x8
	sizeofLongLong = 0x8
)

type (
	_C_short     int16
	_C_int       int32
	_C_long      int64
	_C_long_long int64
)

type Timespec struct {
	Sec  int64
	Nsec int64
}

type Timeval struct {
	Sec  int64
	Usec int64
}

type Timeval32 struct {
	Sec  int32
	Usec int32
}

type Rusage struct {
	Utime    Timeval
	Stime    Timeval
	Maxrss   int64
	Ixrss    int64
	Idrss    int64
	Isrss    int64
	Minflt   int64
	Majflt   int64
	Nswap    int64
	Inblock  int64
	Oublock  int64
	Msgsnd   int64
	Msgrcv   int64
	Nsignals int64
	Nvcsw    int64
	Nivcsw   int64
}

type Rlimit struct {
	Cur uint64
	Max uint64
}

type _Gid_t uint32

const (
	S_IFMT   = 0xf000
	S_IFIFO  = 0x1000
	S_IFCHR  = 0x2000
	S_IFDIR  = 0x4000
	S_IFBLK  = 0x6000
	S_IFREG  = 0x8000
	S_IFLNK  = 0xa000
	S_IFSOCK = 0xc000
	S_ISUID  = 0x800
	S_ISGID  = 0x400
	S_ISVTX  = 0x200
	S_IRUSR  = 0x100
	S_IWUSR  = 0x80
	S_IXUSR  = 0x40
)

type Stat_t struct {
	Dev       uint64
	Ino       uint64
	Mode      uint32
	Nlink     uint32
	Uid       uint32
	Gid       uint32
	Rdev      uint64
	Size      int64
	Atim      Timespec
	Mtim      Timespec
	Ctim      Timespec
	Blksize   int32
	Pad_cgo_0 [4]byte
	Blocks    int64
	Fstype    [16]int8
}

type Flock_t struct {
	Type      int16
	Whence    int16
	Pad_cgo_0 [4]byte
	Start     int64
	Len       int64
	Sysid     int32
	Pid       int32
	Pad       [4]int64
}

type Dirent struct {
	Ino       uint64
	Off       int64
	Reclen    uint16
	Name      [1]int8
	Pad_cgo_0 [5]byte
}

type RawSockaddrInet4 struct {
	Family uint16
	Port   uint16
	Addr   [4]byte /* in_addr */
	Zero   [8]int8
}

type RawSockaddrInet6 struct {
	Family         uint16
	Port           uint16
	Flowinfo       uint32
	Addr           [16]byte /* in6_addr */
	Scope_id       uint32
	X__sin6_src_id uint32
}

type RawSockaddrUnix struct {
	Family uint16
	Path   [108]int8
}

type RawSockaddrDatalink struct {
	Family uint16
	Index  uint16
	Type   uint8
	Nlen   uint8
	Alen   uint8
	Slen   uint8
	Data   [244]int8
}

type RawSockaddr struct {
	Family uint16
	Data   [14]int8
}

type RawSockaddrAny struct {
	Addr RawSockaddr
	Pad  [236]int8
}

type _Socklen uint32

type Linger struct {
	Onoff  int32
	Linger int32
}

type Iovec struct {
	Base *int8
	Len  uint64
}

type IPMreq struct {
	Multiaddr [4]byte /* in_addr */
	Interface [4]byte /* in_addr */
}

type IPv6Mreq struct {
	Multiaddr [16]byte /* in6_addr */
	Interface uint32
}

type Msghdr struct {
	Name         *byte
	Namelen      uint32
	Pad_cgo_0    [4]byte
	Iov          *Iovec
	Iovlen       int32
	Pad_cgo_1    [4]byte
	Accrights    *int8
	Accrightslen int32
	Pad_cgo_2    [4]byte
}

type Cmsghdr struct {
	Len   uint32
	Level int32
	Type  int32
}

type Inet6Pktinfo struct {
	Addr    [16]byte /* in6_addr */
	Ifindex uint32
}

type IPv6MTUInfo struct {
	Addr RawSockaddrInet6
	Mtu  uint32
}

type ICMPv6Filter struct {
	X__icmp6_filt [8]uint32
}

const (
	SizeofSockaddrInet4    = 0x10
	SizeofSockaddrInet6    = 0x20
	SizeofSockaddrAny      = 0xfc
	SizeofSockaddrUnix     = 0x6e
	SizeofSockaddrDatalink = 0xfc
	SizeofLinger           = 0x8
	SizeofIPMreq           = 0x8
	SizeofIPv6Mreq         = 0x14
	SizeofMsghdr           = 0x30
	SizeofCmsghdr          = 0xc
	SizeofInet6Pktinfo     = 0x14
	SizeofIPv6MTUInfo      = 0x24
	SizeofICMPv6Filter     = 0x20
)

type FdSet struct {
	Bits [1024]int64
}

const (
	SizeofIfMsghdr  = 0x54
	SizeofIfData    = 0x44
	SizeofIfaMsghdr = 0x14
	SizeofRtMsghdr  = 0x4c
	SizeofRtMetrics = 0x28
)

type IfMsghdr struct {
	Msglen    uint16
	Version   uint8
	Type      uint8
	Addrs     int32
	Flags     int32
	Index     uint16
	Pad_cgo_0 [2]byte
	Data      IfData
}

type IfData struct {
	Type       uint8
	Addrlen    uint8
	Hdrlen     uint8
	Pad_cgo_0  [1]byte
	Mtu        uint32
	Metric     uint32
	Baudrate   uint32
	Ipackets   uint32
	Ierrors    uint32
	Opackets   uint32
	Oerrors    uint32
	Collisions uint32
	Ibytes     uint32
	Obytes     uint32
	Imcasts    uint32
	Omcasts    uint32
	Iqdrops    uint32
	Noproto    uint32
	Lastchange Timeval32
}

type IfaMsghdr struct {
	Msglen    uint16
	Version   uint8
	Type      uint8
	Addrs     int32
	Flags     int32
	Index     uint16
	Pad_cgo_0 [2]byte
	Metric    int32
}

type RtMsghdr struct {
	Msglen    uint16
	Version   uint8
	Type      uint8
	Index     uint16
	Pad_cgo_0 [2]byte
	Flags     int32
	Addrs     int32
	Pid       int32
	Seq       int32
	Errno     int32
	Use       int32
	Inits     uint32
	Rmx       RtMetrics
}

type RtMetrics struct {
	Locks    uint32
	Mtu      uint32
	Hopcount uint32
	Expire   uint32
	Recvpipe uint32
	Sendpipe uint32
	Ssthresh uint32
	Rtt      uint32
	Rttvar   uint32
	Pksent   uint32
}

const (
	SizeofBpfVersion = 0x4
	SizeofBpfStat    = 0x80
	SizeofBpfProgram = 0x10
	SizeofBpfInsn    = 0x8
	SizeofBpfHdr     = 0x14
)

type BpfVersion struct {
	Major uint16
	Minor uint16
}

type BpfStat struct {
	Recv    uint64
	Drop    uint64
	Capt    uint64
	Padding [13]uint64
}

type BpfProgram struct {
	Len       uint32
	Pad_cgo_0 [4]byte
	Insns     *BpfInsn
}

type BpfInsn struct {
	Code uint16
	Jt   uint8
	Jf   uint8
	K    uint32
}

type BpfTimeval struct {
	Sec  int32
	Usec int32
}

type BpfHdr struct {
	Tstamp    BpfTimeval
	Caplen    uint32
	Datalen   uint32
	Hdrlen    uint16
	Pad_cgo_0 [2]byte
}

type Termios struct {
	Iflag     uint32
	Oflag     uint32
	Cflag     uint32
	Lflag     uint32
	Cc        [19]uint8
	Pad_cgo_0 [1]byte
}
