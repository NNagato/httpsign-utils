package authenticator

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var defaultMaxTimeGap = time.Second * 30

//ErrNonceNotInRange error when nonce not in acceptable range.
var ErrNonceNotInRange = errors.New("nonce submit is not in aceptable range")

// NonceValidator checking validate by time range
type NonceValidator struct {
	// MaxTimeGap is max time different between client submit timestamp
	// and server time that considered valid. The time precision is millisecond.
	MaxTimeGap time.Duration
}

// NonceValidatorOption is the option of Nonce Validator constructor.
type NonceValidatorOption func(*NonceValidator)

// NonceValidatorWithMaxTimeGap is the option to create NonceValidator with custom time gap.
func NonceValidatorWithMaxTimeGap(gap time.Duration) NonceValidatorOption {
	return func(validator *NonceValidator) {
		validator.MaxTimeGap = gap
	}
}

// NewNonceValidator return NonceValidator with default value (30 second)
func NewNonceValidator(options ...NonceValidatorOption) *NonceValidator {
	v := &NonceValidator{
		MaxTimeGap: defaultMaxTimeGap,
	}
	for _, option := range options {
		option(v)
	}
	if v.MaxTimeGap == 0 {
		v.MaxTimeGap = defaultMaxTimeGap
	}
	return v
}

// Validate return error when checking if header date is valid or not
func (v *NonceValidator) Validate(r *http.Request) error {
	nonce, err := strconv.ParseInt(r.Header.Get("nonce"), 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse nonce in header. Error: %s", err.Error())
	}

	clientTime := time.Unix(0, nonce*int64(time.Millisecond))
	gap := time.Now().Sub(clientTime)
	if gap < 0 {
		gap = -gap
	}
	if gap > v.MaxTimeGap {
		return ErrNonceNotInRange
	}
	return nil
}
