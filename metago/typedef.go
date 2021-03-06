// Copyright 2015 Ian Dawes. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/nu7hatch/gouuid"
)

type typeID struct {
	pkg *uuid.UUID
	typ uint32
}

func (t *typeID) String() string {
	return fmt.Sprintf("%s:%d", t.pkg, t.typ)
}

func (t *typeID) Compare(o *typeID) int {
	if t.pkg != o.pkg {
		return bytes.Compare((*t.pkg)[:], (*o.pkg)[:])
	}
	return int(t.typ) - int(o.typ)
}

type typedef struct {
	pkg                 string
	typeID              typeID
	name                string
	abstract            bool
	extendsName         string
	extends             *typedef
	srcfile             string
	srcline             int
	attrDefsByID        map[int]attrDef
	attrDefsByIDInOrder attrdefList
	abstractMethods     []string
	methods             []*methodDef
	imports             map[string]struct{}
	persistent          bool
}

type attrdefList []attrDef

func (l attrdefList) Len() int           { return len(l) }
func (l attrdefList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l attrdefList) Less(i, j int) bool { return l[i].attributeID() < l[j].attributeID() }

/*
	Main internal entry point for parsing a type definition.

	type <typeID> <typeName> <concrete|abstract> [extends <typeName>] {
		attributes {
			<attributedef>*
		}
		methods {
			<import|method>*
		}

	If this type extends another type, the supertype must be defined in the same package. Any method defined on a subtype that has the same name as a method in the supertype must also
	have the same signature.

	Imports may occur at any point in the method block, allowing them to be specified just above the method that requires it. Imports may be duplicated, they will only be present once in the generated
	output.
*/
func parseTypedef(pkg string, pkgUUID *uuid.UUID, r *reader) (*typedef, error) {
	pkgParts := strings.Split(pkg, "/")
	pkgTail := "\"" + pkgParts[len(pkgParts)-1] + "\""
	t, err := parseTypedefHeader(pkgTail, pkgUUID, r)
	if err != nil {
		return nil, err
	}
	if err := t.parseImportBlock(r); err != nil {
		return nil, err
	}
	if err := t.parseAttributeBlock(r); err != nil {
		return nil, err
	}
	if err := t.parseMethodBlock(r); err != nil {
		return nil, err
	}
	fields, err := r.read()
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("incomplete type specification, line %d of file %s", r.line, r.f.Name())
		}
		return nil, err
	}
	if len(fields) != 1 && fields[0] != "}" {
		return nil, fmt.Errorf("missing closing \"}\", line %d of file %s", r.line, r.f.Name())
	}
	return t, nil
}

/*
	typedef header:
		type <typeID> <typeName> <concrete|abstract> [extends <typeName>] {
*/
func parseTypedefHeader(pkg string, pkgUUID *uuid.UUID, r *reader) (*typedef, error) {
	fields, err := r.read()
	if err != nil {
		return nil, err
	}
	if *veryVerbose {
		fmt.Printf("    Parsing %q as type def header from line %d of file %s\n", fields, r.line, r.f.Name())
	}
	if len(fields) < 4 {
		return nil, fmt.Errorf("invalid type specification, line %d of file %s", r.line, r.f.Name())
	}
	if fields[0] != "type" {
		return nil, fmt.Errorf("expecting \"type\", found \"%s\", line %d o file %s", fields[0], r.line, r.f.Name())
	}
	t := &typedef{pkg: pkg, srcfile: r.f.Name(), srcline: r.line}
	subID, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil, fmt.Errorf("expecting integer type number, found \"%s\", line %d of file %s", fields[1], r.line, r.f.Name())
	}
	t.typeID = typeID{pkg: pkgUUID, typ: uint32(subID)}
	t.name = fields[2]
	if fields[3] == "abstract" {
		t.abstract = true
	} else if fields[3] == "concrete" {
		t.abstract = false
	} else {
		return nil, fmt.Errorf("expecting \"abstract\" or \"concrete\", found \"%s\", line %d of file %s", fields[3], r.line, r.f.Name())
	}
	if fields[len(fields)-1] != "{" {
		return nil, fmt.Errorf("missing \"{\", line %d of file %s", r.line, r.f.Name())
	}
	if len(fields) == 5 {
		return t, nil
	}
	if fields[4] != "extends" {
		return nil, fmt.Errorf("expecting \"extends\", found \"%s\", line %d of file %s", fields[5], r.line, r.f.Name())
	}
	if len(fields) < 7 {
		return nil, fmt.Errorf("missing base type name, line %d of file %s", r.line, r.f.Name())
	}
	t.extendsName = fields[5]

	return t, nil
}

