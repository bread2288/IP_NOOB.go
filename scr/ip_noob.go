package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

/*
	IP_NOOB.go - простой сканер портов
	написал: burger
*/

/*
	\
*/

/*
Настройки
*/

var hosts []string // IP-адреса
var ports []string // Порты
var timeout = 1    // Время ожидания (sec)

/*\
   \
	\
1.  |
	|
   /
*/

func is_valid(ip string, ports []string, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, port := range ports {
		address := ip + ":" + port
		conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
		if err == nil {
			conn.Close()
			// сообщение если активно
			names, err := net.LookupAddr(ip)
			if err != nil {
				color.Yellow("[?] -> IP: %s Не имеет имени.", ip)
			}
			for _, name := range names {
				color.Green("[+] -> IP: %s | ИМЯ: %s", ip, name)
			}
			color.Green("[+] -> %s:%s - Активный.", ip, port)
			write("result", "["+ip+":"+port+"] - Активный ("+time.Now().Format("02-01-2006 15:04:05")+")")
		}
	}
}

/*
	|
*/

func manage_ip(ip net.IP) net.IP {
	for j := len(ip) - 1; j >= 0; j-- {
		if ip[j] < 255 {
			ip[j]++
			return ip
		}
		ip[j] = 0
	}
	return ip
}

/*
	|
*/

func read(filename string) []string {
	// Открытие файла
	fike, err := os.Open(filename)
	if err != nil {
		color.Red("[!] -> Не найден файл '%s'", filename)
		color.Red("[!] -> Вам необходимо создать файл '%s'", filename)
		os.Exit(0)
	}
	defer fike.Close()

	// Читение
	var lines []string
	scanner := bufio.NewScanner(fike)
	for scanner.Scan() {
		if filename == "hosts" {
			if strings.Contains(scanner.Text(), "-") {
				line := strings.Split(scanner.Text(), "-")
				start_ip := net.ParseIP(line[0])
				end_ip := net.ParseIP(line[1])

				for ip := start_ip; !ip.Equal(end_ip); ip = manage_ip(ip) {
					lines = append(lines, ip.String())
				}
				lines = append(lines, string(end_ip))
			} else {
				lines = append(lines, scanner.Text())
			}
		} else if filename == "ports" {
			if strings.Contains(scanner.Text(), "-") {
				line := strings.Split(scanner.Text(), "-")
				start_p, err1 := strconv.Atoi(line[0])
				if err1 != nil {
					panic(err1)
				}
				end_p, err2 := strconv.Atoi(line[1])
				if err2 != nil {
					panic(err2)
				}

				for i := start_p; i < end_p+1; i++ {
					lines = append(lines, strconv.Itoa(i))
				}
			} else {
				lines = append(lines, scanner.Text())
			}
		}
	}

	return lines
}

/*
	\
*/

func write(filename string, text string) {
	// Открытие файла
	fike, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		color.Red("[!] -> Не найден файл '%s'", filename)
		color.Red("[!] -> Вам необходимо создать файл '%s'", filename)
		os.Exit(0)
	}
	defer fike.Close()

	// Запись
	_, err = fike.WriteString(text + "\n")
	if err != nil {
		color.Red("[!] -> Не удалось записать новое значение.")
	}
}

/*
	|
*/

func printer_use(ip string, port string, text string) {
	/*
		PRINTER_USE - отправка TCP запросов на печать.
	*/
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		color.Red("[!] -> Не удалось подключится.")
		panic(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(text + "\n"))
	if err != nil {
		color.Red("[!] -> Не удалось отправить текст.")
		panic(err)
	}

	color.Green("[+] -> Текст отправлен.")
}

/*\
   \
	\
2.  |
	|
   /
*/

func main() {
	action := flag.String("action", "scan", "")
	flag.Parse()

	if *action == "scan" {
		// Начальное сообщение
		color.Blue("#| IP_NOOB.go🗿 |#")
		color.Blue("написал: burger")
		color.Blue("")
		color.Blue("[/] -> Загрузка файлов...")

		// Перезапись файла (очистка)
		_, err := os.Create("result")
		if err != nil {
			color.Red("[!] -> Не удалось перезаписать файл.")
		}
		hosts = read("hosts") // загрузка IPs
		ports = read("ports") // загрузка порты

		color.Green("[+] -> Файлы загружены.")
		color.Blue("")
		color.Blue("[√] -> Найдено IP-адресов: %d", len(hosts))
		color.Blue("[√] -> Найдено портов: %d", len(ports))
		color.Blue("[√] -> Время ожидания: %ds", timeout)
		color.White("Нажмите 'Enter' для запуска...")
		fmt.Scanln()

		// Работа
		var wg sync.WaitGroup
		for _, ip := range hosts {

			wg.Add(1)
			go is_valid(ip, ports, &wg)

		}

		wg.Wait()
	} else if *action == "printer" {
		var printer_ip string
		var printer_port string
		var send_text string
		var count int

		color.Blue("###| PRINTER_USE🖨 |###")
		color.White("Введите IP принтера: ")
		fmt.Scanln(&printer_ip)

		color.White("Введите PORT принтера: ")
		fmt.Scanln(&printer_port)

		color.White("Введите Текст (тест): ")
		fmt.Scanln(&send_text)

		color.White("Сколько отправить?: ")
		fmt.Scanln(&count)

		for i := 0; i < count; i++ {
			printer_use(printer_ip, printer_port, send_text)
		}
	}
}
