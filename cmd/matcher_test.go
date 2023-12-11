package cmd_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kusshi94/gipp/cmd"
)

func TestParseIp(t *testing.T) {
	testCases := []struct {
		description string
		ipStr       string
		expectedIP  cmd.IPAddress
		expectedErr error
	}{
		{
			description: "Not Compressed IPv6 Address",
			ipStr:       "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			expectedIP: cmd.IPv6Address{IP: [16]byte{
				0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00,
				0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34,
			}},
			expectedErr: nil,
		},
		{
			description: "Compressed IPv6 Address",
			ipStr:       "2001:db8::abcd:01ff:fe00:0",
			expectedIP: cmd.IPv6Address{IP: [16]byte{
				0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
				0xab, 0xcd, 0x01, 0xff, 0xfe, 0x00, 0x00, 0x00,
			}},
			expectedErr: nil,
		},
		{
			description: "Compressed IPv6 Address",
			ipStr:       "2001:db8::50",
			expectedIP: cmd.IPv6Address{IP: [16]byte{
				0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x50,
			}},
			expectedErr: nil,
		},
		{
			description: "Compressed IPv6 Address",
			ipStr:       "::1",
			expectedIP: cmd.IPv6Address{IP: [16]byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
			}},
			expectedErr: nil,
		},
		{
			description: "Compressed IPv6 Address",
			ipStr:       "2001:db8::",
			expectedIP: cmd.IPv6Address{IP: [16]byte{
				0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			}},
			expectedErr: nil,
		},
		{
			description: "Shortest Compressed IPv6 Address",
			ipStr:       "::",
			expectedIP: cmd.IPv6Address{IP: [16]byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			}},
			expectedErr: nil,
		},
		{
			description: "Longest Compressed IPv6 Address",
			ipStr:       "2001:db8::1:abcd:01ff:fe00:0",
			expectedIP: cmd.IPv6Address{IP: [16]byte{
				0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x01,
				0xab, 0xcd, 0x01, 0xff, 0xfe, 0x00, 0x00, 0x00,
			}},
			expectedErr: nil,
		},
		{
			description: "Too Long IPv6 Address",
			ipStr:       "2001:0db8:85a3:0000:0000:8a2e:0370:7334:abcd",
			expectedIP:  nil,
			expectedErr: cmd.ErrInvalidIP,
		},
		{
			description: "Too Short IPv6 Address",
			ipStr:       "2001:0db8:85a3:0000:0000:8a2e:0370",
			expectedIP:  nil,
			expectedErr: cmd.ErrInvalidIP,
		},
		{
			description: "Not IPv6 Address",
			ipStr:       "202222:::1:12321:::1:1:21:1:1:4",
			expectedIP:  nil,
			expectedErr: cmd.ErrInvalidIP,
		},
		{
			description: "IPv4 Address",
			ipStr:       "192.168.0.1",
			expectedIP:  cmd.IPv4Address{IP: [4]byte{192, 168, 0, 1}},
			expectedErr: nil,
		},
		{
			description: "Too Large IPv4 Address",
			ipStr:       "192.168.0.256",
			expectedIP:  nil,
			expectedErr: cmd.ErrInvalidIP,
		},
		{
			description: "Too Long IPv4 Address",
			ipStr:       "10.0.0.0.1",
			expectedIP:  nil,
			expectedErr: cmd.ErrInvalidIP,
		},
		{
			description: "Too Short IPv4 Address",
			ipStr:       "10.0.0",
			expectedIP:  nil,
			expectedErr: cmd.ErrInvalidIP,
		},
		{
			description: "Not IPv4 Address",
			ipStr:       "111...12.321.321.4",
			expectedIP:  nil,
			expectedErr: cmd.ErrInvalidIP,
		},
	}

	for _, tc := range testCases {
		ip, err := cmd.ParseIp(tc.ipStr)
		if !reflect.DeepEqual(ip, tc.expectedIP) {
			t.Errorf("expected IP: %v, got: %v", tc.expectedIP, ip)
		}
		if err != tc.expectedErr {
			t.Errorf("expected error: %v, got: %v", tc.expectedErr, err)
		}
	}
}

