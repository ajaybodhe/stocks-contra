# stocks-contra
##permissions
sudo chmod 0777 -R /tmp/

##lynx
sudo apt-get install lynx

##mysql
import mysql dump, username=root password=password

##mysql login
mysql --local-infile -uroot -ppassword NSE

##go installation steps
1.Download latest tar from https://golang.org/dl/
  tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz

2.vim $HOME/.bashrc
  Add the following line
  export PATH=$PATH:/usr/local/go/bin
  export GOROOT=/usr/local/go
  export GOBIN=/usr/local/go/bin
  export GOPATH=/root/Go/
  
3.Check that Go is installed correctly by building a simple program, as follows.

Create a file named hello.go and put the following program in it:

package main
import "fmt"
func main() {
    fmt.Printf("hello, world\n")
}

Then run it with the go tool:
$ go run hello.go
hello, world

If you see the "hello, world" message then your Go installation is working.

##go dependencies
go get github.com/go-sql-driver/mysql
go get golang.org/x/crypto
go get github.com/golang/glog
go get github.com/scalingdata/gcfg
go get gopkg.in/gomail.v2
go get github.com/PuerkitoBio/goquery
go get github.com/willf/bitset
go get github.com/cdipaolo/sentiment
