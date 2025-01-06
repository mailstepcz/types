package types

import (
	"reflect"
	"testing"

	"github.com/mailstepcz/slice"
	"github.com/stretchr/testify/require"
)

func TestRelevantFields(t *testing.T) {
	req := require.New(t)

	type person struct {
		Name     string
		Age      int
		Excluded bool
	}

	fields := RelevantFields(reflect.TypeFor[person](), func(f reflect.StructField) bool {
		return f.Name != "Excluded"
	})

	req.Equal([]string{"Name", "Age"}, slice.Fmap(func(f reflect.StructField) string { return f.Name }, fields))
}
