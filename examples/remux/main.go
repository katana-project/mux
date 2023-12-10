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

	var (
		streams       = inCtx.Streams()
		streamMapping = make([]int, len(streams))
	)
	for i, inStream := range streams {
		codec := inStream.Codec()
		if !mp4Muxer.SupportsCodec(codec) {
			streamMapping[i] = -1
			continue
		}

		outStream := outCtx.NewStream(codec)
		if err := inStream.CopyParameters(outStream); err != nil {
			fmt.Printf("Error copying stream %d parameters: %s\n", i, err.Error())
			return
		}

		streamMapping[i] = outStream.Index()
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
		if remapId := streamMapping[streamIdx]; remapId > -1 {
			pkt.SetStreamIndex(remapId)

			pkt.Rescale(
				inCtx.Stream(streamIdx).TimeBase(),
				outCtx.Stream(remapId).TimeBase(),
			)
			pkt.ResetPos()
			if err := outCtx.WriteFrame(pkt); err != nil {
				fmt.Printf("Error writing frame: %s\n", err.Error())
				break
			}
		} else {
			if err := pkt.Clear(); err != nil {
				fmt.Printf("Error clearing frame: %s\n", err.Error())
				break
			}
		}
	}
}
