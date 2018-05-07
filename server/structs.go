package main

type AnnotationMapValue struct {
	AnnotaitionID int    `json:"annotation_id"`
	Key           string `json:"name"`
	Value         string `json:"value"`
	Index         int    `json:"index"`
}
type Action struct {
	Table  string `json:"table"`
	Action string `json:"action"`
	//Data   interface{} `json:"data"`
	Data AnnotationMapValue `json:"data"`
}
