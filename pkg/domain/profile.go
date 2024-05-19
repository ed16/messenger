package domain

type Profile struct {
	UserId      int64
	Description string
	PhotoURL    string
}

func NewProfile() *Profile {
	return &Profile{}
}
