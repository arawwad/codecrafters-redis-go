package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (s *Stream) Append(entryKey RespType, entries []StreamEntry) (BulkString, SimpleError) {
	var timeStamp time.Time
	var sequenceNumber int
	emptyKey := BulkString("")

	str, ok := entryKey.Str()
	if !ok || str == "*-*" {
		return emptyKey, InvalidStreamKeyError
	}

	if str == "0-0" {
		return emptyKey, MinimumStreamKeyError
	}

	if str == "*" {
		timeStamp = time.Now()
		sequenceNumber = 0
	} else {
		idParts := strings.Split(str, "-")
		if len(idParts) != 2 {
			return emptyKey, InvalidStreamKeyError
		}

		if idParts[0] == "*" {
			timeStamp = time.Now()
		} else {
			timeStampPart, err := strconv.Atoi(idParts[0])
			if err != nil {
				return emptyKey, InvalidStreamKeyError
			}
			timeStamp = time.UnixMilli(int64(timeStampPart))
		}

		if timeStamp.Before(s.lastTimeStamp) {
			return emptyKey, InvalidOrderOfStreamKey
		}

		if idParts[1] == "*" {
			if timeStamp != s.lastTimeStamp {
				if timeStamp.UnixMilli() == 0 {
					sequenceNumber = 1
				} else {
					sequenceNumber = 0
				}
			} else {
				sequenceNumber = s.lastSequenceNumber + 1
			}
		} else {
			sequenceNumberPart, err := strconv.Atoi(idParts[1])
			if err != nil {
				return emptyKey, InvalidStreamKeyError
			} else {
				sequenceNumber = sequenceNumberPart
			}
		}

		if timeStamp == s.lastTimeStamp && sequenceNumber <= s.lastSequenceNumber {
			return emptyKey, InvalidOrderOfStreamKey
		}

		s.keys = append(s.keys, entryKey)
		s.lastSequenceNumber = sequenceNumber
		s.lastTimeStamp = timeStamp
		s.values[entryKey] = entries
	}

	return BulkString(fmt.Sprintf("%d-%d", timeStamp.UnixMilli(), sequenceNumber)), EmptySimpleError
}

type StreamEntry struct {
	Key   RespType
	Value RespType
}
