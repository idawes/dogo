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
	return newAttrDef(&baseAttrDef{parentType: t, attributeID: id, name: name, attrType: fields[2], persistence: p, srcline: r.line, srcfile: r.f.Name()})
}

func newAttrDef(b *baseAttrDef) (attrDef, error) {
	c, err := getClass(b.attrType)
	if err != nil {
		return nil, fmt.Errorf("unknown attribute class %s, line %d of file %s", b.attrType, b.srcline, b.srcfile)
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
	return nil, fmt.Errorf("Unknown attribute class %s, line %d of file %s", c, b.srcline, b.srcfile)
}

func getClass(n string) (attrClass, error) {
	switch n {
	case "byte", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64", "string":
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
	AttributeID() int
	PersistenceClass() persistenceClass
	Srcline() int
	Srcfile() string
	Name() string
	Type() string
	GenerateEquals(w *writer, levelID string)
	GenerateDiff(w *writer, levelID string)
	GenerateIns(w *writer, levelID string)
	GenerateDel(w *writer, levelID string)
}

type baseAttrDef struct {
	parentType  *typedef
	attributeID int
	name        string
	attrType    string
	persistence persistenceClass
	srcline     int
	srcfile     string
}

func (a *baseAttrDef) AttributeID() int {
	return a.attributeID
}

func (a *baseAttrDef) PersistenceClass() persistenceClass {
	return a.persistence
}

func (a *baseAttrDef) Srcline() int {
	return a.srcline
}

func (a *baseAttrDef) Srcfile() string {
	return a.srcfile
}

func (a *baseAttrDef) Name() string {
	return a.name
}

func (a *baseAttrDef) Type() string {
	return a.attrType
}

func (a *baseAttrDef) CheckLevel0Hdr(w *writer, levelID string) {
	if levelID == "" {
		w.printf("    {\n")
		w.printf("        va, vb := o1.%[1]s, o2.%[1]s\n", a.name)
	}
}

func (a *baseAttrDef) CheckLevel0Ftr(w *writer, levelID string) {
	if levelID == "" {
		w.printf("    }\n")
	}
}

func (a *baseAttrDef) GenerateEquals(w *writer, levelID string) {
	a.CheckLevel0Hdr(w, levelID)
	w.printf("  if va%[1]s != vb%[1]s {\n    return false\n  }\n", levelID)
	a.CheckLevel0Ftr(w, levelID)
}

func (a *baseAttrDef) GenerateDiff(w *writer, levelID string) {
	a.CheckLevel0Hdr(w, levelID)
	format := `  if va%[1]s != vb%[1]s {
		d%[1]s.Add(metago.New%[2]sChg(&%[3]s%[4]sSREF, vb%[1]s, va%[1]s))
	}
`
	w.printf(format, levelID, strings.Title(a.Type()), a.parentType.name, a.name)
	a.CheckLevel0Ftr(w, levelID)
}

func (a *baseAttrDef) GenerateIns(w *writer, levelID string) {
	w.printf("d%[1]s.Add(metago.New%[2]sChg(&%[3]s%[4]sSREF, va%[1]s))\n", levelID, strings.Title(a.Type()), a.parentType.name, a.name)
}

func (a *baseAttrDef) GenerateDel(w *writer, levelID string) {
	w.printf("d%[1]s.Add(metago.New%[2]sChg(&%[3]s%[4]sSREF, vb%[1]s))\n", levelID, strings.Title(a.Type()), a.parentType.name, a.name)
}

/************************************************************************/
/************************** Time Attribute ******************************/
type timeAttrDef struct {
	baseAttrDef
}

func (a *timeAttrDef) GenerateEquals(w *writer, levelID string) {
	a.CheckLevel0Hdr(w, levelID)
	w.printf("  if va%[1]s.Equal(vb%[1]s) {\n    return false\n  }\n", levelID)
	a.CheckLevel0Ftr(w, levelID)
}

// parameters:
// 1: current level id
// 2: type name
// 3: attribute name
const TimeDiffTemplate = `  if va%[1]s.Equal(vb%[1]s) {
		d%[1]s.Add(metago.NewTimeChg(&%[2]s%[3]sSREF, vb%[1]s, va%[1]s))
	}
`

func (a *timeAttrDef) GenerateDiff(w *writer, levelID string) {
	a.CheckLevel0Hdr(w, levelID)
	w.printf(TimeDiffTemplate, levelID, a.parentType.name, a.name)
	a.CheckLevel0Ftr(w, levelID)
}

func (a *timeAttrDef) GenerateIns(w *writer, levelID string) {
	w.printf("d%[1]s.Add(metago.NewTimeChg(&%[2]s%[3]sSREF, va%[1]s))\n", levelID, a.parentType.name, a.name)
}

func (a *timeAttrDef) GenerateDel(w *writer, levelID string) {
	w.printf("d%[1]s.Add(metago.NewTimeChg(&%[2]s%[3]sSREF, vb%[1]s))\n", levelID, a.parentType.name, a.name)
}

/************************************************************************/
/**************************** Slice Attribute ***************************/
type sliceAttrDef struct {
	baseAttrDef
	valType string
	valAttr attrDef
}

func newSliceAttrDef(b *baseAttrDef) (*sliceAttrDef, error) {
	valType := b.attrType[2:]
	valAttr, err := newAttrDef(&baseAttrDef{parentType: b.parentType, name: b.name, attrType: valType})
	if err != nil {
		return nil, fmt.Errorf("invalid slice attribute specification %s, line %d of file %s", b.attrType, b.srcline, b.srcfile)
	}
	return &sliceAttrDef{baseAttrDef: *b, valType: valType, valAttr: valAttr}, nil
}

func (a *sliceAttrDef) GenerateEquals(w *writer, levelID string) {
	nextLevelID := fmt.Sprintf("%s1", levelID)
	a.CheckLevel0Hdr(w, levelID)
	format := `    if len(va%[1]s) != len(vb%[1]s) {
        return false
    }
    for idx%[1]s, va%[2]s := range va%[1]s {
		vb%[2]s := vb%[1]s[idx%[1]s]
`
	w.printf(format, levelID, nextLevelID)
	a.valAttr.GenerateEquals(w, nextLevelID)
	w.printf("  }\n")
	a.CheckLevel0Ftr(w, levelID)
}

func (a *sliceAttrDef) GenerateDiff(w *writer, levelID string) {
	nextLevelID := fmt.Sprintf("%s1", levelID)
	a.CheckLevel0Hdr(w, levelID)
	format := `    for idx%[1]s, va%[2]s := range va%[1]s {
        if idx%[1]s  < len(vb%[1]s) {
			vb%[2]s := vb%[1]s[idx%[1]s]
			d%[2]s := &metago.Diff{} 
`
	w.printf(format, levelID, nextLevelID)
	a.valAttr.GenerateDiff(w, nextLevelID)
	format = `            if len(d%[2]s.Changes) != 0 {
				d%[1]s.Changes = append(d%[1]s.Changes, metago.NewSliceChg(&%[3]s%[4]sSREF, idx%[1]s, metago.ChangeTypeModify, d%[2]s))
			}
		} else {
			d%[2]s := &metago.Diff{}
`
	w.printf(format, levelID, nextLevelID, a.parentType.name, a.name)
	a.valAttr.GenerateIns(w, nextLevelID)
	format = `            if len(d%[2]s.Changes) != 0 {
				d%[1]s.Changes = append(d%[1]s.Changes, metago.NewSliceChg(&%[3]s%[4]sSREF, idx%[1]s, metago.ChangeTypeInsert, d%[2]s))
			}
		}
	}
	for idx%[1]s, vb%[2]s := range vb%[1]s {
	    d%[2]s := &metago.Diff{}
`
	w.printf(format, levelID, nextLevelID, a.parentType.name, a.name)
	a.valAttr.GenerateDel(w, nextLevelID)
	format = `        if len(d%[2]s.Changes) != 0 {
		    d%[1]s.Changes = append(d%[1]s.Changes, metago.NewSliceChg(&%[3]s%[4]sSREF, idx%[1]s, metago.ChangeTypeDelete, d%[2]s))
        }
	}
`
	w.printf(format, levelID, nextLevelID, a.parentType.name, a.name)
	a.CheckLevel0Ftr(w, levelID)
}

func (a *sliceAttrDef) GenerateIns(w *writer, levelID string) {
	nextLevelID := fmt.Sprintf("%s1", levelID)
	format := `    for idx%[1]s, va%[2]s := range va%[1]s {
		d%s := &metago.Diff{}
`
	w.printf(format, levelID, nextLevelID)
	a.valAttr.GenerateIns(w, nextLevelID)
	format = `        if len(d%[2]s.Changes) != 0 {
		    d%[1]s.Changes = append(d%[1]s.Changes, metago.NewSliceChg(&%[3]s%[4]sSREF, idx%[1]s, metago.ChangeTypeInsert, d%[2]s))
	    }
	}
`
	w.printf(format, levelID, nextLevelID, a.parentType.name, a.name)
}

func (a *sliceAttrDef) GenerateDel(w *writer, levelID string) {
	nextLevelID := fmt.Sprintf("%s1", levelID)
	format := `    for idx%[1]s, va%[2]s := range va%[1]s {
		d%s := &metago.Diff{}
`
	w.printf(format, levelID, nextLevelID)
	a.valAttr.GenerateDel(w, nextLevelID)
	format = `        if len(d%[2]s.Changes) != 0 {
		    d%[1]s.Changes = append(d%[1]s.Changes, metago.NewSliceChg(&%[3]s%[4]sSREF, idx%[1]s, metago.ChangeTypeDelete, d%[2]s))
	    }
	}
`
	w.printf(format, levelID, nextLevelID, a.parentType.name, a.name)
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
	i := strings.Index(b.attrType, "[")
	j := strings.Index(b.attrType, "]")
	if i == -1 || j == -1 {
		return nil, fmt.Errorf("invalid map attribute specification %s, line %d of file %s", b.attrType, b.srcline, b.srcfile)
	}
	keyType := b.attrType[i+1 : j]
	keyAttr, err := newAttrDef(&baseAttrDef{attrType: keyType})
	if err != nil {
		return nil, fmt.Errorf("invalid map attribute specification %s, line %d of file %s", b.attrType, b.srcline, b.srcfile)
	}
	valType := b.attrType[j+1:]
	valAttr, err := newAttrDef(&baseAttrDef{attrType: valType})
	if err != nil {
		return nil, fmt.Errorf("invalid map attribute specification %s, line %d of file %s", b.attrType, b.srcline, b.srcfile)
	}
	return &mapAttrDef{baseAttrDef: *b, keyType: keyType, keyAttr: keyAttr, valType: valType, valAttr: valAttr}, nil
}

/************************************************************************/
/**************************** Diff Obj Attribute ************************/
type structAttrDef struct {
	baseAttrDef
	packageName string
}

func newStructAttrDef(b *baseAttrDef) (*structAttrDef, error) {
	return &structAttrDef{baseAttrDef: *b}, nil
}
