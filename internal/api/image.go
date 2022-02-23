package api

type Image struct {
	ID   string   `jsonapi:"primary,image"`
	URL  string   `jsonapi:"attr,url"`
	Tags []string `jsonapi:"attr,tags"`
}
