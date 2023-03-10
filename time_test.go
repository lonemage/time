package time_test

import (
	"fmt"
	"testing"

	"github.com/lonemage/time"
	. "github.com/smartystreets/goconvey/convey"
)

func TestXtime(t *testing.T) {
	Convey("TestXtime", t, func() {
		time.Conf()

		n := time.Now()
		n3 := time.Now()
		time.Pass(time.Hour)
		n2 := time.Now()
		fmt.Println(n, n2, n3)
		So(n2.Sub(n), ShouldAlmostEqual, time.Hour, time.Second)

		time.In(time.UTC)
		So(time.GetLocation(), ShouldEqual, time.UTC)
	})
}
