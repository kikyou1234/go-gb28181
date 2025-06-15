package sip

import (
	"fmt"
	"strings"
	"time"

	"go-sip/utils"
)

// DefaultProtocol DefaultProtocol
var DefaultProtocol = "udp"

// DefaultSipVersion DefaultSipVersion
var DefaultSipVersion = "SIP/2.0"

// Port number
type Port uint16

// NewPort NewPort
func NewPort(port int) *Port {
	newPort := Port(port)
	return &newPort
}

// Clone clone
func (port *Port) Clone() *Port {
	if port == nil {
		return nil
	}
	newPort := *port
	return &newPort
}

func (port *Port) String() string {
	if port == nil {
		return ""
	}
	return fmt.Sprintf("%d", *port)
}

// Equals Equals
func (port *Port) Equals(other interface{}) bool {
	if p, ok := other.(*Port); ok {
		return Uint16PtrEq((*uint16)(port), (*uint16)(p))
	}

	return false
}

// MaybeString  wrapper
type MaybeString interface {
	String() string
	Equals(other interface{}) bool
}

// String string
type String struct {
	Str string
}

func (str String) String() string {
	return str.Str
}

// Equals Equals
func (str String) Equals(other interface{}) bool {
	if v, ok := other.(String); ok {
		return str.Str == v.Str
	}

	return false
}

// ContentTypeSDP SDP contenttype
var ContentTypeSDP = ContentType("application/sdp")

// ContentTypeXML XML contenttype
var ContentTypeXML = ContentType("Application/MANSCDP+xml")

var ContentTypeRTSP = ContentType("Application/MANSRTSP")

var (
	// CatalogXML 获取设备列表xml样式
	CatalogXML = `<?xml version="1.0" encoding="GB2312"?>
<Query>
<CmdType>Catalog</CmdType>
<SN>%d</SN>
<DeviceID>%s</DeviceID>
</Query>
`
	// RecordInfoXML 获取录像文件列表xml样式
	RecordInfoXML = `<?xml version="1.0" encoding="GB2312"?>
<Query>
<CmdType>RecordInfo</CmdType>
<SN>%d</SN>
<DeviceID>%s</DeviceID>
<StartTime>%s</StartTime>
<EndTime>%s</EndTime>
<Secrecy>0</Secrecy>
<Type>all</Type>
</Query>
`
	// DeviceInfoXML 查询设备详情xml样式
	DeviceInfoXML = `<?xml version="1.0" encoding="GB2312"?>
<Query>
<CmdType>DeviceInfo</CmdType>
<SN>%d</SN>
<DeviceID>%s</DeviceID>
</Query>
`
)

func PausePlayContent() []byte {
	var content strings.Builder
	content.Grow(200) // 预分配内存（可选优化）

	// 添加协议头
	content.WriteString("PAUSE RTSP/1.0\r\n")

	// 添加 CSeq 头（假设 getInfoCseq() 返回 int）
	cseq := utils.GetInfoCseq()
	content.WriteString(fmt.Sprintf("CSeq: %d\r\n", cseq))

	content.WriteString("PauseTime: now\r\n")

	// 获取最终字符串
	result := content.String()
	return []byte(result)
}

func ResumePlayContent() []byte {
	var content strings.Builder
	content.Grow(200) // 预分配内存（可选优化）

	// 添加协议头
	content.WriteString("PLAY RTSP/1.0\r\n")

	// 添加 CSeq 头（假设 getInfoCseq() 返回 int）
	cseq := utils.GetInfoCseq()
	content.WriteString(fmt.Sprintf("CSeq: %d\r\n", cseq))

	content.WriteString("Range: npt=now-\r\n")

	// 获取最终字符串
	result := content.String()
	return []byte(result)
}

func SpeedPlayContent(speed float64) []byte {
	var content strings.Builder
	content.Grow(200) // 预分配内存（可选优化）

	// 添加协议头
	content.WriteString("PLAY RTSP/1.0\r\n")

	// 添加 CSeq 头（假设 getInfoCseq() 返回 int）
	cseq := utils.GetInfoCseq()
	content.WriteString(fmt.Sprintf("CSeq: %d\r\n", cseq))

	content.WriteString(fmt.Sprintf("Scale: %.6f\r\n", speed))

	// 获取最终字符串
	result := content.String()
	return []byte(result)
}

