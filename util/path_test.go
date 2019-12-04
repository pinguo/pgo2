package util

import (
    "testing"
)

func TestCleanPath(t *testing.T) {
    t.Run("empty", func(t *testing.T) {
        if CleanPath("") != "/" {
            t.FailNow()
        }
    })

    if CleanPath("../../aa/b/./") != "/aa/b/" {
        t.FailNow()
    }
}
