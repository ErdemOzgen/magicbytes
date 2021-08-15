package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/ErdemOzgen/magicbytes/magicbytes"
	//"github.com/ErdemOzgen/magicbytes"
)

//=========================Color Const ==========================================
const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[1;32m%s\033[0m"
)

//========================Color Const End=========================================
// ================== Health check and profile tools==============================

func cpuprofileCheck() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()
}

func memoryprofileCheck() {
	f, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	runtime.GC()    // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

func httpServerRuntimeCheck() {
	fmt.Printf(DebugColor, "wget -O trace.out http://localhost:6060/debug/pprof/trace?seconds=5 || go tool trace trace.out")
	fmt.Println("")
	fmt.Printf(NoticeColor, "======================================================================================================")
	fmt.Println("")
	log.Println(http.ListenAndServe("localhost:6060", nil))
}

// ================== Health check and profile tools END==============================
//============================Banner==================================================
type ascii struct {
	banner     string
	bannerline string
}

func (d *ascii) printBanner() {
	fmt.Printf(InfoColor, d.banner)
	fmt.Println("")
}

func (d *ascii) printBannerline() {
	fmt.Printf(NoticeColor, d.bannerline)
	fmt.Println("")
}

//=====================BANNER END========================================================
func main() {
	//=====================BANNER Struct=================================================
	welcomeScreenBanner := ascii{
		banner: `
		 ███    ███  █████   ██████  ██  ██████ ██████  ██    ██ ████████ ███████ ███████ 
		 ████  ████ ██   ██ ██       ██ ██      ██   ██  ██  ██     ██    ██      ██      
		 ██ ████ ██ ███████ ██   ███ ██ ██      ██████    ████      ██    █████   ███████ 
		 ██  ██  ██ ██   ██ ██    ██ ██ ██      ██   ██    ██       ██    ██           ██ 
		 ██      ██ ██   ██  ██████  ██  ██████ ██████     ██       ██    ███████ ███████ `,
		bannerline: "======================================================================================================",
	}
	welcomeScreenBanner.printBannerline()
	welcomeScreenBanner.printBanner()
	welcomeScreenBanner.printBannerline()

	//=====================BANNER Struct END========================================================

	var print = fmt.Println // print shortcut for python like
	var directorypath string
	healthCheck := false // set true for profiling
	//healthCheck := true

	flag.StringVar(&directorypath, "directorypath", "src/", "Default path is setted as current directory")

	fmt.Printf(DebugColor, "pwd ===> ")
	fmt.Println(directorypath)
	ctx, cancel := context.WithCancel(context.Background())
	//print("flag:", healthCheck)
	go func() {
		//log.Println(http.ListenAndServe("localhost:6060", nil))
		if healthCheck == true {
			//https://blog.golang.org/pprof
			fmt.Println("Entering profiling mod")
			//welcomeScreenBanner.printBannerline()
			go cpuprofileCheck()
			go memoryprofileCheck()
			go httpServerRuntimeCheck()
		} else {
			fmt.Printf(DebugColor, "Staring without profiling mod....")
			print()
			//welcomeScreenBanner.printBannerline()
		}
	}()
	//time.Sleep(time.Millisecond * 1)
	m := []*magicbytes.Meta{
		//	{Type: "png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}},
		//   89 50 4E 47 0D 0A 1A 0A
		//   https://cryptii.com/pipes/integer-encoder
		{Type: "png", Offset: 0, Bytes: []byte{0x89, 0x50, 0x4E, 0x47}},
		{Type: "jpeg", Offset: 0, Bytes: []byte{0xff, 0xd8, 0xff, 0xe0}},
		{Type: "pcap", Offset: 0, Bytes: []byte{0xa1, 0xb2, 0xc3, 0xd4}},
		{Type: "pcap2", Offset: 0, Bytes: []byte{0xd4, 0xc3, 0xb2, 0xa1}},
		{Type: "pdf", Offset: 0, Bytes: []byte{0x25, 0x50, 0x44, 0x46, 0x2d}},
		{Type: "tar", Offset: 0, Bytes: []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x00, 0x30, 0x30}},
		//{Type: "tar2", Offset: 0, Bytes: []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x20, 0x20, 0x00}},
		//{Type: "DICOM", Offset: 0x80, Bytes: []byte{0x44, 0x49, 0x43, 0x4D}},
		//{Type: "jpg", Offset: 0, Bytes: []byte{0xFF, 0xD8, 0xFF, 0xDB}},
		{Type: "jpg", Offset: 0, Bytes: []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01}},
		nil, //must handle nil value
	}
	welcomeScreenBanner.printBannerline()
	if err := magicbytes.Search(ctx, directorypath, m, func(path, metaType string) bool {
		fmt.Println(path)

		return false
	}); err != nil {
		log.Fatal(err)
	}
	fmt.Println(m)

	//Add defer when removing the below #change
	cancel()

	/* #remove
	fmt.Println("Waiting for input:")
	var input string
	fmt.Scanln(&input)
	*/
	chanForExit := make(chan os.Signal)
	signal.Notify(chanForExit, os.Interrupt, syscall.SIGTERM)
	func() {
		<-chanForExit

		fmt.Println("\r- Ctrl+C pressed in Terminal")
		fmt.Println("Serve has been shut down")
		os.Exit(0)
	}()
}
