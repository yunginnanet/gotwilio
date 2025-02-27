package gotwilio

import (
	"fmt"
	"testing"
)

var params map[string]string
var ServiceSid = ""
var SessionSid = ""

func init() {
	params = make(map[string]string)
	params["SID"] = "AC0f30491286ab4abb4a108abefbd05d8a"
	params["TOKEN"] = "1dcf52d7a1f3853ed78f0ee20d056dd0"
	params["FROM"] = "+15005550006"
	params["TO"] = "+19135551234"
	params["WHATSAPP"] = "" // if empty, will skip WhatsApp tests
	// setup a WhatsApp sandbox in the Twilio console
	// and add the setup the WHATSAPP_FROM adn WHATSAPP_TO numbers
	params["WHATSAPP_FROM"] = ""
	params["WHATSAPP_TO"] = ""
}

func TestSMS(t *testing.T) {
	msg := "Welcome to gotwilio"
	twilio := NewTwilioClient(params["SID"], params["TOKEN"])
	_, exc, err := twilio.SendSMS(params["FROM"], params["TO"], msg, "", "")
	if err != nil {
		t.Fatal(err)
	}

	if exc.Code == 90010 {
		t.Skip("skipping SMS test, account not active")
	}

	if exc != nil {
		t.Fatal(exc)
	}
}

func TestMMS(t *testing.T) {
	msg := "Welcome to gotwilio"
	twilio := NewTwilioClient(params["SID"], params["TOKEN"])
	file := []string{"https://www.google.com/images/logo.png"}
	_, exc, err := twilio.SendMMS(params["FROM"], params["TO"], msg, file, "", "")
	if err != nil {
		t.Fatal(err)
	}

	if exc.Code == 90010 {
		t.Skip("skipping SMS test, account not active")
	}

	if exc != nil {
		t.Fatal(exc)
	}
}

func TestMMSMultipleFiles(t *testing.T) {
	msg := "Welcome to gotwilio"
	twilio := NewTwilioClient(params["SID"], params["TOKEN"])
	files := []string{"https://www.google.com/images/logo.png", "https://www.google.com/images/logo.png"}
	_, exc, err := twilio.SendMMS(params["FROM"], params["TO"], msg, files, "", "")
	if err != nil {
		t.Fatal(err)
	}

	if exc.Code == 90010 {
		t.Skip("skipping SMS test, account not active")
	}

	if exc != nil {
		t.Fatal(exc)
	}
}

func TestMMSTooManyFiles(t *testing.T) {
	msg := "Welcome to gotwilio"
	twilio := NewTwilioClient(params["SID"], params["TOKEN"])
	var files []string
	for i := 0; i < 11; i++ {
		files = append(files, "https://www.google.com/images/logo.png")
	}
	_, exc, err := twilio.SendMMS(params["FROM"], params["TO"], msg, files, "", "")
	if err != nil {
		t.Fatal(err)
	}

	if exc.Code == 90010 {
		t.Skip("skipping SMS test, account not active")
	}

	// Test for code for too many files
	if exc == nil || int(exc.Code) != 21623 {
		t.Fatal(exc)
	}
}

func TestWhatsApp(t *testing.T) {
	if len(params["WHATSAPP"]) == 0 {
		t.Skip("skipping WhatsApp test")
	}

	msg := "Welcome to gotwilio from WhatsApp"
	twilio := NewTwilioClient(params["SID"], params["TOKEN"])
	_, exc, err := twilio.SendWhatsApp(params["WHATSAPP_FROM"], params["WHATSAPP_TO"], msg, "", "")
	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}
}

func TestWhatsAppMedia(t *testing.T) {
	if len(params["WHATSAPP"]) == 0 {
		t.Skip("skipping WhatsApp test")
	}

	msg := "Welcome to gotwilio from WhatsApp Media, here's a cute cat picture"
	twilio := NewTwilioClient(params["SID"], params["TOKEN"])
	file := []string{
		"https://bit.ly/3bZY1oV", // cute cat photo
	}
	_, exc, err := twilio.SendWhatsAppMedia(params["WHATSAPP_FROM"], params["WHATSAPP_TO"], msg, file, "", "")
	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}
}

