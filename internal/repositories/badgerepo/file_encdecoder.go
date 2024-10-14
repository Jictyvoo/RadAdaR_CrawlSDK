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
		FileMIME:      decodeFileMIME(protoFileInfo.FileMime),
		Envelope:      decodeFileEnvelope(protoFileInfo.Envelope),
		Content:       protoFileInfo.GetContent(),
		Checksum:      protoFileInfo.GetChecksum(),
		CreatedAt:     protoFileInfo.CreatedAt.AsTime(),
		ModifiedAt:    protoFileInfo.ModifiedAt.AsTime(),
		ExtraMetadata: protoFileInfo.GetExtraMetadata(),
	}

	return fileInfo, nil
}

func decodeFileMIME(protoMime *protodtos.FileMIME) cacheproxy.FileMIME {
	return cacheproxy.FileMIME{
		Name:      protoMime.GetName(),
		Extension: protoMime.GetExtension(),
		MimeType:  protoMime.GetMimeType(),
	}
}

func decodeFileEnvelope(protoEnvelope *protodtos.Envelope) cacheproxy.FileEnvelope {
	newEnvelope := cacheproxy.FileEnvelope{
		Headers: make(map[string][]string, len(protoEnvelope.GetHeaders())),
		Status:  uint16(protoEnvelope.GetStatus()),
	}

	for key, values := range protoEnvelope.GetHeaders() {
		newEnvelope.Headers[key] = strings.Split(values, ",")
	}

	return newEnvelope
}
