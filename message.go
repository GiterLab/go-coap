package coap

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// CType represents the message type.
type CType uint8

const (
	// Confirmable messages require acknowledgements.
	Confirmable CType = 0
	// NonConfirmable messages do not require acknowledgements.
	NonConfirmable CType = 1
	// Acknowledgement is a message indicating a response to confirmable message.
	Acknowledgement CType = 2
	// Reset indicates a permanent negative acknowledgement.
	Reset CType = 3
)

var typeNames = [256]string{
	Confirmable:     "Confirmable",
	NonConfirmable:  "NonConfirmable",
	Acknowledgement: "Acknowledgement",
	Reset:           "Reset",
}

func init() {
	for i := range typeNames {
		if typeNames[i] == "" {
			typeNames[i] = fmt.Sprintf("Unknown (0x%x)", i)
		}
	}
}

func (t CType) String() string {
	return typeNames[t]
}

// CCode is the type used for both request and response codes.
type CCode uint8

// Request Codes
const (
	GET    CCode = 1
	POST   CCode = 2
	PUT    CCode = 3
	DELETE CCode = 4
)

// Response Codes
const (
	Created               CCode = 65
	Deleted               CCode = 66
	Valid                 CCode = 67
	Changed               CCode = 68
	Content               CCode = 69
	BadRequest            CCode = 128
	Unauthorized          CCode = 129
	BadOption             CCode = 130
	Forbidden             CCode = 131
	NotFound              CCode = 132
	MethodNotAllowed      CCode = 133
	NotAcceptable         CCode = 134
	PreconditionFailed    CCode = 140
	RequestEntityTooLarge CCode = 141
	UnsupportedMediaType  CCode = 143
	InternalServerError   CCode = 160
	NotImplemented        CCode = 161
	BadGateway            CCode = 162
	ServiceUnavailable    CCode = 163
	GatewayTimeout        CCode = 164
	ProxyingNotSupported  CCode = 165

	// All Code values are assigned by sub-registries according to the
	// following ranges:
	//   0.00      Indicates an Empty message (see Section 4.1).
	//   0.01-0.31 Indicates a request.  Values in this range are assigned by
	//             the "CoAP Method Codes" sub-registry (see Section 12.1.1).
	//   1.00-1.31 Reserved
	//   2.00-5.31 Indicates a response.  Values in this range are assigned by
	//             the "CoAP Response Codes" sub-registry (see
	//             Section 12.1.2).
	//   6.00-7.31 Reserved
	// 6.00-6.31
	GiterlabErrnoOk              = 192 // 正常响应  [PV1/PV2]
	GiterlabErrnoParamConfigure  = 193 // 有新的配置参数 [PV2]
	GiterlabErrnoFirmwareUpdate  = 194 // 有新的固件可以更新 [PV2]
	GiterlabErrnoUserCommand     = 195 // 有用户命令需要执行 [PV2]
	GiterlabErrnoEnterFlightMode = 220 // 进入飞行模式[PV2]

	// 7.00-7.31
	GiterlabErrnoIllegalKey                  = 224 //    KEY错误，设备激活码错误 [PV1/PV2]
	GiterlabErrnoDataError                   = 225 //    数据错误 [PV1/PV2]
	GiterlabErrnoDeviceNotExist              = 226 //    设备不存在或设备传感器类型匹配错误 [PV1/PV2]
	GiterlabErrnoTimeExpired                 = 227 //    时间过期 [PV1/PV2]
	GiterlabErrnoNotSupportProtocolVersion   = 228 //    不支持的协议版本 [PV1/PV2]
	GiterlabErrnoProtocolParsingErrors       = 229 //    议解析错误 [PV1/PV2]
	GiterlabErrnoRequestTimeout              = 230 // [*]请求超时 [PV1/PV2]
	GiterlabErrnoOptProtocolParsingErrors    = 231 //    可选附加头解析错误 [PV1/PV2]
	GiterlabErrnoNotSupportAnalyticalMethods = 232 //    不支持的可选附加头解析方法 [PV1/PV2]
	GiterlabErrnoNotSupportPacketType        = 233 //    不支持的包类型 [PV1/PV2]
	GiterlabErrnoDataDecodingError           = 234 //    数据解码错误 [PV1/PV2]
	GiterlabErrnoPackageLengthError          = 235 //    数据包长度字段错误 [PV1/PV2]
	GiterlabErrnoDuoxieyunServerRequestBusy  = 236 // [*]多协云服务器请求失败 [PV1过时了]
	GiterlabErrnoSluanServerRequestBusy      = 237 // [*]石峦服务器请求失败 [PV2过时了]
	GiterlabErrnoCacheServiceErrors          = 238 // [*]缓存服务出错 [PV1/PV2]
	GiterlabErrnoTableStoreServiceErrors     = 239 // [*]表格存储服务出错 [PV1/PV2]
	GiterlabErrnoDatabaseServiceErrors       = 240 // [*]数据库存储出错 [PV1/PV2]
	GiterlabErrnoNotSupportEncodingType      = 241 //    不支持的编码类型 [PV1/PV2]
	GiterlabErrnoDeviceRepeatRegistered      = 242 //    设备重复注册 [PV2]
	GiterlabErrnoDeviceSimCardUsed           = 243 //    设备手机卡重复使用 [PV2]
	GiterlabErrnoDeviceSimCardIllegal        = 244 //    设备手机卡未登记，非法的SIM卡 [PV2]
	GiterlabErrnoDeviceUpdateForcedFailed    = 245 //    强制更新设备信息失败 [PV2]
)

