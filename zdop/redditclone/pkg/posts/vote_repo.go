package posts

import "sync"

type VoteMemoryRepository struct {
	data  []Vote
	mutex *sync.Mutex
}

func NewVoteMemoryRepo() *VoteMemoryRepository {
	return &VoteMemoryRepository{
		data:  make([]Vote, 0, 10),
		mutex: &sync.Mutex{},
	}
}

func (repo *VoteMemoryRepository) GetAll() ([]*Vote, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	votes := make([]*Vote, 0, len(repo.data))
	for _, vote := range repo.data {
		votes = append(votes, &vote)
	}

	return votes, nil
}

func (repo *VoteMemoryRepository) GetScore() (int, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	score := 0
	for _, vote := range repo.data {
		score += vote.Vote
	}
	return score, nil
}

func (repo *VoteMemoryRepository) GetUpvotePercentage() (int, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	votes := len(repo.data)
	if votes == 0 {
		return 0, nil
	}

	upvotes := 0
	for _, vote := range repo.data {
		if vote.Vote == 1 {
			upvotes++
		}
	}

	return upvotes * 100 / votes, nil
}

func (repo *VoteMemoryRepository) Upvote(userID string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for _, vote := range repo.data {
		if vote.UserID == userID {
			vote.Vote = 1
			return nil
		}
	}
	repo.data = append(repo.data, Vote{UserID: userID, Vote: 1})

	return nil
}

func (repo *VoteMemoryRepository) Downvote(userID string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for _, vote := range repo.data {
		if vote.UserID == userID {
			vote.Vote = -1
			return nil
		}
	}
	repo.data = append(repo.data, Vote{UserID: userID, Vote: -1})

	return nil
}

func (repo *VoteMemoryRepository) Unvote(userID string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for idx, vote := range repo.data {
		if vote.UserID == userID {
			repo.data = append(repo.data[:idx], repo.data[idx+1:]...)
			return nil
		}
	}
	return nil
}
