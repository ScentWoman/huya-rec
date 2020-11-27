package huya

import (
	"fmt"
	"testing"
)

func TestAPI(t *testing.T) {
	info, e := GetInfo("https://www.huya.com/5801xiaowei")
	if e != nil {
		fmt.Println(e)
		t.FailNow()
	}

	fmt.Println(info.Name, "is on:", info.On, info.Title)
	fmt.Println(info.Stream)
	t.Fail()
}
