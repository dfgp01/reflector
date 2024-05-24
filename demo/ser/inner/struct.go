package inner

type (
	User struct {
		Age      int     `json:"age,omitempty"`
		Name     string  `json:"name,omitempty"`
		Children []*User `json:"children,omitempty"`
		Parent   *User   `json:"parent,omitempty"`
	}
)
