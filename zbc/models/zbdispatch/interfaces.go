package zbdispatch

import (
	"encoding/binary"
	"io"
)

// SimpleBinaryEncoding interface is abstraction over all SimpleBinaryEncoding Messages.
type SimpleBinaryEncoding interface {
	Encode(writer io.Writer, order binary.ByteOrder, doRangeCheck bool) error
	Decode(reader io.Reader, order binary.ByteOrder, actingVersion uint16, blockLength uint16, doRangeCheck bool) error
}
