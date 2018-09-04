package config

import (
	"testing"
	"github.com/dispatchlabs/samples/common-util/helper"
)

func TestUpdateDisgoExecutable(t *testing.T) {
	RefreshDisgoExecutable(helper.GetDefaultDirectory())
}

func TestDeleteDB(t *testing.T) {
	ClearDB(helper.GetDefaultDirectory())
}

