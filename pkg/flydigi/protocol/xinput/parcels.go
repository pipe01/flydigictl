package xinput

import (
	"bytes"
	"math"
)

const (
	perPackageMaxLength  = 15
	perPackageDataLength = 10
)

const (
	startKMPC       = 37
	writeDataKMPC   = 36
	specifyConfigId = 160
	writeLed        = 41
	startWriteLed   = 42
)

func splitParcels(start func(parcelCount int) []byte, innerCmd, configID byte, data []byte) [][]byte {
	parcelCount := int(math.Ceil(float64(len(data)) / perPackageDataLength))

	parcels := make([][]byte, 0)

	header := make([]byte, perPackageMaxLength)
	copy(header, start(parcelCount))
	parcels = append(parcels, crcData(header))

	var parcel bytes.Buffer

	for i := 0; i < parcelCount; i++ {
		parcel.Reset()

		parcel.WriteByte(165)
		parcel.WriteByte(innerCmd)

		for j := 0; j < perPackageDataLength; j++ {
			dataIdx := i*perPackageDataLength + j

			if i < parcelCount-1 {
				parcel.WriteByte(data[dataIdx])
			} else {
				if dataIdx < len(data) {
					parcel.WriteByte(data[dataIdx])
				} else {
					parcel.WriteByte(0)
				}
			}
		}

		parcel.WriteByte(160)
		parcel.WriteByte(byte(i))
		parcel.WriteByte(0)

		// buf.Write(crcData(parcel.Bytes()))
		parcels = append(parcels, crcData(parcel.Bytes()))
	}

	return parcels
}

// FormatterWriteFullData
func getConfigDataParcels(data []byte, configID byte) [][]byte {
	return splitParcels(func(parcelCount int) []byte {
		return []byte{
			165,
			startKMPC,
			byte(parcelCount),
			specifyConfigId,
			configID,
		}
	}, writeDataKMPC, configID, data)
}

// FormatterWriteFullLEDData
func getLEDConfigDataParcels(data []byte, configID byte) [][]byte {
	return splitParcels(func(parcelCount int) []byte {
		return []byte{
			165,
			startWriteLed,
			configID,
			byte(parcelCount),
		}
	}, writeLed, configID, data)
}

func crcData(pkg []byte) []byte {
	var sum byte

	out := make([]byte, len(pkg))
	copy(out, pkg)

	for i, v := range pkg {
		if i == len(pkg)-1 {
			out[len(pkg)-1] = sum
		} else {
			sum += v
		}
	}

	return out
}
