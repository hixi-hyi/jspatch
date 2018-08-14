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
type jsonPatchDocument []jsonPatch

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

func (jpdoc jsonPatchDocument) CheckAllowOperation(m interface{}) error {
	allowlist := map[string]string{} // path : ops

	v := reflect.Indirect(reflect.ValueOf(m))
	t := v.Type()
	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		tf := t.Field(i)
		jt := tf.Tag.Get("json")
		jp := tf.Tag.Get("patch")
		allowlist[jt] = jp
	}

	for _, jp := range jpdoc {
		al := allowlist[jp.Path]
		if al != "" {
			if al == "-" {
				return fmt.Errorf("%s: %s operator is not allowed.", jp.Path, jp.Op)
			}
			if !strings.Contains(al, jp.Op) {
				return fmt.Errorf("%s: %s operator is not allowed. allow [%s]", jp.Path, jp.Op, al)
			}
		}
	}

	return nil
}

func (jpdoc *jsonPatchDocument) Apply(m interface{}) error {

	if err := jpdoc.CheckAllowOperation(m); err != nil {
		return err
	}

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

	err = json.Unmarshal(applied, m)
	if err != nil {
		return err
	}

	return nil
}
