package protoutil

import (
	"io"
	"os"

	"google.golang.org/protobuf/proto"
)

func UnmarshalFromFile(m proto.Message, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	bs, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	if err := proto.Unmarshal(bs, m); err != nil {
		return err
	}
	return nil
}