var codeNames = [256]string{
	GET:                   "GET",
	POST:                  "POST",
	PUT:                   "PUT",
	DELETE:                "DELETE",
	Created:               "Created",
	Deleted:               "Deleted",
	Valid:                 "Valid",
	Changed:               "Changed",
	Content:               "Content",
	BadRequest:            "BadRequest",
	Unauthorized:          "Unauthorized",
	BadOption:             "BadOption",
	Forbidden:             "Forbidden",
	NotFound:              "NotFound",
	MethodNotAllowed:      "MethodNotAllowed",
	NotAcceptable:         "NotAcceptable",
	PreconditionFailed:    "PreconditionFailed",
	RequestEntityTooLarge: "RequestEntityTooLarge",
	UnsupportedMediaType:  "UnsupportedMediaType",
	InternalServerError:   "InternalServerError",
	NotImplemented:        "NotImplemented",
	BadGateway:            "BadGateway",
	ServiceUnavailable:    "ServiceUnavailable",
	GatewayTimeout:        "GatewayTimeout",
	ProxyingNotSupported:  "ProxyingNotSupported",

	GiterlabErrnoOk:             "giterlabErrnoOk:",
	GiterlabErrnoParamConfigure: "giterlabErrnoParamConfigure",
	GiterlabErrnoFirmwareUpdate: "giterlabErrnoFirmwareUpdate",

	GiterlabErrnoIllegalKey:                  "GiterlabErrnoIllegalKey",
	GiterlabErrnoDataError:                   "GiterlabErrnoDataError",
	GiterlabErrnoDeviceNotExist:              "GiterlabErrnoDeviceNotExist",
	GiterlabErrnoTimeExpired:                 "GiterlabErrnoTimeExpired",
	GiterlabErrnoNotSupportProtocolVersion:   "GiterlabErrnoNotSupportProtocolVersion",
	GiterlabErrnoProtocolParsingErrors:       "GiterlabErrnoProtocolParsingErrors",
	GiterlabErrnoRequestTimeout:              "GiterlabErrnoRequestTimeout",
	GiterlabErrnoOptProtocolParsingErrors:    "GiterlabErrnoOptProtocolParsingErrors",
	GiterlabErrnoNotSupportAnalyticalMethods: "GiterlabErrnoNotSupportAnalyticalMethods",
	GiterlabErrnoNotSupportPacketType:        "GiterlabErrnoNotSupportPacketType",
	GiterlabErrnoDataDecodingError:           "GiterlabErrnoDataDecodingError",
	GiterlabErrnoPackageLengthError:          "GiterlabErrnoPackageLengthError",
	GiterlabErrnoDuoxieyunServerRequestBusy:  "GiterlabErrnoDuoxieyunServerRequestBusy",
	GiterlabErrnoSluanServerRequestBusy:      "GiterlabErrnoSluanServerRequestBusy",
	GiterlabErrnoCacheServiceErrors:          "GiterlabErrnoCacheServiceErrors",
	GiterlabErrnoTableStoreServiceErrors:     "GiterlabErrnoTableStoreServiceErrors",
	GiterlabErrnoDatabaseServiceErrors:       "GiterlabErrnoDatabaseServiceErrors",
	GiterlabErrnoNotSupportEncodingType:      "GiterlabErrnoNotSupportEncodingType",
	GiterlabErrnoDeviceRepeatRegistered:      "GiterlabErrnoDeviceRepeatRegistered",
	GiterlabErrnoDeviceSimCardUsed:           "GiterlabErrnoDeviceSimCardUsed",
	GiterlabErrnoDeviceSimCardIllegal:        "GiterlabErrnoDeviceSimCardIllegal",
	GiterlabErrnoDeviceUpdateForcedFailed:    "GiterlabErrnoDeviceUpdateForcedFailed",
}

