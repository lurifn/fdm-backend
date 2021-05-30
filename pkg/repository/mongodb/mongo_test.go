package mongodb

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestUnmarshalBson(t *testing.T) {
	obj := `{"key1":"value1","key2": 2}`
	result, err := unmarshalBson(obj)
	expected := bson.M{"key1": "value1", "key2": int32(2)}

	assert.NoError(t, err)
	assert.EqualValues(t, expected, result)
}
