package main

import (
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"MSET":    mset,
	"MGET":    mget,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

func ping(args []Value) Value {
	if len(args) != 0 {
		return Value{typ: "string", str: args[0].bulk}
	}
	return Value{typ: "string", str: "PONG"}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func mset(args []Value) Value {
	if len(args)%2 != 0 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	for i := 0; i < len(args); i += 2 {
		key := args[i].bulk
		val := args[i+1].bulk

		SETsMu.Lock()
		SETs[key] = val
		SETsMu.Unlock()
	}
	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	SETsMu.RLock()
	value, exist := SETs[key]
	SETsMu.RUnlock()

	if !exist {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func mget(args []Value) Value {
	val := Value{}
	val.typ = "array"
	val.array = make([]Value, 0)

	for i := 0; i < len(args); i++ {
		key := args[i].bulk
		SETsMu.RLock()
		value, exist := SETs[key]
		SETsMu.RUnlock()

		if exist {
			val.array = append(val.array, Value{typ: "bulk", bulk: value})
		} else {
			val.array = append(val.array, Value{typ: "null"})
		}
	}
	return val
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}
	hash, key, val := args[0].bulk, args[1].bulk, args[2].bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = make(map[string]string)
	}
	HSETs[hash][key] = val
	HSETsMu.Unlock()
	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	values := []Value{}

	for k, v := range value {
		values = append(values, Value{typ: "bulk", bulk: k})
		values = append(values, Value{typ: "bulk", bulk: v})
	}
	return Value{typ: "array", array: values}
}