func init() {
	for i := range codeNames {
		if codeNames[i] == "" {
			codeNames[i] = fmt.Sprintf("Unknown (0x%x)", i)
		}
	}
}

func (c CCode) String() string {
	return codeNames[c]
}

// Message encoding errors.
var (
	ErrInvalidTokenLen   = errors.New("invalid token length")
	ErrOptionTooLong     = errors.New("option is too long")
	ErrOptionGapTooLarge = errors.New("option gap too large")
)

// OptionID identifies an option in a message.
type OptionID uint32

/*
   +-----+----+---+---+---+----------------+--------+--------+---------+
   | No. | C  | U | N | R | Name           | Format | Length | Default |
   +-----+----+---+---+---+----------------+--------+--------+---------+
   |   1 | x  |   |   | x | If-Match       | opaque | 0-8    | (none)  |
   |   3 | x  | x | - |   | Uri-Host       | string | 1-255  | (see    |
   |     |    |   |   |   |                |        |        | below)  |
   |   4 |    |   |   | x | ETag           | opaque | 1-8    | (none)  |
   |   5 | x  |   |   |   | If-None-Match  | empty  | 0      | (none)  |
   |   7 | x  | x | - |   | Uri-Port       | uint   | 0-2    | (see    |
   |     |    |   |   |   |                |        |        | below)  |
   |   8 |    |   |   | x | Location-Path  | string | 0-255  | (none)  |
   |  11 | x  | x | - | x | Uri-Path       | string | 0-255  | (none)  |
   |  12 |    |   |   |   | Content-Format | uint   | 0-2    | (none)  |
   |  14 |    | x | - |   | Max-Age        | uint   | 0-4    | 60      |
   |  15 | x  | x | - | x | Uri-Query      | string | 0-255  | (none)  |
   |  17 | x  |   |   |   | Accept         | uint   | 0-2    | (none)  |
   |  20 |    |   |   | x | Location-Query | string | 0-255  | (none)  |
   |  35 | x  | x | - |   | Proxy-Uri      | string | 1-1034 | (none)  |
   |  39 | x  | x | - |   | Proxy-Scheme   | string | 1-255  | (none)  |
   |  60 |    |   | x |   | Size1          | uint   | 0-4    | (none)  |
   +-----+----+---+---+---+----------------+--------+--------+---------+
*/

// Option IDs.
const (
	IfMatch       OptionID = 1
	URIHost       OptionID = 3
	ETag          OptionID = 4
	IfNoneMatch   OptionID = 5
	Observe       OptionID = 6
	URIPort       OptionID = 7
	LocationPath  OptionID = 8
	URIPath       OptionID = 11
	ContentFormat OptionID = 12
	MaxAge        OptionID = 14
	URIQuery      OptionID = 15
	Accept        OptionID = 17
	LocationQuery OptionID = 20
	ProxyURI      OptionID = 35
	ProxyScheme   OptionID = 39
	Size1         OptionID = 60

	// The IANA policy for future additions to this sub-registry is split
	// into three tiers as follows.  The range of 0..255 is reserved for
	// options defined by the IETF (IETF Review or IESG Approval).  The
	// range of 256..2047 is reserved for commonly used options with public
	// specifications (Specification Required).  The range of 2048..64999 is
	// for all other options including private or vendor-specific ones,
	// which undergo a Designated Expert review to help ensure that the
	// option semantics are defined correctly.  The option numbers between
	// 65000 and 65535 inclusive are reserved for experiments.  They are not
	// meant for vendor-specific use of any kind and MUST NOT be used in
	// operational deployments.
	GiterLabID    OptionID = 65000
	GiterLabKey   OptionID = 65001
	AccessID      OptionID = 65002
	AccessKey     OptionID = 65003
	CheckCRC32    OptionID = 65004
	EncoderType   OptionID = 65005
	EncoderID     OptionID = 65006
	Flags         OptionID = 65007
	PackageNumber OptionID = 65100
)

