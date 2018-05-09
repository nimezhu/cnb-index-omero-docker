package main

type App struct {
	Appname string
	Version string
}
type AnnotationMapValue struct {
	AnnotationID int    `json:"annotation_id"`
	Name         string `json:"name"`
	Value        string `json:"value"`
	Index        int    `json:"index"`
	ParentType   string `json:"parent_type"`
	ParentID     int    `json:"parent_id"`
}
type RawAnnotationMapValue struct {
	AnnotationID int    `json:"annotation_id"`
	Name         string `json:"name"`
	Value        string `json:"value"`
	Index        int    `json:"index"`
}
type Action struct {
	Table  string `json:"table"`
	Action string `json:"action"`
	//Data   interface{} `json:"data"`
	Data RawAnnotationMapValue `json:"data"`
}
type BedURI struct {
	genome   string
	chr      string
	start    int
	end      int
	id       string
	color    string
	parentID string //parent should be AnnotaitonMapValue
}

func (b *BedURI) SetId(value string) {
	b.id = value
}
func (b *BedURI) SetColor(value string) {
	b.color = value
}
func (b *BedURI) Color() string {
	return b.color
}
func (b *BedURI) Genome() string {
	return b.genome
}
func (b *BedURI) Chr() string {
	return b.chr
}
func (b *BedURI) Start() int {
	return b.start
}
func (b *BedURI) End() int {
	return b.end
}
func (b *BedURI) Id() string {
	return b.id
}
func (b *BedURI) ParentID() string {
	return b.parentID
}
