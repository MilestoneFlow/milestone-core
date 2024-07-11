package rewards

type RewardError struct {
	HttpCode int
	Message  string
}

func (e *RewardError) Error() string {
	return e.Message
}

func NewRewardError(httpCode int, message string) *RewardError {
	return &RewardError{
		HttpCode: httpCode,
		Message:  message,
	}
}

var Errors = struct {
	NotFound       *RewardError
	KeyExistsError *RewardError
	InvalidFormat  *RewardError
	NotFoundByKey  *RewardError
}{
	NotFound:       NewRewardError(404, "Reward not found"),
	KeyExistsError: NewRewardError(400, "Reward with this key already exists"),
	InvalidFormat:  NewRewardError(400, "Invalid reward"),
	NotFoundByKey:  NewRewardError(404, "Reward not found by key"),
}
