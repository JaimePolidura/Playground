package broadcast

func ZeroArray(bytes *[]byte) {
	for i := 0; i < len(*bytes); i++ {
		(*bytes)[i] = 0x00
	}
}