func (t *typedef) parseImportBlock(r *reader) error {
	t.imports = make(map[string]struct{})
	fields, err := r.read()
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("incomplete type specification, line %d of file %s", r.line, r.f.Name())
		}
		return err
	}
	if len(fields) != 2 || fields[0] != "imports" || fields[1] != "{" {
		if *veryVerbose {
			fmt.Println("      No imports specified")
		}
		r.unread(fields)
		return nil
	}
	for {
		fields, err := r.read()
		if err != nil {
			return err
		}
		if len(fields) != 1 {
			return fmt.Errorf("invalid import block entry, line %d of file %s", r.line, r.f.Name())
		}
		if fields[0] == "}" { // naked "}" means end of block.
			if *veryVerbose {
				fmt.Printf("      Found imports: %v\n", t.imports)
			}
			return nil
		}
		t.imports[strings.Replace(fields[0], "\"", "", -1)] = struct{}{}
	}
}

/*
	attribute block:
		attributes {
			<attributedef>*
		}
*/
func (t *typedef) parseAttributeBlock(r *reader) error {
	t.attrDefsByID = make(map[int]attrDef, 10)
	fields, err := r.read()
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("incomplete type specification, line %d of file %s", r.line, r.f.Name())
		}
		return err
	}
	if len(fields) != 2 || fields[0] != "attributes" || fields[1] != "{" {
		if *veryVerbose {
			fmt.Println("      No attributes specified")
		}
		r.unread(fields)
		return nil
	}
	for {
		fields, err := r.read()
		if err != nil {
			return err
		}
		if len(fields) == 1 && fields[0] == "}" { // naked "}" means end of block.
			// generate sorted attribute list before returning
			if len(t.attrDefsByID) > 0 {
				t.attrDefsByIDInOrder = make([]attrDef, 0)
				for _, a := range t.attrDefsByID {
					t.attrDefsByIDInOrder = append(t.attrDefsByIDInOrder, a)
				}
				sort.Sort(t.attrDefsByIDInOrder)
			}
			return nil
		}
		r.unread(fields)
		a, err := parseAttribute(t, r)
		if err != nil {
			return err
		}
		if old, present := t.attrDefsByID[a.attributeID()]; present {
			return fmt.Errorf("duplicate defintion of attribute id %d on line %d of file %s\n   It is also defined on line %d of file %s", a.attributeID(), a.srcline(), a.srcfile(), old.srcline(), old.srcfile())
		}
		t.attrDefsByID[a.attributeID()] = a
		if *veryVerbose {
			fmt.Printf("      Found attribute spec: %v\n", a)
		}
	}
}

/*
	method block:
		methods {
			<import|method>*
		}

	import:
		standard go import spec.

	method:
		standard golang method decl, except leaving out the receiver spec, which will be automatically generated. Special syntax ##super##(<args>) can be used to invoke function of same name in supertype.
*/
func (t *typedef) parseMethodBlock(r *reader) error {
	t.methods = make([]*methodDef, 0)
	fields, err := r.read()
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("incomplete type specification, line %d of file %s", r.line, r.f.Name())
		}
		return err
	}
	if len(fields) != 2 || fields[0] != "methods" || fields[1] != "{" {
		if *veryVerbose {
			fmt.Println("      No methods specified")
		}
		r.unread(fields)
		return nil
	}
	for {
		fields, err := r.read()
		if err != nil {
			return err
		}
		if len(fields) == 1 && fields[0] == "}" { // naked "}" means end of block.
			return nil
		}
		if len(fields) > 0 && fields[0] == "func" {
			r.unread(fields)
			m, err := parseMethod(t, r)
			if err != nil {
				return err
			}
			t.methods = append(t.methods, m)
			continue
		}
		return fmt.Errorf("invalid keyword: \"%s\", line %d of file %s", fields[0], r.line, r.f.Name())
	}
}

func (t *typedef) validateTypeHierarchy(typesByName map[string]*typedef) error {
	if t.extendsName == "" {
		return nil
	}
	super, ok := typesByName[t.extendsName]
	if !ok {
		return fmt.Errorf("couldn't find type %s while validating hierarchy for type %s, defined on line %d of file %s", t.extendsName, t.name, t.srcline, t.srcfile)
	}
	t.extends = super
	return nil
}

func (t *typedef) generate(w *writer) {
	if t.abstract {
		return
	}
	w.printf("import (\n")
	t.generateImports(w)
	w.printf(")\n\ntype %s struct {\n", t.name)
	t.generateAttributes(w)
	w.printf("}\n")
	t.generateMethods(w)
}

func (t *typedef) generateImports(w *writer) {
	t.coalesceImports(t.imports)
	for im := range t.imports {
		w.printf("  \"%s\"\n", im)
	}
}

