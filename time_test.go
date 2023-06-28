package time_test

import (
	"fmt"
	"testing"

	"github.com/lonemage/time"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLocaltime(t *testing.T) {
	Convey("TestLocaltime", t, func() {
		err := time.Conf(time.WithLocal())
		So(err, ShouldBeNil)

		tim, err := time.Parse(time.LayoutYmdHms, "2020-01-01 00:00:00")
		So(err, ShouldBeNil)
		So(tim.Location(), ShouldEqual, time.Local)
		y, m, d := tim.Date()
		So(y, ShouldEqual, 2020)
		So(m, ShouldEqual, 1)
		So(d, ShouldEqual, 1)

		du := time.Since(tim)
		So(du, ShouldBeGreaterThan, 0)

		n := time.Now()
		n3 := time.Now()
		time.Pass(time.Hour)
		n2 := time.Now()
		fmt.Println(n, n2, n3)
		So(n2.Sub(n), ShouldAlmostEqual, time.Hour, time.Second)

		err = time.SetLocalLocation(time.UTC)
		So(err, ShouldBeNil)
		So(time.LocalLocation(), ShouldEqual, time.UTC)

	})
}

func TestXtime(t *testing.T) {
	Convey("TestXtime", t, func() {
		// time.Conf()

		// n := time.Now()
		// n3 := time.Now()
		// time.Pass(time.Hour)
		// n2 := time.Now()
		// fmt.Println(n, n2, n3)
		// So(n2.Sub(n), ShouldAlmostEqual, time.Hour, time.Second)

		// time.In(time.UTC)
		// So(time.GetLocation(), ShouldEqual, time.UTC)
	})
}
