package crushftp

import (
	"encoding/xml"
	"time"
)

type unixDateTime struct {
	time.Time
}

func (c *unixDateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(time.UnixDate, v)
	if err != nil {
		return err
	}
	*c = unixDateTime{parse}
	return nil
}
