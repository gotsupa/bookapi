package book

type Books struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Author     string `json:"author"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}
