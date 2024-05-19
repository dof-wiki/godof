package binary_helper

import "encoding/binary"

func DecryptCrc(data []byte, crc uint32) []byte {
	key := uint32(0x81A79011)
	xor := crc ^ key
	intNum := len(data) / 4

	// Create key_all
	keyAll := make([]byte, 4*intNum)
	for i := 0; i < intNum; i++ {
		binary.LittleEndian.PutUint32(keyAll[i*4:], xor)
	}

	// Convert data and key_all to uint32 slices for XOR operation
	dataInts := make([]uint32, intNum)
	keyInts := make([]uint32, intNum)
	for i := 0; i < intNum; i++ {
		dataInts[i] = binary.LittleEndian.Uint32(data[i*4:])
		keyInts[i] = binary.LittleEndian.Uint32(keyAll[i*4:])
	}

	// Perform XOR
	valueXoredAll := make([]uint32, intNum)
	for i := 0; i < intNum; i++ {
		valueXoredAll[i] = dataInts[i] ^ keyInts[i]
	}

	// Define masks
	mask1 := uint32(0b00000000_00000000_00000000_00111111)
	mask2 := uint32(0b11111111_11111111_11111111_11000000)

	// Apply masks and shift
	result := make([]byte, 4*intNum)
	for i := 0; i < intNum; i++ {
		value1 := valueXoredAll[i] & mask1
		value2 := valueXoredAll[i] & mask2
		value := (value1 << 26) | (value2 >> 6)
		binary.LittleEndian.PutUint32(result[i*4:], value)
	}

	return result
}
