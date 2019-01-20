package cache

import (
	"compress/gzip"
	"encoding/gob"
	"context"
	"bytes"
)

type DefaultMarshal struct {
	V         interface{}
	GzipLevel int
}

func (o *DefaultMarshal) Unmarshal(ctx context.Context, b []byte) error {

	//span, ctx := opentracing.StartSpanFromContext(ctx, "Unmarshal")
	//defer span.Finish()

	buf := bytes.NewReader(b)
	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(gzipReader)
	err = dec.Decode(o.V)
	if err != nil {
		return err
	}
	return nil
}

func (o *DefaultMarshal) Marshal(ctx context.Context) ([]byte, error) {

	//span, ctx := opentracing.StartSpanFromContext(ctx, "Marshal")
	//defer span.Finish()

	var buf bytes.Buffer
	gzipWriter, err := gzip.NewWriterLevel(&buf, o.GzipLevel)
	if err != nil {
		return nil, err
	}
	enc := gob.NewEncoder(gzipWriter)
	err = enc.Encode(o.V)
	if err != nil {
		return nil, err
	}
	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
