package json

import (
	"encoding/json"
	"shortner/core"

	errs "github.com/pkg/errors"
)

type Redirect struct{}

// decode json to serializer type
func (r *Redirect) Decode(input []byte) (*core.Redirect, error) {
	redirect := &core.Redirect{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errs.Wrap(err, "serializer.Redirect.Decode")
	}
	return redirect, nil
}

// Encode serialzier object to json
func (r *Redirect) Encode(input *core.Redirect) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errs.Wrap(err, "serializer.Redirect.Encode")
	}
	return rawMsg, nil
}
