package datastore

type VoteStore interface {
	RecordVote(itemID, itemType, action string, botID string) error
	HasVoted(itemID, itemType, botID string) (bool, string, error)
	Close() error
}