// Option value format (RFC7252 section 3.2)
type valueFormat uint8

const (
	valueUnknown valueFormat = iota
	valueEmpty
	valueOpaque
	valueUint
	valueString
)

type optionDef struct {
	valueFormat valueFormat
	minLen      int
	maxLen      int
}

var optionDefs = [65536]optionDef{
	IfMatch:       {valueFormat: valueOpaque, minLen: 0, maxLen: 8},
	URIHost:       {valueFormat: valueString, minLen: 1, maxLen: 255},
	ETag:          {valueFormat: valueOpaque, minLen: 1, maxLen: 8},
	IfNoneMatch:   {valueFormat: valueEmpty, minLen: 0, maxLen: 0},
	Observe:       {valueFormat: valueUint, minLen: 0, maxLen: 3},
	URIPort:       {valueFormat: valueUint, minLen: 0, maxLen: 2},
	LocationPath:  {valueFormat: valueString, minLen: 0, maxLen: 255},
	URIPath:       {valueFormat: valueString, minLen: 0, maxLen: 255},
	ContentFormat: {valueFormat: valueUint, minLen: 0, maxLen: 2},
	MaxAge:        {valueFormat: valueUint, minLen: 0, maxLen: 4},
	URIQuery:      {valueFormat: valueString, minLen: 0, maxLen: 255},
	Accept:        {valueFormat: valueUint, minLen: 0, maxLen: 2},
	LocationQuery: {valueFormat: valueString, minLen: 0, maxLen: 255},
	ProxyURI:      {valueFormat: valueString, minLen: 1, maxLen: 1034},
	ProxyScheme:   {valueFormat: valueString, minLen: 1, maxLen: 255},
	Size1:         {valueFormat: valueUint, minLen: 0, maxLen: 4},

	// GiterLab: add private options
	GiterLabID:    {valueFormat: valueString, minLen: 0, maxLen: 255},
	GiterLabKey:   {valueFormat: valueString, minLen: 0, maxLen: 255},
	AccessID:      {valueFormat: valueString, minLen: 0, maxLen: 255},
	AccessKey:     {valueFormat: valueString, minLen: 0, maxLen: 255},
	CheckCRC32:    {valueFormat: valueUint, minLen: 0, maxLen: 4},
	EncoderType:   {valueFormat: valueUint, minLen: 0, maxLen: 4},
	EncoderID:     {valueFormat: valueUint, minLen: 0, maxLen: 4},
	Flags:         {valueFormat: valueUint, minLen: 0, maxLen: 4},
	PackageNumber: {valueFormat: valueUint, minLen: 0, maxLen: 4},
}

// MediaType specifies the content type of a message.
type MediaType uint16

// Content types.
const (
	TextPlain     MediaType = 0  // text/plain;charset=utf-8
	AppLinkFormat MediaType = 40 // application/link-format
	AppXML        MediaType = 41 // application/xml
	AppOctets     MediaType = 42 // application/octet-stream
	AppExi        MediaType = 47 // application/exi
	AppJSON       MediaType = 50 // application/json
)

type option struct {
	ID    OptionID
	Value interface{}
}

func encodeInt(v uint32) []byte {
	switch {
	case v == 0:
		return nil
	case v < 256:
		return []byte{byte(v)}
	case v < 65536:
		rv := []byte{0, 0}
		binary.BigEndian.PutUint16(rv, uint16(v))
		return rv
	case v < 16777216:
		rv := []byte{0, 0, 0, 0}
		binary.BigEndian.PutUint32(rv, uint32(v))
		return rv[1:]
	default:
		rv := []byte{0, 0, 0, 0}
		binary.BigEndian.PutUint32(rv, uint32(v))
		return rv
	}
}

func decodeInt(b []byte) uint32 {
	tmp := []byte{0, 0, 0, 0}
	copy(tmp[4-len(b):], b)
	return binary.BigEndian.Uint32(tmp)
}

func (o option) toBytes() []byte {
	var v uint32

	switch i := o.Value.(type) {
	case string:
		return []byte(i)
	case []byte:
		return i
	case MediaType:
		v = uint32(i)
	case int:
		v = uint32(i)
	case int32:
		v = uint32(i)
	case uint:
		v = uint32(i)
	case uint32:
		v = i
	default:
		panic(fmt.Errorf("invalid type for option %x: %T (%v)",
			o.ID, o.Value, o.Value))
	}

	return encodeInt(v)
}

