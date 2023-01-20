package main

import (
	"fmt"
	"server/connector"
	"server/server"
)

func main() {
	//app := &cli.App{
	//	Name:    "UroborosQ's simple distributed file system - server",
	//	Usage:   "If you want to, you can use it",
	//	Version: "v0.0.1",
	//	Commands: []*cli.Command{
	//		{
	//			Name:     "add",
	//			Usage:    "Add file to the volume",
	//			Category: "file",
	//			Flags: []cli.Flag{
	//				&cli.StringFlag{
	//					Name:  "source",
	//					Usage: "Path to your local file",
	//				},
	//				&cli.StringFlag{
	//					Name:  "target",
	//					Usage: "Partial path on the volume.",
	//				},
	//			},
	//		},
	//		{
	//			Name:     "remove",
	//			Usage:    "Remove file from the volume",
	//			Category: "file",
	//			Flags: []cli.Flag{
	//				&cli.StringFlag{
	//					Name:  "target",
	//					Usage: "Partial path on the volume",
	//				},
	//			},
	//		},
	//		{
	//			Name:     "add",
	//			Usage:    "Add node to the system",
	//			Category: "node",
	//			Flags: []cli.Flag{
	//				&cli.StringFlag{
	//					Name:  "ip",
	//					Usage: "ip address of the node",
	//				},
	//				&cli.IntFlag{
	//					Name:  "port",
	//					Usage: "port of the node",
	//				},
	//				&cli.Int64Flag{
	//					Name:  "size",
	//					Usage: "Max amount of bytes, which client can use on the node",
	//				},
	//			},
	//		},
	//		{
	//			Name:     "remove",
	//			Usage:    "Remove node from the system",
	//			Category: "node",
	//		},
	//		{
	//			Name:     "clean",
	//			Usage:    "Move all files from one node to others",
	//			Category: "node",
	//		},
	//		{
	//			Name:  "balance",
	//			Usage: "balance files between nodes",
	//		},
	//	},
	//	Action: func(cCtx *cli.Context) error {
	//		fmt.Println("Welcome to the distributed file system client!")
	//		return nil
	//	},
	//}
	//
	//if err := app.Run(os.Args); err != nil {
	//	log.Fatal(err)
	//}
	s := server.CreateServer("/home/uroborosq/lol")
	_, err := s.AddNode("localhost", "12345", 123, connector.Tcp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = s.AddFile("/home/uroborosq/Изображения/Снимки экрана/Снимок экрана от 2022-12-20 13-58-34.png", "1.png" +
		"")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
