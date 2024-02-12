package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"time"
)
import "fmt"

type ResToken struct {
	Token string `json:"token"`
}

func catch(err error) {
	if err != nil {
		log.Println(err)
	}
}

const sleepConst = 100
const sleepBetweenConst = 5

func printRequest(conn net.Conn, reqType string) {
	res, err := bufio.NewReader(conn).ReadString('}')
	catch(err)
	fmt.Println("Request (" + reqType + "): " + res)
	time.Sleep(sleepConst * time.Millisecond)
}

func printBroadcast(conn2 net.Conn) {
	res2, err := bufio.NewReader(conn2).ReadString('}')
	catch(err)
	fmt.Println("\nBroadcast: " + res2 + "\n")
	//time.Sleep(sleepConst * time.Millisecond)
}

func registerAndEntergame(conn net.Conn, conn2 net.Conn, usernames []string) []string {
	tokens := []string{}
	for _, u := range usernames {
		// Register
		fmt.Fprintf(conn, "{\"method\": \"register\", \"username\": \""+u+"\", \"password\": \"parol123\"}\n")
		res, err := bufio.NewReader(conn).ReadString('}')
		catch(err)
		fmt.Println("Request (register): ...")
		time.Sleep(sleepConst * time.Millisecond)
		resToken := ResToken{}
		err = json.Unmarshal([]byte(res), &resToken)
		catch(err)
		token := resToken.Token
		tokens = append(tokens, token)

		// Entergame
		fmt.Fprintf(conn, "{\"method\": \"entergame\", \"token\": \""+token+"\"}\n")
		printRequest(conn, "entergame")
		time.Sleep(2 * time.Second)

		go printBroadcast(conn2)
	}
	return tokens
}

func saveAnswers(conn net.Conn, conn2 net.Conn, tokens []string) {
	for qn := range []int{0, 1} {
		for i, t := range tokens {
			// Get question
			fmt.Fprintf(conn, "{\"method\": \"getquestion\", \"token\": \""+t+"\"}\n")
			printRequest(conn, "getquestion")

			// Save answer
			fmt.Fprintf(conn, "{\"method\": \"saveanswer\", \"token\": \""+t+"\", \"answer\": \"ans "+strconv.Itoa(qn)+"."+strconv.Itoa(i)+"!\"}\n")
			printRequest(conn, "saveanswer")

			go printBroadcast(conn2)
		}
	}
}

func sendVotes(conn net.Conn, conn2 net.Conn, tokens []string) {

	for i, t := range tokens {
		fmt.Println("----- Voting:", i)
		// Get duel
		fmt.Fprintf(conn, "{\"method\": \"getduel\", \"token\": \""+t+"\"}\n")
		printRequest(conn, "getduel")

		// Save vote
		fmt.Fprintf(conn, "{\"method\": \"savevote\", \"vote\": 1, \"token\": \""+t+"\"}\n")
		printRequest(conn, "savevote")

		go printBroadcast(conn2)
	}
}

func main() {
	// Подключаемся к сокету
	fmt.Println("Start client")
	connReq, err := net.Dial("tcp", "6.tcp.eu.ngrok.io:16581")
	catch(err)
	connBrcast, err := net.Dial("tcp", "4.tcp.eu.ngrok.io:10138")
	catch(err)
	fmt.Println("CONNS:", connReq, connBrcast)

	tokens := registerAndEntergame(connReq, connBrcast, []string{"dovolniy", "Yuriy", "Posevin", "user4"})
	time.Sleep(20 * time.Second)

	saveAnswers(connReq, connBrcast, tokens)
	time.Sleep(20 * time.Second)

	sendVotes(connReq, connBrcast, tokens)
	time.Sleep(20 * time.Second)

	fmt.Fprint(connReq, "{\"method\": \"getroundresult\", \"token\": \""+tokens[0]+"\"}\n")
	printRequest(connReq, "getroundresult")
	printBroadcast(connBrcast)

}
