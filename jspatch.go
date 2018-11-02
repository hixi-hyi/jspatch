package jspatch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
)

type jsonPatch struct {
	Op    string      `json:"op" binding:"required"`
	Path  string      `json:"path" binding:"required"`
	Value interface{} `json:"value"`
}
type jsonPatchDocument []*jsonPatch

var JSONPatch = &jsonPatch{}
var JSONPatchDocument = jsonPatchDocument{}

func (_ *jsonPatch) New() *jsonPatch {
	return &jsonPatch{}
}

func (_ *jsonPatchDocument) New() jsonPatchDocument {
	return jsonPatchDocument{}
}

func (m *jsonPatchDocument) Bind(req *http.Request) error {
	d := json.NewDecoder(req.Body)
	d.UseNumber()
	err := d.Decode(&m)
	return err
}

func (m *jsonPatchDocument) Add(jp *jsonPatch) {
	*m = append(*m, jp)
}

func (jpdoc jsonPatchDocument) CheckWithPrefix(m interface{}, prefix string) error {
	allowlist := map[string]string{} // path : ops

	v := reflect.Indirect(reflect.ValueOf(m))
	t := v.Type()
	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		tf := t.Field(i)
		switch tf.Type.Kind() {
		case reflect.Struct:
			jt := tf.Tag.Get("json")
			if err := jpdoc.CheckWithPrefix(v.Field(i).Interface(), prefix+jt+"/"); err != nil {
				return err
			}
		default:
			jt := tf.Tag.Get("json")
			jp := tf.Tag.Get("patch")
			allowlist[prefix+jt] = jp
		}
	}

	for _, jp := range jpdoc {
		al := allowlist[jp.Path]
		if al != "" {
			if al == "-" {
				return fmt.Errorf("'%s' operation is not allowed. [%s]", jp.Op, jp.Path)
			}
			if !strings.Contains(al, jp.Op) {
				return fmt.Errorf("'%s' operation is not allowed. [%s]. allow [%s]", jp.Op, jp.Path, al)
			}
		}
	}

	return nil
}

func (jpdoc jsonPatchDocument) Check(m interface{}) error {
	return jpdoc.CheckWithPrefix(m, "/")
}

func (jpdoc *jsonPatchDocument) CheckApply(m interface{}) error {
	if err := jpdoc.Check(m); err != nil {
		return err
	}
	return jpdoc.Apply(m)
}

func (jpdoc *jsonPatchDocument) Apply(m interface{}) error {
	doc, err := json.Marshal(m)
	if err != nil {
		return err
	}

	rawjpdoc, err := json.Marshal(jpdoc)
	if err != nil {
		return err
	}
	obj, err := jsonpatch.DecodePatch(rawjpdoc)
	if err != nil {
		return err
	}

	applied, err := obj.Apply(doc)
	if err != nil {
		return err
	}

	// clear
	p := reflect.ValueOf(m).Elem()
	p.Set(reflect.Zero(p.Type()))

	err = json.Unmarshal(applied, m)
	if err != nil {
		return err
	}

	return nil
}