func TestVoice(t *testing.T) {
	callback := NewCallbackParameters("http://example.com")
	twilio := NewTwilioClient(params["SID"], params["TOKEN"])
	_, exc, err := twilio.CallWithUrlCallbacks(params["FROM"], params["TO"], callback)
	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}
}

func TestProxyServiceCRUD(t *testing.T) {
	t.Skip("Skipping as default test")

	twilio := NewTwilioClient(params["SID"], params["TOKEN"])

	req := ProxyServiceRequest{
		UniqueName:  "Test Service Name",
		CallbackURL: "https://www.example.com/",
	}

	resp, exc, err := twilio.NewProxyService(req)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	serviceID := resp.Sid

	resp, exc, err = twilio.GetProxyService(serviceID)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	req.OutOfSessionCallbackURL = "https://www.example.com/out"
	resp, exc, err = twilio.UpdateProxyService(serviceID, req)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	exc, err = twilio.DeleteProxyService(serviceID)
	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}

}

func TestProxySessionCRUD(t *testing.T) {
	t.Skip("Skipping as default test")

	twilio := NewTwilioClient(params["SID"], params["TOKEN"])
	// New Service to attach Session
	sreq := ProxyServiceRequest{
		UniqueName:  "Test Service Name",
		CallbackURL: "https://www.example.com/",
	}

	service, exc, err := twilio.NewProxyService(sreq)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	// Create Session
	req := ProxySessionRequest{
		UniqueName: "First Session Name",
		Mode:       "message-only",
	}

	session, exc, err := twilio.NewProxySession(service.Sid, req)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	// Get Session
	session, exc, err = twilio.GetProxySession(service.Sid, session.Sid)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	// Update Session
	req.Status = "closed"
	session, exc, err = twilio.UpdateProxySession(service.Sid, session.Sid, req)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	// Delete Session
	exc, err = twilio.DeleteProxySession(service.Sid, session.Sid)
	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}

	// Cleanup Service
	exc, err = twilio.DeleteProxyService(service.Sid)
	if err != nil {
		t.Fatal(err)
	}

	if exc != nil {
		t.Fatal(exc)
	}

}

func TestParticipantCreationAndMessage(t *testing.T) {
	t.Skip("Skipping as default test")

	twilio := NewTwilioClient(params["SID"], params["TOKEN"])

	service, exc, err := twilio.GetProxyService(ServiceSid)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	// Create ProxySession
	// req := ProxySessionRequest{
	// 	UniqueName: "Session#1234",
	// 	Mode:       "message-only",
	// }

	// session, exc, err := twilio.NewProxySession(service.Sid, req)
	session, exc, err := twilio.GetProxySession(service.Sid, SessionSid)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	// List Participants
	participants, exc, err := session.ListParticipants()
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}
	fmt.Println("Participant Count", len(participants))
	fmt.Println("")

	// Add Participant
	var participant Participant
	if len(participants) == 0 {
		u := ParticipantRequest{
			FriendlyName: "John Doe",
			Identifier:   "+15715551234",
		}

		participant, exc, err = session.AddParticipant(u)
		if err != nil {
			t.Fatal(err)
		}
		if exc != nil {
			t.Fatal(exc)
		}
	} else {
		// Get first one
		participant, exc, err = session.GetParticipant(participants[0].Sid)
		if err != nil {
			t.Fatal(err)
		}
		if exc != nil {
			t.Fatal(exc)
		}
		fmt.Printf("Participant: %v\n\n", participant)
	}

	msg := ProxyMessage{
		Body: "A follow up message",
	}
	resp, exc, err := session.CreateInteraction(participant.Sid, msg)
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}

	fmt.Printf("Message Interaction: %v\n\n", resp)

	interactions, exc, err := session.GetInteractions()
	if err != nil {
		t.Fatal(err)
	}
	if exc != nil {
		t.Fatal(exc)
	}
	fmt.Println("Interaction Count", len(interactions.Interactions))

	// participant, exc, err := session.DeleteParticipant(participant.ID)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if exc != nil {
	// 	t.Fatal(exc)
	// }

	// // Delete Session
	// exc, err = twilio.DeleteProxySession(service.Sid, session.Sid)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// if exc != nil {
	// 	t.Fatal(exc)
	// }

}
