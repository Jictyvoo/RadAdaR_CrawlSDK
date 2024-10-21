package badgerepo

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jictyvoo/radadar_crawlsdk/internal/protodtos"
	"github.com/jictyvoo/radadar_crawlsdk/pkg/cacheproxy"
)

func fixtureFileInfo() cacheproxy.FileInformation {
	return cacheproxy.FileInformation{
		FileMIME: cacheproxy.FileMIME{
			Name:      "example",
			Extension: "txt",
			MimeType:  "text/plain",
		},
		Envelope: cacheproxy.FileEnvelope{
			Headers: map[string][]string{
				"Content-Type": {"text/plain"},
				"X-Custom":     {"custom1", "custom2"},
			},
			Status: 200,
		},
		Content:       []byte("Hello, world!"),
		Checksum:      []byte("315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"),
		CreatedAt:     time.Now().UTC(),
		ModifiedAt:    time.Now().UTC(),
		ExtraMetadata: map[string]string{"author": "John Doe"},
	}
}

func compareProtoWithDTO(
	t *testing.T,
	fileInfo cacheproxy.FileInformation,
	decoded *protodtos.FileInformation,
) {
	// Check that encoded values match the original data
	if fileInfo.FileMIME.Name != decoded.FileMime.GetName() {
		t.Errorf(
			"Expected FileMIME Name %v, got %v",
			fileInfo.FileMIME.Name,
			decoded.FileMime.GetName(),
		)
	}
	if fileInfo.FileMIME.Extension != decoded.FileMime.GetExtension() {
		t.Errorf(
			"Expected FileMIME Extension %v, got %v",
			fileInfo.FileMIME.Extension,
			decoded.FileMime.GetExtension(),
		)
	}
	if fileInfo.FileMIME.MimeType != decoded.FileMime.GetMimeType() {
		t.Errorf(
			"Expected FileMIME MimeType %v, got %v",
			fileInfo.FileMIME.MimeType,
			decoded.FileMime.GetMimeType(),
		)
	}
	if uint32(fileInfo.Envelope.Status) != decoded.Envelope.Status {
		t.Errorf(
			"Expected Envelope Status %v, got %v",
			fileInfo.Envelope.Status,
			decoded.Envelope.Status,
		)
	}
	for key, value := range fileInfo.Envelope.Headers {
		compareValues := strings.SplitN(decoded.Envelope.Headers[key], ",", len(value))
		if !reflect.DeepEqual(value, compareValues) {
			t.Errorf(
				"Expected Envelope Header[%s] to be `%+v` = `%+v`",
				key, value, compareValues,
			)
		}
	}
	if !reflect.DeepEqual(fileInfo.Content, decoded.GetContent()) {
		t.Errorf("Expected Content %v, got %v", fileInfo.Content, decoded.GetContent())
	}
	if !reflect.DeepEqual(fileInfo.Checksum, decoded.GetChecksum()) {
		t.Errorf("Expected Checksum %v, got %v", fileInfo.Checksum, decoded.GetChecksum())
	}
}

func TestEncodeFileInfo_Success(t *testing.T) {
	fileInfo := fixtureFileInfo()

	encoded, err := EncodeFileInfo(fileInfo)
	if err != nil {
		t.Fatalf("EncodeFileInfo failed: %v", err)
	}
	if len(encoded) == 0 {
		t.Fatalf("Encoded data is empty")
	}

	// Decode back to verify
	var decoded protodtos.FileInformation
	err = proto.Unmarshal(encoded, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal encoded data: %v", err)
	}

	compareProtoWithDTO(t, fileInfo, &decoded)
}

func TestEncodeFileInfo_EmptyHeaders(t *testing.T) {
	fileInfo := fixtureFileInfo()
	fileInfo.Envelope.Headers = nil
	fileInfo.ExtraMetadata = nil

	encoded, err := EncodeFileInfo(fileInfo)
	if err != nil {
		t.Fatalf("EncodeFileInfo failed: %v", err)
	}
	if len(encoded) == 0 {
		t.Fatalf("Encoded data is empty")
	}
}

func TestDecodeFileInfo_Success(t *testing.T) {
	// Prepare a valid protodtos.FileInformation protobuf object
	protoFileInfo := &protodtos.FileInformation{
		FileMime: &protodtos.FileMIME{
			Name:      "example",
			Extension: "txt",
			MimeType:  "text/plain",
		},
		Envelope: &protodtos.Envelope{
			Headers: map[string]string{
				"Content-Type": "text/plain",
				"X-Custom":     "custom1,custom2",
			},
			Status: 200,
		},
		Content:       []byte("Hello, world!"),
		Checksum:      []byte("abc123"),
		CreatedAt:     timestamppb.Now(),
		ModifiedAt:    timestamppb.Now(),
		ExtraMetadata: map[string]string{"author": "John Doe"},
	}

	// Encode it to bytes
	encoded, err := proto.Marshal(protoFileInfo)
	if err != nil {
		t.Fatalf("EncodeFileInfo failed: %v", err)
	}
	if len(encoded) == 0 {
		t.Fatalf("Encoded data is empty")
	}

	// Decode using DecodeFileInfo
	var decodedFileInfo cacheproxy.FileInformation
	decodedFileInfo, err = DecodeFileInfo(encoded)
	if err != nil {
		t.Fatalf("DecodeFileInfo failed: %v", err)
	}

	compareProtoWithDTO(t, decodedFileInfo, protoFileInfo)
}

func TestDecodeFileInfo_Error(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		shouldErr bool
	}{
		{
			name:      "Invalid Data",
			input:     []byte("invalid protobuf data"),
			shouldErr: true,
		},
		{
			name:      "Empty Data",
			input:     []byte{},
			shouldErr: false,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeFileInfo(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("expected an error, but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("did not expect an error, but got: %v", err)
			}
		})
	}
}
