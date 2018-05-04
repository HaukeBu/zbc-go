package zbmsgpack

import (
	"fmt"
	"strings"
)

type RaftMember struct {
	Host string `msgpack:"host"`
	Port int    `msgpack:"port"`
}

func (m RaftMember) String() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

type RaftEvent struct {
	Members []RaftMember `msgpack:"members"`
}

func (e RaftEvent) String() string {
	var members []string
	for _, member := range e.Members {
		members = append(members, member.String())
	}
	return fmt.Sprintf("Raft Event [members: [%s]]", strings.Join(members, ", "))
}
