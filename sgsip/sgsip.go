// SIPExer Generic SIP Parsing Library
//
package sgsip

import (
	"net"
	"strconv"
	"strings"
)

// return and error code values
const (
	SGSIPRetOK = 0
	// generic errors
	SGSIPRetErr      = -1
	SGSIPRetNotFound = -2

	// first line parse errors
	SGSIPRetErrFLineShort          = -100
	SGSIPRetErrFLineFormat         = -101
	SGSIPRetErrFLineResponseShort  = -102
	SGSIPRetErrFLineResponseFormat = -103
	SGSIPRetErrFLineResponseCode   = -104
	SGSIPRetErrFLineRequestFormat  = -120

	// sip params parse errors
	SGSIPRetErrParamFormat = -150
)

const (
	ProtoNONE = iota
	ProtoUDP
	ProtoTCP
	ProtoTLS
	ProtoSCTP
	ProtoWS
	ProtoWSS
)

const (
	AFNONE = 0
	AFIPv4 = 4
	AFIPv6 = 6
	AFHost = 8
)

const (
	SchemaNONE = iota
	SchemaSIP
	SchemaSIPS
	SchemaTEL
)

const (
	ParamValNone = iota
	ParamValBare
	ParamValQuoted
)

const (
	FLineNone = iota
	FLineRequest
	FLineResponse
)

type SGSIPSocketAddress struct {
	Val     string
	Proto   string
	Addr    string
	Port    string
	PortNo  int
	AType   int
	ProtoId int
}

type SGSIPURI struct {
	Val      string
	Schema   string
	SchemaId int
	User     string
	Addr     string
	Port     string
	PortNo   int
	Params   string
	UParams  string
	Proto    string
	ProtoId  int
	AType    int
}

type SGSIPParam struct {
	Name  string
	Value string
	PMode int
}

type SGSIPFirstLine struct {
	Val      string
	MType    int
	Proto    string
	Method   string
	MethodId int
	URI      string
	Code     int    // status code number
	CodeVal  string // status code string
	Reason   string
}

type SGSIPHeader struct {
	Name  string
	Body  string
	HType int
}

type SGSIPBody struct {
	Content     string
	ContentLen  int
	ContentType string
}

type SGSIPMessage struct {
	Data    string
	FLine   SGSIPFirstLine
	RURI    SGSIPURI
	Headers []SGSIPHeader
	Body    SGSIPBody
	MFlags  int
}

// Quick detection of ip/address type
func SGAddrType(addr string) int {
	if net.ParseIP(addr) == nil {
		return AFHost
	}
	for i := 0; i < len(addr); i++ {
		switch addr[i] {
		case '.':
			return AFIPv4
		case ':':
			return AFIPv6
		}
	}
	return AFNONE
}

// Quick detection of ip/address type also with IPv6 square brackets
func SGAddrTypeEx(addr string) int {
	if addr[0:1] == "[" {
		// assuming only IPv6 address
		if addr[len(addr)-1:] != "]" {
			return AFNONE
		}
		if SGAddrType(addr[1:len(addr)-1]) != AFIPv6 {
			return AFNONE
		}
		return AFIPv6
	} else {
		return SGAddrType(addr)
	}
}

// SGSIPSetProto --
func SGSIPSetProto(protostr string, protoval *string, protoid *int) int {
	switch protostr {
	case "udp", "UDP":
		*protoid = ProtoUDP
		*protoval = "udp"
		return SGSIPRetOK
	case "tcp", "TCP":
		*protoid = ProtoTCP
		*protoval = "tcp"
		return SGSIPRetOK
	case "tls", "TLS":
		*protoid = ProtoTLS
		*protoval = "tls"
		return SGSIPRetOK
	case "sctp", "SCTP":
		*protoid = ProtoSCTP
		*protoval = "sctp"
		return SGSIPRetOK
	case "ws", "WS":
		*protoid = ProtoWS
		*protoval = "ws"
		return SGSIPRetOK
	case "wss", "WSS":
		*protoid = ProtoWSS
		*protoval = "wss"
		return SGSIPRetOK
	}
	return SGSIPRetErr
}

