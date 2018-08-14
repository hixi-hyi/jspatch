# Synopsis
```
package main

type User struct {
    ID        int64     `json:"id"        patch:"test"`
    Name      string    `json:"name"      patch:"test,replace"`
    CreatedAt time.Time `json:"-"         patch:"test"`
    UpdatedAt time.Time `json:"-"         patch:"test"`
}

func main() {
    // http serve
}

func PatchSelf(w http.ResponseWriter, r *http.Request) {
    var jpobj = patch.JSONPatchDocument.New()
    if err := jpobj.Bind(r); err != nil {
        return
    }

    m := GetUser()
    if err := jpobj.CheckApply(m); err != nil {
        return
    }
    PutUser(m)

    // success
}
```
