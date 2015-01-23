//
// AUTO-GENERATED by metago. DO NOT EDIT!
//

package test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/myUser/myPkg"
	"time"
)

type ExtendedObject struct {
	ByteField         byte
	U8Field           uint8
	U16Field          uint16
	U32Field          uint32
	U64Field          uint64
	S8Field           int8
	S16Field          int16
	S32Field          int32
	S64Field          int64
	StringField       string
	TimeField         time.Time
	ExtendedByteField byte
	StructField       myPkg.MyStruct
}

func (this *ExtendedObject) Dump() string {
	this.ExtendedByteField = 4
	return this.Dump_super()
}

func (this *ExtendedObject) ConditionalDump(t bool) string {
	return this.ConditionalDump_super(!t)
}

// from: BasicAttrTypesObject
func (this *ExtendedObject) Dump_super() (string, error) {
	return spew.Sdump(*this)
}

// from: BasicAttrTypesObject
func (this *ExtendedObject) ConditionalDump_super(t bool) string {
	if t {
		return this.Dump()
	}
	return ""
}

func (o1 *ExtendedObject) Equals(other interface{}) bool {
	switch o2 := other.(type) {
	case *ExtendedObject:
		return o1.equals(o2)
	case ExtendedObject:
		return o1.equals(&o2)
	}
	return false
}

func (o1 *ExtendedObject) equals(o2 *ExtendedObject) bool {

	//---------  comparison for ByteField ----------------------------------/
	if o1.ByteField != o2.ByteField {
		return false
	}

	//---------  comparison for U8Field ----------------------------------/
	if o1.U8Field != o2.U8Field {
		return false
	}

	//---------  comparison for U16Field ----------------------------------/
	if o1.U16Field != o2.U16Field {
		return false
	}

	//---------  comparison for U32Field ----------------------------------/
	if o1.U32Field != o2.U32Field {
		return false
	}

	//---------  comparison for U64Field ----------------------------------/
	if o1.U64Field != o2.U64Field {
		return false
	}

	//---------  comparison for S8Field ----------------------------------/
	if o1.S8Field != o2.S8Field {
		return false
	}

	//---------  comparison for S16Field ----------------------------------/
	if o1.S16Field != o2.S16Field {
		return false
	}

	//---------  comparison for S32Field ----------------------------------/
	if o1.S32Field != o2.S32Field {
		return false
	}

	//---------  comparison for S64Field ----------------------------------/
	if o1.S64Field != o2.S64Field {
		return false
	}

	//---------  comparison for StringField ----------------------------------/
	if o1.StringField != o2.StringField {
		return false
	}

	//---------  comparison for TimeField ----------------------------------/
	if !o1.TimeField.Equal(o2.TimeField) {
		return false
	}

	//---------  comparison for ExtendedByteField ----------------------------------/
	if o1.ExtendedByteField != o2.ExtendedByteField {
		return false
	}

	//---------  comparison for StructField ----------------------------------/
	if o1.StructField != o2.StructField {
		return false
	}
	return true
}