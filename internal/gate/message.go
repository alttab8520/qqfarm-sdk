package gate

import "github.com/alttab8520/qqfarm-sdk/internal/pb"

const (
	TypeRequest  = 1
	TypeResponse = 2
	TypeNotify   = 3
)

var Services = map[string]string{
	"User":   "gamepb.userpb.UserService",
	"Plant":  "gamepb.plantpb.PlantService",
	"Visit":  "gamepb.visitpb.VisitService",
	"Friend": "gamepb.friendpb.FriendService",
}

type Message struct {
	ServiceName  string
	MethodName   string
	MessageType  int32
	ClientSeq    int64
	ServerSeq    int64
	ErrorCode    int64
	ErrorMessage string
	Body         []byte
	Token        string
}

func EncodeRequest(serviceKey, method string, body []byte, token string, clientSeq, serverSeq int64) []byte {
	name := Services[serviceKey]
	if name == "" {
		name = serviceKey
	}
	meta := pb.NewEncoder()
	meta.String(1, name)
	meta.String(2, method)
	meta.Int(3, int64(TypeRequest))
	meta.Int(4, clientSeq)
	meta.Int(5, serverSeq)
	msg := pb.NewEncoder()
	msg.Message(1, meta.Bytes())
	msg.BytesField(2, body)
	msg.String(3, token)
	return msg.Bytes()
}

func Decode(data []byte) (Message, error) {
	var out Message
	fields, err := pb.Walk(data)
	if err != nil {
		return out, err
	}
	for _, f := range fields {
		switch f.Num {
		case 1:
			meta, err := pb.FieldMap(f.Bytes)
			if err != nil {
				return out, err
			}
			out.ServiceName = pb.StringField(meta, 1)
			out.MethodName = pb.StringField(meta, 2)
			out.MessageType = int32(pb.IntField(meta, 3))
			out.ClientSeq = pb.IntField(meta, 4)
			out.ServerSeq = pb.IntField(meta, 5)
			out.ErrorCode = pb.IntField(meta, 6)
			out.ErrorMessage = pb.StringField(meta, 7)
		case 2:
			out.Body = f.Bytes
		case 3:
			out.Token = string(f.Bytes)
		}
	}
	return out, nil
}