// SGSIPSetSchema --
func SGSIPSetSchema(schemastr string, schemaval *string, schemaid *int) int {
	switch schemastr {
	case "sip", "SIP":
		*schemaval = "sip"
		*schemaid = SchemaSIP
		return SGSIPRetOK
	case "sips", "SIPS":
		*schemaval = "sips"
		*schemaid = SchemaSIPS
		return SGSIPRetOK
	case "tel", "TEL":
		*schemaval = "tel"
		*schemaid = SchemaTEL
		return SGSIPRetOK
	default:
		return SGSIPRetErr
	}
}

// SGSIPParseSocketAddress --
func SGSIPParseSocketAddress(sockstr string, sockaddr *SGSIPSocketAddress) int {
	if sockstr[0:1] == "[" && sockstr[len(sockstr)-1:] == "]" {
		// assuming only IPv6 address -- fill with defaults
		sockaddr.AType = SGAddrTypeEx(sockaddr.Addr)
		if sockaddr.AType != AFIPv6 {
			return SGSIPRetErr
		}
		sockaddr.Val = sockstr
		sockaddr.Proto = "udp"
		sockaddr.ProtoId = ProtoUDP
		sockaddr.Addr = sockstr
		sockaddr.Port = "5060"
		sockaddr.PortNo = 5060
		return SGSIPRetOK
	}
	strArray := strings.SplitN(sockstr, ":", 2)
	if len(strArray) == 1 {
		// only host address -- fill with defaults
		sockaddr.Val = sockstr
		sockaddr.Proto = "udp"
		sockaddr.ProtoId = ProtoUDP
		sockaddr.Addr = sockstr
		sockaddr.Port = "5060"
		sockaddr.PortNo = 5060
		sockaddr.AType = SGAddrType(sockaddr.Addr)
		return SGSIPRetOK
	}
	strProto := strArray[0]
	strAddrPort := strArray[1]

	ret := SGSIPSetProto(strProto, &sockaddr.Proto, &sockaddr.ProtoId)
	if ret != SGSIPRetOK {
		// first token is not proto - assume addr:port
		sockaddr.Proto = "udp"
		sockaddr.ProtoId = ProtoUDP
		strAddrPort = sockstr
		strProto = ""
	}
	if strAddrPort[0:1] == "[" {
		strArray = strings.SplitN(strAddrPort, "]", 2)
		if strProto == "" && strArray[1][0:1] != ":" {
			// no port and only IPv6 tested before
			return SGSIPRetErr
		}
		sockaddr.Port = strArray[1][1:]
		i, err := strconv.Atoi(sockaddr.Port)
		if err != nil {
			return SGSIPRetErr
		}
		sockaddr.PortNo = i
		sockaddr.Addr = strArray[0] + "]"
		sockaddr.AType = SGAddrTypeEx(sockaddr.Addr)
		if sockaddr.AType != AFIPv6 {
			return SGSIPRetErr
		}
	} else {
		strArray = strings.SplitN(strAddrPort, ":", 2)
		if len(strArray) > 1 {
			sockaddr.Port = strArray[1]
			i, err := strconv.Atoi(sockaddr.Port)
			if err != nil {
				return SGSIPRetErr
			}
			sockaddr.PortNo = i
		} else {
			sockaddr.Port = "5060"
			sockaddr.PortNo = 5060
		}
		sockaddr.Addr = strArray[0]
		sockaddr.AType = SGAddrType(sockaddr.Addr)
	}
	sockaddr.Val = sockstr
	return SGSIPRetOK
}

