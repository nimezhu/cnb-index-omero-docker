package main

type AnnotationMapValue struct {
	AnnotaitionID int    `json:"annotation_id"`
	Name          string `json:"name"`
	Value         string `json:"value"`
	Index         int    `json:"index"`
}
type Action struct {
	Table  string `json:"table"`
	Action string `json:"action"`
	//Data   interface{} `json:"data"`
	Data AnnotationMapValue `json:"data"`
}
type Bed4 struct {
	chr   string
	start int
	end   int
	id    string
}

func (b *Bed4) Chr() string {
	return b.chr
}
func (b *Bed4) Start() int {
	return b.start
}
func (b *Bed4) End() int {
	return b.end
}
func (b *Bed4) Id() string {
	return b.id
}
