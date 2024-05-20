package convertor

type OldItem struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	Count    int        `json:"count"`
	Children *[]OldItem `json:"children"`
}
