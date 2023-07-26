package arr

import (
	"sort"
	"strconv"

	"github.com/svcodestore/goutils/rand"
	"github.com/svcodestore/goutils/snowflake"
)

func ListToTree(list []map[string]any, idKey, pidKey, childrenKey, idType string, isGenId bool) (tree []map[string]any) {
	listLen := len(list)
	if listLen == 0 {
		return
	}

	if idType == "string" {
		return ListToTreeWithStrID(list, idKey, pidKey, childrenKey, isGenId)
	}

	if idType == "int" {
		return ListToTreeWithIntID(list, idKey, pidKey, childrenKey, isGenId)
	}

	return
}

func ListToTreeWithStrID(list []map[string]any, idKey, pidKey, childrenKey string, isGenId bool) (tree []map[string]any) {
	listLen := len(list)
	if listLen == 0 {
		return
	}

	var idMap = make(map[string]map[string]any, listLen)
	var pidMap = make(map[string]string, listLen)

	var ids []string
	var sortedList []map[string]any

	for _, m := range list {
		id := (m[idKey]).(string)
		ids = append(ids, id)
		idMap[id] = m

		if pid, ok := m[pidKey]; ok {
			pidMap[pid.(string)] = pid.(string)
		} else {
			return list
		}
	}

	for _, m := range list {
		id := (m[idKey]).(string)
		if pidMap[id] != "" {
			delete(pidMap, id)
		}
	}

	sort.Slice(ids, func(i, j int) bool {
		numA, _ := strconv.Atoi(ids[i])
		numB, _ := strconv.Atoi(ids[j])
		return numA < numB
	})

	for _, id := range ids {
		sortedList = append(sortedList, idMap[id])
	}

	for _, m := range sortedList {
		pid := (m[pidKey]).(string)
		if isGenId {
			m[idKey] = snowflake.SnowflakeId(rand.RandomInt(1024)).String()
		}
		if pidMap[pid] != "" {
			tree = append(tree, m)
		} else {
			if isGenId {
				id := (m[idKey]).(string)
				if pidMap[id] != "" {
					idMap[id][pidKey] = m[idKey]
				}
			}
			parent := idMap[pid]
			if parent[childrenKey] == nil {
				parent[childrenKey] = []map[string]any{}

			}
			if pid != "0" {
				m[pidKey] = parent[idKey]
			}
			children := parent[childrenKey].([]map[string]any)
			children = append(children, m)
			parent[childrenKey] = children
		}
	}

	return
}

func ListToTreeWithIntID(list []map[string]any, idKey, pidKey, childrenKey string, isGenId bool) (tree []map[string]any) {
	listLen := len(list)
	if listLen == 0 {
		return
	}

	var idMap = make(map[int]map[string]any, listLen)
	var pidMap = make(map[int]int, listLen)

	var ids []int
	var sortedList []map[string]any

	for _, m := range list {
		id := int((m[idKey]).(float64))
		ids = append(ids, id)
		idMap[id] = m

		if pid, ok := m[pidKey]; ok {
			pidMap[int(pid.(float64))] = int(pid.(float64))
		} else {
			return list
		}
	}

	for _, m := range list {
		id := int((m[idKey]).(float64))
		delete(pidMap, id)
	}

	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})

	for _, id := range ids {
		sortedList = append(sortedList, idMap[id])
	}

	for _, m := range sortedList {
		pid := int((m[pidKey]).(float64))
		if isGenId {
			m[idKey] = snowflake.SnowflakeId(rand.RandomInt(1024)).String()
		}
		if _, ok := pidMap[pid]; ok {
			tree = append(tree, m)
		} else {
			if isGenId {
				id := int((m[idKey]).(float64))
				if pidMap[id] != 0 {
					idMap[id][pidKey] = m[idKey]
				}
			}
			parent := idMap[pid]
			if parent[childrenKey] == nil {
				parent[childrenKey] = []map[string]any{}

			}
			if pid != 0 {
				m[pidKey] = parent[idKey]
			}
			children := parent[childrenKey].([]map[string]any)
			children = append(children, m)
			parent[childrenKey] = children
		}
	}

	return
}
