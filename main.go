package main

/*
botões grande:
NoteOn at channel (0-7): value: 127 controller: 35
NoteOn at channel (0-7): value: 0 controller: 35

Canais analógicos:
NoteOn at channel (0-7): value: (0-126) controller: 15

Mute:
NoteOn at channel (0-6): value: 127 controller: 27
NoteOn at channel (0-6): value: 0 controller: 27

Listem:
NoteOn at channel (0-7): value: 127 controller: 24
NoteOn at channel (0-7): value: 0 controller: 24

*/

import (
	"fmt"
	"github.com/micmonay/keybd_event"
	"gitlab.com/gomidi/rtmididrv"
	"log"
	"os/exec"
	"runtime"
	"time"

	"gitlab.com/gomidi/midi"
	. "gitlab.com/gomidi/midi/midimessage/channel" // (Channel Messages)
	"gitlab.com/gomidi/midi/reader"
)

// This example reads from the first input port
func main() {
	drv, err := rtmididrv.New()
	must(err)

	// make sure to close the driver at the end


	ins, err := drv.Ins()
	must(err)

	in := ins[0]

	fmt.Printf("opening MIDI IN Port %v\n", in)
	must(in.Open())

	// to disable logging, pass mid.NoLogger() as option
	rd := reader.New(
		reader.NoLogger(),
		// print every message
		reader.Each(func(pos *reader.Position, msg midi.Message) {

			// inspect
			//fmt.Println(msg)

			switch v := msg.(type) {
			case Aftertouch:
				fmt.Printf("Aftertouch at channel %v: Pressure: %v\n", v.Channel(), v.Pressure())
			case ControlChange:
				fmt.Printf("NoteOn at channel %v: value: %v controller: %v\n", v.Channel(), v.Value(), v.Controller())
			case PolyAftertouch:
				fmt.Printf("NoteOn at channel %v: key: %v pressure: %v\n", v.Channel(), v.Key(), v.Pressure())
			case ProgramChange:
				fmt.Printf("NoteOn at channel %v: program: %v\n", v.Channel(), v.Program())
			case Pitchbend:
				fmt.Printf("NoteOn at channel %v: value: %v abs value: %v\n", v.Channel(), v.Value(), v.AbsValue())
			case NoteOn:

				b, e := exec.Command("/usr/bin/open", "-a", "GoLand-M1" ).Output()
				if e != nil {
					panic(e)
				}
				fmt.Printf("%s", b)

				kb, err := keybd_event.NewKeyBonding()
				if err != nil {
					log.Fatal(err)
				}

				// For linux, need to wait 2 seconds
				if runtime.GOOS == "linux" {
					time.Sleep(2 * time.Second)
				}

				time.Sleep(5*time.Second)
				//set keys
				kb.SetKeys(keybd_event.VK_H)
				//kb.HasCTRL(true)
				kb.HasSuper(true)


				//launch
				err = kb.Launching()
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("NoteOn at channel %v: key: %v velocity: %v\n", v.Channel(), v.Key(), v.Velocity())
			case NoteOff:
				fmt.Printf("NoteOff at channel %v: key: %v\n", v.Channel(), v.Key())
			case NoteOffVelocity:
				fmt.Printf("NoteOff at channel %v: key: %v velocity: %v\n", v.Channel(), v.Key(), v.Velocity())
			}
		}),
	)

	// listen for MIDI
	err = rd.ListenTo(in)
	must(err)

	time.Sleep(60 * 60 * time.Second)
	err = in.StopListening()
	must(err)

	err = in.Close()
	must(err)

	err = drv.Close()
	must(err)

	fmt.Printf("closing MIDI IN Port %v\n", in)
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
