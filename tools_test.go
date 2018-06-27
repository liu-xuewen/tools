package tools

import (
	"fmt"
)

import (
	"testing"
)

type data struct {
	input string
	want  string
	id    int
}

/*----------------------
	UnitTest
----------------------*/
func TestStringToTime(t *testing.T) {
	datas := []data{
		{"2016-02-25 09:17:00",
			"success",
			1},
		{"2016-02-25 09:17:0",
			"fail",
			2},
		{"2016-02-25 09:17:",
			"fail",
			3},
		{"2016-02-25 09:17:61",
			"fail",
			4},
		{"2016-02-25 09:17:a1",
			"fail",
			5},
	}

	for _, d := range datas {
		r, err := StringToTime(d.input)
		if err != nil && d.want != "fail" {
			t.Errorf("测试ID: %d,StringToTime(%s) -> %s ,wanted: %s, result: %s ",
				d.id, d.input, r.String(), d.want, "fail")
		} else if err ==nil && d.want != "success" {
			t.Errorf("测试ID: %d,StringToTime(%s) -> %s ,wanted: %s ,result: %s",
				d.id, d.input, r.String(), d.want, "success")
		}
	}
}

/*----------------------
	Benchmark
----------------------*/
func BenchmarkStringToTime(b *testing.B) {
	b.ResetTimer()
	
	for i:=0;i<b.N;i++ {
		StringToTime("2016-02-25 09:17:00")
	}
}
/*----------------------
	Example
----------------------*/

//Example code for RandString function
func ExampleRandString() {
	fmt.Println(RandString())
	// Output: mv9PnyYs8-jMXHAYcNYv6xKg0wqdwDklLt1F6RFkmBc=
}

//Example code for StringToTime function
func ExampleStringToTime() {
	t, _ := StringToTime("2016-02-25 09:17:00")
	fmt.Println(t.String())
	// Output: 2016-02-25 09:17:00 +0800 CST
}
