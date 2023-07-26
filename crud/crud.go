package crud

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/svcodestore/goutils/arr"
	"github.com/svcodestore/goutils/str"
)

type CrudRequestData struct {
	Creates []map[string]any `json:"A"`
	Updates []map[string]any `json:"U"`
	Deletes []int            `json:"D"`
}

func ParseBatchCrudFormData(values url.Values) (adds []map[string]any, updates []map[string]any, deletes []string) {
	for k, val := range values {
		data := str.SplitByPairSymbol(k, "[", "]")
		if strings.HasPrefix(k, "D") {
			deletes = append(deletes, val...)
		} else {
			i, e := strconv.Atoi(data[0])
			if e != nil {
				continue
			}
			v := ""
			if len(val) > 0 {
				v = val[0]
			}
			if strings.HasPrefix(k, "A") {
				if len(adds) == 0 {
					adds = append(adds, map[string]any{})
				}
				if len(adds) > i {
					if v != "" {
						adds[i][data[1]] = v
					}
				} else {
					// problem
					adds = append(adds, map[string]any{})
					if v != "" {
						adds[i][data[1]] = v
					}
				}
			}
			if strings.HasPrefix(k, "U") {
				if len(updates) == 0 {
					updates = append(updates, map[string]any{})
				}
				if len(updates) > i {
					updates[i][data[2]] = v
					updates[i]["id"] = data[1]
				} else {
					// problem
					updates = append(updates, map[string]any{})
					updates[i][data[2]] = v
					updates[i]["id"] = data[1]
				}
			}
		}
	}
	return
}

func ExecFormCrudBatch(values url.Values, addCb func(b []byte) (err error), updateCb func(b []byte) (err error), deleteCb func(ids []string) (err error)) (err error) {
	adds, updates, deletes := ParseBatchCrudFormData(values)
	if adds == nil && updates == nil && deletes == nil {
		return
	}

	addLen := len(adds)
	if addLen > 0 {
		for i := 0; i < addLen; i++ {
			b, e := json.Marshal(adds[i])
			if e != nil {
				err = e
				return
			}
			e = addCb(b)
			if e != nil {
				err = e
				return
			}
		}
	}

	updateLen := len(updates)
	if updateLen > 0 {
		for i := 0; i < updateLen; i++ {
			b, e := json.Marshal(updates[i])
			if e != nil {
				err = e
				return
			}
			e = updateCb(b)
			if e != nil {
				err = e
				return
			}
		}
	}

	deleteLen := len(deletes)
	if deleteLen > 0 {
		e := deleteCb(deletes)
		if e != nil {
			err = e
			return
		}
	}

	return
}

func ExecJsonCrudBatch(data *CrudRequestData, addCb func(b []byte) (id any, err error), updateCb func(b []byte) (err error), deleteCb func(ids []int) (err error)) (err error) {
	adds := data.Creates
	updates := data.Updates
	deletes := data.Deletes
	if adds == nil && updates == nil && deletes == nil {
		return
	}

	if adds != nil {
		addLen := len(adds)
		if addLen > 0 {
			tree := arr.ListToTree(adds, "id", "pid", "children", "int", false)
			err = addTreeCb(nil, tree, addCb)
			if err != nil {
				return
			}
		}
	}

	if updates != nil {
		updateLen := len(updates)
		if updateLen > 0 {
			for i := 0; i < updateLen; i++ {
				b, e := json.Marshal(updates[i])
				if e != nil {
					err = e
					return
				}
				e = updateCb(b)
				if e != nil {
					err = e
					return
				}
			}
		}
	}

	if deletes != nil {
		deleteLen := len(deletes)
		if deleteLen > 0 {
			e := deleteCb(deletes)
			if e != nil {
				err = e
				return
			}
		}
	}

	return
}

func addTreeCb(pid any, tree []map[string]any, addCb func(b []byte) (id any, err error)) (err error) {
	for _, m := range tree {
		m["id"] = 0
		if pid != nil {
			m["pid"] = pid
		}
		children := m["children"]
		if children != nil {
			delete(m, "children")
		}

		b, e := json.Marshal(m)
		if e != nil {
			err = e
			return
		}
		id, e := addCb(b)

		if e != nil {
			return e
		}

		if children != nil && len(children.([]map[string]any)) > 0 {
			err = addTreeCb(id, children.([]map[string]any), addCb)
		}
	}

	return
}
