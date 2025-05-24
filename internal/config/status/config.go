package status

import "github.com/idelchi/godyl/internal/config/common"

type Status struct {
	common.Tracker `json:"-"    mapstructure:"-"`
	Tags           []string `json:"tags" mapstructure:"tags"`
}

func (s Status) ToCommon() common.Common {
	return common.Common{
		Tracker: s.Tracker,
	}
}
