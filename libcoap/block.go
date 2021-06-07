package libcoap

/*
#cgo LDFLAGS: -lcoap-3-openssl
#include <coap3/coap.h>
*/
import "C"
import "math"
import "strings"
import "strconv"
import log "github.com/sirupsen/logrus"

type Block struct {
    NUM   int
    M     int
    SZX   int
}

const (
    BLOCK_COMPONENT_NUMBER = 3
    SZX_BIT_NUMBER         = 3
    M_BIT_NUMBER           = 1
    ADDITIONAL_SZX         = 4
    MORE_BLOCK             = 1
    LAST_BLOCK             = 0
    LARGEST_BLOCK_SIZE     = 6
)

/*
 * Convert value type String to type Block
 */
func StringToBlock(val string) (*Block) {
    if val == "" { return nil }

    blocks := strings.Split(val, "/")
	if len(blocks) != BLOCK_COMPONENT_NUMBER {
        log.Warn("Block2 is invalid. Please enter the value follow to format: NUM/M/SZX")
        return nil
	} else {
        num, err := strconv.Atoi(blocks[0])
        if err != nil {
            log.Warnf("Block2 NUM is not number: %+v", blocks[0])
            return nil
        }

        m, err := strconv.Atoi(blocks[1])
        if err != nil {
            log.Warnf("Block2 M is not number: %+v", blocks[1])
            return nil
        }

        size, err := strconv.ParseFloat(blocks[2], 64)
        if err != nil {
            log.Warnf("Block2 SZX is not number: %+v", blocks[2])
            return nil
        }
        szx := int(math.Log2(size) - ADDITIONAL_SZX)

		return &Block{ num, m, szx }
	}
}

/*
 * Convert value type Block to type String
 */
func (block *Block) ToString() string {
    ret := strconv.Itoa(block.NUM) + "/" + strconv.Itoa(block.M) + "/" + strconv.Itoa(1 << (uint(block.SZX) + ADDITIONAL_SZX))
    return ret
}

/*
 * Convert value type Int to type Block
 */
func IntToBlock(val int) (*Block) {
    if val >= 0 {
        if val == 0 {
            return &Block{}
        } else {
            num := val >> (SZX_BIT_NUMBER + M_BIT_NUMBER)
            m := (val - num << (SZX_BIT_NUMBER + M_BIT_NUMBER)) >> (SZX_BIT_NUMBER)
            szx := (val - num << (SZX_BIT_NUMBER + M_BIT_NUMBER) - m << SZX_BIT_NUMBER)
            return &Block{ num, m, szx }
        }
    }
    return nil
}

/*
 * Convert value type Block to type Int
 */
func (block *Block) ToInt() int {
    ret := block.NUM << (SZX_BIT_NUMBER + M_BIT_NUMBER) + block.M << SZX_BIT_NUMBER + block.SZX
    return ret
}