func parseOptionValue(optionID OptionID, valueBuf []byte) interface{} {
	def := optionDefs[optionID]
	if def.valueFormat == valueUnknown {
		// Skip unrecognized options (RFC7252 section 5.4.1)
		return nil
	}
	if len(valueBuf) < def.minLen || len(valueBuf) > def.maxLen {
		// Skip options with illegal value length (RFC7252 section 5.4.3)
		return nil
	}
	switch def.valueFormat {
	case valueUint:
		intValue := decodeInt(valueBuf)
		if optionID == ContentFormat || optionID == Accept {
			return MediaType(intValue)
		}
		return intValue
	case valueString:
		return string(valueBuf)
	case valueOpaque, valueEmpty:
		return valueBuf
	}
	// Skip unrecognized options (should never be reached)
	return nil
}

type options []option

func (o options) Len() int {
	return len(o)
}

func (o options) Less(i, j int) bool {
	if o[i].ID == o[j].ID {
		return i < j
	}
	return o[i].ID < o[j].ID
}

func (o options) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o options) Minus(oid OptionID) options {
	rv := options{}
	for _, opt := range o {
		if opt.ID != oid {
			rv = append(rv, opt)
		}
	}
	return rv
}

// Message is a CoAP message.
type Message struct {
	Type      CType
	Code      CCode
	MessageID uint16

	Token, Payload []byte

	opts options
}

// IsConfirmable returns true if this message is confirmable.
func (m Message) IsConfirmable() bool {
	return m.Type == Confirmable
}

// Options gets all the values for the given option.
func (m Message) Options(o OptionID) []interface{} {
	var rv []interface{}

	for _, v := range m.opts {
		if o == v.ID {
			rv = append(rv, v.Value)
		}
	}

	return rv
}

// Option gets the first value for the given option ID.
func (m Message) Option(o OptionID) interface{} {
	for _, v := range m.opts {
		if o == v.ID {
			return v.Value
		}
	}
	return nil
}

func (m Message) optionStrings(o OptionID) []string {
	var rv []string
	for _, o := range m.Options(o) {
		rv = append(rv, o.(string))
	}
	return rv
}

// Path gets the Path set on this message if any.
func (m Message) Path() []string {
	return m.optionStrings(URIPath)
}

// PathString gets a path as a / separated string.
func (m Message) PathString() string {
	return strings.Join(m.Path(), "/")
}

// SetPathString sets a path by a / separated string.
func (m *Message) SetPathString(s string) {
	for s[0] == '/' {
		s = s[1:]
	}
	m.SetPath(strings.Split(s, "/"))
}

// SetPath updates or adds a URIPath attribute on this message.
func (m *Message) SetPath(s []string) {
	m.SetOption(URIPath, s)
}

// RemoveOption removes all references to an option
func (m *Message) RemoveOption(opID OptionID) {
	m.opts = m.opts.Minus(opID)
}

// AddOption adds an option.
func (m *Message) AddOption(opID OptionID, val interface{}) {
	iv := reflect.ValueOf(val)
	if (iv.Kind() == reflect.Slice || iv.Kind() == reflect.Array) &&
		iv.Type().Elem().Kind() == reflect.String {
		for i := 0; i < iv.Len(); i++ {
			m.opts = append(m.opts, option{opID, iv.Index(i).Interface()})
		}
		return
	}
	m.opts = append(m.opts, option{opID, val})
}

// SetOption sets an option, discarding any previous value
func (m *Message) SetOption(opID OptionID, val interface{}) {
	m.RemoveOption(opID)
	m.AddOption(opID, val)
}

const (
	extoptByteCode   = 13
	extoptByteAddend = 13
	extoptWordCode   = 14
	extoptWordAddend = 269
	extoptError      = 15
)