// SGSIPParseURI --
func SGSIPParseURI(uristr string, uri *SGSIPURI) int {
	strArray := strings.SplitN(uristr, ":", 2)

	if len(strArray) < 2 {
		return SGSIPRetErr
	}
	ret := SGSIPSetSchema(strArray[0], &uri.Schema, &uri.SchemaId)
	if ret != SGSIPRetOK {
		return ret
	}
	atPos := strings.Index(strArray[1], "@")
	colPos := strings.Index(strArray[1], ":")
	scPos := strings.Index(strArray[1], ";")
	if atPos == 0 {
		// empty user part
		return SGSIPRetErr
	}
	if atPos < 0 && colPos < 0 && scPos < 0 {
		// no user, no port, no parameters
		uri.Addr = strArray[1]
		uri.Proto = "udp"
		uri.ProtoId = ProtoUDP
		uri.Port = "5060"
		uri.PortNo = 5060
		uri.AType = SGAddrType(uri.Addr)
		uri.Val = uristr
		return SGSIPRetOK
	}
	pHostPP := ""
	if atPos > 0 {
		pUser := strArray[1][0:atPos]
		pHostPP = strArray[1][atPos+1:]
		uScPos := strings.Index(pUser, ";")
		if uScPos == 0 {
			// empty user part
			return SGSIPRetErr
		}
		if uScPos < 0 {
			uri.User = pUser
		} else {
			uri.User = pUser[0 : uScPos+1]
			uri.UParams = pUser[uScPos+1:]
		}
	} else {
		pHostPP = strArray[1]
	}
	if colPos < 0 && scPos < 0 {
		// no port, no params
		uri.Addr = pHostPP
		uri.Proto = "udp"
		uri.ProtoId = ProtoUDP
		uri.Port = "5060"
		uri.PortNo = 5060
		uri.AType = SGAddrType(uri.Addr)
		uri.Val = uristr
		return SGSIPRetOK
	}
	pPortParams := ""
	if pHostPP[0:1] == "[" {
		if pHostPP[len(pHostPP)-1:] == "]" {
			// only IPv6 address
			uri.Addr = pHostPP
			uri.Proto = "udp"
			uri.ProtoId = ProtoUDP
			uri.Port = "5060"
			uri.PortNo = 5060
			uri.AType = SGAddrTypeEx(uri.Addr)
			if uri.AType != AFIPv6 {
				return SGSIPRetErr
			}
			uri.Val = uristr
			return SGSIPRetOK
		}
		strArray = strings.SplitN(pHostPP, "]", 2)
		uri.Addr = strArray[0] + "]"
		uri.AType = SGAddrTypeEx(uri.Addr)
		if uri.AType != AFIPv6 {
			return SGSIPRetErr
		}
		pPortParams = strArray[1]
	} else {
		scPos = strings.IndexAny(pHostPP, ":;")
		uri.Addr = pHostPP[0:scPos]
		uri.AType = SGAddrType(uri.Addr)
		pPortParams = pHostPP[scPos:]
	}
	pParams := ""
	if pPortParams[0:1] == ":" {
		// port
		pPort := ""
		scPos = strings.Index(pPortParams, ";")
		if scPos < 0 {
			pPort = pPortParams[1:]
		} else {
			pPort = pPortParams[1:scPos]
		}
		i, err := strconv.Atoi(pPort)
		if err != nil || i <= 0 {
			return SGSIPRetErr
		}
		uri.Port = pPort
		uri.PortNo = i
		if scPos < 0 {
			uri.Proto = "udp"
			uri.ProtoId = ProtoUDP
			uri.Val = uristr
			return SGSIPRetOK
		}
		pParams = pPortParams[scPos:]
	} else if pPortParams[0:1] == ";" {
		pParams = pPortParams
	} else {
		return SGSIPRetErr
	}
	uri.Proto = "udp"
	uri.ProtoId = ProtoUDP
	if len(pParams) > 0 {
		uri.Params = pParams[1:]
		strArray = strings.Split(pParams, ";transport=")
		if len(strArray) == 1 {
			uri.Val = uristr
			return SGSIPRetOK
		}
		scPos = strings.Index(strArray[1], ";")
		pProto := ""
		if scPos < 0 {
			pProto = strArray[1]
		} else {
			pProto = strArray[1][0:scPos]
		}
		ret := SGSIPSetProto(pProto, &uri.Proto, &uri.ProtoId)
		if ret != SGSIPRetOK {
			return SGSIPRetErr
		}
	}
	uri.Val = uristr
	return SGSIPRetOK
}

