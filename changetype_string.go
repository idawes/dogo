// generated by stringer -type=ChangeType; DO NOT EDIT

package metago

import "fmt"

const _ChangeType_name = "ChangeTypeInsertChangeTypeDeleteChangeTypeModify"

var _ChangeType_index = [...]uint8{16, 32, 48}

func (i ChangeType) String() string {
	if i < 0 || i >= ChangeType(len(_ChangeType_index)) {
		return fmt.Sprintf("ChangeType(%d)", i)
	}
	hi := _ChangeType_index[i]
	lo := uint8(0)
	if i > 0 {
		lo = _ChangeType_index[i-1]
	}
	return _ChangeType_name[lo:hi]
}
