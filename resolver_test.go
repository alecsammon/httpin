package httpin

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type MissingDecoderName struct {
	Birthday time.Time `in:"form=birthday;decoder;required"`
}

type UnknownDecoderName struct {
	Birthday time.Time `in:"form=birthday;decoder=decodeBirthday"`
}

func TestFieldResolver(t *testing.T) {
	Convey("Build resolver tree", t, func() {
		resolver, err := buildResolverTree(reflect.TypeOf(ProductQuery{}))
		So(err, ShouldBeNil)
		So(resolver, ShouldNotBeNil)
		r, _ := http.NewRequest("GET", "https://example.com", nil)
		r.Form = make(url.Values)
		r.Form.Set("created_at", time.Now().Format(time.RFC3339))
		r.Form.Set("color", "red")
		r.Form.Set("is_soldout", "true")
		r.Form.Add("sort_by", "id")
		r.Form.Add("sort_by", "quantity")
		r.Form.Add("sort_desc", "0")
		r.Form.Add("sort_desc", "true")
		r.Form.Set("page", "1")
		r.Form.Set("per_page", "20")
		r.Header.Set("x-api-token", "cad979df-5e40-4bfd-b31d-f870ca2c14ea")
		rv, err := resolver.resolve(r)
		So(err, ShouldBeNil)
		So(rv.Elem().Interface(), ShouldHaveSameTypeAs, ProductQuery{})
		bs, _ := json.Marshal(rv.Interface())
		t.Logf("ProductQuery: %s\n", bs)
	})

	Convey("Parse decoder directive in struct tags, missing decoder name", t, func() {
		_, err := buildResolverTree(reflect.TypeOf(MissingDecoderName{}))
		So(err, ShouldNotBeNil)
		So(errors.Is(err, ErrMissingDecoderName), ShouldBeTrue)
	})

	Convey("Parse decoder directive in struct tags, unknown name of decoder", t, func() {
		_, err := buildResolverTree(reflect.TypeOf(UnknownDecoderName{}))
		So(err, ShouldNotBeNil)
		So(errors.Is(err, ErrDecoderNotFound), ShouldBeTrue)
	})
}