func TestParseIPPattern(t *testing.T) {
	testCases := []struct {
		description     string
		pattern         string
		expectedPattern cmd.Pattern
		expectedErr     error
	}{
		{
			description: "IPv6 No Masks Pattern",
			pattern:     "2001:db8::abcd:01ff:fe00:0",
			expectedPattern: cmd.Pattern{
				IP: cmd.IPv6Address{IP: [16]byte{
					0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
					0xab, 0xcd, 0x01, 0xff, 0xfe, 0x00, 0x00, 0x00,
				}},
				MaskEnd:   128,
				MaskStart: 0,
			},
			expectedErr: nil,
		},
		{
			description: "IPv6 Prefix Pattern",
			pattern:     "fe80::/10",
			expectedPattern: cmd.Pattern{
				IP: cmd.IPv6Address{IP: [16]byte{
					0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				}},
				MaskEnd:   10,
				MaskStart: 0,
			},
			expectedErr: nil,
		},
		{
			description: "IPv6 Suffix Pattern",
			pattern:     "::100/-9",
			expectedPattern: cmd.Pattern{
				IP: cmd.IPv6Address{IP: [16]byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00,
				}},
				MaskEnd:   128,
				MaskStart: 119,
			},
			expectedErr: nil,
		},
		{
			description: "IPv6 Prefix and Suffix Pattern",
			pattern:     "::abcd:01ff:fe00:0/-64/104",
			expectedPattern: cmd.Pattern{
				IP: cmd.IPv6Address{IP: [16]byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0xab, 0xcd, 0x01, 0xff, 0xfe, 0x00, 0x00, 0x00,
				}},
				MaskEnd:   104,
				MaskStart: 64,
			},
			expectedErr: nil,
		},
		{
			description: "IPv6 Invalid Pattern",
			pattern:     "::abcd:01ff:fe00:0/-64/129",
			expectedPattern: cmd.Pattern{
				IP:        nil,
				MaskEnd:   0,
				MaskStart: 0,
			},
			expectedErr: cmd.ErrInvalidPattern,
		},
		{
			description: "IPv6 Invalid Pattern",
			pattern:     "::abcd:01ff:fe00:0/-129",
			expectedPattern: cmd.Pattern{
				IP:        nil,
				MaskEnd:   0,
				MaskStart: 0,
			},
			expectedErr: cmd.ErrInvalidPattern,
		},
		{
			description: "IPv4 No Masks Pattern",
			pattern:     "192.168.1.100",
			expectedPattern: cmd.Pattern{
				IP:        cmd.IPv4Address{IP: [4]byte{192, 168, 1, 100}},
				MaskEnd:   32,
				MaskStart: 0,
			},
			expectedErr: nil,
		},
		{
			description: "IPv4 Prefix Pattern",
			pattern:     "192.168.1.0/24",
			expectedPattern: cmd.Pattern{
				IP:        cmd.IPv4Address{IP: [4]byte{192, 168, 1, 0}},
				MaskEnd:   24,
				MaskStart: 0,
			},
			expectedErr: nil,
		},
		{
			description: "IPv4 Suffix Pattern",
			pattern:     "0.0.0.1/-8",
			expectedPattern: cmd.Pattern{
				IP:        cmd.IPv4Address{IP: [4]byte{0, 0, 0, 1}},
				MaskEnd:   32,
				MaskStart: 24,
			},
			expectedErr: nil,
		},
		{
			description: "IPv4 Prefix and Suffix Pattern",
			pattern:     "0.0.100.0/-16/24",
			expectedPattern: cmd.Pattern{
				IP:        cmd.IPv4Address{IP: [4]byte{0, 0, 100, 0}},
				MaskEnd:   24,
				MaskStart: 16,
			},
			expectedErr: nil,
		},
		{
			description: "IPv4 Invalid Pattern",
			pattern:     "192.168.1.0/-33",
			expectedPattern: cmd.Pattern{
				IP:        nil,
				MaskEnd:   0,
				MaskStart: 0,
			},
			expectedErr: cmd.ErrInvalidPattern,
		},
	}

	for _, tc := range testCases {
		pattern, err := cmd.ParsePattern(tc.pattern)
		if !reflect.DeepEqual(pattern, tc.expectedPattern) {
			t.Errorf("expected pattern: %v, got: %v", tc.expectedPattern, pattern)
		}
		if err != tc.expectedErr {
			t.Errorf("expected error: %v, got: %v", tc.expectedErr, err)
		}
	}
}

