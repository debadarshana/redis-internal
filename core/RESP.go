package core

func readSimpleString(data []byte) (string, int) {
	// fmt.Println("Reading Simple String...")
	pos := 1
	for data[pos] != '\r' {
		pos++
	}
	result := string(data[1:pos])
	// fmt.Printf("Simple String parsed: %q\n", result)
	return result, pos + 2
}

func readArray(data []byte) ([]interface{}, int) {
	// fmt.Println("Reading Array...")
	pos := 1
	// Read array length
	strLen := 0
	for data[pos] != '\r' {
		strLen = strLen*10 + int(data[pos]-'0')
		pos++
	}
	pos += 2 // Skip \r\n
	// fmt.Printf("Array length: %d\n", strLen)

	var elems []interface{} = make([]interface{}, strLen)
	for i := 0; i < strLen; i++ {
		ele, delta := DecodeOne(data[pos:])
		elems[i] = ele
		pos += delta
	}
	// fmt.Printf("Array parsed: %v\n", elems)
	return elems, pos
}

func readBulkString(data []byte) (string, int) {
	// fmt.Println("Reading Bulk String...")
	pos := 1
	// Read string length
	strLen := 0
	for data[pos] != '\r' {
		strLen = strLen*10 + int(data[pos]-'0')
		pos++
	}
	pos += 2 // Skip \r\n after length
	// fmt.Printf("Bulk String length: %d\n", strLen)

	// Extract the string
	result := string(data[pos : pos+strLen])
	pos += strLen + 2 // Skip string content + \r\n
	// fmt.Printf("Bulk String parsed: %q\n", result)
	return result, pos
}

func readError(data []byte) (string, int) {
	// fmt.Println("Reading Error...")
	result, delta := readSimpleString(data)
	// fmt.Printf("Error parsed: %q\n", result)
	return result, delta
}

func readInt64(data []byte) (int64, int) {
	// fmt.Println("Reading Integer...")
	pos := 1
	var value int64 = 0
	for data[pos] != '\r' {
		value = value*10 + int64(data[pos]-'0')
		pos++
	}
	// fmt.Printf("Integer parsed: %d\n", value)
	return value, pos + 2
}

func DecodeOne(data []byte) (interface{}, int) {
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '*':
		return readArray(data)
	case '$':
		return readBulkString(data)
	case '-':
		return readError(data)
	case ':':
		return readInt64(data)
	default:
		// fmt.Printf("Unknown RESP type: %c\n", data[0])
		return string(data), len(data)
	}
}

func RESPParser(data []byte) interface{} {
	// fmt.Printf("Parsing RESP data: %q\n", string(data))
	value, _ := DecodeOne(data)
	return value
}

// This function will read the Byte array the client sent
// and will return the array of strings into tokens
func DecodeCmd(ReadBuffer []byte) ([]string, error) {
	// fmt.Printf("Decoding command from buffer: %q\n", string(ReadBuffer))
	value, _ := DecodeOne(ReadBuffer)
	ts := value.([]interface{})
	tokens := make([]string, len(ts))
	for i := range tokens {
		tokens[i] = ts[i].(string)
	}
	// fmt.Printf("Decoded command tokens: %v\n", tokens)
	return tokens, nil
}
