// Copyright 2015 Ian Dawes. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strconv"
	"strings"
)

/*
	attribute definition:
	<id> <name> <type> [<persistenceType>]
*/

func parseAttribute(t *typedef, r *reader) (attrDef, error) {
	fields, err := r.read()
	if len(fields) < 3 || len(fields) > 4 {
		return nil, fmt.Errorf("invalid attribute specification, line %d of file %s", r.line, r.f.Name())
	}
	id, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("expecting an integer attribute id, found \"%s\", line %d of file %s", fields[0], r.line, r.f.Name())
	}
	name := fields[1]
	p := persistenceClassPersistent
	if len(fields) > 3 {
		switch fields[3] {
		case "persistent":
			p = persistenceClassPersistent
		case "non-persistent":
			p = persistenceClassNonPersistent
		case "ephemeral":
			p = persistenceClassEphemeral
		default:
			return nil, fmt.Errorf("unrecognized persistence type \"%s\", line %d of file %s", fields[3], r.line, r.f.Name())
		}
	}
	return newAttrDef(&baseAttrDef{parentType: t, attrID: id, nm: name, typ: fields[2], persistence: p, sline: r.line, sfile: r.f.Name()})
}

func newAttrDef(b *baseAttrDef) (attrDef, error) {
	c, err := getClass(b.typ)
	if err != nil {
		return nil, fmt.Errorf("unknown attribute class %s, line %d of file %s", b.typ, b.sline, b.sfile)
	}
	switch c {
	case attrClassBuiltin:
		return b, nil
	case attrClassTime:
		return &timeAttrDef{baseAttrDef: *b}, nil
	case attrClassMap:
		return newMapAttrDef(b)
	case attrClassSlice:
		return newSliceAttrDef(b)
	case attrClassStruct:
		return newStructAttrDef(b)
	}
	return nil, fmt.Errorf("Unknown attribute class %s, line %d of file %s", c, b.sline, b.sfile)
}

func getClass(n string) (attrClass, error) {
	switch n {
	case "byte", "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64", "float32", "float64", "string":
		return attrClassBuiltin, nil
	case "time.Time":
		return attrClassTime, nil
	}

	switch {
	case strings.HasPrefix(n, "[]"):
		return attrClassSlice, nil
	case strings.HasPrefix(n, "map["):
		return attrClassMap, nil
	default:
		return attrClassStruct, nil
	}
}

//go:generate stringer -type=attrClass
type attrClass int

const (
	attrClassBuiltin attrClass = iota
	attrClassTime
	attrClassSlice
	attrClassMap
	attrClassStruct
)

//go:generate stringer -type=persistenceClass
type persistenceClass int

const (
	persistenceClassPersistent    persistenceClass = iota // serialized to disk and wire
	persistenceClassNonPersistent                         // serialized to wire
	persistenceClassEphemeral                             // computed or temporary storage - not serialized
)

type attrDef interface {
	attributeID() int
	persistenceClass() persistenceClass
	srcline() int
	srcfile() string
	name() string
	typeName() string
	isMultiValued() bool
	generateEquals(w *writer, levelID string)
	generateDiff(w *writer, levelID string)
	generateIns(w *writer, levelID string)
	generateDel(w *writer, levelID string)
	generateApply(w *writer, levelID string)
	generateMapModify(w *writer, levelID string)
}

type baseAttrDef struct {
	parentType  *typedef
	attrID      int
	nm          string
	typ         string
	persistence persistenceClass
	sline       int
	sfile       string
}

func (a *baseAttrDef) attributeID() int {
	return a.attrID
}

func (a *baseAttrDef) persistenceClass() persistenceClass {
	return a.persistence
}

func (a *baseAttrDef) srcline() int {
	return a.sline
}

func (a *baseAttrDef) srcfile() string {
	return a.sfile
}

func (a *baseAttrDef) name() string {
	return a.nm
}

func (a *baseAttrDef) typeName() string {
	return a.typ
}

func (a *baseAttrDef) isMultiValued() bool {
	return false
}

func (a *baseAttrDef) checkLevel0Hdr(w *writer, levelID string) {
	if levelID == "" {
		w.printf("    {\n")
		w.printf("        va, vb := o1.%[1]s, o2.%[1]s\n", a.nm)
	}
}

func (a *baseAttrDef) checkLevel0Ftr(w *writer, levelID string) {
	if levelID == "" {
		w.printf("    }\n")
	}
}

func (a *baseAttrDef) generateEquals(w *writer, levelID string) {
	a.checkLevel0Hdr(w, levelID)
	w.printf("  if va%[1]s != vb%[1]s {\n    return false\n  }\n", levelID)
	a.checkLevel0Ftr(w, levelID)
}