func TestIPPatternMatch(t *testing.T) {
	testCases := []struct {
		description string
		pattern     string
		ip          string
		expected    bool
	}{
		{
			description: "IPv6 No Masks Pattern",
			pattern:     "2001:db8::abcd:1ff:fe00:0",
			ip:          "2001:db8::abcd:1ff:fe00:0",
			expected:    true,
		},
		{
			description: "IPv6 /128 Pattern",
			pattern:     "2001:db8::abcd:1ff:fe00:0/128",
			ip:          "2001:db8::abcd:1ff:fe00:0",
			expected:    true,
		},
		{
			description: "IPv6 /64 Pattern",
			pattern:     "2001:db8::/64",
			ip:          "2001:db8::abcd:1ff:fe00:0",
			expected:    true,
		},
		{
			description: "IPv6 /-64 Pattern",
			pattern:     "0::abcd:1ff:fe00:0/-64",
			ip:          "2001:db8::abcd:1ff:fe00:0",
			expected:    true,
		},
		{
			description: "IPv6 EUI-64 Pattern",
			pattern:     "::abcd:1ff:fe00:0/-64/104",
			ip:          "2001:db8::abcd:1ff:fe00:0",
			expected:    true,
		},
		{
			description: "IPv6 No Masks and No Match Pattern",
			pattern:     "2001:db8::abcd:1ff:fe00:0",
			ip:          "2001:db8::abcd:1ff:fe00:1",
			expected:    false,
		},
		{
			description: "IPv6 /128 and No Match Pattern",
			pattern:     "2001:db8::abcd:1ff:fe00:0/128",
			ip:          "2001:db8::abcd:1ff:fe00:1",
			expected:    false,
		},
		{
			description: "IPv6 /64 and No Match Pattern",
			pattern:     "2001:db8:100::/64",
			ip:          "2001:db8:200::abcd:1ff:fe00:1",
			expected:    false,
		},
		{
			description: "IPv6 /-64 and No Match Pattern",
			pattern:     "::ef01:1ff:fe00:0/-64",
			ip:          "2001:db8::abcd:1ff:fe00:1",
			expected:    false,
		},
		{
			description: "IPv6 EUI-64 and No Match Pattern",
			pattern:     "::ef01:1ff:fe00:0/-64/104",
			ip:          "2001:db8::abcd:1ff:fe00:1",
			expected:    false,
		},
		{
			description: "IPv4 No Masks Pattern",
			pattern:     "192.168.100.1",
			ip:          "192.168.100.1",
			expected:    true,
		},
		{
			description: "IPv4 /32 Pattern",
			pattern:     "192.168.100.1/32",
			ip:          "192.168.100.1",
			expected:    true,
		},
		{
			description: "IPv4 /24 Pattern",
			pattern:     "192.168.100.0/24",
			ip:          "192.168.100.1",
			expected:    true,
		},
		{
			description: "IPv4 /-8 Pattern",
			pattern:     "0.0.0.101/-8",
			ip:          "192.168.100.101",
			expected:    true,
		},
		{
			description: "IPv4 No Masks and No Match Pattern",
			pattern:     "192.168.100.1",
			ip:          "10.0.0.1",
			expected:    false,
		},
		{
			description: "IPv4 /32 and No Match Pattern",
			pattern:     "192.168.100.1/32",
			ip:          "10.0.0.1",
			expected:    false,
		},
		{
			description: "IPv4 /24 and No Match Pattern",
			pattern:     "192.168.100.0/24",
			ip:          "10.0.0.1",
			expected:    false,
		},
		{
			description: "IPv4 /-24 and No Match Pattern",
			pattern:     "0.0.0.101/-24",
			ip:          "10.0.0.1",
			expected:    false,
		},
	}

	for _, tc := range testCases {
		fmt.Println(tc.description)
		pattern, err := cmd.ParsePattern(tc.pattern)
		if err != nil {
			t.Errorf("parse pattern: unexpected error: %v", err)
		}
		ip, err := cmd.ParseIp(tc.ip)
		if err != nil {
			t.Errorf("parse ip: unexpected error: %v", err)
		}
		if pattern.Match(ip) != tc.expected {
			t.Errorf("expected: %v, got: %v", tc.expected, pattern.Match(ip))
		}
	}
}
