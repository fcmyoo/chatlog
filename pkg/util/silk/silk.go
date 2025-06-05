package silk

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/sjzar/go-lame"
	"github.com/sjzar/go-silk"
)

func Silk2MP3(data []byte) ([]byte, error) {

	sd := silk.SilkInit()
	defer sd.Close()

	pcmdata := sd.Decode(data)
	if len(pcmdata) == 0 {
		return nil, fmt.Errorf("silk decode failed")
	}

	le := lame.Init()
	defer le.Close()

	le.SetInSamplerate(24000)
	le.SetOutSamplerate(24000)
	le.SetNumChannels(1)
	le.SetBitrate(16)
	// IMPORTANT!
	le.InitParams()

	mp3data := le.Encode(pcmdata)
	if len(mp3data) == 0 {
		return nil, fmt.Errorf("mp3 encode failed")
	}

	return mp3data, nil
}

// Silk2Wav 将silk音频解码为wav格式字节流
func Silk2Wav(data []byte) ([]byte, error) {
	sd := silk.SilkInit()
	defer sd.Close()

	pcmdata := sd.Decode(data)
	if len(pcmdata) == 0 {
		return nil, fmt.Errorf("silk decode failed")
	}

	// 构造wav头部
	buf := new(bytes.Buffer)
	// 写RIFF头
	buf.WriteString("RIFF")
	// 先写0, 后面再补文件长度
	binary.Write(buf, binary.LittleEndian, uint32(36+len(pcmdata)*2))
	buf.WriteString("WAVEfmt ")
	binary.Write(buf, binary.LittleEndian, uint32(16))      // PCM块大小
	binary.Write(buf, binary.LittleEndian, uint16(1))       // PCM格式
	binary.Write(buf, binary.LittleEndian, uint16(1))       // 单声道
	binary.Write(buf, binary.LittleEndian, uint32(24000))   // 采样率
	binary.Write(buf, binary.LittleEndian, uint32(24000*2)) // 字节率
	binary.Write(buf, binary.LittleEndian, uint16(2))       // 每采样点字节数
	binary.Write(buf, binary.LittleEndian, uint16(16))      // 位深
	buf.WriteString("data")
	binary.Write(buf, binary.LittleEndian, uint32(len(pcmdata)*2))
	// 写PCM数据
	for _, v := range pcmdata {
		binary.Write(buf, binary.LittleEndian, v)
	}
	return buf.Bytes(), nil
}
