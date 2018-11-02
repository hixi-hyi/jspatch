package jspatch

import "testing"

// TODO ちゃんとしたテスト書く
func TestCheckOperation(t *testing.T) {

	type Model struct {
		ID    int64  `json:"id" patch:"-"`
		Value string `json:"value" patch:"remove"`
	}
	jpdoc := JSONPatchDocument.New()
	jp := JSONPatch.New()
	jp.Op = "remove"
	jp.Path = "/value"
	jpdoc.Add(jp)

	model := &Model{
		ID:    100,
		Value: "hixi",
	}

	t.Logf("%#v", model)

	err := jpdoc.CheckApply(model)
	t.Logf("%#v", err)

	t.Logf("%#v", model)
}
