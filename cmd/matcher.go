package cmd

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidIP      = errors.New("invalid ip")
	ErrInvalidPattern = errors.New("invalid pattern")
)

type IPAddress interface {
	Bytes() []byte
	Version() int
}

type IPv6Address struct {
	IP [16]byte
}

func (ip IPv6Address) Bytes() []byte {
	return ip.IP[:]
}

func (ip IPv6Address) Version() int {
	return 6
}

type IPv4Address struct {
	IP [4]byte
}

func (ip IPv4Address) Bytes() []byte {
	return ip.IP[:]
}

func (ip IPv4Address) Version() int {
	return 4
}

func ParseIp(ip string) (IPAddress, error) {
	for i := 0; i < len(ip); i++ {
		if ip[i] == '.' {
			return parseIPv4(ip)
		}
		if ip[i] == ':' {
			return parseIPv6(ip)
		}
	}
	return nil, ErrInvalidIP
}

func parseIPv6(ip string) (IPAddress, error) {
	var err error
	// 略記を展開する
	ip, err = extendIPv6(ip)
	if err != nil {
		return nil, err
	}

	// コロンで分割する
	blocks := strings.Split(ip, ":")
	// ブロックの数が8でない場合はエラー
	if len(blocks) != 8 {
		return nil, ErrInvalidIP
	}

	// ブロックを16進数に変換する
	var ipBytes [16]byte
	for i := 0; i < len(blocks); i++ {
		// ブロックが空の場合はエラー
		if blocks[i] == "" {
			return nil, ErrInvalidIP
		}
		// ブロックが4桁を超えている場合はエラー
		if len(blocks[i]) > 4 {
			return nil, ErrInvalidIP
		}
		// ブロックを16進数に変換する
		block, err := hexToBytes(blocks[i])
		if err != nil {
			return nil, err
		}
		// ブロックを挿入する
		copy(ipBytes[i*2:], block)
	}

	return IPv6Address{IP: ipBytes}, nil
}

// 略記を展開し、すべてのブロックを4桁にする
func extendIPv6(ipv6 string) (string, error) {
	// コロン2つによる略記が複数ある場合はエラー
	if strings.Count(ipv6, "::") > 1 {
		return "", ErrInvalidIP
	}

	// コロン2つによる略記を展開する
	if strings.Contains(ipv6, "::") {
		// コロン2つのみの場合
		if ipv6 == "::" {
			return "0000:0000:0000:0000:0000:0000:0000:0000", nil
		}

		// コロン2つが先頭でも末尾でもない場合
		if !strings.HasPrefix(ipv6, "::") && !strings.HasSuffix(ipv6, "::") {
			// コロン2つの位置を取得する
			idx := strings.Index(ipv6, "::")
			// コロン2つの位置で分割する
			head := ipv6[:idx]
			tail := ipv6[idx+1:]
			// 全体のコロンの数を数える
			colonCount := strings.Count(ipv6, ":")
			// 追加するブロックの数 = 8 - 全体のコロンの数
			addedBlockCount := 8 - colonCount
			// 追加するブロックを作成する
			addedBlock := strings.Repeat(":0000", addedBlockCount)
			// 追加するブロックを挿入する
			ipv6 = head + addedBlock + tail
		}

		// コロン2つが先頭にある場合
		if strings.HasPrefix(ipv6, "::") {
			// 先頭のコロン2つを削除する
			ipv6 = ipv6[2:]
			// 全体のコロンの数を数える
			colonCount := strings.Count(ipv6, ":")
			// 追加するブロックの数 = 8 - (全体のコロンの数 + 1)
			addedBlockCount := 8 - (colonCount + 1)
			// 追加するブロックを作成する
			addedBlock := strings.Repeat("0000:", addedBlockCount)
			// 追加するブロックを挿入する
			ipv6 = addedBlock + ipv6
		}

		// コロン2つが末尾にある場合
		if strings.HasSuffix(ipv6, "::") {
			// 末尾のコロン2つを削除する
			ipv6 = ipv6[:len(ipv6)-2]
			// 全体のコロンの数を数える
			colonCount := strings.Count(ipv6, ":")
			// 追加するブロックの数 = 8 - (全体のコロンの数 + 1)
			addedBlockCount := 8 - (colonCount + 1)
			// 追加するブロックを作成する
			addedBlock := strings.Repeat(":0000", addedBlockCount)
			// 追加するブロックを挿入する
			ipv6 = ipv6 + addedBlock
		}
	}

	// すべてのブロックが4桁になるように0を追加する
	blocks := strings.Split(ipv6, ":")
	for i := 0; i < len(blocks); i++ {
		// ブロックが4桁を超えている場合はエラー
		if len(blocks[i]) > 4 {
			return "", ErrInvalidIP
		}
		blocks[i] = strings.Repeat("0", 4-len(blocks[i])) + blocks[i]
	}
	return strings.Join(blocks, ":"), nil
}

