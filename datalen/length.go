package datalen

import "errors"

//	+----------+
//	| 0xxxxxxx |
//	+----------+
//	max: 2^(7*1) - 1

//	+----------+----------+
//	| 10xxxxxx | xxxxxxxx |
//	+----------+----------+
//	max: 2^(7*2) - 1

//	+----------+----------+----------+
//	| 110xxxxx | xxxxxxxx | xxxxxxxx |
//	+----------+----------+----------+
//	max: 2^(7*3) - 1

//	+----------+----------+----------+----------+
//	| 1110xxxx | xxxxxxxx | xxxxxxxx | xxxxxxxx |
//	+----------+----------+----------+----------+
//	max: 2^(7*4) - 1

const (
	max1 = (1 << (7 * 1)) - 1
	max2 = (1 << (7 * 2)) - 1
	max3 = (1 << (7 * 3)) - 1
	max4 = (1 << (7 * 4)) - 1

	tag1 = 0x00 // 0000 0000
	tag2 = 0x80 // 1000 0000
	tag3 = 0xC0 // 1100 0000
	tag4 = 0xE0 // 1110 0000

	t1_mask = 0x80 // 1000 0000
	t2_mask = 0xC0 // 1100 0000
	t3_mask = 0xE0 // 1110 0000
	t4_mask = 0xF0 // 1111 0000

	v1_mask = 0x7F // 0111 1111
	v2_mask = 0x3F // 0011 1111
	v3_mask = 0x1F // 0001 1111
	v4_mask = 0x0F // 0000 1111
)

const (
	MaxValue = max4 // Maximum length value
	MaxSize  = 4    // Maximum number of bytes for encoding length
)

var (
	errNegativeValue = errors.New("datalen: length has a negative value")
	errMoreThanMax   = errors.New("datalen: length more than max value")
	errInsufData     = errors.New("datalen: insufficient data length")
	errFirstByte     = errors.New("datalen: invalid first byte")
)

func Encode(length int, data []byte) (int, error) {
	if length < 0 {
		return 0, errNegativeValue
	}
	var x = uint(length)
	if x <= max1 {
		if len(data) < 1 {
			return 0, errInsufData
		}
		data[0] = byte(x)
		return 1, nil
	}
	if x <= max2 {
		if len(data) < 2 {
			return 0, errInsufData
		}
		data[0] = tag2 | byte(x>>8)
		data[1] = byte(x)
		return 2, nil
	}
	if x <= max3 {
		if len(data) < 3 {
			return 0, errInsufData
		}
		data[0] = tag3 | byte(x>>16)
		data[1] = byte(x >> 8)
		data[2] = byte(x)
		return 3, nil
	}
	if x <= max4 {
		if len(data) < 4 {
			return 0, errInsufData
		}
		data[0] = tag4 | byte(x>>24)
		data[1] = byte(x >> 16)
		data[2] = byte(x >> 8)
		data[3] = byte(x)
		return 4, nil
	}
	return 0, errMoreThanMax
}

func Decode(data []byte, length *int) (int, error) {
	if len(data) == 0 {
		return 0, errInsufData
	}
	first := data[0]
	if (first & t1_mask) == tag1 {
		x := uint(first & v1_mask)
		*length = int(x)
		return 1, nil
	}
	if (first & t2_mask) == tag2 {
		if len(data) < 2 {
			return 0, errInsufData
		}
		x := uint(first & v2_mask)
		x = (x << 8) | uint(data[1])
		*length = int(x)
		return 2, nil
	}
	if (first & t3_mask) == tag3 {
		if len(data) < 3 {
			return 0, errInsufData
		}
		x := uint(first & v3_mask)
		x = (x << 8) | uint(data[1])
		x = (x << 8) | uint(data[2])
		*length = int(x)
		return 3, nil
	}
	if (first & t4_mask) == tag4 {
		if len(data) < 4 {
			return 0, errInsufData
		}
		x := uint(first & v4_mask)
		x = (x << 8) | uint(data[1])
		x = (x << 8) | uint(data[2])
		x = (x << 8) | uint(data[3])
		*length = int(x)
		return 4, nil
	}
	return 0, errFirstByte
}
