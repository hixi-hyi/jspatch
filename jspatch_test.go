package jspatch

import "testing"

// TODO ちゃんとしたテスト書く
func TestSimple(t *testing.T) {
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

	err := jpdoc.CheckApply(model)
	if err != nil {
		t.Errorf("error occured: %#v", err)
	}
	if model.Value != "" {
		t.Errorf("got: %s want: %s", model.Value, "")
	}
}

func TestCheckSucceed(t *testing.T) {
	type Model struct {
		ID    int64  `json:"id" patch:"-"`
		Value string `json:"value" patch:"replace"`
	}
	jpdoc := JSONPatchDocument.New()
	jp := JSONPatch.New()
	jp.Op = "replace"
	jp.Path = "/value"
	jp.Value = "hyi"
	jpdoc.Add(jp)

	model := &Model{
		ID:    100,
		Value: "hixi",
	}

	err := jpdoc.Check(model)
	if err != nil {
		t.Errorf("error occured: %#v", err)
	}
}
func TestCheckFailure(t *testing.T) {
	type Model struct {
		ID    int64  `json:"id" patch:"-"`
		Value string `json:"value" patch:"remove"`
	}
	jpdoc := JSONPatchDocument.New()
	jp := JSONPatch.New()
	jp.Op = "replace"
	jp.Path = "/value"
	jp.Value = "hyi"
	jpdoc.Add(jp)

	model := &Model{
		ID:    100,
		Value: "hixi",
	}

	err := jpdoc.Check(model)
	if err == nil {
		t.Errorf("error occured: %#v", err)
	}
	t.Logf("expected error occured: %s", err)
}

func TestNestedCheckSucceed(t *testing.T) {
	type User struct {
		Name string `json:"name" patch:"replace"`
	}
	type Model struct {
		ID   int64 `json:"id" patch:"-"`
		User User  `json:"user" patch:"-"`
	}
	jpdoc := JSONPatchDocument.New()
	jp := JSONPatch.New()
	jp.Op = "replace"
	jp.Path = "/user/name"
	jp.Value = "hyi"
	jpdoc.Add(jp)

	model := &Model{
		ID: 100,
		User: User{
			Name: "hixi",
		},
	}

	err := jpdoc.Check(model)
	if err != nil {
		t.Errorf("error occured: %#v", err)
	}
}

func Test3NestedCheckSucceed(t *testing.T) {
	type Company struct {
		Address string `json:"address" patch:"replace"`
	}
	type User struct {
		Name    string  `json:"name" patch:"replace"`
		Company Company `json:"company" patch:"replace"`
	}
	type Model struct {
		ID   int64 `json:"id" patch:"-"`
		User User  `json:"user" patch:"-"`
	}
	// replace /user/company/address
	{
		jpdoc := JSONPatchDocument.New()
		jp := JSONPatch.New()
		jp.Op = "replace"
		jp.Path = "/user/company/address"
		jp.Value = "Tokyo"
		jpdoc.Add(jp)

		model := &Model{
			ID: 100,
			User: User{
				Name: "hixi",
				Company: Company{
					Address: "Toyama",
				},
			},
		}

		err := jpdoc.Check(model)
		if err != nil {
			t.Errorf("error occured: %#v", err)
		}
	}
	// replace /user/company/-
	{
		jpdoc := JSONPatchDocument.New()
		jp := JSONPatch.New()
		jp.Op = "replace"
		jp.Path = "/user/company"
		jp.Value = `{"address":"Tokyo"}`
		jpdoc.Add(jp)

		model := &Model{
			ID: 100,
			User: User{
				Name: "hixi",
				Company: Company{
					Address: "Toyama",
				},
			},
		}

		err := jpdoc.Check(model)
		if err != nil {
			t.Errorf("error occured: %#v", err)
		}
	}
}
func Test3NestedCheckFailure(t *testing.T) {
	type Company struct {
		Address string `json:"address" patch:"replace"`
	}
	type User struct {
		Name    string  `json:"name" patch:"replace"`
		Company Company `json:"company" patch:"replace"`
	}
	type Model struct {
		ID   int64 `json:"id" patch:"-"`
		User User  `json:"user" patch:"-"`
	}
	// replace /user/company/address
	{
		jpdoc := JSONPatchDocument.New()
		jp := JSONPatch.New()
		jp.Op = "remove"
		jp.Path = "/user/company/address"
		jpdoc.Add(jp)

		model := &Model{
			ID: 100,
			User: User{
				Name: "hixi",
				Company: Company{
					Address: "Toyama",
				},
			},
		}

		err := jpdoc.Check(model)
		if err == nil {
			t.Errorf("error occured: %#v", err)
		}
		t.Logf("expected error occured: %s", err)
	}
	// replace /user/company/-
	{
		jpdoc := JSONPatchDocument.New()
		jp := JSONPatch.New()
		jp.Op = "remove"
		jp.Path = "/user/company"
		jpdoc.Add(jp)

		model := &Model{
			ID: 100,
			User: User{
				Name: "hixi",
				Company: Company{
					Address: "Toyama",
				},
			},
		}

		err := jpdoc.Check(model)
		if err == nil {
			t.Errorf("error occured: %#v", err)
		}
		t.Logf("expected error occured: %s", err)
	}
}

func TestNestedCheckFailure(t *testing.T) {
	type User struct {
		Name string `json:"name" patch:"replace"`
	}
	type Model struct {
		ID   int64 `json:"id" patch:"-"`
		User User  `json:"user" patch:"-"`
	}
	jpdoc := JSONPatchDocument.New()
	jp := JSONPatch.New()
	jp.Op = "remove"
	jp.Path = "/user/name"
	jpdoc.Add(jp)

	model := &Model{
		ID: 100,
		User: User{
			Name: "hixi",
		},
	}

	err := jpdoc.Check(model)
	if err == nil {
		t.Errorf("error occured: %#v", err)
	}
	t.Logf("expected error occured: %s", err)
}
