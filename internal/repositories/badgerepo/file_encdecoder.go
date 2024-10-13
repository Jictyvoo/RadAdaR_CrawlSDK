package badgerepo

import (
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jictyvoo/tcg_deck-resolver/internal/protodtos"
	"github.com/jictyvoo/tcg_deck-resolver/pkg/cacheproxy"
)

func EncodeFileInfo(fileInfo cacheproxy.FileInformation) ([]byte, error) {
	protoFileInfo := &protodtos.FileInformation{
		FileMime: &protodtos.FileMIME{
			Name:      fileInfo.FileMIME.Name,
			Extension: fileInfo.FileMIME.Extension,
			MimeType:  fileInfo.FileMIME.MimeType,
		},
		Envelope: &protodtos.Envelope{
			Headers: make(map[string]string, len(fileInfo.Envelope.Headers)),
			Status:  uint32(fileInfo.Envelope.Status),
		},
		Content:       fileInfo.Content,
		Checksum:      fileInfo.Checksum,
		CreatedAt:     timestamppb.New(fileInfo.CreatedAt),
		ModifiedAt:    timestamppb.New(fileInfo.ModifiedAt),
		ExtraMetadata: fileInfo.ExtraMetadata,
	}

	// Convert map of strings to Protobuf map
	for key, values := range fileInfo.Envelope.Headers {
		protoFileInfo.Envelope.Headers[key] = strings.Join(values, ",")
	}

	return proto.Marshal(protoFileInfo)
}

func DecodeFileInfo(bytes []byte) (cacheproxy.FileInformation, error) {
	protoFileInfo := protodtos.FileInformation{}
	if err := proto.Unmarshal(bytes, &protoFileInfo); err != nil {
		return cacheproxy.FileInformation{}, err
	}

	fileInfo := cacheproxy.FileInformation{
		FileMIME: cacheproxy.FileMIME{
			Name:      protoFileInfo.FileMime.Name,
			Extension: protoFileInfo.FileMime.Extension,
			MimeType:  protoFileInfo.FileMime.MimeType,
		},
		Envelope: cacheproxy.FileEnvelope{
			Headers: make(map[string][]string, len(protoFileInfo.Envelope.Headers)),
			Status:  uint16(protoFileInfo.Envelope.Status),
		},
		Content:       protoFileInfo.GetContent(),
		Checksum:      protoFileInfo.GetChecksum(),
		CreatedAt:     protoFileInfo.CreatedAt.AsTime(),
		ModifiedAt:    protoFileInfo.ModifiedAt.AsTime(),
		ExtraMetadata: protoFileInfo.GetExtraMetadata(),
	}

	for _, values := range protoFileInfo.Envelope.Headers {
		fileInfo.Envelope.Headers[values] = strings.Split(values, ",")
	}

	return fileInfo, nil
}