func SeekPlayContent(sub_time int64) []byte {
	var content strings.Builder
	content.Grow(200) // 预分配内存（可选优化）

	// 添加协议头
	content.WriteString("PLAY RTSP/1.0\r\n")

	// 添加 CSeq 头（假设 getInfoCseq() 返回 int）
	cseq := utils.GetInfoCseq()
	content.WriteString(fmt.Sprintf("CSeq: %d\r\n", cseq))

	content.WriteString(fmt.Sprintf("Range: npt= %d-\r\n", sub_time))

	// 获取最终字符串
	result := content.String()
	return []byte(result)
}

// GetDeviceInfoXML 获取设备详情指令
func GetDeviceInfoXML(id string) []byte {
	return []byte(fmt.Sprintf(DeviceInfoXML, utils.RandInt(100000, 999999), id))
}

// GetCatalogXML 获取NVR下设备列表指令
func GetCatalogXML(id string) []byte {
	return []byte(fmt.Sprintf(CatalogXML, utils.RandInt(100000, 999999), id))
}

// GetRecordInfoXML 获取录像文件列表指令
func GetRecordInfoXML(id string, sceqNo int, start, end int64) []byte {
	return []byte(fmt.Sprintf(RecordInfoXML, sceqNo, id, time.Unix(start, 0).Format("2006-01-02T15:04:05"), time.Unix(end, 0).Format("2006-01-02T15:04:05")))
}

func GetPausePlayContent() []byte {
	return PausePlayContent()
}

func GetResumePlayContent() []byte {
	return ResumePlayContent()
}

func GetSeekPlayContent(sub_time int64) []byte {
	return SeekPlayContent(sub_time)
}

func GetSpeedPlayContent(speed float64) []byte {
	return SpeedPlayContent(speed)
}

func GenerateBroadcastXML(ServerID string, channelID string) string {
	var builder strings.Builder
	builder.Grow(200) // 预分配内存

	sn := utils.RandInt(100000, 999999) // [0,900000) + 100000 = [100000,999999]

	// 构建 XML
	builder.WriteString(fmt.Sprintf("<?xml version=\"1.0\" encoding=\"%s\"?>\r\n", "GB2312"))
	builder.WriteString("<Notify>\r\n")
	builder.WriteString("<CmdType>Broadcast</CmdType>\r\n")
	builder.WriteString(fmt.Sprintf("<SN>%d</SN>\r\n", sn))
	builder.WriteString(fmt.Sprintf("<SourceID>%s</SourceID>\r\n", ServerID))
	builder.WriteString(fmt.Sprintf("<TargetID>%s</TargetID>\r\n", channelID))
	builder.WriteString("</Notify>\r\n")

	return builder.String()
}

func GenerateDeviceControl(DeviceID, cmdStr string) string {
	var ptzXml strings.Builder
	ptzXml.Grow(200) // 类似 Java 的 new StringBuilder(200)

	sn := utils.RandInt(100000, 999999) // [0,900000) + 100000 = [100000,999999]

	ptzXml.WriteString(fmt.Sprintf("<?xml version=\"1.0\" encoding=\"%s\"?>\r\n", "GB2312"))
	ptzXml.WriteString("\r\n<Control>\r\n")
	ptzXml.WriteString("<CmdType>DeviceControl</CmdType>\r\n")
	ptzXml.WriteString(fmt.Sprintf("<SN>%d</SN>\r\n", sn))
	ptzXml.WriteString(fmt.Sprintf("<DeviceID>%s</DeviceID>\r\n", DeviceID))
	ptzXml.WriteString(fmt.Sprintf("<PTZCmd>%s</PTZCmd>\r\n", cmdStr))
	ptzXml.WriteString("<Info>\r\n")
	ptzXml.WriteString("<ControlPriority>5</ControlPriority>\r\n")
	ptzXml.WriteString("</Info>\r\n")
	ptzXml.WriteString("</Control>\r\n")

	return ptzXml.String()
}

// RFC3261BranchMagicCookie RFC3261BranchMagicCookie
const RFC3261BranchMagicCookie = "z9hG4bK"

// GenerateBranch returns random unique branch ID.
func GenerateBranch() string {
	return strings.Join([]string{
		RFC3261BranchMagicCookie,
		utils.RandString(32),
	}, "")
}
