package citylib_test

import (
	"github.com/dimroc/urbanevents/cityservice/citylib"
	//. "github.com/dimroc/urbanevents/cityservice/utils"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewGeoEventFromTweet(t *testing.T) {
	Convey("Given a city and an anaconda.Tweet with a poi", t, func() {
		city := Fixture.GetCity()

		Convey("and an anaconda.Tweet with a poi", func() {
			tweet := Fixture.GetPoiTweet()

			Convey("it should create a geoevent", func() {
				geoevent, err := citylib.NewGeoEventFromTweet(city, tweet)

				So(err, ShouldBeNil)
				So(geoevent.LocationType, ShouldEqual, "poi")
				So(geoevent.Place, ShouldEqual, "Hill Country Chicken")
				So(geoevent.GeoJson.Type, ShouldEqual, "Polygon")
			})
		})

		Convey("and an anaconda.Tweet with a coordinate", func() {
			tweet := Fixture.GetCoordinateTweet()

			Convey("it should create a geoevent", func() {
				geoevent, err := citylib.NewGeoEventFromTweet(city, tweet)

				So(err, ShouldBeNil)
				So(geoevent.LocationType, ShouldEqual, "coordinate")
				So(geoevent.GeoJson.Type, ShouldEqual, "point")
			})
		})

		Convey("and an anaconda.Tweet with a video", func() {
			tweet := Fixture.GetVideoTweet()

			Convey("it should create a geoevent", func() {
				geoevent, err := citylib.NewGeoEventFromTweet(city, tweet)

				So(err, ShouldBeNil)
				So(geoevent.MediaType, ShouldEqual, "video")
				So(geoevent.ThumbnailUrl, ShouldContainSubstring, "thumb")
				So(geoevent.MediaUrl, ShouldEndWith, "mp4")
			})
		})

		Convey("and an anaconda.Tweet with an instagram link", func() {
			tweet := Fixture.GetInstagramTweet()

			Convey("it should create a text geoevent", func() {
				geoevent, err := citylib.NewGeoEventFromTweet(city, tweet)

				So(err, ShouldBeNil)
				So(geoevent.MediaType, ShouldEqual, "text")
				So(geoevent.ThumbnailUrl, ShouldBeEmpty)
				So(geoevent.MediaUrl, ShouldBeEmpty)
				So(geoevent.ExpandedUrl, ShouldContainSubstring, "instagram")

				// This geoevent is later enriched by the InstagramTweetEnricher
			})
		})
	})
}
