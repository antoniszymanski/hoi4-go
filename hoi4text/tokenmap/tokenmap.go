package tokenmap

import (
	"bufio"
	"compress/flate"
	"encoding/binary"
	"errors"
	"io"
	"maps"
	"slices"
	"unsafe"
)

func Encode(w io.Writer, m map[uint16]string) (n int, err error) {
	fw, _ := flate.NewWriter(w, flate.BestCompression)
	write := func(b []byte) error {
		written, writeErr := fw.Write(b)
		n += written
		return writeErr
	}
	var writeBuf [binary.MaxVarintLen64]byte
	writeUvarint := func(x uint64) error {
		i := binary.PutUvarint(writeBuf[:], x)
		return write(writeBuf[:i])
	}
	writeUint16 := func(x uint16) error {
		binary.LittleEndian.PutUint16(writeBuf[:], x)
		return write(writeBuf[:2])
	}
	writeString := func(s string) error {
		return write(unsafe.Slice(unsafe.StringData(s), len(s)))
	}
	var allValuesLen uint64
	for key, value := range m {
		if value == "" {
			delete(m, key)
			continue
		}
		allValuesLen += uint64(len(value))
	}
	if err = writeUvarint(allValuesLen); err != nil {
		return
	}
	keys := make([]uint16, 0, len(m))
	keys = slices.AppendSeq(keys, maps.Keys(m))
	slices.Sort(keys)
	if err = writeUvarint(uint64(len(keys))); err != nil {
		return
	}
	for _, key := range keys {
		value := m[key]
		if err = writeUint16(key); err != nil {
			return
		}
		if err = writeUvarint(uint64(len(value))); err != nil {
			return
		}
		if err = writeString(value); err != nil {
			return
		}
	}
	err = fw.Close()
	return
}

func Decode(r io.Reader) (map[uint16]string, error) {
	fr := flate.NewReader(r)
	br := bufio.NewReader(fr)
	allValuesLen, err := binary.ReadUvarint(br)
	if err != nil {
		return nil, err
	}
	mapLen, err := binary.ReadUvarint(br)
	if err != nil {
		return nil, err
	}
	allValues := make([]byte, allValuesLen)
	m := make(map[uint16]string, mapLen)
	keyBuf := make([]byte, 2)
	for range mapLen {
		if _, err := io.ReadFull(br, keyBuf); err != nil {
			return nil, err
		}
		key := binary.LittleEndian.Uint16(keyBuf)
		valueLen, err := binary.ReadUvarint(br)
		if err != nil {
			return nil, err
		} else if uint64(len(allValues)) < valueLen {
			return nil, errors.New("value length exceeds remaining buffer capacity")
		}
		valueBuf := allValues[:valueLen]
		allValues = allValues[valueLen:]
		if _, err := io.ReadFull(br, valueBuf); err != nil {
			return nil, err
		}
		value := unsafe.String(unsafe.SliceData(valueBuf), len(valueBuf))
		m[key] = value
	}
	if err := fr.Close(); err != nil {
		return nil, err
	}
	return m, nil
}
