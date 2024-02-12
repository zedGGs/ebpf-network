package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"

	manager "github.com/gojue/ebpfmanager"
)

//go:embed ebpf/bin/probe.o
var _bytecode []byte


func trigger() {
	fmt.Println("Generating some network traffic to the probes ...")
	_, _ = http.Get("https://example.com")
}

func main() {
	m := &manager.Manager{
		Probes: []*manager.Probe{
			{
				Section:       "xdp/ingress",
				EbpfFuncName:  "egress_cls_func",
				Ifname:        "wlp3s0",
				XDPAttachMode: manager.XdpAttachModeSkb,
			},
		},
	}
	err := m.Init(bytes.NewReader(_bytecode))
	if err != nil {
		fmt.Println(err)
		return
	}
	// Start the manager
	if err := m.Start(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("successfully started, head over to /sys/kernel/debug/tracing/trace_pipe")

	trigger()

	// Close the manager
	if err := m.Stop(manager.CleanAll); err != nil {
		fmt.Println(err)
	}
}