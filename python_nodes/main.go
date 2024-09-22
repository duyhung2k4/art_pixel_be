package pythonnodes

import (
	"log"
	"os/exec"
)

func RunPythonServer(port string) {
	// Truyền port làm argument cho lệnh Python
	cmd := exec.Command("python3", "python_nodes/index.py", port)

	// Bắt đầu chạy lệnh
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start Python server: %v", err)
	}

	// Chờ Python server kết thúc
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Python server exited with error: %v", err)
	}
}