// 16進数の文字列4桁をバイト列に変換する
func hexToBytes(s string) ([]byte, error) {
	if len(s) != 4 {
		return nil, ErrInvalidIP
	}

	var b [2]byte
	for i := 0; i < len(s); i++ {
		var n byte
		switch {
		case '0' <= s[i] && s[i] <= '9':
			n = s[i] - '0'
		case 'a' <= s[i] && s[i] <= 'f':
			n = s[i] - 'a' + 10
		case 'A' <= s[i] && s[i] <= 'F':
			n = s[i] - 'A' + 10
		default:
			return nil, ErrInvalidIP
		}
		if i%2 == 0 {
			b[i/2] = n << 4
		} else {
			b[i/2] |= n
		}
	}
	return b[:], nil
}

func parseIPv4(ip string) (IPAddress, error) {
	// ドットで分割する
	blocks := strings.Split(ip, ".")
	// ブロックの数が4でない場合はエラー
	if len(blocks) != 4 {
		return nil, ErrInvalidIP
	}

	// ブロックを10進数に変換する
	var ipBytes [4]byte
	for i := 0; i < len(blocks); i++ {
		// ブロックが空の場合はエラー
		if blocks[i] == "" {
			return nil, ErrInvalidIP
		}
		// ブロックが3桁を超えている場合はエラー
		if len(blocks[i]) > 3 {
			return nil, ErrInvalidIP
		}
		// ブロックを10進数に変換する
		block, err := strconv.Atoi(blocks[i])
		if err != nil {
			return nil, ErrInvalidIP
		}
		// ブロックが0~255の範囲外の場合はエラー
		if block < 0 || block > 255 {
			return nil, ErrInvalidIP
		}
		// ブロックを挿入する
		ipBytes[i] = byte(block)
	}

	return IPv4Address{IP: ipBytes}, nil
}

type Pattern interface {
	Match(ip IPAddress) bool
}

type IPv6Pattern struct {
	IP        IPv6Address
	MaskEnd   int
	MaskStart int
}

func (p IPv6Pattern) Match(ip IPAddress) bool {
	if ip.Version() != 6 {
		return false
	}
	ipBytes := ip.Bytes()
	for i := p.MaskStart; i < p.MaskEnd; i++ {
		if ipBytes[i] != p.IP.IP[i] {
			return false
		}
	}
	return true
}

type IPv4Pattern struct {
	IP   IPv4Address
	Mask [4]byte
}

func (p IPv4Pattern) Match(ip IPAddress) bool {
	if ip.Version() != 4 {
		return false
	}
	ipBytes := ip.Bytes()
	for i := 0; i < 32; i++ {
		if ipBytes[i] != p.IP.IP[i]&p.Mask[i] {
			return false
		}
	}
	return true
}

func ParsePattern(s string) (Pattern, error) {
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			return parseIPv4Pattern(s)
		}
		if s[i] == ':' {
			return parseIPv6Pattern(s)
		}
	}
	return nil, ErrInvalidIP
}

func parseIPv4Pattern(s string) (Pattern, error) {
	return nil, nil
}

func parseIPv6Pattern(s string) (Pattern, error) {
	// IPアドレスの部分を取り出す
	var ipPart string
	if strings.Contains(s, "/") {
		idx := strings.Index(s, "/")
		ipPart = s[:idx]
	} else {
		ipPart = s
	}
	ip, err := ParseIp(ipPart)
	if err != nil {
		return nil, err
	}

	// マスクの部分を取り出す
	var maskPart string
	if strings.Contains(s, "/") {
		idx := strings.Index(s, "/")
		maskPart = s[idx:]
	} else {
		maskPart = ""
	}
	// マスクを分割する
	masks := strings.Split(maskPart, "/")

	// マスクを適用する
	mask := [16]byte{}
	for i := 0; i < len(mask); i++ {
		mask[i] = 0xff
	}

	maskStart := 0
	maskEnd := 128
	for i := 0; i < len(masks); i++ {
		if masks[i] == "" {
			continue
		}
		masklen, err := strconv.Atoi(masks[i])
		if err != nil {
			return nil, ErrInvalidPattern
		}
		if masklen < -128 || masklen > 128 || masklen == 0 {
			return nil, ErrInvalidPattern
		}

		// Prefix指定の場合
		if masklen > 0 {
			maskEnd = masklen
		}
		// Suffix指定の場合
		if masklen < 0 {
			maskStart = 128 + masklen
		}
	}

	return IPv6Pattern{
		IP:        ip.(IPv6Address),
		MaskEnd:   maskEnd,
		MaskStart: maskStart,
	}, nil
}
