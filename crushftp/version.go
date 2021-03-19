package crushftp

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type version struct {
	Major  int64
	Minor  int64
	Bugfix int64
}

func (v *version) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	d.DecodeElement(&s, &start)

	plainVersion := strings.TrimPrefix(s, "Version ")
	versionParts := strings.Split(plainVersion, ".")

	var major, minor, bugfix int64
	var err error

	if major, err = strconv.ParseInt(versionParts[0], 10, 64); err != nil {
		return nil
	}

	if minor, err = strconv.ParseInt(versionParts[1], 10, 64); err != nil {
		return nil
	}
	if bugfix, err = strconv.ParseInt(versionParts[2], 10, 64); err != nil {
		return nil
	}

	*v = version{
		Major:  major,
		Minor:  minor,
		Bugfix: bugfix,
	}

	return nil
}

func (v version) String() string {
	return fmt.Sprintf(
		"%d.%d.%d",
		v.Major,
		v.Minor,
		v.Bugfix,
	)
}
