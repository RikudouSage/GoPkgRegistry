package types

import (
	"database/sql/driver"
	"fmt"
	"net/url"

	"github.com/samber/lo"
)

type URI url.URL

func (receiver *URI) UnmarshalText(text []byte) error {
	return receiver.Scan(string(text))
}

func (receiver *URI) MarshalText() (text []byte, err error) {
	return []byte(receiver.String()), nil
}

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

func (receiver *URI) Value() (driver.Value, error) {
	return receiver.String(), nil
}

func (receiver *URI) AsGoURL() *url.URL {
	return (*url.URL)(receiver)
}

func (receiver *URI) String() string {
	return lo.Must(url.PathUnescape(receiver.AsGoURL().String()))
}
