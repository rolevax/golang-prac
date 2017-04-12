package messages

import (
	"bufio"
	"io"

	"github.com/gogo/protobuf/proto"
)

func ReadPb(reader *bufio.Reader, pb proto.Message) error {
	p, err := readBytes(reader)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(p, pb)
	if err != nil {
		return err
	}
	return nil
}

func readBytes(reader *bufio.Reader) ([]byte, error) {
	size, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	p := make([]byte, size)
	for i := 0; i < int(size); i++ {
		p[i], err = reader.ReadByte()
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func WritePb(conn io.Writer, pb proto.Message) error {
	bytes, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	conn.Write([]byte{byte(len(bytes))}) // assume size < 256
	conn.Write(bytes)
	return nil
}
