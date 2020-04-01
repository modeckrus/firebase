package postpubliched

//Post is actualy post whic will in collection post(publiched)
type Post struct {
	UID         string   `json:"uid"`         //User id(who publiched this post)
	PostID      string   `json:"postId"`      //Post Identifer
	Nick        string   `json:"nick"`        //Displayable name
	Title       string   `json:"title"`       //Title of this post
	Body        string   `json:"body"`        //Body of this post
	Hasattach   bool     `json:"hasattach"`   //Has attachment
	Images      []string `json:"images"`      //Images array of path to firestorage
	Likes       int      `json:"likes"`       //Value of likes
	BestComment string   `json:"bestcomment"` //Best comment about this post
}
