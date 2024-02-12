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
const sleepBetweenConst = 2

func printRequest(conn net.Conn, reqType string, token string) {
	res, err := bufio.NewReader(conn).ReadString('}')
	catch(err)
	fmt.Println("Request " + token[len(token)-3:] + " (" + reqType + "): " + res)
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
		printRequest(conn, "entergame", token)

		go printBroadcast(conn2)
	}
	return tokens
}

func saveAnswers(conn net.Conn, conn2 net.Conn, tokens []string) {
	for qn := range []int{0, 1} {
		for i, t := range tokens {
			// Get question
			fmt.Fprintf(conn, "{\"method\": \"getquestion\", \"token\": \""+t+"\"}\n")
			printRequest(conn, "getquestion", t)

			// Save answer
			fmt.Fprintf(conn, "{\"method\": \"saveanswer\", \"token\": \""+t+"\", \"answer\": \"ans "+strconv.Itoa(qn)+"."+strconv.Itoa(i)+"!\"}\n")
			printRequest(conn, "saveanswer", t)

			go printBroadcast(conn2)
		}
	}
}

func sendVotes(conn net.Conn, conn2 net.Conn, tokens []string) {
	//for i := 0; i < len(tokens) - 2; i++
	for i := range tokens {
		fmt.Println("----- Voting:", i)
		for _, t := range tokens {
			// Get duel
			fmt.Fprintf(conn, "{\"method\": \"getduel\", \"token\": \""+t+"\"}\n")
			printRequest(conn, "getduel", t)

			// Save vote
			fmt.Fprintf(conn, "{\"method\": \"savevote\", \"vote\": 1, \"token\": \""+t+"\"}\n")
			printRequest(conn, "savevote", t)
		}

		go printBroadcast(conn2)
		//time.Sleep(sleepBetweenConst * time.Second)
	}

	fmt.Fprint(conn, "{\"method\": \"getroundresult\", \"token\": \""+tokens[0]+"\"}\n")
	printRequest(conn, "getroundresult", tokens[0])
}

func getGameResult(conn net.Conn, conn2 net.Conn, tokens []string) {
	fmt.Fprint(conn, "{\"method\": \"getgameresult\", \"token\": \""+tokens[0]+"\"}\n")
	printRequest(conn, "getgameresult", tokens[0])
}

func main() {
	// Подключаемся к сокету
	fmt.Println("Start client")
	connReq, err := net.Dial("tcp", "127.0.0.1:8081")
	catch(err)
	connBrcast, err := net.Dial("tcp", "127.0.0.1:8082")
	catch(err)

	tokens := registerAndEntergame(connReq, connBrcast, []string{"dovolniy", "Yuriy", "Andrew", "user4", "user5"})
	time.Sleep(3 * time.Second)

	// Игра 1
	saveAnswers(connReq, connBrcast, tokens)

	sendVotes(connReq, connBrcast, tokens)

	printBroadcast(connBrcast)

	getGameResult(connReq, connBrcast, tokens)

	time.Sleep(7 * time.Second)

	tokens = registerAndEntergame(connReq, connBrcast, []string{"MOLODOY", "stariy", "kek", "cheburek", "hohotunchik"})
	time.Sleep(3 * time.Second)

	// Игра 2
	saveAnswers(connReq, connBrcast, tokens)

	sendVotes(connReq, connBrcast, tokens)

	printBroadcast(connBrcast)

	getGameResult(connReq, connBrcast, tokens)

}
