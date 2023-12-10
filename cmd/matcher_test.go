package cmd_test

import (
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
			expectedIP:  cmd.IPv6Address{IP: [16]byte{
				0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00,
				0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34,
			}},
			expectedErr: nil,
		},
		{
			description: "Compressed IPv6 Address",
			ipStr:       "2001:db8::abcd:01ff:fe00:0",
			expectedIP:  cmd.IPv6Address{IP: [16]byte{
				0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
				0xab, 0xcd, 0x01, 0xff, 0xfe, 0x00, 0x00, 0x00,
			}},
			expectedErr: nil,
		},
		{
			description: "Compressed IPv6 Address",
			ipStr:       "2001:db8::50",
			expectedIP:  cmd.IPv6Address{IP: [16]byte{
				0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x50,
			}},
			expectedErr: nil,
		},
		{
			description: "Compressed IPv6 Address",
			ipStr:       "::1",
			expectedIP:  cmd.IPv6Address{IP: [16]byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
			}},
			expectedErr: nil,
		},
		{
			description: "Compressed IPv6 Address",
			ipStr:       "2001:db8::",
			expectedIP:  cmd.IPv6Address{IP: [16]byte{
				0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			}},
			expectedErr: nil,
		},
		{
			description: "Shortest Compressed IPv6 Address",
			ipStr:       "::",
			expectedIP:  cmd.IPv6Address{IP: [16]byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			}},
			expectedErr: nil,
		},
		{
			description: "Longest Compressed IPv6 Address",
			ipStr:       "2001:db8::1:abcd:01ff:fe00:0",
			expectedIP:  cmd.IPv6Address{IP: [16]byte{
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
			description: "IPv6 No Mask Pattern",
			pattern:     "2001:db8::abcd:01ff:fe00:0",
			expectedPattern: cmd.Pattern{
				IP: cmd.IPv6Address{IP: [16]byte{
					0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
					0xab, 0xcd, 0x01, 0xff, 0xfe, 0x00, 0x00, 0x00,
				}},
				MaskEnd:  128,
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
				MaskEnd:  10,
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
				MaskEnd:  128,
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
				MaskEnd:  104,
				MaskStart: 64,
			},
			expectedErr: nil,
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
		if err != nil {
			t.Errorf("expected error: %v, got: %v", nil, err)
		}
	}
}
