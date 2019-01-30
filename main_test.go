package main

import (
    "testing"
    "net"
    "net/http"
    "strconv"
    "log"
    "os"
   
   // "os/exec"
    
    
)
var goroutine chan int 

// Run the server on  a goroutine to pass router Test
func setup() chan int {
	ch := make(chan int)
	go func() {
		port1 := 8080
		p,ln,err := setPort(&port1)
		if err != nil {
			os.Exit(1)
		}
	srv := setServer(strconv.Itoa(*p))
	gracefullshutdown(srv)
	log.Fatal(srv.Serve(ln))
	
	}()
	return ch
}
func init() {
	goroutine = setup()
}
// Check that our Splitter function  give us the expected result
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
// Test Busy port and out of range port declaration
func TestPort(t *testing.T){
	port := 70000
	p,l,err := setPort(&port)
	if err == nil {
		t.Error("port" + strconv.Itoa(*p) + " out of range doe'snt raise error")
	} 
	if l != nil {
	l.Close()
	}
	ln, err := net.Listen("tcp", ":" + strconv.Itoa(9000))
	port = 9000
	p,l,err = setPort(&port)
	if err == nil {
		t.Error("busy port does'nt raise error")
	}
	if l != nil {
	l.Close()
	}
	ln.Close()
	
}
// Check that setENV func parse correctly the flags
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
// Test of versionString does'nt work
// go test  -v -ldflags=-X=main.VersionString=`git rev-parse HEAD` does'nt compile the versionstring
/*func TestVersion(t *testing.T){

	var	(
		cmdOut []byte 
		err error
	)
	cmdArgs := []string{"rev-parse", "HEAD"}
	if cmdOut,err = exec.Command("git",cmdArgs...).Output(); err != nil {
		 t.Error(err) 
		
	}
	cmd := string(cmdOut)
	if VersionString != cmd {
		 t.Fatal(cmd + "Versionstring mismatch" + VersionString )
	}
	
}
*/
// Testing of the general routing on the  
func TestRouter(t *testing.T) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/helloworld", nil)
    if err != nil {
        t.Fatal("Creating  request failed!")
    } else {
	    res, err := http.DefaultClient.Do(req)
	    if err != nil {
	            t.Error(err) //Something is wrong while sending request
	        }
		if res.StatusCode != 200  {
	            t.Errorf("Page not found : %d ", res.StatusCode) 
		}
    }
	req, err = http.NewRequest("GET", "http://127.0.0.1:8080/helloworld?name=poeaizA", nil)
    if err != nil {
        t.Fatal("Creating 'GET /questions/1/SC' request failed!")
    } else { 
	    res, err := http.DefaultClient.Do(req)
	    if err != nil {
	            t.Error(err) //Something is wrong while sending request
	        }
		if res.StatusCode != 200  {
	            t.Errorf("Page not found : %d", res.StatusCode) 
		}
    }
	req, err = http.NewRequest("GET", "http://127.0.0.1:8080/version", nil)
    if err != nil {
        t.Fatal("Creating 'GET /questions/1/SC' request failed!")
    } else {
	    res, err := http.DefaultClient.Do(req)
	    if err != nil {
	            t.Error(err) //Something is wrong while sending request
	        }
		if res.StatusCode != 200  {
	            t.Errorf("Page not found : %d", res.StatusCode) 
		}
    }
	
	req, err = http.NewRequest("GET", "http://127.0.0.1:8080/redirect", nil)
    if err != nil {
        t.Fatal("Creating 'GET /questions/1/SC' request failed!")
    } else {
	    res, err := http.DefaultClient.Do(req)
	    if err != nil {
	            t.Error(err) //Something is wrong while sending request
	        }
		if res.StatusCode != 404  {
	            t.Errorf("Server does'nt trigger 404 but %d", res.StatusCode) 
		}
    }
	
}
// Did'nt manage to send sigkill to the goroutine
// Maybe try to  get the process via ps and send the signal 
// Still cleaning up  the server 
func TestShutdown(t *testing.T) {
	close(goroutine)
}
