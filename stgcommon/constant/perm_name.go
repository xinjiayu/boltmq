package constant

import (
	"strings"
)

// PERM 读写权限
// @author gaoyanlei
// @since 2017/8/9
const (
	PERM_PRIORITY = 0x1 << 3
	PERM_READ     = 0x1 << 2
	PERM_WRITE    = 0x1 << 1
	PERM_INHERIT  = 0x1 << 0
)

func isReadable(perm int) bool {
	return (perm & PERM_READ) == PERM_READ
}

func isWriteable(perm int) bool {
	return (perm & PERM_WRITE) == PERM_WRITE
}

func isInherited(perm int) bool {
	return (perm & PERM_INHERIT) == PERM_INHERIT
}

func Perm2String(perm int) string {
	str := "---"
	if isReadable(perm) {
		str = strings.Replace(str, "---", "R--", -1)
	}

	if isWriteable(perm) {
		str = strings.Replace(str, "---", "-W-", -1)
	}

	if isInherited(perm) {
		str = strings.Replace(str, "---", "--X", -1)
	}
	return str
}
