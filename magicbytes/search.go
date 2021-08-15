package magicbytes

import (
	"bytes"
	"context"
	"fmt"
	"os"
)

func searchMetasInAFile(ctx context.Context, walkingPath string, metaslice *[]*Meta) (*Meta, error) {
	select {
	case <-ctx.Done():
		//time.Sleep(time.Second * 1)
		if status, ok := FileStatus[walkingPath]; ok && status {
			return nil, ctx.Err()
		}
	default:
		//fmt.Print("*default*")
	}

	f, err := os.Open(walkingPath)
	if err != nil {
		fmt.Println("2")
		//log.Println(err, walkingPath, &MyCustomError{ErrDes: "cant open file error=>path", StatusCode: 953})
		return nil, err

	}

	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()

	stat, err := f.Stat()
	if err != nil {

		//log.Println(err, walkingPath, &MyCustomError{ErrDes: "cant stat file error=>path", StatusCode: 952})
		return nil, err
	}

	fileSize := stat.Size()

	for _, meta := range *metaslice {
		select {
		case <-ctx.Done():
			if status, ok := FileStatus[walkingPath]; ok && status {
				return nil, ctx.Err()
			}
		default:
		}

		if meta == nil {
			continue
		}

		//Offset
		lenBytes := len(meta.Bytes)
		if meta.Offset+int64(lenBytes) > fileSize {

			continue
		}

		_, e := f.Seek(meta.Offset, 0)
		if e != nil {

			//log.Println(err, walkingPath, &MyCustomError{ErrDes: "cant seek on file error=>", StatusCode: 951})
			continue
		}

		mb := make([]byte, lenBytes)

		n, err := f.Read(mb)
		if err != nil {
			//log.Println(err, walkingPath, &MyCustomError{ErrDes: "cant read on file error=>", StatusCode: 949})

			continue
		}

		if n == lenBytes && bytes.Equal(mb, meta.Bytes) {
			return meta, nil
		}

	}

	return nil, nil
}
