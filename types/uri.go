package types

import (
	"database/sql/driver"
	"fmt"
	"net/url"

	"github.com/samber/lo"
)

// URI wraps [url.URL] with text and SQL serialization support.
type URI url.URL

// UnmarshalText parses a textual URI into receiver.
func (receiver *URI) UnmarshalText(text []byte) error {
	return receiver.Scan(string(text))
}

// MarshalText returns the receiver's unescaped textual representation.
func (receiver *URI) MarshalText() (text []byte, err error) {
	return []byte(receiver.String()), nil
}

// Scan implements the database/sql.Scanner interface. It accepts URI values
// stored as strings; scanning nil leaves receiver unchanged.
func (receiver *URI) Scan(src any) error {
	if src == nil {
		return nil
	}

	if str, ok := src.(string); ok {
		goUrl, err := url.Parse(str)
		if err != nil {
			return fmt.Errorf("failed parsing url: %w", err)
		}

		*receiver = URI(*goUrl)
		return nil
	}

	return fmt.Errorf("the value must be a string, %T given", src)
}

// Value implements [driver.Valuer].
func (receiver *URI) Value() (driver.Value, error) {
	return receiver.String(), nil
}

// AsGoURL returns receiver as a standard-library URL.
func (receiver *URI) AsGoURL() *url.URL {
	return (*url.URL)(receiver)
}

// String returns the URI string with percent-escaped path characters decoded.
func (receiver *URI) String() string {
	if receiver == nil {
		return ""
	}

	return lo.Must(url.PathUnescape(receiver.AsGoURL().String()))
}