func (t *typedef) coalesceImports(imports map[string]struct{}) {
	if t.extends != nil {
		t.extends.coalesceImports(imports)
	}
	for k, v := range t.imports {
		imports[k] = v
	}
	imports["github.com/idawes/metago"] = struct{}{}
}

func (t *typedef) generateAttributes(w *writer) {
	if t.extends != nil {
		t.extends.generateAttributes(w)
	}
	for _, a := range t.attrDefsByIDInOrder {
		w.printf("  %s %s\n", a.name(), a.typeName())
	}
}

func (t *typedef) generateMethods(w *writer) {
	methods := make(map[string]*methodDef)
	t.resolveMethods(methods)
	s := make([]string, 0)
	for n, _ := range methods {
		s = append(s, n)
	}
	sort.Strings(s)
	for _, n := range s {
		m := methods[n]
		if strings.Contains(m.name, "_super") {
			w.printf("\n// from: %s", m.parentType.name)
		}
		w.printf("\nfunc (this *%s) %s(%s) %s {\n%s}\n", t.name, m.name, m.params, m.returns, m.body)
	}
	t.generateEquals(w)
	t.generateDiff(w)
	t.generateApply(w)
}

func (t *typedef) resolveMethods(methods map[string]*methodDef) {
	for _, m := range t.methods {
		name := m.name
		super := fmt.Sprintf("this.%s_super", name)
		for {
			if _, exists := methods[name]; exists {
				name = fmt.Sprintf("%s_super", name)
				super = fmt.Sprintf("%s_super", super)
			} else {
				break
			}
		}
		if m.name != name {
			mcopy := *m
			m = &mcopy
		}
		m.name = name
		m.body = strings.Replace(m.body, "##super##", super, -1)
		methods[m.name] = m
	}
	if t.extends != nil {
		t.extends.resolveMethods(methods)
	}
}

func (t *typedef) generateEquals(w *writer) {
	w.printf("\nfunc (o1 %[1]s) Equals(o2 %[1]s) bool {\n", t.name)
	t.generateAttrEquals(w)
	w.printf("    return true\n}\n")
}

func (t *typedef) generateAttrEquals(w *writer) {
	if t.extends != nil {
		t.extends.generateAttrEquals(w)
	}
	for _, a := range t.attrDefsByIDInOrder {
		w.printf("\n")
		a.generateEquals(w, "")
	}
}

func (t *typedef) generateDiff(w *writer) {
	w.printf("\n// Diff returns a record of the differences between o1 and o2 such that applying the generated record to o1 would make it equal to o2")
	w.printf("\nfunc (o1 %[1]s) Diff(o2 %[1]s) metago.Diff {\n", t.name)
	w.printf("    chgs := make([]metago.Chg, 0)\n")
	t.generateAttrDiffs(w)
	w.printf("    return metago.Diff{Chgs: chgs}\n}\n")

}

func (t *typedef) generateAttrDiffs(w *writer) {
	if t.extends != nil {
		t.extends.generateAttrDiffs(w)
	}
	for _, a := range t.attrDefsByIDInOrder {
		w.printf("\n")
		a.generateDiff(w, "")
	}
}

func (t *typedef) generateApply(w *writer) {
	format := `
func (orig *%[1]s) Apply(d metago.Diff) error {
	for _, c := range d.Chgs {
		switch c.AttributeID() {
`
	w.printf(format, t.name)
	t.generateAttrApply(w)
	format = `		}
	}
	return nil
}
`
	w.printf(format)
}

func (t *typedef) generateAttrApply(w *writer) {
	if t.extends != nil {
		t.extends.generateAttrApply(w)
	}
	for _, a := range t.attrDefsByIDInOrder {
		w.printf("\n")
		a.generateApply(w, "")
	}
}

func (t *typedef) generateSchema(w *writer) {
	if t.abstract {
		return
	}
	if *verbose {
		fmt.Printf("generating Schema for %s\n", t.name)
	}
	w.printf("\n  %sTID  = metago.TypeID{Pkg: &MetagoPackageUUID, PkgName: %s, Typ: %d, Name: \"%s\"}\n", t.name, t.pkg, t.typeID.typ, t.name)
	for _, a := range t.attrDefsByIDInOrder {
		w.printf("  %[1]s%[2]sAID = metago.AttrID{TypeID: &%[1]sTID, Attr: %[3]d, Name: \"%[4]s\"}\n", t.name, a.name(), a.attributeID(), a.name())
		w.printf("  %[1]s%[2]sSREF = metago.Attrdef{ID: &%[1]s%[2]sAID, Persistence: metago.%[3]s}\n", t.name, a.name(), strings.Title(a.persistenceClass().String()))
	}
}
