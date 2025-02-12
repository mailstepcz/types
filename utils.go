package types

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/fealsamh/go-utils/dbutils"
	"github.com/google/uuid"
	"github.com/mailstepcz/enums"
	"github.com/mailstepcz/maybe"
	"github.com/mailstepcz/serr"
	"github.com/mailstepcz/types/iface"
	"github.com/mailstepcz/validate"
	"github.com/oklog/ulid/v2"
	"github.com/rickb777/date/v2"
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
	Date         = reflect.TypeFor[date.Date]()
	Time         = reflect.TypeFor[time.Time]()
	TimePtr      = reflect.TypeFor[*time.Time]()
	NullTime     = reflect.TypeFor[sql.NullTime]()
	TimestampPtr = reflect.TypeFor[*timestamppb.Timestamp]()
	Decimal      = reflect.TypeFor[decimal.Decimal]()
	Copiable     = reflect.TypeFor[iface.Copiable]()
	LanguageTag  = reflect.TypeFor[language.Tag]()
	ClosedEnum   = reflect.TypeFor[enums.ClosedEnum]()
	ULID         = reflect.TypeFor[ulid.ULID]()
	StructpbPtr  = reflect.TypeFor[*structpb.Struct]()
	Maybe        = maybe.IfaceType
	Required     = validate.RequiredIfaceType
)

var (
	// ErrNoTypeAttributes indicates that no attributes were found for a Postgres type.
	ErrNoTypeAttributes = errors.New("no attributes for type")
)

// CheckAlignment checks whether two types are well aligned fieldwise.
func CheckAlignment(t1, t2 reflect.Type) error {
	fields1 := RelevantFields(t1, nil)
	fields2 := RelevantFields(t2, nil)
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
func RelevantFields(t reflect.Type, filter func(reflect.StructField) bool) []reflect.StructField {
	var fields []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath == "" && !f.Anonymous {
			if filter == nil || filter(f) {
				fields = append(fields, f)
			}
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

	if len(fields) == 0 {
		return nil, serr.Wrap("", ErrNoTypeAttributes, serr.String("type", typname))
	}

	return fields, nil
}