// MarshalBinary produces the binary form of this Message.
func (m *Message) MarshalBinary() ([]byte, error) {
	tmpbuf := []byte{0, 0}
	binary.BigEndian.PutUint16(tmpbuf, m.MessageID)

	/*
	     0                   1                   2                   3
	    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |Ver| T |  TKL  |      Code     |          Message ID           |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |   Token (if any, TKL bytes) ...
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |   Options (if any) ...
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |1 1 1 1 1 1 1 1|    Payload (if any) ...
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	*/

	buf := bytes.Buffer{}
	buf.Write([]byte{
		(1 << 6) | (uint8(m.Type) << 4) | uint8(0xf&len(m.Token)),
		byte(m.Code),
		tmpbuf[0], tmpbuf[1],
	})
	buf.Write(m.Token)

	/*
	     0   1   2   3   4   5   6   7
	   +---------------+---------------+
	   |               |               |
	   |  Option Delta | Option Length |   1 byte
	   |               |               |
	   +---------------+---------------+
	   \                               \
	   /         Option Delta          /   0-2 bytes
	   \          (extended)           \
	   +-------------------------------+
	   \                               \
	   /         Option Length         /   0-2 bytes
	   \          (extended)           \
	   +-------------------------------+
	   \                               \
	   /                               /
	   \                               \
	   /         Option Value          /   0 or more bytes
	   \                               \
	   /                               /
	   \                               \
	   +-------------------------------+

	   See parseExtOption(), extendOption()
	   and writeOptionHeader() below for implementation details
	*/

	extendOpt := func(opt int) (int, int) {
		ext := 0
		if opt >= extoptByteAddend {
			if opt >= extoptWordAddend {
				ext = opt - extoptWordAddend
				opt = extoptWordCode
			} else {
				ext = opt - extoptByteAddend
				opt = extoptByteCode
			}
		}
		return opt, ext
	}

	writeOptHeader := func(delta, length int) {
		d, dx := extendOpt(delta)
		l, lx := extendOpt(length)

		buf.WriteByte(byte(d<<4) | byte(l))

		tmp := []byte{0, 0}
		writeExt := func(opt, ext int) {
			switch opt {
			case extoptByteCode:
				buf.WriteByte(byte(ext))
			case extoptWordCode:
				binary.BigEndian.PutUint16(tmp, uint16(ext))
				buf.Write(tmp)
			}
		}

		writeExt(d, dx)
		writeExt(l, lx)
	}

	sort.Stable(&m.opts)

	prev := 0

	for _, o := range m.opts {
		b := o.toBytes()
		writeOptHeader(int(o.ID)-prev, len(b))
		buf.Write(b)
		prev = int(o.ID)
	}

	if len(m.Payload) > 0 {
		buf.Write([]byte{0xff})
	}

	buf.Write(m.Payload)

	return buf.Bytes(), nil
}

// ParseMessage extracts the Message from the given input.
func ParseMessage(data []byte) (Message, error) {
	rv := Message{}
	return rv, rv.UnmarshalBinary(data)
}

// UnmarshalBinary parses the given binary slice as a Message.
func (m *Message) UnmarshalBinary(data []byte) error {
	if len(data) < 4 {
		return errors.New("short packet")
	}

	if data[0]>>6 != 1 {
		return errors.New("invalid version")
	}

	m.Type = CType((data[0] >> 4) & 0x3)
	tokenLen := int(data[0] & 0xf)
	if tokenLen > 8 {
		return ErrInvalidTokenLen
	}

	m.Code = CCode(data[1])
	m.MessageID = binary.BigEndian.Uint16(data[2:4])

	if tokenLen > 0 {
		m.Token = make([]byte, tokenLen)
	}
	if len(data) < 4+tokenLen {
		return errors.New("truncated")
	}
	copy(m.Token, data[4:4+tokenLen])
	b := data[4+tokenLen:]
	prev := 0

	parseExtOpt := func(opt int) (int, error) {
		switch opt {
		case extoptByteCode:
			if len(b) < 1 {
				return -1, errors.New("truncated")
			}
			opt = int(b[0]) + extoptByteAddend
			b = b[1:]
		case extoptWordCode:
			if len(b) < 2 {
				return -1, errors.New("truncated")
			}
			opt = int(binary.BigEndian.Uint16(b[:2])) + extoptWordAddend
			b = b[2:]
		}
		return opt, nil
	}

	for len(b) > 0 {
		if b[0] == 0xff {
			b = b[1:]
			break
		}

		delta := int(b[0] >> 4)
		length := int(b[0] & 0x0f)

		if delta == extoptError || length == extoptError {
			return errors.New("unexpected extended option marker")
		}

		b = b[1:]

		delta, err := parseExtOpt(delta)
		if err != nil {
			return err
		}
		length, err = parseExtOpt(length)
		if err != nil {
			return err
		}

		if len(b) < length {
			return errors.New("truncated")
		}

		oid := OptionID(prev + delta)
		opval := parseOptionValue(oid, b[:length])
		b = b[length:]
		prev = int(oid)

		if opval != nil {
			m.opts = append(m.opts, option{ID: oid, Value: opval})
		}
	}
	m.Payload = b
	return nil
}
