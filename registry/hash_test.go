package registry

import "testing"
import "fmt"

func TestHash(t *testing.T) {

	b := &blob{
		content: []byte("AAAAAAAAAAAAAAAA"),
	}

	b.hash()

	t.Log("testing")

	t.Log(fmt.Sprintf("this should be something%+v\n", b))

    // total := Sum(5, 5)
    // if total != 10 {
    //    t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
    // }
}