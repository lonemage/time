package time_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"

	time "github.com/lonemage/xtime"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFakeDriver(t *testing.T) {
	Convey("TestFakeDriver", t, func() {
		d := time.NewFakeDriver()
		v1, err := d.Get()
		So(err, ShouldBeNil)

		_, n1, l1 := time.GetValue(v1)
		_, err = d.Set(time.MakeValue(0, n1+time.Hour.Nanoseconds(), l1))
		So(err, ShouldBeNil)

		v2, err2 := d.Get()
		_, n2, l2 := time.GetValue(v2)
		So(err2, ShouldBeNil)
		So(n2-n1, ShouldAlmostEqual, time.Hour, time.Second)
		So(l2, ShouldEqual, l1)
	})
}

var fsPort = rand.Int()%10000 + 3000

func FakeServer() error {
	c := time.NewFakeDriver()

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		// v,err:=c.Get()
		// if err!=nil{

		// }
	})
	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		c.Set(nil)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", fsPort), nil)
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
