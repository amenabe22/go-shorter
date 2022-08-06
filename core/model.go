package core

// define core redirection model and types as a struct
type Redirect struct {
	Code      string `json:"code" bson:"code"`
	URL       string `json:"url" bson:"url" validate:"empty=false & format=url`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}
