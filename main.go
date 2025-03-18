package main

import (
        "fmt"
        "io"
        "os/exec"
        "runtime"
        "strings"
        "time"
)

func main() {
        w, h := 80, 24
        if runtime.GOOS == "windows" {
                if out, err := exec.Command("powershell", "(Get-Host).UI.RawUI.WindowSize.Width").Output(); err == nil {
                        fmt.Sscanf(string(out), "%d", &w)
                }
                if out, err := exec.Command("powershell", "(Get-Host).UI.RawUI.WindowSize.Height").Output(); err == nil {
                        fmt.Sscanf(string(out), "%d", &h)
                }
        } else if out, err := exec.Command("stty", "size").Output(); err == nil {
                fmt.Sscanf(string(out), "%d %d", &h, &w)
        }

        cmd := exec.Command("ffmpeg", "-i", "bad_apple.mp4", "-vf", fmt.Sprintf("fps=30,scale=%d:%d", w/2, h),
                "-pix_fmt", "gray", "-f", "image2pipe", "-vcodec", "rawvideo", "-")
        pipe, _ := cmd.StdoutPipe()
        cmd.Start()

        chars := " .:-=+*#%@"
        buffer := make([]byte, w/2*h)
        fmt.Print("\033[?25l\033[2J")
        start := time.Now()
        frame := 0

        for {
                if _, err := io.ReadFull(pipe, buffer); err != nil {
                        break
                }

                var sb strings.Builder
                sb.Grow(w * (h + 1))
                sb.WriteString("\033[H")

                for y := 0; y < h; y++ {
                        for x := 0; x < w/2; x++ {
                                idx := int(buffer[y*w/2+x]) * (len(chars) - 1) / 255
                                sb.WriteByte(chars[idx])
                                sb.WriteByte(chars[idx])
                        }
                        sb.WriteByte('\n')
                }

                fmt.Print(sb.String())
                frame++

                target := start.Add(time.Duration(frame) * time.Second / 30)
                if sleep := target.Sub(time.Now()); sleep > 0 {
                        time.Sleep(sleep)
                }
        }

        fmt.Print("\033[?25h")
        cmd.Wait()
}
