package core

import (
	"fmt"
	"strconv"
	"time"
)

/* implement OK and NIL */
var RESP_OK []byte = []byte("+OK\r\n")
var RESP_NIL []byte = []byte("$-1\r\n")

// RedisCmd represents a parsed Redis command
type RedisCmd struct {
	Cmd  string
	Args []string
}

func Encode(value interface{}, isSimple bool) []byte {
	// fmt.Printf("Encoding value: %v, isSimple: %t\n", value, isSimple)
	switch v := value.(type) {
	case string:
		if isSimple {
			result := []byte(fmt.Sprintf("+%s\r\n", v))
			// fmt.Printf("Encoded as Simple String: %q\n", string(result))
			return result
		}
		result := []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
		// fmt.Printf("Encoded as Bulk String: %q\n", string(result))
		return result
	case int, int8, int16, int32, int64:
		result := []byte(fmt.Sprintf(":%d\r\n", v))
		return result
	}
	// fmt.Println("Unknown type, returning empty")
	return []byte{}
}
func evalEXPIRE(Args []string) []byte {
	//EXPIRE key time in sec
	if len(Args) != 2 {
		return []byte("-ERR wrong number of arguments for 'expire' command\r\n")
	}
	//Get the key,timeout duration
	var key string = Args[0]
	expireDurationSec, err := strconv.ParseInt(Args[1], 10, 64)
	if err != nil {
		return []byte("-ERR value is not an integer or out of range\r\n")
	}

	//get the key
	obj := Get(key)
	if obj == nil {
		//return 0 if key is invalid
		return Encode(0, false)
	}
	obj.ExpiresAt = time.Now().UnixMilli() + expireDurationSec*1000 //store in mili second

	//return 1 success
	return Encode(1, false)

}
func evalDEL(Args []string) []byte {
	//DEL k1,k2,..
	if len(Args) < 1 {
		return []byte("-ERR wrong number of arguments for 'del' command\r\n")

	}
	var del_cnt int = 0
	for i := 0; i < len(Args); i++ {
		key := Args[i]
		if Del(key) == true {
			del_cnt++
		}
	}
	return Encode(del_cnt, false)
}
func evalTTL(Args []string) []byte {
	if len(Args) != 1 {
		return []byte("-ERR wrong number of arguments for 'TTL' command\r\n")
	}
	var key string = Args[0]

	obj := Get(key)

	if obj == nil {
		return []byte(":-2\r\n")
	}
	if obj.ExpiresAt == -1 {
		return []byte(":-1\r\n")
	}

	durationMS := obj.ExpiresAt - time.Now().UnixMilli()

	if durationMS < 0 {
		return []byte(":-2\r\n")
	}
	return Encode(durationMS/1000, false)
}
func evalGET(Args []string) []byte {
	if len(Args) != 1 {
		return []byte("-ERR wrong number of arguments for 'GET' command\r\n")
	}
	var key string = Args[0]

	//Get the object
	obj := Get(key)
	// key not exist

	if obj == nil {
		return RESP_NIL
	}
	return Encode(obj.Value, false)
}
func evalSET(Args []string) []byte {
	//check the size
	if len(Args) <= 1 {
		return []byte("-ERR wrong number of arguments for 'SET' command\r\n")
	}

	key, value := Args[0], Args[1]
	var exDurationMs int64 = -1
	for i := 2; i < len(Args); i++ {
		switch Args[i] {
		case "EX", "ex":
			i++
			if i == len(Args) {
				// means no argument provided for EX time
				return []byte("-ERR syntax error\r\n")
			}
			//calculate the expiration time and change to ms
			exSeconds, err := strconv.ParseInt(Args[i], 10, 64)
			if err != nil {
				return []byte("-ERR value is not an integer or out of range\r\n")
			}
			exDurationMs = exSeconds * 1000 //convert to ms
		default:
			return []byte(fmt.Sprintf("-ERR unknown Argument '%s'\r\n", Args[i]))
			//store the k,v pair
		}
	}
	Put(key, NewObj(value, exDurationMs))
	return RESP_OK

}
func evalTIME(Args []string) []byte {
	if len(Args) > 0 {
		return []byte("-ERR wrong number of arguments for 'TIME' command\r\n")
	}

	now := time.Now()
	seconds := now.Unix()
	microseconds := now.Nanosecond() / 1000

	// TIME command returns an array with 2 elements: [seconds, microseconds]
	// Format: *2\r\n$len\r\nseconds\r\n$len\r\nmicroseconds\r\n
	secondsStr := fmt.Sprintf("%d", seconds)
	microsecondsStr := fmt.Sprintf("%d", microseconds)

	result := fmt.Sprintf("*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n",
		len(secondsStr), secondsStr,
		len(microsecondsStr), microsecondsStr)

	return []byte(result)
}
func evalECHO(Args []string) []byte {
	if len(Args) != 1 {
		return []byte("-ERR wrong number of arguments for 'echo' command\r\n")
	} else {
		return Encode(Args[0], false)
	}
}

func evalPING(Args []string) []byte {
	// fmt.Printf("Evaluating PING command with %d args: %v\n", len(Args), Args)

	if len(Args) >= 2 {
		// fmt.Println("Too many arguments for PING")
		return []byte("-ERR wrong number of arguments for 'ping' command\r\n")
	}

	if len(Args) == 0 {
		// fmt.Println("PING with no args, returning PONG")
		return Encode("PONG", true)
	} else {
		// fmt.Printf("PING with arg: %s\n", Args[0])
		return Encode(Args[0], false)
	}
}

func EvalAndResponse(Command *RedisCmd) []byte {
	// fmt.Printf("Evaluating command: %s with args: %v\n", Command.Cmd, Command.Args)

	switch Command.Cmd {
	case "PING":
		return evalPING(Command.Args)
	case "ECHO":
		return evalECHO(Command.Args)
	case "TIME":
		return evalTIME(Command.Args)
	case "SET":
		return evalSET(Command.Args)
	case "GET":
		return evalGET(Command.Args)
	case "TTL":
		return evalTTL(Command.Args)
	case "DEL":
		return evalDEL(Command.Args)
	case "EXPIRE":
		return evalEXPIRE(Command.Args)
	default:
		// fmt.Printf("Command %s not supported\n", Command.Cmd)
		return []byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", Command.Cmd))
	}
}
