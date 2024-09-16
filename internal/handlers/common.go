package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func renderPartial(c *gin.Context, partial string, title string) {

	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	tmpl, err := template.ParseFiles("web/templates/" + partial)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering partial: %v", err)
		return
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, gin.H{"User": user}); err != nil {
		c.String(http.StatusInternalServerError, "Error executing template: %v", err)
		return
	}
	c.HTML(http.StatusOK, "base.html", gin.H{"Title": title, "User": user, "Content": template.HTML(buf.String())})
}

// HandleTerm ...
func HandleTerm(c *gin.Context) {
	wsConn, err := upgradeWebSocket(c)
	if err != nil {
		log.Println(err)
		return
	}
	defer wsConn.Close()

	config, err := createSSHConfig()
	if err != nil {
		log.Println(err)
		return
	}

	sshConn, err := establishSSHConnection(config)
	if err != nil {
		log.Println(err)
		return
	}
	defer sshConn.Close()

	if err := handleSSHSession(wsConn, sshConn); err != nil {
		log.Println(err)
	}
}

func upgradeWebSocket(c *gin.Context) (*websocket.Conn, error) {
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, fmt.Errorf("upgrade error: %v", err)
	}
	return wsConn, nil
}

func createSSHConfig() (*ssh.ClientConfig, error) {
	key, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/id_rsa")
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: os.Getenv("USER"), // Assuming the environment variable USER is set to the SSH user
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return config, nil
}

func establishSSHConnection(config *ssh.ClientConfig) (*ssh.Client, error) {
	// TODO: change to config
	host := "localhost"
	port := "22"
	hostport := fmt.Sprintf("%s:%s", host, port)

	sshConn, err := ssh.Dial("tcp", hostport, config)
	if err != nil {
		return nil, fmt.Errorf("SSH dial error: %v", err)
	}
	return sshConn, nil
}

func handleSSHSession(wsConn *websocket.Conn, sshConn *ssh.Client) error {
	session, err := sshConn.NewSession()
	if err != nil {
		return fmt.Errorf("SSH session error: %v", err)
	}
	defer session.Close()

	sshOut, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("STDOUT pipe error: %v", err)
	}

	sshIn, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("STDIN pipe error: %v", err)
	}

	if err := session.RequestPty("xterm", 80, 40, ssh.TerminalModes{}); err != nil {
		return fmt.Errorf("request PTY error: %v", err)
	}

	if err := session.Shell(); err != nil {
		return fmt.Errorf("start shell error: %v", err)
	}

	go func() {
		defer session.Close()
		buf := make([]byte, 1024)
		for {
			n, err := sshOut.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Println("Read from SSH stdout error:", err)
				}
				return
			}
			if n > 0 {
				err = wsConn.WriteMessage(websocket.BinaryMessage, buf[:n])
				if err != nil {
					log.Println("Write to WebSocket error:", err)
					return
				}
			}
		}
	}()

	for {
		messageType, p, err := wsConn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Println("Read from WebSocket error:", err)
			}
			return err
		}
		if messageType == websocket.BinaryMessage || messageType == websocket.TextMessage {
			_, err = sshIn.Write(p)
			if err != nil {
				log.Println("Write to SSH stdin error:", err)
				return err
			}
		}
	}
}
