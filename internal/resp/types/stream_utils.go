package types

import (
	"strconv"
	"strings"
	"time"
)

type StreamEntryKey struct {
	TimeStamp      time.Time
	SequenceNumber int
	Value          RespType
}

func (s *Stream) Append(entryKey RespType, entries []StreamEntry) SimpleError {
	var timeStamp time.Time
	var sequenceNumber int

	str, ok := entryKey.Str()
	if !ok || str == "*-*" {
		return InvalidStreamKeyError
	}

	if str == "0-0" {
		return MinimumStreamKeyError
	}

	if str == "*" {
		timeStamp = time.Now()
		sequenceNumber = s.lastSequenceNumber + 1
	} else {

		idParts := strings.Split(str, "-")
		if len(idParts) != 2 {
			return InvalidStreamKeyError
		}

		if idParts[0] == "*" {
			timeStamp = time.Now()
		} else {
			timeStampPart, err := strconv.Atoi(idParts[0])
			if err != nil {
				return InvalidStreamKeyError
			}
			timeStamp = time.Unix(int64(timeStampPart), 0)
		}

		if idParts[1] == "*" {
			sequenceNumber = s.lastSequenceNumber + 1
		} else {
			sequenceNumberPart, err := strconv.Atoi(idParts[1])
			if err != nil {
				return InvalidStreamKeyError
			} else {
				sequenceNumber = sequenceNumberPart
			}
		}

		if timeStamp.Before(s.lastTimeStamp) || (timeStamp == s.lastTimeStamp && sequenceNumber <= s.lastSequenceNumber) {
			return InvalidOrderOfStreamKey
		}

		s.keys = append(s.keys, entryKey)
		s.lastSequenceNumber = sequenceNumber
		s.lastTimeStamp = timeStamp
		s.values[entryKey] = entries
	}

	return EmptySimpleError
}

type StreamEntry struct {
	Key   RespType
	Value RespType
}
