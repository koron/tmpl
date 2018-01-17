package gosrc

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"strconv"
)

// WarnLog is log destination for warning messages
var WarnLog = log.New(os.Stderr, "W ", log.LstdFlags)

func warnf(s string, args ...interface{}) {
	if WarnLog != nil {
		WarnLog.Printf(s, args...)
	}
}

// DebugLog is log destination for debug messages
var DebugLog = log.New(os.Stderr, "D ", 0)

func debugf(s string, args ...interface{}) {
	if DebugLog != nil {
		DebugLog.Printf(s, args...)
	}
}

// baseTypeName returns the name of the base type of x (or "")
// and whether the type is imported or not.
//
func baseTypeName(x ast.Expr) (name string, imported bool) {
	switch t := x.(type) {
	case *ast.Ident:
		return t.Name, false
	case *ast.SelectorExpr:
		if _, ok := t.X.(*ast.Ident); ok {
			// only possible for qualified type names;
			// assume type is imported
			return t.Sel.Name, true
		}
	case *ast.StarExpr:
		return baseTypeName(t.X)
	}
	return
}

func typeString(x ast.Expr) string {
	switch t := x.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		if _, ok := t.X.(*ast.Ident); ok {
			return typeString(t.X) + "." + t.Sel.Name
		}
	case *ast.StarExpr:
		return "*" + typeString(t.X)
	}
	return ""
}

func firstName(names []*ast.Ident) string {
	if len(names) == 0 {
		return ""
	}
	return names[0].Name
}

// Package represents a go pacakge.
type Package struct {
	Name string

	Imports []*Import

	Values []*Value
	Funcs  map[string]*Func
	Types  map[string]*Type
}

func (p *Package) putValue(v *Value) {
	p.Values = append(p.Values, v)
}

func (p *Package) putType(typ *Type) {
	if p.Types == nil {
		p.Types = make(map[string]*Type)
	}
	p.Types[typ.Name] = typ
}

func (p *Package) assureType(name string) *Type {
	if typ, ok := p.Types[name]; ok {
		return typ
	}
	typ := &Type{Name: name}
	p.putType(typ)
	return typ
}

func (p *Package) putFunc(fun *Func) {
	if p.Funcs == nil {
		p.Funcs = make(map[string]*Func)
	}
	p.Funcs[fun.Name] = fun
}

// Import represents an import.
type Import struct {
	Name string
	Path string
}

// Var represents a variable.
type Var struct {
	Name string
	Type string
}

// Field represents a variable.
type Field struct {
	Name string
	Type string
	Tag  *Tag
}

// Tag represents a tag for field
type Tag struct {
	Raw    string
	Values map[string]*TagValue
}

func parseTag(tag string) *Tag {
	raw := tag
	values := map[string]*TagValue{}
	// parse tag: partially copied from reflect.StructTag.Lookup()
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			break
		}
		values[name] = parseTagValue(value)
	}
	return &Tag{
		Raw:    raw,
		Values: values,
	}
}

func (tag *Tag) match(name string, value *string) bool {
	for k, v := range tag.Values{
		if k == name {
			if value == nil || v.has(*value) {
				return true
			}
		}
	}
	return false
}

var tagValueRx = regexp.MustCompile(`\s+`)

// TagValue represents content of a tag.
type TagValue struct {
	Raw    string
	Values []string
}

func parseTagValue(s string) *TagValue {
	return &TagValue{
		Raw:    s,
		Values: tagValueRx.Split(s, -1),
	}
}

func (tv *TagValue) has(value string) bool {
	for _, v := range tv.Values  {
		if v == value {
			return true
		}
	}
	return false
}

// Func represents a function.
type Func struct {
	Name    string
	Params  []*Var
	Results []*Var
}

// Type represents a function.
type Type struct {
	Name     string
	Embedded map[string]struct{}
	Methods  map[string]*Func
	Fields   map[string]*Field
	IsStruct bool

	defined bool
}

