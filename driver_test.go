package time_test

import (
	"testing"

	"github.com/lonemage/time"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLocalDriver(t *testing.T) {
	Convey("TestFakeDriver", t, func() {
		ns := ""
		d := time.NewLocalDriver()

		n1, err := d.Sync(ns)
		So(err, ShouldBeNil)
		So(n1.Local(), ShouldEqual, time.Local)

		err = time.Set(time.Now().Add(time.Hour))
		So(err, ShouldBeNil)

		n2, err := d.Sync(ns)
		So(err, ShouldBeNil)
		So(n2.Sub(n1).Milliseconds(), ShouldEqual, 0)

		err = time.In(time.UTC)
		So(err, ShouldBeNil)

		n3, err := d.Sync(ns)
		So(err, ShouldBeNil)
		So(n3.Location(), ShouldEqual, time.UTC)
	})
}
