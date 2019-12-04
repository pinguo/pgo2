package util

import (
    "testing"
)

func TestNewHashRing(t *testing.T) {
    var obj interface{}
    obj = NewHashRing("127.0.0.1:6379", "127.0.0.1:6379", 16, HashSha1Crc32)
    if _, ok := obj.(*HashRing); ok == false {
        t.FailNow()
    }
}

func TestHashRing_AddNode(t *testing.T) {
    h := NewHashRing("127.0.0.1:6379", "127.0.0.1:6379", 32, HashSha1Crc32)
    h.AddNode("127.0.0.1:6379")
    if h.weights["127.0.0.1:6379"] != 3 {
        t.FailNow()
    }
}

func TestHashRing_GetNode(t *testing.T) {
    h := NewHashRing("127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381", 32, HashSha1Crc32)
    k := "test"
    if h.GetNode(k) != "127.0.0.1:6380" {
        t.FailNow()
    }
}

func TestHashRing_DelNode(t *testing.T) {
    h := NewHashRing("127.0.0.1:6379", "127.0.0.1:6379", 32, HashSha1Crc32)
    h.DelNode("127.0.0.1:6379")
    if len(h.items) != 0 {
        t.FailNow()
    }
}