func (t *Type) putMethod(fun *Func) {
	if t.Methods == nil {
		t.Methods = make(map[string]*Func)
	}
	t.Methods[fun.Name] = fun
}

func (t *Type) putField(f *Field) {
	if t.Fields == nil {
		t.Fields = make(map[string]*Field)
	}
	t.Fields[f.Name] = f
}

func (t *Type) putEmbedded(typeName string) {
	if t.Embedded == nil {
		t.Embedded = make(map[string]struct{})
	}
	t.Embedded[typeName] = struct{}{}
}

// FieldsByTag collects fields which match with query.
// The query's format is "{tagName}" or "{tagName}:{value}".
func (t *Type) FieldsByTag(tagQuery string) []*Field {
	var hits []*Field
	var name string
	var value *string
	for _, f := range t.Fields {
		if f.Tag != nil && f.Tag.match(name, value) {
			hits = append(hits, f)
		}
	}
	return hits
}

// Value represents a value or const
type Value struct {
	Name    string
	Type    string
	IsConst bool
}

// Parser is a parser for go source files.
type Parser struct {
	Package *Package
}

func (p *Parser) readImport(s *ast.ImportSpec) error {
	path, err := strconv.Unquote(s.Path.Value)
	if err != nil {
		return err
	}
	name := ""
	if s.Name != nil {
		name = s.Name.Name
	}
	p.Package.Imports = append(p.Package.Imports, &Import{
		Name: name,
		Path: path,
	})
	return nil
}

func (p *Parser) readValue(d *ast.GenDecl) error {
	prev := ""
	for _, spec := range d.Specs {
		s, ok := spec.(*ast.ValueSpec)
		if !ok {
			warnf("readValue not support: %T", spec)
			continue
		}
		// determine var/const typeName
		typeName := ""
		var isConst bool
		switch {
		case s.Type == nil:
			if n, imp := baseTypeName(s.Type); !imp {
				typeName = n
			}
		case d.Tok == token.CONST:
			typeName = prev
			isConst = true
		}
		for _, n := range s.Names {
			p.Package.putValue(&Value{
				Name:    n.Name,
				Type:    typeName,
				IsConst: isConst,
			})
		}
	}
	return nil
}

func (p *Parser) readType(spec *ast.TypeSpec) error {
	name := spec.Name.Name
	typ := p.Package.assureType(name)
	typ.defined = true
	return p.readTypeFields(spec.Type, typ)
}

func (p *Parser) readTypeFields(expr ast.Expr, typ *Type) error {
	switch t := expr.(type) {
	case *ast.StructType:
		typ.IsStruct = true
		if t.Fields == nil || len(t.Fields.List) == 0 {
			break
		}
		for _, astField := range t.Fields.List {
			f, err := p.toField(astField)
			if err != nil {
				return err
			}
			if f.Name == "" {
				typ.putEmbedded(f.Type)
				break
			}
			typ.putField(f)
		}
	case *ast.InterfaceType:
		if t.Methods == nil || len(t.Methods.List) == 0 {
			break
		}
		for _, astField := range t.Methods.List {
			name := firstName(astField.Names)
			if name == "" {
				// TODO: support embedded interface
				return errors.New("embedded interface not supported")
			}
			switch ft := astField.Type.(type) {
			case *ast.FuncType:
				f, err := p.toFunc(name, ft)
				if err != nil {
					return err
				}
				typ.putMethod(f)
			default:
				return fmt.Errorf("unsupported interface method type: %T", ft)
			}
		}
	}
	return nil
}

