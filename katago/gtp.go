package katago

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// GTPWriter gtp writer
type GTPWriter struct {
	// NumOfTransmitMoves number of transmit moves
	NumOfTransmitMoves   int
	MinRefreshCentSecond int
	Compression          bool
	writer               io.Writer
	buffer               *bytes.Buffer
	latestInfoWriteAt    *time.Time
	firstWrite           bool
}

type infoData struct {
	info   string
	visits int
}

// NewGTPWriter new gtp writer
func NewGTPWriter(writer io.Writer) *GTPWriter {
	return &GTPWriter{
		NumOfTransmitMoves:   15,
		writer:               writer,
		MinRefreshCentSecond: 30,
		Compression:          false,
	}
}

func (writer *GTPWriter) Write(buf []byte) {
	if writer.buffer == nil {
		writer.buffer = bytes.NewBuffer(buf)
	} else {
		writer.buffer.Write(buf)
	}
	//log.Printf("DEBUG got new buffer[%v]\n", string(buf))
	// split the whole buffer by lines
	content := string(writer.buffer.Bytes())
	// log.Printf("DEBUG content[%v]\n", content)
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		// empty string
		return
	}
	lastLine := lines[len(lines)-1]
	totalLines := len(lines)
	// last line is not a complete line
	writer.buffer = bytes.NewBuffer([]byte(lastLine))
	totalLines--
	// log.Printf("DEBUG lines found: %v, lines sent: %v\n", len(lines), totalLines)

	bufToWrite := writer.processLines(lines[:totalLines])

	if len(bufToWrite) > 0 {
		writer.writer.Write(bufToWrite)
	}
}

func toGZipBuffer(buf []byte) []byte {
	resultBuffer := bytes.NewBuffer(make([]byte, 0))
	zippedBuffer := bytes.NewBuffer(make([]byte, 0))
	gw, _ := gzip.NewWriterLevel(zippedBuffer, gzip.DefaultCompression)
	gw.Write(buf)
	gw.Close()

	zipped := zippedBuffer.Bytes()
	// log.Printf("zipped, origin len: %d, zipped len: %d\n", len(buf), len(zipped))
	binary.Write(resultBuffer, binary.LittleEndian, uint32(len(zipped)))
	resultBuffer.Write(zipped)
	return resultBuffer.Bytes()
}

func (writer *GTPWriter) processLines(lines []string) []byte {
	var buffer bytes.Buffer
	now := time.Now()
	for _, line := range lines {
		processedLine := writer.processLine(line)
		if strings.HasPrefix(processedLine, "info") {
			// info line too fast, ignore this line
			if writer.MinRefreshCentSecond > 0 && writer.latestInfoWriteAt != nil && writer.latestInfoWriteAt.After(now.Add(time.Millisecond*time.Duration(writer.MinRefreshCentSecond*-10))) {
				// too fast, ignore
				// log.Printf("DEBUG too fast, ignored info")
			} else {
				if writer.Compression {
					buffer.Write([]byte{0xff})
					buffer.Write(toGZipBuffer([]byte(processedLine)))
				} else {
					buffer.WriteString(processedLine)
				}
				buffer.WriteString("\n")
				writer.latestInfoWriteAt = &now
			}
		} else {
			// write directly
			writer.latestInfoWriteAt = nil
			buffer.WriteString(processedLine)
			buffer.WriteString("\n")
		}
	}
	return buffer.Bytes()
}

func (writer *GTPWriter) processLine(line string) string {
	if !strings.HasPrefix(line, "info") {
		return line
	}
	if writer.NumOfTransmitMoves == 0 {
		return line
	}
	ownershipIndex := strings.Index(line, "ownership")
	var infos []string
	if ownershipIndex >= 0 {
		infos = strings.Split(line[:ownershipIndex], "info")
	} else {
		infos = strings.Split(line, "info")
	}
	// log.Printf("DEBUG infos found: %v\n", len(infos))
	infoDatas := make([]infoData, len(infos))
	m := regexp.MustCompile(`visits ([0-9]+)`)
	for i, info := range infos {
		visits := 0
		match := m.FindStringSubmatch(info)
		if len(match) > 1 {
			v, err := strconv.Atoi(match[1])
			if err != nil {
				v = 0
			}
			visits = v
		}
		infoDatas[i] = infoData{
			visits: visits,
			info:   info,
		}
	}
	sort.SliceStable(infoDatas, func(i, j int) bool {
		return infoDatas[i].visits > infoDatas[j].visits
	})
	var buffer bytes.Buffer

	for i, infoData := range infoDatas {
		if i >= writer.NumOfTransmitMoves {
			break
		}
		buffer.WriteString("info")
		buffer.WriteString(infoData.info)
	}
	if ownershipIndex >= 0 {
		buffer.WriteString(line[ownershipIndex:])
	}
	return buffer.String()
}
