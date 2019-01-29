package main

import (
    "testing"
    "net"
    "net/http"
    "strconv"
    "log"
    "os"
    //"syscall"
    
)
var goroutine chan int 

func setup() chan int {
	ch := make(chan int)
	go func() {
	srv := setServer("8080")
	gracefullshutdown()
	log.Fatal(srv.ListenAndServe())
	
	}()
	return ch
}
func init() {
	goroutine = setup()
}

func TestSplitterBasic(t *testing.T) {
		if paramSplit("ABC") != "A B C" {
			
			t.Error("ABC != A B C splitter is failing")
		}
		if paramSplit(`#{[|\^@]}æâ€êþÿÿÿûîœôŀïüð’‘ëßä«»© ↓¬¿×÷¡≤P`) != `#{[|\^@]}æâ€êþÿÿÿûîœôŀïüð’‘ëßä«»© ↓¬¿×÷¡≤ P` {
			
			t.Error("ABC != A B C splitter is failing" + paramSplit("#{[|\\^@]}æâ€êþÿÿÿûîœôŀïüð’‘ëßä«»© ↓¬¿×÷¡≤ P"))
		}
		if paramSplit(`AlfredENeuman`) != `Alfred E Neuman` {
			
			t.Error("ABC != A B C splitter is failing" + paramSplit("#{[|\\^@]}æâ€êþÿÿÿûîœôŀïüð’‘ëßä«»© ↓¬¿×÷¡≤ P"))
		}
	
}
func TestPort(t *testing.T){
	port := 70000
	p,err := setPort(&port)
	if err == nil {
		t.Error("port" + strconv.Itoa(*p) + " out of range doe'snt raise error")
	} 
	ln, err := net.Listen("tcp", ":" + strconv.Itoa(9000))
	port = 9000
	p,err = setPort(&port)
	if err == nil {
		t.Error("busy port does'nt raise error")
	}
	ln.Close()
	
}
func TestsetEnv(t *testing.T){
	env := []string{"TEST=/bin"}
	setEnv(env)
	if os.Getenv("TEST") != "/bin" {
		t.Errorf("setting environnement var TEST failed %s ", os.Getenv("TEST") ) 
	}
	env = []string{"TEST==/bin"}
	err := setEnv(env)
	if err == nil {
		t.Errorf("set Env Does'nt throw an error ")
	} 
	env = []string{"TEST=/bin=TEST2=/BIN"}
	err = setEnv(env)
	if err == nil {
		t.Errorf("set Env Does'nt throw an error ")
	} 
	env = []string{"TEST=/bin","TEST2=/BIN"}
	setEnv(env)
	if  os.Getenv("TEST") != "/bin" && os.Getenv("TEST2") != "/BIN" {
		t.Errorf("Does'nt handle more than one env var")
	}
	
}
func TestRouter(t *testing.T) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/helloworld", nil)
    if err != nil {
        t.Fatal("Creating  request failed!")
    }
    res, err := http.DefaultClient.Do(req)
    if err != nil {
            t.Error(err) //Something is wrong while sending request
        }
	if res.StatusCode != 200  {
            t.Errorf("Page not found : %d ", res.StatusCode) 
	}
	req, err = http.NewRequest("GET", "http://127.0.0.1:8080/helloworld?name=poeaizA", nil)
    if err != nil {
        t.Fatal("Creating 'GET /questions/1/SC' request failed!")
    }
    res, err = http.DefaultClient.Do(req)
    if err != nil {
            t.Error(err) //Something is wrong while sending request
        }
	if res.StatusCode != 200  {
            t.Errorf("Page not found : %d", res.StatusCode) 
	}
	req, err = http.NewRequest("GET", "http://127.0.0.1:8080/version", nil)
    if err != nil {
        t.Fatal("Creating 'GET /questions/1/SC' request failed!")
    }
    res, err = http.DefaultClient.Do(req)
    if err != nil {
            t.Error(err) //Something is wrong while sending request
        }
	if res.StatusCode != 200  {
            t.Errorf("Page not found : %d", res.StatusCode) 
	}
	req, err = http.NewRequest("GET", "http://127.0.0.1:8080/redirect", nil)
    if err != nil {
        t.Fatal("Creating 'GET /questions/1/SC' request failed!")
    }
    res, err = http.DefaultClient.Do(req)
    if err != nil {
            t.Error(err) //Something is wrong while sending request
        }
	if res.StatusCode != 404  {
            t.Errorf("Server does'nt trigger 404 but %d", res.StatusCode) 
	}
	
}
func TestShutdown(t *testing.T) {
	close(goroutine)
}