func (a *baseAttrDef) generateDiff(w *writer, levelID string) {
	a.checkLevel0Hdr(w, levelID)
	format := `  if va%[1]s != vb%[1]s {
		chgs%[1]s = append(chgs%[1]s, metago.New%[2]sChg(&%[3]s%[4]sSREF, vb%[1]s, va%[1]s))
	}
`
	w.printf(format, levelID, strings.Title(a.typ), a.parentType.name, a.nm)
	a.checkLevel0Ftr(w, levelID)
}

func (a *baseAttrDef) generateIns(w *writer, levelID string) {
	w.printf("chgs%[1]s = append(chgs%[1]s, metago.New%[2]sChg(&%[3]s%[4]sSREF, va%[1]s))\n", levelID, strings.Title(a.typ), a.parentType.name, a.nm)
}

func (a *baseAttrDef) generateDel(w *writer, levelID string) {
	w.printf("chgs%[1]s = append(chgs%[1]s, metago.New%[2]sChg(&%[3]s%[4]sSREF, vb%[1]s))\n", levelID, strings.Title(a.typ), a.parentType.name, a.nm)
}

func (a *baseAttrDef) checkLevel0ApplyHdr(w *writer, levelID string) {
	if levelID == "" {
		w.printf("    {\n")
		w.printf("        v := &o.%[1]s\n", a.nm)
	}
}

func (a *baseAttrDef) checkLevel0ApplyFtr(w *writer, levelID string) {
	if levelID == "" {
		w.printf("    }\n")
	}
}
func (a *baseAttrDef) generateApply(w *writer, levelID string) {
	w.printf("    case &%[1]s%[2]sAID:\n", a.parentType.name, a.nm)
	a.checkLevel0ApplyHdr(w, levelID)
	w.printf("        *v%[1]s = c.(*metago.%[2]sChg).NewValue\n", levelID, strings.Title(a.typ))
	a.checkLevel0ApplyFtr(w, levelID)
}

func (a *baseAttrDef) generateMapModify(w *writer, levelID string) {
	nextLevelID := fmt.Sprintf("%s1", levelID)
	w.printf("              m%[1]s[key%[1]s] = mc%[1]s.Chgs[0].(*metago.%[3]sChg).NewValue\n", levelID, nextLevelID, strings.Title(a.typ))
}

/************************************************************************/
/************************** Time Attribute ******************************/
type timeAttrDef struct {
	baseAttrDef
}

func (a *timeAttrDef) unqualifiedTypeName() string {
	return "Time"
}

func (a *timeAttrDef) generateEquals(w *writer, levelID string) {
	a.checkLevel0Hdr(w, levelID)
	w.printf("  if va%[1]s.Equal(vb%[1]s) {\n    return false\n  }\n", levelID)
	a.checkLevel0Ftr(w, levelID)
}

func (a *timeAttrDef) generateDiff(w *writer, levelID string) {
	a.checkLevel0Hdr(w, levelID)
	format := `  if va%[1]s.Equal(vb%[1]s) {
		chgs%[1]s = append(chgs%[1]s, metago.NewTimeChg(&%[2]s%[3]sSREF, vb%[1]s, va%[1]s))
	}
`
	w.printf(format, levelID, a.parentType.name, a.nm)
	a.checkLevel0Ftr(w, levelID)
}

func (a *timeAttrDef) generateIns(w *writer, levelID string) {
	w.printf("chgs%[1]s = append(chgs%[1]s, metago.NewTimeChg(&%[2]s%[3]sSREF, va%[1]s))\n", levelID, a.parentType.name, a.nm)
}

func (a *timeAttrDef) generateDel(w *writer, levelID string) {
	w.printf("chgs%[1]s = append(chgs%[1]s, metago.NewTimeChg(&%[2]s%[3]sSREF, vb%[1]s))\n", levelID, a.parentType.name, a.nm)
}

func (a *timeAttrDef) generateApply(w *writer, levelID string) {

}

func (a *timeAttrDef) generateMapModify(w *writer, levelID string) {
	nextLevelID := fmt.Sprintf("%s1", levelID)
	w.printf("              m%[1]s[key%[1]s] = mc%[1]s.Chgs[0].(*metago.TimeChg).NewValue\n", levelID, nextLevelID, strings.Title(a.typ))
}

/************************************************************************/
/**************************** Slice Attribute ***************************/
type sliceAttrDef struct {
	baseAttrDef
	valType string
	valAttr attrDef
}

