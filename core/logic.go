package core

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type port struct {
}

func (p port) Upload(filename string) error {

	shell := fmt.Sprintf("-v error -select_streams v:0 -show_entries stream=width,height -of csv=s=x:p=0 %s.mp4", "input")

	cmd := exec.Command("ffprobe", strings.Split(shell, " ")...)
	stdout, err := cmd.Output()

	if err != nil {
		return err
	}

	size := strings.Trim(string(stdout), "\n")

	dim := strings.Split(size, "x")

	W, _ := strconv.Atoi(dim[0])
	H, _ := strconv.Atoi(dim[1])

	fmt.Println(size)

	var scale []string
	switch {
	case H > 850:
		scale = append(scale, ":1080")
		scale = append(scale, splitter(W, H, 720.0))
		scale = append(scale, splitter(W, H, 480.0))
		scale = append(scale, splitter(W, H, 240.0))
	case H > 650:
		scale = append(scale, ":720")
		scale = append(scale, splitter(W, H, 480.0))
		scale = append(scale, splitter(W, H, 240.0))
	case H > 350:
		scale = append(scale, ":480")
		scale = append(scale, splitter(W, H, 240.0))
	case H > 180:
		scale = append(scale, ":240")
	}

	fmt.Println(scale)

	if len(scale) == 0 {
		fmt.Println("No scale")
		return errors.New("no scale")
	}

	if err := os.Mkdir("../video/"+filename, 0755); err != nil {
		fmt.Println(err)
		return err
	}

	for i, x := range scale {
		go doHLS(filename, x, i)
	}

	return nil
}

func splitter(w, h int, div float32) string {

	rem := float32(h) / div

	w = int(float32(w) / rem)

	if w % 2 != 0 {
		w += 1
	}

	return fmt.Sprintf("%d:%d", w, int(div))
}

func doHLS(f, x string, i int)  {

	var shell string

	dim := strings.Split(x, ":")[1]

	if i == 0 {
		shell = fmt.Sprintf("-re -i %s.mp4 -codec copy -map 0 -f segment -segment_list ../video/%s/%s-list.m3u8 -segment_list_flags +live -segment_time 10 -segment_format mpegts ../video/%s/%s-%s.ts", "input", f, dim, f, dim, "%03d")
	} else {
		shell = fmt.Sprintf("-re -i %s.mp4 -vf scale=%s -preset veryslow -crf 28 -map 0 -f segment -segment_list ../video/%s/%s-list.m3u8 -segment_list_flags +live -segment_time 10 -segment_format mpegts ../video/%s/%s-%s.ts",
			"input", x, f, dim, f, dim, "%03d")
	}

	fmt.Println(shell, "\n", i, "Started", x, "...")

	cmd := exec.Command("ffmpeg", strings.Split(shell, " ")...)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(string(stdout), err)
		return
	}

	fmt.Println(i, string(stdout))
}

func (p port) Stream() {
	panic("implement me")
}

func NewService() Service {

	return &port{}
}
