package jwt

import "time"

func Renew[T any](token *string, exp *time.Time) (*string, error) {
	if time.Until(*exp).Seconds() < float64(*conf.duration)/2 {
		tokenData := new(T)
		if err := Parse(*token, tokenData); err != nil {
			return nil, err
		}
		return Create(tokenData)
	}
	return token, nil
}
