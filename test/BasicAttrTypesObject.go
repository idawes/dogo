//
// AUTO-GENERATED by metago. DO NOT EDIT!
//

package test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/idawes/metago"
	"time"
)

type BasicAttrTypesObject struct {
	ByteField   byte
	UbbField    uint
	U8Field     uint8
	U16Field    uint16
	U32Field    uint32
	U64Field    uint64
	SField      int
	S8Field     int8
	S16Field    int16
	S32Field    int32
	S64Field    int64
	F32Field    float32
	F64Field    float64
	StringField string
	TimeField   time.Time
}

func (this *BasicAttrTypesObject) ConditionalDump(t bool) string {
	if t {
		return this.Dump()
	}
	return ""
}

func (this *BasicAttrTypesObject) Dump() string {
	return spew.Sdump(*this)
}

func (o1 *BasicAttrTypesObject) Equals(o2 *BasicAttrTypesObject) bool {

	{
		va, vb := o1.ByteField, o2.ByteField
		if va != vb {
			return false
		}
	}

	{
		va, vb := o1.UbbField, o2.UbbField
		if va != vb {
			return false
		}
	}

	{
		va, vb := o1.U8Field, o2.U8Field
		if va != vb {
			return false
		}
	}

	{
		va, vb := o1.U16Field, o2.U16Field
		if va != vb {
			return false
		}
	}

	{
		va, vb := o1.U32Field, o2.U32Field
		if va != vb {
			return false
		}
	}

	{
		va, vb := o1.U64Field, o2.U64Field
		if va != vb {
			return false
		}
	}

	{
		va, vb := o1.SField, o2.SField
		if va != vb {
			return false
		}
	}

	{
		va, vb := o1.S8Field, o2.S8Field
		if va != vb {
			return false
		}
	}

	{
		va, vb := o1.S16Field, o2.S16Field
		if va != vb {
			return false
		}
	}

	{
		va, vb := o1.S32Field, o2.S32Field
		if va != vb {
			return false
		}
	}

	{
		va, vb := o1.S64Field, o2.S64Field
		if va != vb {
			return false
		}
	}

	{
		va, vb := o1.F32Field, o2.F32Field
		if va != vb {
			return false
		}
	}

	{
		va, vb := o1.F64Field, o2.F64Field
		if va != vb {
			return false
		}
	}

	{
		va, vb := o1.StringField, o2.StringField
		if va != vb {
			return false
		}
	}

	{
		va, vb := o1.TimeField, o2.TimeField
		if va.Equal(vb) {
			return false
		}
	}
	return true
}