func newSliceAttrDef(b *baseAttrDef) (*sliceAttrDef, error) {
	valType := b.typ[2:]
	valAttr, err := newAttrDef(&baseAttrDef{parentType: b.parentType, nm: b.nm, typ: valType})
	if err != nil {
		return nil, fmt.Errorf("invalid slice attribute specification %s, line %d of file %s", b.typ, b.sline, b.sfile)
	}
	return &sliceAttrDef{baseAttrDef: *b, valType: valType, valAttr: valAttr}, nil
}

func (a *sliceAttrDef) isMultiValued() bool {
	return true
}

func (a *sliceAttrDef) generateEquals(w *writer, levelID string) {
	nextLevelID := fmt.Sprintf("%s1", levelID)
	a.checkLevel0Hdr(w, levelID)
	format := `    if len(va%[1]s) != len(vb%[1]s) {
        return false
    }
    for idx%[1]s, va%[2]s := range va%[1]s {
		vb%[2]s := vb%[1]s[idx%[1]s]
`
	w.printf(format, levelID, nextLevelID)
	a.valAttr.generateEquals(w, nextLevelID)
	w.printf("  }\n")
	a.checkLevel0Ftr(w, levelID)
}

func (a *sliceAttrDef) generateDiff(w *writer, levelID string) {
	nextLevelID := fmt.Sprintf("%s1", levelID)
	a.checkLevel0Hdr(w, levelID)
	format := `    for idx%[1]s, va%[2]s := range va%[1]s {
        if idx%[1]s  < len(vb%[1]s) {
			vb%[2]s := vb%[1]s[idx%[1]s]
			chgs%[2]s := make([]metago.Chg, 0)
`
	w.printf(format, levelID, nextLevelID)
	a.valAttr.generateDiff(w, nextLevelID)
	format = `            if len(chgs%[2]s) != 0 {
				chgs%[1]s = append(chgs%[1]s, metago.NewSliceChg(&%[3]s%[4]sSREF, idx%[1]s, metago.ChangeTypeModify, chgs%[2]s))
			}
		} else {
			chgs%[2]s := make([]metago.Chg, 0)
`
	w.printf(format, levelID, nextLevelID, a.parentType.name, a.nm)
	a.valAttr.generateIns(w, nextLevelID)
	format = `            if len(chgs%[2]s) != 0 {
				chgs%[1]s = append(chgs%[1]s, metago.NewSliceChg(&%[3]s%[4]sSREF, idx%[1]s, metago.ChangeTypeInsert, chgs%[2]s))
			}
		}
	}
	for idx%[1]s, vb%[2]s := range vb%[1]s {
		chgs%[2]s := make([]metago.Chg, 0)
`
	w.printf(format, levelID, nextLevelID, a.parentType.name, a.nm)
	a.valAttr.generateDel(w, nextLevelID)
	format = `            if len(chgs%[2]s) != 0 {
				chgs%[1]s = append(chgs%[1]s, metago.NewSliceChg(&%[3]s%[4]sSREF, idx%[1]s, metago.ChangeTypeDelete, chgs%[2]s))
        }
	}
`
	w.printf(format, levelID, nextLevelID, a.parentType.name, a.nm)
	a.checkLevel0Ftr(w, levelID)
}

func (a *sliceAttrDef) generateIns(w *writer, levelID string) {
	nextLevelID := fmt.Sprintf("%s1", levelID)
	format := `    for idx%[1]s, va%[2]s := range va%[1]s {
		chgs%[2]s := make([]metago.Chg, 0)
`
	w.printf(format, levelID, nextLevelID)
	a.valAttr.generateIns(w, nextLevelID)
	format = `            if len(chgs%[2]s) != 0 {
				chgs%[1]s = append(chgs%[1]s, metago.NewSliceChg(&%[3]s%[4]sSREF, idx%[1]s, metago.ChangeTypeInsert, chgs%[2]s))
		}
	}
`
	w.printf(format, levelID, nextLevelID, a.parentType.name, a.nm)
}

func (a *sliceAttrDef) generateDel(w *writer, levelID string) {
	nextLevelID := fmt.Sprintf("%s1", levelID)
	format := `    for idx%[1]s, vb%[2]s := range vb%[1]s {
		chgs%[2]s := make([]metago.Chg, 0)
`
	w.printf(format, levelID, nextLevelID)
	a.valAttr.generateDel(w, nextLevelID)
	format = `            if len(chgs%[2]s) != 0 {
				chgs%[1]s = append(chgs%[1]s, metago.NewSliceChg(&%[3]s%[4]sSREF, idx%[1]s, metago.ChangeTypeDelete, chgs%[2]s))
		}
	}
`
	w.printf(format, levelID, nextLevelID, a.parentType.name, a.nm)
}

