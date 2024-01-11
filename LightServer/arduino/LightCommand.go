//
//  LightCommand.go
//  GitHubLight
//
//  Created by Karl Kraft on 12/30/2023
//  Copyright 2023-2024 Karl Kraft. All rights reserved
//

package arduino

import (
	"fmt"
	"net"
)

type LightCommand struct {
	Red1   uint8
	Green1 uint8
	Blue1  uint8

	Red2   uint8
	Green2 uint8
	Blue2  uint8

	Red3   uint8
	Green3 uint8
	Blue3  uint8
}

func (packet LightCommand) Send(connection net.Conn) {
	buffer := make([]byte, 12)
	buffer[0] = 0x34
	buffer[1] = 0x12

	buffer[2] = packet.Red1
	buffer[3] = packet.Green1
	buffer[4] = packet.Blue1

	buffer[5] = packet.Red2
	buffer[6] = packet.Green2
	buffer[7] = packet.Blue2

	buffer[8] = packet.Red3
	buffer[9] = packet.Green3
	buffer[10] = packet.Blue3

	_, err := connection.Write(buffer)
	if err != nil {
		fmt.Printf("Couldn't send LightCommand packet %v", err)
	}

}