// SGSIPURIToSocketAddress --
func SGSIPURIToSocketAddress(uri *SGSIPURI, sockaddr *SGSIPSocketAddress) int {
	if len(uri.Proto) > 0 {
		sockaddr.Proto = uri.Proto
		sockaddr.ProtoId = uri.ProtoId
	} else {
		sockaddr.Proto = "udp"
		sockaddr.ProtoId = ProtoUDP
	}
	if len(uri.Addr) > 0 {
		sockaddr.Addr = uri.Addr
	} else {
		sockaddr.Addr = "127.0.0.1"
	}
	sockaddr.AType = SGAddrTypeEx(sockaddr.Addr)

	if len(uri.Port) > 0 {
		sockaddr.Port = uri.Port
		sockaddr.PortNo = uri.PortNo
	} else {
		sockaddr.Port = "5060"
		sockaddr.PortNo = 5060
	}
	sockaddr.Val = sockaddr.Proto + ":" + sockaddr.Addr + ":" + sockaddr.Port
	return SGSIPRetOK
}

// SGSocketAddressToSIPURI --
func SGSocketAddressToSIPURI(sockaddr *SGSIPSocketAddress, user string, tmode int, uri *SGSIPURI) int {
	if len(sockaddr.Proto) > 0 {
		uri.Proto = sockaddr.Proto
		uri.ProtoId = sockaddr.ProtoId
	} else {
		uri.Proto = "udp"
		uri.ProtoId = ProtoUDP
	}
	upart := ""
	if len(user) > 0 {
		uri.User = user
		upart = user + "@"
	}
	if len(sockaddr.Addr) > 0 {
		uri.Addr = sockaddr.Addr
	} else {
		uri.Addr = "127.0.0.1"
	}
	uri.AType = SGAddrTypeEx(uri.Addr)

	if len(sockaddr.Port) > 0 {
		uri.Port = sockaddr.Port
		uri.PortNo = sockaddr.PortNo
	} else {
		uri.Port = "5060"
		uri.PortNo = 5060
	}
	uri.Schema = "sip"
	uri.SchemaId = SchemaSIP

	if tmode == 0 && uri.ProtoId == ProtoUDP {
		uri.Val = uri.Schema + ":" + upart + sockaddr.Addr + ":" + sockaddr.Port
	} else {
		uri.Val = uri.Schema + ":" + upart + sockaddr.Addr + ":" + sockaddr.Port + ";transport=" + sockaddr.Proto
	}

	return SGSIPRetOK
}

// SGSIPParamsGet --
func SGSIPParamsGet(paramStr string, paramName string, vmode int, paramVal *SGSIPParam) int {
	if len(paramStr) < len(paramName) {
		return SGSIPRetNotFound
	}
	pStr := paramStr
	if pStr[0:1] != ";" {
		pStr = ";" + pStr
	}
	if pStr[len(pStr)-1:] != ";" {
		pStr = pStr + ";"
	}

	if strings.Index(pStr, ";"+paramName+";") >= 0 {
		// parameter without value
		paramVal.Name = paramName
		paramVal.Value = ""
		paramVal.PMode = ParamValBare
		return SGSIPRetOK
	}

	strArray := strings.Split(pStr, ";"+paramName+"=")
	if len(strArray) == 1 {
		return SGSIPRetNotFound
	}
	scPos := -1
	qVal := 0
	if strArray[1][0:1] == "\"" {
		if vmode == 0 {
			return SGSIPRetErrParamFormat
		}
		scPos = strings.Index(strArray[1], "\";")
		paramVal.PMode = ParamValQuoted
		qVal = 1
	} else {
		paramVal.PMode = ParamValBare
		scPos = strings.Index(strArray[1], ";")
	}
	if scPos < 0 {
		paramVal.Value = strArray[1]
	} else {
		paramVal.Value = strArray[1][0 : scPos+qVal]
	}
	paramVal.Name = paramName
	return SGSIPRetOK
}

