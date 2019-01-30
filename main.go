package main

import (
	"context"
	"net"
    "net/http"
    "net/url"
    "strconv"
	"os"
	"flag"
    "fmt"
    "log"
    "unicode"
    "strings"
    "text/template"
   // "crypto/tls"
    "time"
    "errors"
    "os/signal"
       "syscall"
      
)
// Version string compiled with -ldflags "-X=main.VersionString=`git rev-parse HEAD`"
var VersionString = "unset"

// We defined our logwritter to ovveride default log format
type logWriter struct {
	Request    *http.Request
	URL        url.URL
	TimeStamp  time.Time
	StatusCode int
	Size int
}
func (writer logWriter) Write(bytes []byte) (int, error) {
    return fmt.Print(time.Now().Format("2006-01-02 15:04:05 ") + string(bytes))
}
// Split a world by it's uppercase and add a space in between the splitted part
// Returns the expected string "ABC" -> "A B C"
func paramSplit(param string) string {
     var words string
     var l int
     var s string
     l = 0 
     for  s = param; s != ""; s = s[l:] {
          l = strings.IndexFunc(s[1:], unicode.IsUpper) + 1
          if l <= 0 {
                l = len(s)
          }
          if l == len(s) {
	          words +=  s[:l] 
          } else {
	          words +=  s[:l] + " "
          }
        } 
     return words
}
// Check port flag and return port value if valid
func setPort(port *int) (*int,net.Listener,error) {
	
	if (*port != 80) && ( *port < 1024 || *port > 65535 ) {
		err := errors.New("port " + strconv.Itoa(*port) + " is out of range.")
		return port,nil,err
	} else {
		ln, err := net.Listen("tcp", ":" + strconv.Itoa(*port))
		if err != nil {
			err := errors.New("port " + strconv.Itoa(*port) + " already in use")
			if ln != nil { 
			ln.Close()
			}
			return port,nil,err
		}
		return  port,ln,err 
	}
}
// Check extra argument and set them  as env var if valid 
func setEnv(env []string ) error{
	
	for i:=0; i < len(env);i++ {
		parts := strings.Split(env[i], "=")
		if (len(parts) != 2){
			err := errors.New("Wrong env var declaration " + env[i])
			return err
		} else {
			err := os.Setenv(parts[0], parts[1])
			if err != nil {
				return  err
			} else {
				log.Print("Setting up  environnement var " + parts[0] +" to " + parts[1])
			}
		}
	}
	return nil 
}
//  Set the server port timeouts and route
func setServer(port string,) (*http.Server){
	mux := http.NewServeMux()
	mux.HandleFunc("/helloworld",paramHandler)
	mux.HandleFunc("/version",respondVersion)
	srv := &http.Server{
		Addr:           ":" + port,
	    ReadTimeout:  5 * time.Second,
	    WriteTimeout: 5 * time.Second,
	    IdleTimeout:  5 * time.Second,
	   // TLSConfig:    tlsConfig,
	    Handler:      middleWare(mux),
	}
	return srv
}
// middleWare ... Here we log access and errors, We can also set user info, auth, coockies...etc  
func middleWare(handler http.Handler) http.Handler {
	
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Right ! not all this does is log like  "github.com/gorilla/handlers"
		if (r.URL.Path != "/helloworld" && r.URL.Path != "/version") {
	        log.Printf("404 Not found %s %s",r.Method, r.URL)
	        // Reference here a handler for a custom 404 page
	        
	    }else {
			log.Printf("200 Ok %s %s",r.Method, r.URL)
		}

        handler.ServeHTTP(w, r)
    })
}
// Handle our helloworld page and it's potential args (name=)
// The function sanitize the params before rewritting them  in  the body
func paramHandler(w http.ResponseWriter, r *http.Request) {
	
	
	names, err := r.URL.Query()["name"]
    
    if !err || len(names[0]) < 1 {
        fmt.Fprintf(w, "<h1>Hello World stranger</h1>")  
        return     
    } else {
    name := names[0]
    name = paramSplit(name)
    //sanitize parameters from xss rewrite
    fmt.Fprintf(w, "<h1>Hello World %s</h1>", template.HTMLEscapeString(name))
	    return
    }
}
// Handle Version page
func respondVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Version %s</h1>",VersionString)
}
// Gracefull shutdown code
func gracefullshutdown(srv *http.Server){
	
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
    signal.Notify(gracefulStop, syscall.SIGINT)
    go func() {
		      
              sig := <-gracefulStop
              log.Printf("caught sig: %+v", sig)
              ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
              defer cancel()
              err := srv.Shutdown(ctx)         
              if err != nil {
              	log.Print(" " + err.Error())
              }
       }()
    
}
//Our Set the log output,parse arg line and  check  their validity ,set the server param 
// And listen and serve on the defined port   
func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
	port := flag.Int("p", 8080, "port")
	flag.Parse()
	env := flag.Args()
	if env != nil {
		err := setEnv(env) 
		if err != nil {
			log.Print("wrong environnement variable! Cannot Parse:" + err.Error())
		}
	}
	port,ln,err := setPort(port)
	if (err != nil) {
		log.Fatal(err)
	} else {
		
		srv := setServer(strconv.Itoa(*port))
		gracefullshutdown(srv)
		
	    log.Fatal(srv.Serve(ln))
	}
}



