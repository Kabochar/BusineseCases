// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

// Package main implements an example TURN server supporting TCP
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pion/turn/v2"
	"github.com/spf13/cast"

	"remote-assistant/internal/config"
)

func main() {
	config.LoadConfig()

	publicIP := os.Getenv("TURN_ADDR")
	port := os.Getenv("TURN_PORT")
	users := fmt.Sprintf("%s=%s", os.Getenv("TURN_ACCOUNT"), os.Getenv("TURN_ACC_PWD"))
	realm := flag.String("realm", "pion.ly", "Realm (defaults to \"pion.ly\")")
	minPortRange, maxPortRange := os.Getenv("TURN_MIN_PORT_RANGE"), os.Getenv("TURN_MAX_PORT_RANGE")
	flag.Parse()

	if len(publicIP) == 0 {
		log.Fatalf("'public-ip' is required")
	} else if len(users) == 0 {
		log.Fatalf("'users' is required")
	}

	// Create a TCP listener to pass into pion/turn
	// pion/turn itself doesn't allocate any TCP listeners, but lets the user pass them in
	// this allows us to add logging, storage or modify inbound/outbound traffic
	tcpListener, err := net.Listen("tcp4", "0.0.0.0:"+cast.ToString(port))
	if err != nil {
		log.Panicf("Failed to create TURN server listener: %s", err)
	}

	// Create a UDP listener to pass into pion/turn
	// pion/turn itself doesn't allocate any UDP sockets, but lets the user pass them in
	// this allows us to add logging, storage or modify inbound/outbound traffic
	udpListener, err := net.ListenPacket("udp4", "0.0.0.0:"+cast.ToString(port))
	if err != nil {
		log.Panicf("Failed to create TURN server listener: %s", err)
	}

	// Cache -users flag for easy lookup later
	// If passwords are stored they should be saved to your DB hashed using turn.GenerateAuthKey
	usersMap := map[string][]byte{}
	// for _, kv := range regexp.MustCompile(`(\w+)=(\w+)`).FindAllStringSubmatch(users, -1) {
	// 	usersMap[kv[1]] = turn.GenerateAuthKey(kv[1], *realm, kv[2])
	// }
	// use custom define account info
	accountInfo := strings.Split(users, "=")
	usersMap[accountInfo[0]] = turn.GenerateAuthKey(accountInfo[0], *realm, accountInfo[1])

	s, err := turn.NewServer(turn.ServerConfig{
		Realm: *realm,
		// Set AuthHandler callback
		// This is called every time a user tries to authenticate with the TURN server
		// Return the key for that user, or false when no user is found
		AuthHandler: func(username string, realm string, srcAddr net.Addr) ([]byte, bool) {
			if key, ok := usersMap[username]; ok {
				return key, true
			}
			return nil, false
		},
		// PacketConnConfigs is a list of UDP Listeners and the configuration around them
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: udpListener,
				RelayAddressGenerator: &turn.RelayAddressGeneratorPortRange{
					RelayAddress: net.ParseIP(publicIP), // Claim that we are listening on IP passed by user (This should be your Public IP)
					Address:      "0.0.0.0",             // But actually be listening on every interface
					MinPort:      cast.ToUint16(minPortRange),
					MaxPort:      cast.ToUint16(maxPortRange),
				},
			},
		},
		// ListenerConfig is a list of Listeners and the configuration around them
		ListenerConfigs: []turn.ListenerConfig{
			{
				Listener: tcpListener,
				RelayAddressGenerator: &turn.RelayAddressGeneratorStatic{
					RelayAddress: net.ParseIP(publicIP),
					Address:      "0.0.0.0",
				},
			},
		},
	})
	if err != nil {
		log.Panic(err)
	}
	log.Println("start listen on", publicIP, port)

	// Block until user sends SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	if err = s.Close(); err != nil {
		log.Panic(err)
	}
}
