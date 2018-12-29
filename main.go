package main

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	pb "github.com/snxq/decode_bazel_build_event/BuildEventStream"

	"github.com/golang/protobuf/proto"
)

var file = flag.String("file", "", "bazel build event file")

var errInvalidVarint = errors.New("invalid varint32 encountered")

func main() {
	flag.Parse()
	if *file == "" {
		fmt.Println("--file must be set!")
	}

	event := &pb.BuildEvent{}

	var outFiles []string

	in, _ := os.Open(*file)
	defer in.Close()
	for {
		if _, err := ReadDelimited(in, event); err != nil {
			if err != io.EOF {
				fmt.Println(err)
				return
			} else {
				break
			}
		}
		payload := reflect.ValueOf(event.Payload)
		payloadType := reflect.TypeOf(event.Payload)
		isFiles := payloadType.ConvertibleTo(reflect.TypeOf(&pb.BuildEvent_NamedSetOfFiles{}))
		if isFiles {
			namedSetOfFiles := payload.Interface().(*pb.BuildEvent_NamedSetOfFiles)
			for _, f := range namedSetOfFiles.NamedSetOfFiles.Files {
				v := reflect.ValueOf(f.File)
				t := reflect.TypeOf(f.File)
				b := t.ConvertibleTo(reflect.TypeOf(&pb.File_Uri{}))
				if b {
					c := v.Interface().(*pb.File_Uri)
					if strings.HasPrefix(c.Uri, "file:///private") {
						outFiles = append(outFiles, c.Uri[7:])
						sum := Sha256Sum(c.Uri[7:])
						if sum != nil {
							fmt.Printf("%x\t%s\n", sum, c.Uri[7:])
						}
					}
				}
			}
		}
	}
}

func ReadDelimited(r io.Reader, m proto.Message) (n int, err error) {
	// Per AbstractParser#parsePartialDelimitedFrom with
	// CodedInputStream#readRawVarint32.
	var headerBuf [binary.MaxVarintLen32]byte
	var bytesRead, varIntBytes int
	var messageLength uint64
	for varIntBytes == 0 { // i.e. no varint has been decoded yet.
		if bytesRead >= len(headerBuf) {
			return bytesRead, errInvalidVarint
		}
		// We have to read byte by byte here to avoid reading more bytes
		// than required. Each read byte is appended to what we have
		// read before.
		newBytesRead, err := r.Read(headerBuf[bytesRead : bytesRead+1])
		if newBytesRead == 0 {
			if err != nil {
				return bytesRead, err
			}
			// A Reader should not return (0, nil), but if it does,
			// it should be treated as no-op (according to the
			// Reader contract). So let's go on...
			continue
		}
		bytesRead += newBytesRead
		// Now present everything read so far to the varint decoder and
		// see if a varint can be decoded already.
		messageLength, varIntBytes = proto.DecodeVarint(headerBuf[:bytesRead])
	}

	messageBuf := make([]byte, messageLength)
	newBytesRead, err := io.ReadFull(r, messageBuf)
	bytesRead += newBytesRead
	if err != nil {
		return bytesRead, err
	}

	return bytesRead, proto.Unmarshal(messageBuf, m)
}

func Sha256Sum(file string) (result []byte) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return
	}
	return h.Sum(nil)
}
