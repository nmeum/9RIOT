package main

import (
	"bytes"
	"github.com/Harvey-OS/ninep/protocol"
)

// Maps strings written by the client to the control socket to
// server replies. Every test function needs an entry in this table.
var ctlcmds = map[string]ServerReply{
	"rversion_success":          {RversionSuccess, protocol.Tversion},
	"rversion_unknown":          {RversionUnknown, protocol.Tversion},
	"rversion_msize_too_big":    {RversionMsizeTooBig, protocol.Tversion},
	"rversion_invalid":          {RversionInvalidVersion, protocol.Tversion},
	"rversion_invalid_len":      {RversionInvalidLength, protocol.Tversion},
	"rversion_version_too_long": {RversionVersionTooLong, protocol.Tversion},

	"rattach_success":     {RattachSuccess, protocol.Tattach},
	"rattach_invalid_len": {RattachInvalidLength, protocol.Tattach},

	"rstat_success": {RstatSuccess, protocol.Tstat},
}

// Replies with the msize and version send by the client. This should always be
// successfully parsed by the client implementation.
func RversionSuccess(b *bytes.Buffer) error {
	TMsize, TVersion, t, err := protocol.UnmarshalTversionPkt(b)
	if err != nil {
		return err
	}

	protocol.MarshalRversionPkt(b, t, TMsize, TVersion)
	return nil
}

// Replies with the version string "unknown".
//
// From version(5):
//   If the server does not understand the client's version
//   string, it should respond with an Rversion message (not
//   Rerror) with the version string the 7 characters
//   ``unknown''.
//
// The client should therefore be able to parse this message but should
// return an error.
func RversionUnknown(b *bytes.Buffer) error {
	TMsize, _, t, err := protocol.UnmarshalTversionPkt(b)
	if err != nil {
		return err
	}

	protocol.MarshalRversionPkt(b, t, TMsize, "unknown")
	return nil
}

// Replies with an msize value equal to the one send by the client plus
// one.
//
// From version(5):
//   The server responds with its own maximum, msize, which must
//   be less than or equal to the client's value.
//
// The client should therefore be able to parse this message but should
// return an error.
func RversionMsizeTooBig(b *bytes.Buffer) error {
	TMsize, TVersion, t, err := protocol.UnmarshalTversionPkt(b)
	if err != nil {
		return err
	}

	protocol.MarshalRversionPkt(b, t, TMsize+1, TVersion)
	return nil
}

// Replies with an invalid version string.
//
// From version(5):
//  After stripping any such period-separated suffix, the server is
//  allowed to respond with a string of the form 9Pnnnn, where nnnn is
//  less than or equal to the digits sent by the client.
//
// Depending on the client implementation the client may not be able to
// parse the R-message. In any case it should return an error since the
// he shouldn't support this version of the 9P protocol.
func RversionInvalidVersion(b *bytes.Buffer) error {
	TMsize, _, t, err := protocol.UnmarshalTversionPkt(b)
	if err != nil {
		return err
	}

	protocol.MarshalRversionPkt(b, t, TMsize, "9P20009P2000")
	return nil
}

// Replies with an length field in the version R-message.
//
// The value of the size field is one byte too short. Thereby causing
// the length field of the version string to report a string size that
// would exceeds the size of the packet itself.
//
// The client should not be able to parse this R-message.
func RversionInvalidLength(b *bytes.Buffer) error {
	TMsize, TVersion, t, err := protocol.UnmarshalTversionPkt(b)
	if err != nil {
		return err
	}

	protocol.MarshalRversionPkt(b, t, TMsize, TVersion)

	var len uint64 = uint64(b.Len()) - 1
	copy(b.Bytes(), []byte{uint8(len)})

	return nil
}

// Replies with a version string that is one byte longer than the
// longest valid version string `unknown`.
//
// The client will most likely be able to parse this messages but might
// reject it if it exceeds a statically allocated buffer.
func RversionVersionTooLong(b *bytes.Buffer) error {
	TMsize, _, t, err := protocol.UnmarshalTversionPkt(b)
	if err != nil {
		return err
	}

	protocol.MarshalRversionPkt(b, t, TMsize, "12345678")
	return nil
}

// Successfully attaches the client. Replying with a vaild qid. The
// client should be able to parse this successfully.
func RattachSuccess(b *bytes.Buffer) error {
	_, _, _, _, t, err := protocol.UnmarshalTattachPkt(b)
	if err != nil {
		return err
	}

	protocol.MarshalRattachPkt(b, t, protocol.QID{})
	return nil
}

// Replies with a modified packet length field causing the packet length
// to be one byte shorter than neccessary. The client should not be able
// to parse this.
func RattachInvalidLength(b *bytes.Buffer) error {
	_, _, _, _, t, err := protocol.UnmarshalTattachPkt(b)
	if err != nil {
		return err
	}

	protocol.MarshalRattachPkt(b, t, protocol.QID{})

	var l uint64 = uint64(b.Len()) - 1
	copy(b.Bytes(), []byte{uint8(l), uint8(l >> 8), uint8(l >> 16), uint8(l >> 24)})

	return nil
}

// Replies with a valid Rstat message. The client must be able to parse
// this R-message.
func RstatSuccess(b *bytes.Buffer) error {
	_, t, err := protocol.UnmarshalTstatPkt(b)
	if err != nil {
		return err
	}

	qid := protocol.QID{
		Type:    23,
		Version: 2342,
		Path:    1337,
	}

	dir := protocol.Dir{
		Type:    9001,
		Dev:     5,
		QID:     qid,
		Mode:    protocol.DMDIR,
		Atime:   1494443596,
		Mtime:   1494443609,
		Length:  2342,
		Name:    "testfile",
		User:    "testuser",
		Group:   "testgroup",
		ModUser: "ken",
	}

	var B bytes.Buffer
	protocol.Marshaldir(&B, dir)

	protocol.MarshalRstatPkt(b, t, B.Bytes())
	return nil
}
