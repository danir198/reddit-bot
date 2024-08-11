package datastore

type VoteStore interface {
	RecordVote(itemID, itemType, action string, botID string) error
	HasVoted(itemID, botID string) (bool, string, error)
	Close() error
}