// SGSIPParseFirstLine --
func SGSIPParseFirstLine(inputStr string, flineVal *SGSIPFirstLine) int {
	strArray := strings.SplitN(inputStr, "\n", 2)
	strFLine := strings.Trim(strArray[0], " \t\r")
	if len(strFLine) < 8 {
		return SGSIPRetErrFLineShort
	}
	if strFLine[0:8] == "SIP/2.0 " {
		flineVal.MType = FLineResponse
	} else if strFLine[len(strFLine)-8:] == " SIP/2.0" {
		flineVal.MType = FLineRequest
	} else {
		return SGSIPRetErrFLineFormat
	}
	flineVal.Val = strFLine
	if flineVal.MType == FLineResponse {
		strCR := strFLine[8:]
		if len(strCR) < 5 {
			return SGSIPRetErrFLineResponseShort
		}
		strArray = strings.SplitN(strCR, " ", 2)
		if len(strArray) < 2 || len(strArray[0]) != 3 {
			return SGSIPRetErrFLineResponseFormat
		}
		i, err := strconv.Atoi(strArray[0])
		if err != nil || i < 100 || i > 999 {
			return SGSIPRetErrFLineResponseCode
		}
		flineVal.Code = i
		flineVal.CodeVal = strArray[0]
		flineVal.Reason = strings.Trim(strArray[1], " \t\r")
		return SGSIPRetOK
	}
	strMU := strFLine[0 : len(strFLine)-8]
	strArray = strings.SplitN(strMU, " ", 2)
	if len(strArray) < 2 || len(strArray[0]) < 3 || len(strArray[1]) < 5 {
		return SGSIPRetErrFLineRequestFormat
	}
	flineVal.Method = strings.Trim(strArray[0], " \t\r")
	flineVal.URI = strings.Trim(strArray[1], " \t\r")
	return SGSIPRetOK
}

// SGSIPValidHeaderName --
func SGSIPHeaderValidName(name string) bool {
	for i, r := range name {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			if i == 0 {
				return false
			}
			if (r < '0' || r > '9') && (r != '-') {
				return false
			}
		}
	}
	return true
}

//
// SGSIPHeaderParseDigestAuthBody - parse www/proxy-authenticate header body.
// Return a map of parameters or nil if the header is not Digest auth header.
func SGSIPHeaderParseDigestAuthBody(hbody string) map[string]string {
	s := strings.SplitN(strings.Trim(hbody, " "), " ", 2)
	if len(s) != 2 || s[0] != "Digest" {
		return nil
	}

	params := map[string]string{}
	for _, kv := range strings.Split(s[1], ",") {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			continue
		}
		params[strings.Trim(parts[0], "\" ")] = strings.Trim(parts[1], "\" ")
	}
	return params
}

// SGSIPParseHeader --
func SGSIPParseHeaders(inputStr string, pMode int, headersList *[]SGSIPHeader) int {
	var strArray []string
	var strHeaders string
	if pMode == 0 {
		strArray = strings.SplitN(inputStr, "\n", 2)
		if len(strArray) < 2 {
			return SGSIPRetErr
		}
		strHeaders = strArray[1]
	} else {
		strHeaders = strings.TrimLeft(inputStr, " \t\r\n")
	}
	if len(strHeaders) == 0 || strHeaders[0:1] == "\r" || strHeaders[0:1] == "\n" {
		// empty or first char is an EoL
		return SGSIPRetErr
	}
	for {
		var hdrItem SGSIPHeader = SGSIPHeader{}
		// split name: body
		strArray = strings.SplitN(strHeaders, ":", 2)
		if len(strArray) < 2 || len(strArray[0]) == 0 || len(strArray[1]) == 0 {
			return SGSIPRetErr
		}
		if !SGSIPHeaderValidName(strArray[0]) {
			return SGSIPRetErr
		}
		hdrItem.Name = strings.TrimRight(strArray[0], " \t")
		hdrItem.Body = ""
		for {
			strArray = strings.SplitN(strArray[1], "\n", 2)
			if len(strArray) < 2 {
				return SGSIPRetErr
			}
			hdrItem.Body += strArray[0]
			if len(strArray[1]) == 0 {
				break
			}
			// check if body spans over next line
			if strArray[1][0:1] != " " && strArray[1][0:1] != "\t" {
				break
			}
		}
		hdrItem.Body = strings.Trim(hdrItem.Body, " \t\r")
		*headersList = append(*headersList, hdrItem)
		strHeaders = strArray[1]
		if len(strHeaders) == 0 || strHeaders[0:1] == "\r" || strHeaders[0:1] == "\n" {
			// EoH
			break
		}
	}
	return SGSIPRetOK
}

