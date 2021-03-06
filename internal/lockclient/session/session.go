package session

import "github.com/SystemBuilders/LocKey/internal/lockclient/id"

// Session captures all necessary parameters necessary to
// describe a session with the lockservice in the lockclient.
type Session interface {
	// SessionID is the unique ID that represents this session.
	// This will be used in every transaction for validating the user.
	SessionID() id.ID
	// ClientID is the ID of the client that will be created when
	// the client is created. This acts as a second layer check along
	// with the sessionID.
	ClientID() id.ID
	// ProcessID the unique ID assigned for the process by the client.
	// This will be the third layer check in the security mechanism.
	ProcessID() id.ID
}