func (o1 *BasicAttrTypesObject) Diff(o2 *BasicAttrTypesObject) *metago.Diff {
	chgs := make([]metago.Chg, 0)

	{
		va, vb := o1.ByteField, o2.ByteField
		if va != vb {
			chgs = append(chgs, metago.NewByteChg(&BasicAttrTypesObjectByteFieldSREF, vb, va))
		}
	}

	{
		va, vb := o1.UbbField, o2.UbbField
		if va != vb {
			chgs = append(chgs, metago.NewUintChg(&BasicAttrTypesObjectUbbFieldSREF, vb, va))
		}
	}

	{
		va, vb := o1.U8Field, o2.U8Field
		if va != vb {
			chgs = append(chgs, metago.NewUint8Chg(&BasicAttrTypesObjectU8FieldSREF, vb, va))
		}
	}

	{
		va, vb := o1.U16Field, o2.U16Field
		if va != vb {
			chgs = append(chgs, metago.NewUint16Chg(&BasicAttrTypesObjectU16FieldSREF, vb, va))
		}
	}

	{
		va, vb := o1.U32Field, o2.U32Field
		if va != vb {
			chgs = append(chgs, metago.NewUint32Chg(&BasicAttrTypesObjectU32FieldSREF, vb, va))
		}
	}

	{
		va, vb := o1.U64Field, o2.U64Field
		if va != vb {
			chgs = append(chgs, metago.NewUint64Chg(&BasicAttrTypesObjectU64FieldSREF, vb, va))
		}
	}

	{
		va, vb := o1.SField, o2.SField
		if va != vb {
			chgs = append(chgs, metago.NewIntChg(&BasicAttrTypesObjectSFieldSREF, vb, va))
		}
	}

	{
		va, vb := o1.S8Field, o2.S8Field
		if va != vb {
			chgs = append(chgs, metago.NewInt8Chg(&BasicAttrTypesObjectS8FieldSREF, vb, va))
		}
	}

	{
		va, vb := o1.S16Field, o2.S16Field
		if va != vb {
			chgs = append(chgs, metago.NewInt16Chg(&BasicAttrTypesObjectS16FieldSREF, vb, va))
		}
	}

	{
		va, vb := o1.S32Field, o2.S32Field
		if va != vb {
			chgs = append(chgs, metago.NewInt32Chg(&BasicAttrTypesObjectS32FieldSREF, vb, va))
		}
	}

	{
		va, vb := o1.S64Field, o2.S64Field
		if va != vb {
			chgs = append(chgs, metago.NewInt64Chg(&BasicAttrTypesObjectS64FieldSREF, vb, va))
		}
	}

	{
		va, vb := o1.F32Field, o2.F32Field
		if va != vb {
			chgs = append(chgs, metago.NewFloat32Chg(&BasicAttrTypesObjectF32FieldSREF, vb, va))
		}
	}

	{
		va, vb := o1.F64Field, o2.F64Field
		if va != vb {
			chgs = append(chgs, metago.NewFloat64Chg(&BasicAttrTypesObjectF64FieldSREF, vb, va))
		}
	}

	{
		va, vb := o1.StringField, o2.StringField
		if va != vb {
			chgs = append(chgs, metago.NewStringChg(&BasicAttrTypesObjectStringFieldSREF, vb, va))
		}
	}

	{
		va, vb := o1.TimeField, o2.TimeField
		if va.Equal(vb) {
			chgs = append(chgs, metago.NewTimeChg(&BasicAttrTypesObjectTimeFieldSREF, vb, va))
		}
	}
	return &metago.Diff{Chgs: chgs}
}

func (o *BasicAttrTypesObject) Apply(d *metago.Diff) error {
	for _, c := range d.Chgs {
		switch c.AttributeID() {

		case &BasicAttrTypesObjectByteFieldAID:
			{
				v := &o.ByteField
				*v = c.(*metago.ByteChg).NewValue
			}

		case &BasicAttrTypesObjectUbbFieldAID:
			{
				v := &o.UbbField
				*v = c.(*metago.UintChg).NewValue
			}

		case &BasicAttrTypesObjectU8FieldAID:
			{
				v := &o.U8Field
				*v = c.(*metago.Uint8Chg).NewValue
			}

		case &BasicAttrTypesObjectU16FieldAID:
			{
				v := &o.U16Field
				*v = c.(*metago.Uint16Chg).NewValue
			}

		case &BasicAttrTypesObjectU32FieldAID:
			{
				v := &o.U32Field
				*v = c.(*metago.Uint32Chg).NewValue
			}

		case &BasicAttrTypesObjectU64FieldAID:
			{
				v := &o.U64Field
				*v = c.(*metago.Uint64Chg).NewValue
			}

		case &BasicAttrTypesObjectSFieldAID:
			{
				v := &o.SField
				*v = c.(*metago.IntChg).NewValue
			}

		case &BasicAttrTypesObjectS8FieldAID:
			{
				v := &o.S8Field
				*v = c.(*metago.Int8Chg).NewValue
			}

		case &BasicAttrTypesObjectS16FieldAID:
			{
				v := &o.S16Field
				*v = c.(*metago.Int16Chg).NewValue
			}

		case &BasicAttrTypesObjectS32FieldAID:
			{
				v := &o.S32Field
				*v = c.(*metago.Int32Chg).NewValue
			}

		case &BasicAttrTypesObjectS64FieldAID:
			{
				v := &o.S64Field
				*v = c.(*metago.Int64Chg).NewValue
			}

		case &BasicAttrTypesObjectF32FieldAID:
			{
				v := &o.F32Field
				*v = c.(*metago.Float32Chg).NewValue
			}

		case &BasicAttrTypesObjectF64FieldAID:
			{
				v := &o.F64Field
				*v = c.(*metago.Float64Chg).NewValue
			}

		case &BasicAttrTypesObjectStringFieldAID:
			{
				v := &o.StringField
				*v = c.(*metago.StringChg).NewValue
			}

		}
	}
	return nil
}
