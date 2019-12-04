package maxmind

import (
    "fmt"
    "net"
    "reflect"
    "testing"

    "github.com/agiledragon/gomonkey"
    "github.com/oschwald/maxminddb-golang"
)

func TestNew(t *testing.T) {

    t.Run("normal", func(t *testing.T) {
        config := make(map[string]interface{})
        config["countryFile"] = "testCountry"
        config["cityFile"] = "testCity"

        patches := gomonkey.ApplyFunc(maxminddb.Open, func(file string) (*maxminddb.Reader, error) {
            t.Log("mock maxminddb.Open")
            return &maxminddb.Reader{}, nil
        })
        defer patches.Reset()

        _, err := New(config)
        if err != nil {
            t.FailNow()
        }
    })

    t.Run("db empty", func(t *testing.T) {
        _, err := New(nil)
        if err == nil {
            t.FailNow()
        }
    })
}

func TestClient_GeoByIp(t *testing.T) {
    config := make(map[string]interface{})
    config["countryFile"] = "testCountry"
    config["cityFile"] = "testCity"

    patches := gomonkey.ApplyFunc(maxminddb.Open, func(file string) (*maxminddb.Reader, error) {
        t.Log("mock maxminddb.Open")
        return &maxminddb.Reader{}, nil
    })
    defer patches.Reset()
    iM, _ := New(config)
    m := iM.(*Client)
    t.Run("invalid arg", func(t *testing.T) {

        _, err := m.GeoByIp("127.0.0.1", []string{})
        if err == nil {
            t.FailNow()
        }
    })

    t.Run("lookup return err", func(t *testing.T) {
        var r *maxminddb.Reader
        patches := gomonkey.ApplyMethod(reflect.TypeOf(r), "Lookup", func(_ *maxminddb.Reader, _ net.IP, result interface{}) error {
            t.Log("mock Client.Lookup")
            return fmt.Errorf("mock err")

        })
        defer patches.Reset()

        _, err := m.GeoByIp("127.0.0.1")
        if err == nil {
            t.FailNow()
        }
    })

    t.Run("normal", func(t *testing.T) {
        mData := make(map[string]interface{})
        mData["continent"] = map[string]interface{}{"names": map[string]interface{}{"en": "continent", "zh-cn": "continent-cn"}}
        mData["country"] = map[string]interface{}{"names": map[string]interface{}{"en": "china", "zh-cn": "china-cn"}, "iso_code": "CN"}
        mData["city"] = map[string]interface{}{"names": map[string]interface{}{"en": "chengdou", "zh-cn": "chengdou-cn"}}
        sub := map[string]interface{}{"names": map[string]interface{}{"en": "subdivisions", "zh-cn": "subdivisions-cn"}}
        subs := make([]interface{}, 0)
        subs = append(subs, sub)
        mData["subdivisions"] = subs

        var r *maxminddb.Reader
        patches := gomonkey.ApplyMethod(reflect.TypeOf(r), "Lookup", func(_ *maxminddb.Reader, _ net.IP, result interface{}) error {
            //result = mData
            rvResult := reflect.ValueOf(result)
            rvResult.Elem().Set(reflect.ValueOf(mData))
            t.Log("mock Client.Lookup")
            return nil

        })
        defer patches.Reset()

        geo, err := m.GeoByIp("127.0.0.1", 1, "zh-cn")
        if err != nil {
            t.Fatal(`err !=nil`)
        }

        if geo.Code != "CN" {
            t.Fatal(`geo.Code != "CN"`)
        }

        if geo.Country != "china" {
            t.Fatal(`geo.Country != "china"`)
        }

        if geo.I18n.Country != "china-cn" {
            t.Fatal(`geo.I18n.Country != "china-cn"`)
        }

        if geo.City != "chengdou" {
            t.Fatal(`geo.City != "chengdou"`)
        }

        if geo.I18n.City != "chengdou-cn" {
            t.Fatal(`geo.I18n.City != "chengdou-cn"`)
        }

        if geo.Province != "subdivisions" {
            t.Fatal(`geo.Province != "subdivisions"`)
        }

        if geo.I18n.Province != "subdivisions-cn" {
            t.Fatal(`geo.I18n.Province != "subdivisions-cn"`)
        }

        if geo.Continent != "continent" {
            t.Fatal(`geo.Continent != "continent"`)
        }

        if geo.I18n.Continent != "continent-cn" {
            t.Fatal(`geo.I18n.Continent != "continent-cn"`)
        }

    })
}
