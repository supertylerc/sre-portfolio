package leader

const CheckInterval = 200

type Leader interface {
	IsCurrentLeader() (bool, error)
	ReadLeader() (string, error)
	WriteLeader() error
}
