package npk

import "fmt"

const (
	NPK_MAGIC = "NeoplePack_Bill"
)

var (
	//NPK_FILENAME_DECORD_FLAG = bytes('puchikon@neople dungeon and fighter %s\x00' % ('DNF' * 73), encoding='ascii')
	NPK_FILENAME_DECORD_FLAG []byte
)

func init() {
	baseString := "puchikon@neople dungeon and fighter %s\\x00"
	repeatedString := ""
	for i := 0; i < 73; i++ {
		repeatedString += "DNF"
	}
	NPK_FILENAME_DECORD_FLAG = []byte(fmt.Sprintf(baseString, repeatedString))
}
