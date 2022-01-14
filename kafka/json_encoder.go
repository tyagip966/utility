package kafka

import "encoding/json"

type JsonEncoder struct {
	bytes            []byte
	marshallingError error
}

func NewJsonEncoder(value interface{}) JsonEncoder {
	bytes, err := json.Marshal(value)
	return JsonEncoder{
		bytes:            bytes,
		marshallingError: err,
	}
}

func (encoder JsonEncoder) Encode() ([]byte, error) {
	return encoder.bytes, encoder.marshallingError
}

func (encoder JsonEncoder) Length() int {
	return len(encoder.bytes)
}
