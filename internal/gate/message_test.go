package gate

import "testing"

func TestRequestRoundTrip(t *testing.T) {
	raw := EncodeRequest("User", "Login", []byte("hi"), "tok", 7, 3)
	msg, err := Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if msg.ServiceName != Services["User"] || msg.MethodName != "Login" {
		t.Fatalf("meta %+v", msg)
	}
	if string(msg.Body) != "hi" || msg.Token != "tok" || msg.ClientSeq != 7 || msg.ServerSeq != 3 {
		t.Fatalf("msg %+v", msg)
	}
	if msg.MessageType != TypeRequest {
		t.Fatalf("type %d", msg.MessageType)
	}
}
