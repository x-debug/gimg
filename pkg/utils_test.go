package pkg

import (
	"os"
	"testing"
)

func TestCalcMd5(t *testing.T) {
	f, _ := os.CreateTemp("", "unittest")
	_, _ = f.WriteString("this is md5 string")
	md5Val, _ := CalcMd5File(f)
	t.Logf("File Md5 Value: %s", md5Val)
	_ = os.Remove(f.Name())
}
