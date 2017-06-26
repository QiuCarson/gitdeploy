package models

type SystemService struct{}

func (this *SystemService) GetPermList() map[string][]Perm {
	var list []Perm
	o.Raw("SELECT * FROM " + tableName("perm")).QueryRows(&list)

	result := make(map[string][]Perm)
	for _, v := range list {
		v.Key = v.Module + "." + v.Action
		if _, ok := result[v.Module]; !ok {
			result[v.Module] = make([]Perm, 0)
		}
		result[v.Module] = append(result[v.Module], v)
	}
	return result
}
