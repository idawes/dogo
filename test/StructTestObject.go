//
// AUTO-GENERATED by metago. DO NOT EDIT!
//

package test

import (
	"github.com/idawes/metago"
)

type StructTestObject struct {
	B  BasicAttrTypesObject
	MB map[int]BasicAttrTypesObject
}

func (o1 *StructTestObject) Equals(o2 *StructTestObject) bool {

	{
		va, vb := o1.B, o2.B
		if va.Equals(&vb) {
			return false
		}
	}

	{
		va, vb := o1.MB, o2.MB
		if len(va) != len(vb) {
			return false
		}
		for key, va1 := range va {
			if vb1, ok := vb[key]; ok {
				if va1.Equals(&vb1) {
					return false
				}
			} else {
				return false // didn't find key in vb
			}
		}
	}
	return true
}

func (o1 *StructTestObject) Diff(o2 *StructTestObject) *metago.Diff {
	d := &metago.Diff{}

	{
		va, vb := o1.B, o2.B
		d.Chgs = append(d.Chgs, metago.NewStructChg(&StructTestObjectBSREF, va.Diff(&vb)))
	}

	{
		va, vb := o1.MB, o2.MB
		for key, va1 := range va {
			if vb1, ok := vb[key]; ok {
				d1 := &metago.Diff{}
				d1.Chgs = append(d1.Chgs, metago.NewStructChg(&StructTestObjectMBSREF, va1.Diff(&vb1)))
				if len(d1.Chgs) != 0 {
					d.Chgs = append(d.Chgs, metago.NewIntMapChg(&StructTestObjectMBSREF, key, metago.ChangeTypeModify, d1))
				}
			} else {
				d1 := &metago.Diff{}
				t := BasicAttrTypesObject{}
				d1.Add(metago.NewStructChg(&StructTestObjectMBSREF, t.Diff(&va1)))
				if len(d1.Chgs) != 0 {
					d.Chgs = append(d.Chgs, metago.NewIntMapChg(&StructTestObjectMBSREF, key, metago.ChangeTypeInsert, d1))
				}
			}
		}
		for key, vb1 := range vb {
			d1 := &metago.Diff{}
			t := BasicAttrTypesObject{}
			d1.Add(metago.NewStructChg(&StructTestObjectMBSREF, vb1.Diff(&t)))
			if len(d1.Chgs) != 0 {
				d.Chgs = append(d.Chgs, metago.NewIntMapChg(&StructTestObjectMBSREF, key, metago.ChangeTypeDelete, d1))
			}
		}
	}
	return d
}

func (o *StructTestObject) Apply(d *metago.Diff) error {
	for _, c := range d.Chgs {
		switch c.AttributeID() {

		case &StructTestObjectMBAID:
			{
				mc := c.(*metago.IntMapChg)
				switch mc.Typ {
				case metago.ChangeTypeModify:
					o1 := o.MB[mc.Key]
					o1.Apply(&mc.Chgs)
				case metago.ChangeTypeInsert:
				case metago.ChangeTypeDelete:
				}
			}
		}
	}
	return nil
}
