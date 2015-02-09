//
// AUTO-GENERATED by metago. DO NOT EDIT!
//

package test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/idawes/metago"
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

func (this *ExtendedObject) ConditionalDump(t bool) string {
	return this.ConditionalDump_super(!t)
}

// from: BasicAttrTypesObject
func (this *ExtendedObject) ConditionalDump_super(t bool) string {
	if t {
		return this.Dump()
	}
	return ""
}

func (this *ExtendedObject) Dump() string {
	this.ExtendedByteField = 4
	return this.Dump_super()
}

// from: BasicAttrTypesObject
func (this *ExtendedObject) Dump_super() (string, error) {
	return spew.Sdump(*this)
}

func (o1 *ExtendedObject) Equals(o2 *ExtendedObject) bool {

	if o1.ByteField != o2.ByteField {
		return false
	}

	if o1.U8Field != o2.U8Field {
		return false
	}

	if o1.U16Field != o2.U16Field {
		return false
	}

	if o1.U32Field != o2.U32Field {
		return false
	}

	if o1.U64Field != o2.U64Field {
		return false
	}

	if o1.S8Field != o2.S8Field {
		return false
	}

	if o1.S16Field != o2.S16Field {
		return false
	}

	if o1.S32Field != o2.S32Field {
		return false
	}

	if o1.S64Field != o2.S64Field {
		return false
	}

	if o1.StringField != o2.StringField {
		return false
	}

	if !o1.TimeField.Equal(o2.TimeField) {
		return false
	}

	if o1.ExtendedByteField != o2.ExtendedByteField {
		return false
	}

	if o1.StructField != o2.StructField {
		return false
	}
	return true
}

func (o1 *ExtendedObject) Diff(o2 *ExtendedObject) bool {

	if o1.ByteField != o2.ByteField {
		d.Add(NewByteDiff(AID_BasicAttrTypesObject_ByteField, true, o1.ByteField, o2.ByteField))
	}

	if o1.U8Field != o2.U8Field {
		d.Add(NewUint8Diff(AID_BasicAttrTypesObject_U8Field, true, o1.U8Field, o2.U8Field))
	}

	if o1.U16Field != o2.U16Field {
		d.Add(NewUint16Diff(AID_BasicAttrTypesObject_U16Field, true, o1.U16Field, o2.U16Field))
	}

	if o1.U32Field != o2.U32Field {
		d.Add(NewUint32Diff(AID_BasicAttrTypesObject_U32Field, true, o1.U32Field, o2.U32Field))
	}

	if o1.U64Field != o2.U64Field {
		d.Add(NewUint64Diff(AID_BasicAttrTypesObject_U64Field, true, o1.U64Field, o2.U64Field))
	}

	if o1.S8Field != o2.S8Field {
		d.Add(NewInt8Diff(AID_BasicAttrTypesObject_S8Field, true, o1.S8Field, o2.S8Field))
	}

	if o1.S16Field != o2.S16Field {
		d.Add(NewInt16Diff(AID_BasicAttrTypesObject_S16Field, true, o1.S16Field, o2.S16Field))
	}

	if o1.S32Field != o2.S32Field {
		d.Add(NewInt32Diff(AID_BasicAttrTypesObject_S32Field, true, o1.S32Field, o2.S32Field))
	}

	if o1.S64Field != o2.S64Field {
		d.Add(NewInt64Diff(AID_BasicAttrTypesObject_S64Field, true, o1.S64Field, o2.S64Field))
	}

	if o1.StringField != o2.StringField {
		d.Add(NewStringDiff(AID_BasicAttrTypesObject_StringField, true, o1.StringField, o2.StringField))
	}

	if !o1.TimeField.Equal(o2.TimeField) {
		return false
	}

	if o1.ExtendedByteField != o2.ExtendedByteField {
		d.Add(NewByteDiff(AID_ExtendedObject_ExtendedByteField, true, o1.ExtendedByteField, o2.ExtendedByteField))
	}

	if o1.StructField != o2.StructField {
		d.Add(NewMyPkg.MyStructDiff(AID_ExtendedObject_StructField, true, o1.StructField, o2.StructField))
	}
	return true
}
