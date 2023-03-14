package time_test

import (
	"math/rand"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var fsPort = rand.Int()%10000 + 3000

func FakeServer() error {
	return nil
}

func TestXDriver(t *testing.T) {
	go func() {
		if err := FakeServer(); err != nil {
			os.Exit(-1)
		}
	}()

	Convey("TestXDriver", t, func() {
		// err := time.Conf(time.WithRemote(fmt.Sprintf("localhost:%d", fsPort), true))
		// So(err, ShouldBeNil)

		// time.Pass(time.Hour)
		// So(time.Now().Sub(time.Now()), ShouldEqual, time.Hour)
	})
}
