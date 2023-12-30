//
//  LightCommand.go
//  GitHubLight
//
//  Created by Karl Kraft on 12/30/2023
//  Copyright 2023 Karl Kraft. All rights reserved.
//

package GitHubLight

import (
	"fmt"
	"net"
)

type LightCommand struct {
	Start  uint8
	Length uint8
	Red    uint8
	Green  uint8
	Blue   uint8
}

func (packet LightCommand) Send(connection net.Conn) {
	buffer := make([]byte, 8)
	buffer[0] = 0x34
	buffer[1] = 0x12
	buffer[2] = packet.Start
	buffer[3] = packet.Length
	buffer[4] = packet.Red
	buffer[5] = packet.Green
	buffer[6] = packet.Blue
	buffer[7] = 0x00

	_, err := connection.Write(buffer)
	if err != nil {
		fmt.Printf("Couldn't send LightCommand packet %v", err)
	}

}
