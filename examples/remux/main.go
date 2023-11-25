package main

import (
	"flag"
	"fmt"
	"github.com/katana-project/mux"
	"io"
)

var mp4Muxer = mux.MustFindMuxer("MP4", "mp4", "video/mp4")

func main() {
	var (
		input  = flag.String("i", "", "Input file")
		output = flag.String("o", "./output.mp4", "Output file")
	)
	flag.Parse()

	if *input == "" { // zero value
		fmt.Println("No input file was provided, exiting.")
		return
	}

	inCtx, err := mux.NewInputContext(*input)
	if err != nil {
		fmt.Printf("Error opening input context: %s\n", err.Error())
		return
	}
	defer inCtx.Close()

	outCtx, err := mux.NewOutputContext(mp4Muxer, *output)
	if err != nil {
		fmt.Printf("Error opening output context: %s\n", err.Error())
		return
	}
	defer outCtx.Close()

	for i, inStream := range inCtx.Streams() {
		outStream := outCtx.NewStream(nil)
		if err := inStream.CopyParameters(outStream); err != nil {
			fmt.Printf("Error copying stream %d parameters: %s\n", i, err.Error())
			return
		}
	}

	pkt := mux.NewPacket()
	defer pkt.Close()

	for {
		err := inCtx.ReadFrame(pkt)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading frame: %s\n", err.Error())
			}

			break
		}

		streamIdx := pkt.StreamIndex()
		pkt.Rescale(
			inCtx.Stream(streamIdx).TimeBase(),
			outCtx.Stream(streamIdx).TimeBase(),
		)
		pkt.ResetPos()

		if err := outCtx.WriteFrame(pkt); err != nil {
			fmt.Printf("Error writing frame: %s\n", err.Error())
			break
		}
	}
}
