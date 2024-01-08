package dinput

import (
	"bytes"
	"flydigi-linux/utils"
	"math"
)

const maxParcelLength = 10

func splitParcels(start func(parcelCount int) []byte, innerCmd, configID byte, data []byte) [][]byte {
	parcelCount := int(math.Ceil(float64(len(data)) / maxParcelLength))

	header := make([]byte, maxParcelLength+4)
	copy(header, start(parcelCount))

	var buf bytes.Buffer
	buf.Write(header)

	for i := 0; i < parcelCount; i++ {
		buf.WriteByte(5)
		buf.WriteByte(innerCmd)

		for j := 0; j < maxParcelLength; j++ {
			dataIdx := i*maxParcelLength + j

			if i < parcelCount-1 {
				buf.WriteByte(data[dataIdx])
			} else {
				if dataIdx < len(data) {
					buf.WriteByte(data[dataIdx])
				} else {
					buf.WriteByte(0)
				}
			}
		}

		buf.WriteByte(160)
		buf.WriteByte(byte(i))
	}

	return utils.SplitSlice(buf.Bytes(), maxParcelLength+4)
}

// FormatterWriteFullData
func getConfigDataParcels(data []byte, configID byte) [][]byte {
	return splitParcels(func(parcelCount int) []byte {
		return []byte{
			5,
			234,
			byte(parcelCount),
			160,
			configID,
		}
	}, 34, configID, data)
}

// FormatterWriteFullLEDData
func getLEDConfigDataParcels(data []byte, configID byte) [][]byte {
	return splitParcels(func(parcelCount int) []byte {
		return []byte{
			5,
			231,
			configID,
			byte(parcelCount),
		}
	}, 51, configID, data)
}