func (a *sliceAttrDef) generateApply(w *writer, levelID string) {

}

/************************************************************************/
/**************************** Map Attribute *****************************/
type mapAttrDef struct {
	baseAttrDef
	keyType string
	keyAttr attrDef
	valType string
	valAttr attrDef
}

func newMapAttrDef(b *baseAttrDef) (*mapAttrDef, error) {
	// map[int]string
	i := strings.Index(b.typ, "[")
	j := strings.Index(b.typ, "]")
	if i == -1 || j == -1 {
		return nil, fmt.Errorf("invalid map attribute specification %s, line %d of file %s", b.typ, b.sline, b.sfile)
	}
	keyType := b.typ[i+1 : j]
	keyAttr, err := newAttrDef(&baseAttrDef{parentType: b.parentType, nm: b.nm, typ: keyType})
	if err != nil {
		return nil, fmt.Errorf("invalid map attribute specification %s, line %d of file %s", b.typ, b.sline, b.sfile)
	}
	valType := b.typ[j+1:]
	valAttr, err := newAttrDef(&baseAttrDef{parentType: b.parentType, nm: b.nm, typ: valType})
	if err != nil {
		return nil, fmt.Errorf("invalid map attribute specification %s, line %d of file %s", b.typ, b.sline, b.sfile)
	}
	return &mapAttrDef{baseAttrDef: *b, keyType: keyType, keyAttr: keyAttr, valType: valType, valAttr: valAttr}, nil
}

func (a *mapAttrDef) isMultiValued() bool {
	return true
}

func (a *mapAttrDef) generateEquals(w *writer, levelID string) {
	nextLevelID := fmt.Sprintf("%s1", levelID)
	a.checkLevel0Hdr(w, levelID)
	format := `    if len(va%[1]s) != len(vb%[1]s) {
        return false
    }
    for key%[1]s, va%[2]s := range va%[1]s {
		if vb%[2]s, ok := vb%[1]s[key%[1]s]; ok {
`
	w.printf(format, levelID, nextLevelID)
	a.valAttr.generateEquals(w, nextLevelID)
	format = `        } else {
			return false // didn't find key%[1]s in vb%[1]s
		}
    }
`
	w.printf(format, levelID)
	a.checkLevel0Ftr(w, levelID)
}

func (a *mapAttrDef) generateDiff(w *writer, levelID string) {
	nextLevelID := fmt.Sprintf("%s1", levelID)
	a.checkLevel0Hdr(w, levelID)
	format := `    for key%[1]s, va%[2]s := range va%[1]s {
		if vb%[2]s, ok := vb%[1]s[key%[1]s]; ok {
			chgs%[2]s := make([]metago.Chg, 0)
`
	w.printf(format, levelID, nextLevelID)
	a.valAttr.generateDiff(w, nextLevelID)
	format = `            if len(chgs%[2]s) != 0 {
				chgs%[1]s = append(chgs%[1]s, metago.New%[3]sMapChg(&%[4]s%[5]sSREF, key%[1]s, metago.ChangeTypeModify, chgs%[2]s))
			}
		} else {
			chgs%[2]s := make([]metago.Chg, 0)
`
	w.printf(format, levelID, nextLevelID, strings.Title(a.keyType), a.parentType.name, a.nm)
	a.valAttr.generateIns(w, nextLevelID)
	format = `            if len(chgs%[2]s) != 0 {
				chgs%[1]s = append(chgs%[1]s, metago.New%[3]sMapChg(&%[4]s%[5]sSREF, key%[1]s, metago.ChangeTypeInsert, chgs%[2]s))
			}
		}
	}
	for key%[1]s, vb%[2]s := range vb%[1]s {
			chgs%[2]s := make([]metago.Chg, 0)
`
	w.printf(format, levelID, nextLevelID, strings.Title(a.keyType), a.parentType.name, a.nm)
	a.valAttr.generateDel(w, nextLevelID)
	format = `            if len(chgs%[2]s) != 0 {
				chgs%[1]s = append(chgs%[1]s, metago.New%[3]sMapChg(&%[4]s%[5]sSREF, key%[1]s, metago.ChangeTypeDelete, chgs%[2]s))
        }
	}
`
	w.printf(format, levelID, nextLevelID, strings.Title(a.keyType), a.parentType.name, a.nm)
	a.checkLevel0Ftr(w, levelID)
}

