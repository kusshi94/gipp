# gipp
A IP Addresses Filtering Tool with Prefix and Suffix

## Installation

with homebrew:

```bash
brew install kusshi94/tap/gipp
```

with `go install`:

```bash
go install github.com/kusshi/gipp@latest
```

## Usage

with files:

```bash
gipp [-e patterns ...] [file ...]
```

with stdin:

```bash
cat file | gipp [-e patterns ...]
```

### Filtering Patterns

#### Prefix

gipp filters IP addresses that have the specified prefix.
The prefix can be written in CIDR notation.

example:

```bash
gipp -e 192.168.1.0/24 input.txt
```

#### Suffix

gipp filters IP addresses that have the specified suffix.
The suffix can be written in CIDR notation-like format.
You can specify the length of the suffix by writing a minus number after the slash.

example:

```bash
gipp -e 0.0.0.1/-8 input.txt
```

#### Both of Prefix and Suffix

If you specify both prefix and suffix, gipp filters IP addresses that have the intersection of the prefix and suffix.

example:

```bash
gipp -e ::ef01:1ff:fe00:0/-64/104 input.txt
```
