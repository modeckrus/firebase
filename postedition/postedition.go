package postedition

//PostEdition is struct represent the post which not beeing published
type PostEdition struct {
	Title     string   `json:"title"` //Title of the post
	Subtitle  string   `json:"subtitle"`
	Body      string   `json:"body"`      //Body of the post
	Hasattach bool     `json:"hasattach"` //Has attachment
	Images    []string //Imagese contains the paths of the images in firestorage
}
