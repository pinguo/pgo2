package util

import (
    "testing"
)

type mockTestMergeFieldA struct {
    Name string
    Id   int
    desc string
}

type mockTestMergeFieldB struct {
    Name string
    Id   string
    desc string
    data string
}

func TestSTMergeSame(t *testing.T) {
    t.Run("param1 not pointer", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        a := "a"
        STMergeSame(a, a)
    })

    t.Run("param2 not pointer", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        a := &mockTestMergeFieldA{Name: "name1", Id: 1, desc: "desc1"}
        b := "a"
        STMergeSame(a, b)
    })

    t.Run("not same type", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        a := &mockTestMergeFieldA{Name: "name1", Id: 1, desc: "desc1"}
        b := &mockTestMergeFieldB{}
        STMergeSame(a, b)
    })

    t.Run("normal", func(t *testing.T) {
        a := &mockTestMergeFieldA{Name: "name1", Id: 1, desc: "desc1"}
        b := &mockTestMergeFieldA{Name: "name2", Id: 2, desc: "desc2"}

        STMergeSame(a, b)

        if a.Name != b.Name {
            t.Fatal(`a.Name!=b.Name`)
        }

        if a.Id != b.Id {
            t.Fatal(`a.Id!=b.Id`)
        }

        if a.desc != a.desc {
            t.Fatal(`a.desc!=a.desc`)
        }
    })

}

func TestSTMergeField(t *testing.T) {
    t.Run("param1 not pointer", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        a := "a"
        STMergeField(a, a)
    })

    t.Run("not struct", func(t *testing.T) {
        defer func() {
            if err := recover(); err != nil {
                return
            }
            t.FailNow()
        }()
        a := "a"
        STMergeField(&a, a)
    })

    t.Run("normal", func(t *testing.T) {
        a := &mockTestMergeFieldA{Name: "name1", Id: 1, desc: "desc1"}
        b := &mockTestMergeFieldB{Name: "name2", Id: "id2", desc: "desc2", data: "data2"}
        STMergeField(a, b)
        if a.Name != b.Name {
            t.Fatal(`a.Name!=b.Name`)
        }

        if a.Id != a.Id {
            t.Fatal(`a.Id!=a.Id`)
        }

        if a.desc != a.desc {
            t.Fatal(`a.desc!=a.desc`)
        }
    })

}
