package posts

type Vote struct {
	UserID string `json:"user"`
	Vote   int    `json:"vote"`
}

type VotesRepo interface {
	GetAll() ([]*Vote, error)
	GetScore() (int, error)
	GetUpvotePercentage() (int, error)
	Upvote(userID string) error
	Downvote(userID string) error
	Unvote(userID string) error
}
