package xtime_test

import (
	"testing"
	"time"

	"github.com/lonemage/xtime"

	. "github.com/smartystreets/goconvey/convey"
)

func TestXtime(t *testing.T) {
	Convey("TestXtime", t, func() {
		xtime.Conf()

		xtime.Pass(xtime.Hour)
		So(xtime.Now().Sub(time.Now()), ShouldEqual, time.Hour)
	})
}