func (p *Parser) readFunc(fun *ast.FuncDecl) error {
	f, err := p.toFunc(fun.Name.Name, fun.Type)
	if err != nil {
		return err
	}
	if fun.Recv != nil {
		if len(fun.Recv.List) == 0 {
			// should not happen (incorrect AST);
			return fmt.Errorf("no receivers: %q", fun.Name.Name)
		}
		recvTypeName, imp := baseTypeName(fun.Recv.List[0].Type)
		if imp {
			// should not happen (incorrect AST);
			return fmt.Errorf("method fro imported receiver: %q", recvTypeName)
		}
		p.Package.assureType(recvTypeName).putMethod(f)
		return nil
	}
	p.Package.putFunc(f)
	return nil
}

func (p *Parser) toFunc(name string, funcType *ast.FuncType) (*Func, error) {
	f := &Func{Name: name}
	if funcType != nil {
		var err error
		f.Params, err = p.toVarArray(funcType.Params)
		if err != nil {
			return nil, err
		}
		f.Results, err = p.toVarArray(funcType.Results)
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}

func (p *Parser) toVarArray(fl *ast.FieldList) ([]*Var, error) {
	if fl == nil || len(fl.List) == 0 {
		return nil, nil
	}
	vars := make([]*Var, 0, len(fl.List))
	for _, f := range fl.List {
		v, err := p.toVar(f)
		if err != nil {
			return nil, err
		}
		vars = append(vars, v)
	}
	return vars, nil
}

func (p *Parser) toVar(f *ast.Field) (*Var, error) {
	return &Var{
		Name: firstName(f.Names),
		Type: typeString(f.Type),
	}, nil
}

func (p *Parser) toField(f *ast.Field) (*Field, error) {
	tag, err := p.toTag(f.Tag)
	if err != nil {
		return nil, err
	}
	return &Field{
		Name: firstName(f.Names),
		Type: typeString(f.Type),
		Tag:  tag,
	}, nil
}

func (p *Parser) toTag(x *ast.BasicLit) (*Tag, error) {
	if x == nil {
		return &Tag{}, nil
	}
	switch x.Kind {
	case token.STRING:
		v, err := strconv.Unquote(x.Value)
		if err != nil {
			return nil, err
		}
		return parseTag(v), nil
	default:
		return nil, fmt.Errorf("unsupported token for tag: %s", x.Kind)
	}
}

func (p *Parser) readFile(file *ast.File) error {
	if p.Package == nil || p.Package.Name != file.Name.Name {
		p.Package = &Package{
			Name: file.Name.Name,
		}
	}
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				for _, spec := range d.Specs {
					if s, ok := spec.(*ast.ImportSpec); ok {
						err := p.readImport(s)
						if err != nil {
							return err
						}
					}
				}
			case token.CONST, token.VAR:
				err := p.readValue(d)
				if err != nil {
					return err
				}
			case token.TYPE:
				if len(d.Specs) == 1 && !d.Lparen.IsValid() {
					if s, ok := d.Specs[0].(*ast.TypeSpec); ok {
						err := p.readType(s)
						if err != nil {
							return err
						}
					}
					break
				}
				for _, spec := range d.Specs {
					if s, ok := spec.(*ast.TypeSpec); ok {
						err := p.readType(s)
						if err != nil {
							return err
						}
					}
				}
			}
		case *ast.FuncDecl:
			err := p.readFunc(d)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ReadFile reads a file as a Package.
func ReadFile(name string) (*Package, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, name, f, 0)
	if err != nil {
		return nil, err
	}
	p := &Parser{}
	err = p.readFile(file)
	if err != nil {
		return nil, err
	}
	return p.Package, nil
}

// ReadDir reads all files in a directory as a Package.
func ReadDir(path string) (*Package, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}
	if len(pkgs) != 1 {
		return nil, fmt.Errorf("multiple packages in directory %s", path)
	}
	p := &Parser{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			err := p.readFile(file)
			if err != nil {
				return nil, err
			}
		}
	}
	return p.Package, nil
}

// Read reads a file or directory as a Package.
func Read(path string) (*Package, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return ReadDir(path)
	}
	return ReadFile(path)
}
