package crushftp

import (
	"encoding/xml"
	"time"
)

type unixDateTime struct {
	time.Time
}

func (t *unixDateTime) Parse(s string) error {
	parsed, err := time.Parse(time.UnixDate, s)

	if err != nil {
		return err
	}

	t.Time = parsed

	return nil
}

func (c *unixDateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	d.DecodeElement(&s, &start)

	err := c.Parse(s)

	if err != nil {
		return err
	}

	return nil
}
