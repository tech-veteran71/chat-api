package main

import "encoding/json"

type BJSON []byte

func (b BJSON) MarshalJSON() ([]byte, error) {
	return b, nil
}

func (b *BJSON) UnmarshalJSON(data []byte) error {
	*b = data
	return nil
}

func (b *BJSON) Update(callback func(j map[string]interface{}) error) error {
	// Decode JSON.
	var raw map[string]interface{}
	err := json.Unmarshal(*b, &raw)
	if err != nil {
		return err
	}

	// Call update function.
	err = callback(raw)
	if err != nil {
		return err
	}

	// Encode JSON.
	bb, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	*b = bb

	// JSON updated.
	return nil
}
