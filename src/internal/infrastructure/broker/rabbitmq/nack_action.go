package rabbitmq

import "errors"

type NackAction int

const (
	NackDiscard NackAction = iota
	NackRequeue
)

type NackError struct {
	Action NackAction
	Err    error
}

func (e *NackError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err == nil {
		return "nack"
	}
	return e.Err.Error()
}

func (e *NackError) Unwrap() error { return e.Err }

func Discard(err error) error { return &NackError{Action: NackDiscard, Err: err} }
func Requeue(err error) error { return &NackError{Action: NackRequeue, Err: err} }

func NackActionFromError(err error) (NackAction, bool) {
	var nackErr *NackError
	if errors.As(err, &nackErr) && nackErr != nil {
		return nackErr.Action, true
	}
	return 0, false
}
