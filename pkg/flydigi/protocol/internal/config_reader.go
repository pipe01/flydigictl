package internal

type ConfigReader struct {
	totalPkgCount, pkgSize int

	recPkgCount int

	buf []byte
}

func NewConfigReader(pkgCount, pkgSize int) *ConfigReader {
	return &ConfigReader{
		totalPkgCount: pkgCount,
		pkgSize:       pkgSize,
		buf:           make([]byte, pkgSize*pkgCount),
	}
}

func (c *ConfigReader) GotPackage(pkgIndex int, data []byte) {
	if pkgIndex > c.totalPkgCount || pkgIndex != c.recPkgCount {
		return
	}

	copy(c.buf[pkgIndex*c.pkgSize:], data)
	c.recPkgCount++
}

func (c *ConfigReader) IsFinished() bool {
	return c.recPkgCount == c.totalPkgCount
}

func (c *ConfigReader) Data() []byte {
	if c.IsFinished() {
		return c.buf
	}

	return nil
}

func (c *ConfigReader) Reset() {
	c.recPkgCount = 0
}
