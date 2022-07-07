package crypto

func Adler32(buffer []byte, offset int, size int) uint32 {
	const adler = 65521
	d := [2]uint32{1, 0}

	p := offset
	for size > 0 {
		var tlen int
		if tlen > 5552 {
			tlen = 5552
		} else {
			tlen = size
		}
		size -= tlen

		for tlen > 0 {
			d[0] += uint32(buffer[p])
			d[1] += d[0]
			tlen--
			p++
		}

		d[0] %= adler
		d[1] %= adler
	}

	return (d[1] << 16) | d[0]
}
