package url_signer

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"hash"
	"sync"
	"time"

	"golang.org/x/crypto/blake2b"
)

// Protect is the main struct used to sign and unsigned data
type Protect struct {
	sync.Mutex
	hash      hash.Hash
	dirty     bool
	timestamp bool
	epoch     int64
}

var (
	ErrInvalidSignature = errors.New("invalid signature")
	ErrShortToken       = errors.New("token is too small to be valid")
)

// New takes a secret key and returns a new Protect struct.  If no Options are provided
// then minimal defaults will be used. NOTE: The key must be 64 bytes or fewer
// . If a larger key is provided it will be truncated to 64 bytes.
func New(key []byte, options ...func(*Protect)) *Protect {

	var err error

	// Create a map for decoding Base58.  This speeds up the process tremendously.
	for i := 0; i < len(encBase58Map); i++ {
		decBase58Map[encBase58Map[i]] = byte(i)
	}

	s := &Protect{}

	for _, opt := range options {
		opt(s)
	}

	s.hash, err = blake2b.New256(key)
	if err != nil {
		// The only possible error that can be returned here is if the key
		// is larger than 64 bytes - which the blake2b hash will not accept.
		// This is a case that is so easily avoidable when using this package
		// and since chaining is convenient for this package.  We're going
		// to do the below to handle this possible case, so we don't have
		// to return an error.
		s.hash, _ = blake2b.New256(key[0:64])
	}

	return s
}

// Epoch is a functional option that can be passed to New() to set the Epoch
// to be used.
func Epoch(e int64) func(*Protect) {
	return func(s *Protect) {
		s.epoch = e
	}
}

// Timestamp is a functional option that can be passed to New() to add a
// timestamp to signatures.
func Timestamp(s *Protect) {
	s.timestamp = true
}

// Sign signs data and returns []byte in the format `data.signature`. Optionally
// add a timestamp and return to the format `data.timestamp.signature`
func (s *Protect) Sign(data []byte) []byte {

	// Build the payload
	el := base64.RawURLEncoding.EncodedLen(s.hash.Size())
	var t []byte

	if s.timestamp {
		ts := time.Now().Unix() - s.epoch
		etl := encodeBase58Len(ts)
		t = make([]byte, 0, len(data)+etl+el+2) // +2 for "." chars
		t = append(t, data...)
		t = append(t, '.')
		t = t[0 : len(t)+etl] // expand for timestamp
		encodeBase58(ts, t)
	} else {
		t = make([]byte, 0, len(data)+el+1)
		t = append(t, data...)
	}

	// Append and encode signature to token
	t = append(t, '.')
	tl := len(t)
	t = t[0 : tl+el]

	// Add the signature to the token
	s.sign(t[tl:], t[0:tl-1])

	// Return the token to the caller
	return t
}

// UnSing validates a signature and if successful returns the data portion of []byte.
// If unsuccessful it will return an error and nil for the data.
func (s *Protect) UnSing(token []byte) ([]byte, error) {

	tl := len(token)
	el := base64.RawURLEncoding.EncodedLen(s.hash.Size())

	// A token must be at least el+2 bytes long to be valid.
	if tl < el+2 {
		return nil, ErrShortToken
	}

	// Get the signature of the payload
	dst := make([]byte, el)
	s.sign(dst, token[0:tl-(el+1)])

	if subtle.ConstantTimeCompare(token[tl-el:], dst) != 1 {
		return nil, ErrInvalidSignature
	}

	return token[0 : tl-(el+1)], nil
}

// This is the map of characters used during base58 encoding.
const encBase58Map = "789924579abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

// Used to create a decoded map, so we can decode base58 fairly fast.
var decBase58Map [256]byte

// sign creates the encoded signature of payload and writes to dst
func (s *Protect) sign(dst, payload []byte) {

	s.Lock()
	if s.dirty {
		s.hash.Reset()
	}
	s.dirty = true
	s.hash.Write(payload)
	h := s.hash.Sum(nil)
	s.Unlock()

	base64.RawURLEncoding.Encode(dst, h)
}

// returns the len of base58 encoded i
func encodeBase58Len(i int64) int {

	var l = 1
	for i >= 58 {
		l++
		i /= 58
	}
	return l
}

// encode time int64 into []byte
func encodeBase58(i int64, b []byte) {
	p := len(b) - 1
	for i >= 58 {
		b[p] = encBase58Map[i%58]
		p--
		i /= 58
	}
	b[p] = encBase58Map[i]
}

// parses a base58 []byte into an int64
func decodeBase58(b []byte) int64 {
	var id int64
	for p := range b {
		id = id*58 + int64(decBase58Map[b[p]])
	}
	return id
}

type Token struct {
	Payload   []byte
	Timestamp time.Time
}

func (s *Protect) Parse(t []byte) Token {

	tl := len(t)
	el := base64.RawURLEncoding.EncodedLen(s.hash.Size())

	token := Token{}

	if s.timestamp {
		for i := tl - (el + 2); i >= 0; i-- {
			if t[i] == '.' {
				token.Payload = t[0:i]
				token.Timestamp = time.Unix(decodeBase58(t[i+1:tl-(el+1)])+s.epoch, 0)
				break
			}
		}
	} else {
		token.Payload = t[0 : tl-(el+1)]
	}

	return token
}
