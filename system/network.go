package system

import (
	"github.com/pytimer/win-netstat"
)

func HasProcessAnyOpenConnections(processId int) (bool, error) {
	conns, err := GetConnectionsWithPid(processId)
	if err != nil {
		return false, err
	}

	return len(conns) > 0, nil
}

func GetConnectionsWithPid(processId int) ([]winnetstat.NetStat, error) {
	var connections []winnetstat.NetStat

	var connTypes = [...]string{"tcp4", "udp4", "tcp6", "udp6"}

	for _, connType := range connTypes {
		conns, err := winnetstat.ConnectionsWithPid(connType, processId)
		if err != nil {
			return nil, err
		}

		connections = append(connections, conns...)
	}

	return connections, nil
}