// SGSIPParseBody --
func SGSIPParseBody(inputStr string, bodyVal *SGSIPBody) int {
	strArray := strings.SplitN(inputStr, "\r\n\r\n", 2)
	if len(strArray) < 2 {
		strArray := strings.SplitN(inputStr, "\n\n", 2)
		if len(strArray) < 2 {
			return SGSIPRetErr
		}
	}
	bodyVal.Content = strArray[1]
	bodyVal.ContentLen = len(strArray[1])
	return SGSIPRetOK
}

// SGSIPMessageHeaderSet --
func SGSIPMessageHeaderSet(msgVal *SGSIPMessage, hname string, hbody string) int {
	for i, hdr := range msgVal.Headers {
		if hdr.Name == hname {
			msgVal.Headers[i].Body = strings.Trim(hbody, " \t\r")
			return SGSIPRetOK
		}
	}
	var hdrItem SGSIPHeader = SGSIPHeader{}
	hdrItem.Name = strings.Trim(hname, " \t\r")
	hdrItem.Body = strings.Trim(hbody, " \t\r")
	msgVal.Headers = append(msgVal.Headers, hdrItem)

	return SGSIPRetOK
}

// SGSIPMessageHeaderGet --
func SGSIPMessageHeaderGet(msgVal *SGSIPMessage, hname string, hbody *string) int {
	for i, hdr := range msgVal.Headers {
		if hdr.Name == hname {
			*hbody = msgVal.Headers[i].Body
			return SGSIPRetOK
		}
	}

	return SGSIPRetErr
}

// SGSIPMessageCSeqUpdate --
func SGSIPMessageCSeqUpdate(msgVal *SGSIPMessage, ival int) int {
	for i, hdr := range msgVal.Headers {
		if hdr.Name == "CSeq" || hdr.Name == "s" {
			slist := strings.SplitN(msgVal.Headers[i].Body, " ", 2)
			if len(slist) != 2 {
				return SGSIPRetErr
			}
			csn, _ := strconv.Atoi(slist[0])

			msgVal.Headers[i].Body = strconv.Itoa(csn+ival) + " " + slist[1]
			return SGSIPRetOK
		}
	}
	return SGSIPRetErr
}

// SGSIPParseMessage --
func SGSIPParseMessage(inputStr string, msgVal *SGSIPMessage) int {
	ret := SGSIPParseFirstLine(inputStr, &msgVal.FLine)
	if ret != SGSIPRetOK {
		return ret
	}
	if msgVal.FLine.MType == FLineRequest {
		ret = SGSIPParseURI(msgVal.FLine.URI, &msgVal.RURI)
		if ret != SGSIPRetOK {
			return ret
		}
	}
	ret = SGSIPParseHeaders(inputStr, 0, &msgVal.Headers)
	if ret != SGSIPRetOK {
		return ret
	}
	ret = SGSIPParseBody(inputStr, &msgVal.Body)
	return SGSIPRetOK
}

// SGSIPMessageToString --
func SGSIPMessageToString(msgVal *SGSIPMessage, outputStr *string) int {
	var sb strings.Builder
	if len(msgVal.FLine.Val) == 0 || len(msgVal.Headers) == 0 {
		return SGSIPRetErr
	}
	sb.WriteString(msgVal.FLine.Val + "\r\n")

	if msgVal.Body.ContentLen > 0 {
		SGSIPMessageHeaderSet(msgVal, "Content-Length", strconv.Itoa(msgVal.Body.ContentLen))
		if len(msgVal.Body.ContentType) > 0 {
			SGSIPMessageHeaderSet(msgVal, "Content-Type", msgVal.Body.ContentType)
		}
	}

	for _, h := range msgVal.Headers {
		sb.WriteString(h.Name + ": " + h.Body + "\r\n")
	}
	sb.WriteString("\r\n")

	if msgVal.Body.ContentLen > 0 {
		sb.WriteString(msgVal.Body.Content)
	}

	*outputStr = sb.String()
	return SGSIPRetOK
}
