package crushftp

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type version struct {
	Major int64
	Minor int64
	Patch int64
}

func (v *version) Parse(s string) error {
	part := strings.Split(s, ".")

	var major, minor, patch int64
	var err error

	if major, err = strconv.ParseInt(part[0], 10, 64); err != nil {
		return err
	}

	if minor, err = strconv.ParseInt(part[1], 10, 64); err != nil {
		return err
	}
	if patch, err = strconv.ParseInt(part[2], 10, 64); err != nil {
		return err
	}

	v.Major = major
	v.Minor = minor
	v.Patch = patch

	return nil
}

func extractVersionFromString(s string) string {
	return strings.TrimPrefix(s, "Version ")
}

func (v *version) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	d.DecodeElement(&s, &start)

	s = extractVersionFromString(s)

	err := v.Parse(s)

	if err != nil {
		return err
	}

	return nil
}

func (v version) String() string {
	return fmt.Sprintf(
		"%d.%d.%d",
		v.Major,
		v.Minor,
		v.Patch,
	)
}
