package core

import (
	"fmt"
)

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
	}
	// fmt.Println("Unknown type, returning empty")
	return []byte{}
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
	default:
		// fmt.Printf("Command %s not supported\n", Command.Cmd)
		return []byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", Command.Cmd))
	}
}