func (a *mapAttrDef) generateIns(w *writer, levelID string) {
	nextLevelID := fmt.Sprintf("%s1", levelID)
	format := `    for key%[1]s, va%[2]s := range va%[1]s {
			chgs%[2]s := make([]metago.Chg, 0)
`
	w.printf(format, levelID, nextLevelID)
	a.valAttr.generateIns(w, nextLevelID)
	format = `            if len(chgs%[2]s) != 0 {
				chgs%[1]s = append(chgs%[1]s, metago.New%[3]sMapChg(&%[4]s%[5]sSREF, key%[1]s, metago.ChangeTypeInsert, chgs%[2]s))
		}
	}
`
	w.printf(format, levelID, nextLevelID, strings.Title(a.keyType), a.parentType.name, a.nm)
}

func (a *mapAttrDef) generateDel(w *writer, levelID string) {
	nextLevelID := fmt.Sprintf("%s1", levelID)
	format := `    for key%[1]s, vb%[2]s := range vb%[1]s {
			chgs%[2]s := make([]metago.Chg, 0)
`
	w.printf(format, levelID, nextLevelID)
	a.valAttr.generateDel(w, nextLevelID)
	format = `            if len(chgs%[2]s) != 0 {
				chgs%[1]s = append(chgs%[1]s, metago.New%[3]sMapChg(&%[4]s%[5]sSREF, key%[1]s, metago.ChangeTypeDelete, chgs%[2]s))
		}
	}
`
	w.printf(format, levelID, nextLevelID, strings.Title(a.keyType), a.parentType.name, a.nm)
}

func (a *mapAttrDef) generateApply(w *writer, levelID string) {
	format := `   case &%[1]s%[2]sAID:
			{
			    m := o.%[2]s
`
	w.printf(format, a.parentType.name, a.nm)
	a.generateApplyBody(w, levelID)
	w.printf("			}\n")
}

func (a *mapAttrDef) generateMapModify(w *writer, levelID string) {
	nextLevelID := fmt.Sprintf("%s1", levelID)
	w.printf("				for _, c%[2]s := range mc%[1]s.Chgs {\n", levelID, nextLevelID)
	w.printf("            m%[2]s := m%[1]s[key%[1]s]\n", levelID, nextLevelID)
	a.generateApplyBody(w, nextLevelID)
	w.printf("              }\n")
}

func (a *mapAttrDef) generateApplyBody(w *writer, levelID string) {
	format := `            mc%[1]s := c%[1]s.(*metago.%[2]sMapChg)
	            key%[1]s := mc%[1]s.Key
				switch mc%[1]s.Typ {
				case metago.ChangeTypeModify:
`
	w.printf(format, levelID, strings.Title(a.keyType))
	a.valAttr.generateMapModify(w, levelID)
	format = `				case metago.ChangeTypeInsert:
`
	w.printf(format)
	a.valAttr.generateMapModify(w, levelID)
	format = `				case metago.ChangeTypeDelete:
				delete(m%[1]s, key%[1]s)
            }
`
	w.printf(format, levelID)
}

/************************************************************************/
/**************************** Struct Attribute **************************/
type structAttrDef struct {
	baseAttrDef
	packageName string
}

func newStructAttrDef(b *baseAttrDef) (*structAttrDef, error) {
	return &structAttrDef{baseAttrDef: *b}, nil
}

func (a *structAttrDef) generateEquals(w *writer, levelID string) {
	a.checkLevel0Hdr(w, levelID)
	w.printf("  if va%[1]s.Equals(&vb%[1]s) {\n    return false\n  }\n", levelID)
	a.checkLevel0Ftr(w, levelID)
}

func (a *structAttrDef) generateDiff(w *writer, levelID string) {
	a.checkLevel0Hdr(w, levelID)
	w.printf("    chgs%[1]s = append(chgs%[1]s, metago.NewStructChg(&%[2]s%[3]sSREF, va%[1]s.Diff(&vb%[1]s)))\n", levelID, a.parentType.name, a.nm)
	a.checkLevel0Ftr(w, levelID)
}

func (a *structAttrDef) generateIns(w *writer, levelID string) {
	w.printf("t := %[1]s{}\n", a.typ)
	w.printf("chgs%[1]s = append(chgs%[1]s, metago.NewStructChg(&%[2]s%[3]sSREF, t.Diff(&va%[1]s)))\n", levelID, a.parentType.name, a.nm)
}

func (a *structAttrDef) generateDel(w *writer, levelID string) {
	w.printf("t := %[1]s{}\n", a.typ)
	w.printf("chgs%[1]s = append(chgs%[1]s, metago.NewStructChg(&%[2]s%[3]sSREF, vb%[1]s.Diff(&t)))\n", levelID, a.parentType.name, a.nm)
}

func (a *structAttrDef) generateApply(w *writer, levelID string) {

}
