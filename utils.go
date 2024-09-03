package types

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/fealsamh/go-utils/dbutils"
	"github.com/google/uuid"
	"github.com/mailstepcz/maybe"
	"github.com/mailstepcz/types/iface"
	"github.com/mailstepcz/validate"
	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Predefined types.
var (
	String       = reflect.TypeFor[string]()
	UUID         = reflect.TypeFor[uuid.UUID]()
	UUIDPtr      = reflect.TypeFor[*uuid.UUID]()
	NullUUID     = reflect.TypeFor[uuid.NullUUID]()
	Time         = reflect.TypeFor[time.Time]()
	TimePtr      = reflect.TypeFor[*time.Time]()
	NullTime     = reflect.TypeFor[sql.NullTime]()
	TimestampPtr = reflect.TypeFor[*timestamppb.Timestamp]()
	Decimal      = reflect.TypeFor[decimal.Decimal]()
	Copiable     = reflect.TypeFor[iface.Copiable]()
	LanguageTag  = reflect.TypeFor[language.Tag]()
	ClosedEnum   = reflect.TypeFor[iface.ClosedEnum]()
	ULID         = reflect.TypeFor[ulid.ULID]()
	StructpbPtr  = reflect.TypeFor[*structpb.Struct]()
	Maybe        = maybe.IfaceType
	Required     = validate.RequiredIfaceType
)

// CheckAlignment checks whether two types are well aligned fieldwise.
func CheckAlignment(t1, t2 reflect.Type) error {
	fields1 := RelevantFields(t1)
	fields2 := RelevantFields(t2)
	if len(fields1) != len(fields2) {
		return fmt.Errorf("not equal number of fields: %d vs %d", len(fields1), len(fields2))
	}
	for i, f1 := range fields1 {
		f2 := fields2[i]
		if f1.Name != f2.Name {
			return fmt.Errorf("not equal names of fields: %s vs %s", f1.Name, f2.Name)
		}
	}
	return nil
}

// RelevantFields returns the list of relevant type fields.
func RelevantFields(t reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath == "" && !f.Anonymous {
			fields = append(fields, f)
		}
	}
	return fields
}

// TypeAttributes returns the attributes of a Postgres type.
func TypeAttributes(db dbutils.Querier, typname string) ([]string, error) {
	rows, err := db.QueryContext(context.Background(), `
		SELECT attname FROM pg_attribute
		INNER JOIN pg_type ON typrelid = attrelid
		WHERE typname = $1 AND attisdropped IS FALSE
		ORDER BY attnum`, typname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		fields = append(fields, name)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fields, nil
}
