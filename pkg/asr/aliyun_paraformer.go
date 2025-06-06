package asr

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// 配置参数
const (
	AliyunASRWebSocketURL = "wss://dashscope.aliyuncs.com/api-ws/v1/inference" // 替换为实际的阿里云 WebSocket 地址
	APIKey                = ""              // 推荐用临时Token，安全性更高
	SampleRate            = 16000                                              // 采样率
	ChunkSize             = 1024                                               // 每次发送的音频字节数
)

// 发送启动任务的消息
func sendRunTask(conn *websocket.Conn, taskID string) error {
	runTaskMsg := map[string]interface{}{
		"header": map[string]interface{}{
			"action":    "run-task",
			"task_id":   taskID,
			"streaming": "duplex",
			"appkey":    APIKey,
		},
		"payload": map[string]interface{}{
			"task_group": "audio",
			"task":       "asr",
			"function":   "recognition",
			"model":      "paraformer-realtime-v2",
			"parameters": map[string]interface{}{
				"sample_rate": SampleRate,
				"format":      "wav",
			},
			"input": map[string]interface{}{},
		},
	}
	return conn.WriteJSON(runTaskMsg)
}

// 发送结束任务的消息
func sendFinishTask(conn *websocket.Conn, taskID string) error {
	finishTaskMsg := map[string]interface{}{
		"header": map[string]interface{}{
			"action":    "finish-task",
			"task_id":   taskID,
			"streaming": "duplex",
		},
		"payload": map[string]interface{}{
			"input": map[string]interface{}{},
		},
	}
	return conn.WriteJSON(finishTaskMsg)
}

// 发送音频流
func sendAudioStream(conn *websocket.Conn, audioFile string) error {
	file, err := os.Open(audioFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buf := make([]byte, ChunkSize)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				return err
			}
			time.Sleep(100 * time.Millisecond) // 按文档建议每100ms发送一次
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// 解析识别结果
func receiveASRResult(conn *websocket.Conn) (string, error) {
	var resultText bytes.Buffer
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return "", err
		}
		// 解析 JSON
		var resp map[string]interface{}
		if err := json.Unmarshal(message, &resp); err != nil {
			continue
		}
		// 检查事件类型
		if header, ok := resp["header"].(map[string]interface{}); ok {
			if action, ok := header["action"].(string); ok {
				if action == "result-generated" {
					if payload, ok := resp["payload"].(map[string]interface{}); ok {
						if output, ok := payload["output"].(map[string]interface{}); ok {
							if sentence, ok := output["sentence"].(map[string]interface{}); ok {
								if text, ok := sentence["text"].(string); ok {
									resultText.WriteString(text)
								}
								// 判断是否为一句话的结尾
								if end, ok := sentence["sentence_end"].(bool); ok && end {
									break
								}
							}
						}
					}
				}
				// 任务结束
				if action == "task-finished" {
					break
				}
			}
		}
	}
	return resultText.String(), nil
}

// 对外主函数：传入音频文件路径，返回识别文本
func AliyunASR(audioFile string) (string, error) {
	taskID := fmt.Sprintf("task-%d", time.Now().UnixNano())
	header := make(map[string][]string)
	header["Authorization"] = []string{APIKey}

	conn, _, err := websocket.DefaultDialer.Dial(AliyunASRWebSocketURL, header)
	if err != nil {
		return "", fmt.Errorf("websocket connect failed: %v", err)
	}
	defer conn.Close()

	// 1. 发送 run-task
	if err := sendRunTask(conn, taskID); err != nil {
		return "", fmt.Errorf("send run-task failed: %v", err)
	}

	// 2. 发送音频流
	if err := sendAudioStream(conn, audioFile); err != nil {
		return "", fmt.Errorf("send audio stream failed: %v", err)
	}

	// 3. 发送 finish-task
	if err := sendFinishTask(conn, taskID); err != nil {
		return "", fmt.Errorf("send finish-task failed: %v", err)
	}

	// 4. 接收识别结果
	text, err := receiveASRResult(conn)
	if err != nil {
		return "", fmt.Errorf("receive asr result failed: %v", err)
	}
	return text, nil
}
