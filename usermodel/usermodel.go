package usermodel

//User is represent the struct of the user in user collection
type User struct {
	UID      string `json:"uid"`      //User id
	Nick     string `json:"nickname"` //Nick name wich will display(for example in posts)
	Name     string `json:"name"`     //Actual Name
	Surname  string `json:"surname"`  //Actual Surname
	Avatar   string `json:"avatar"`   //Path in firebase storage to avatar // For displaying we will use an autogenereted thubnails
	Email    string `json:"email"`    //Email of user
	IsSetted bool   `json:"issetted"` //Is User was setted(specified name avatar and nick)
}
