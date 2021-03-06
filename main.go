package main

import (
        "fmt"
        "strings"
        "strconv"
        "flag"
	"os"
        "./configfile"
)

func loadlist(m map[int]string) {
    c := make(chan string);
    tf := make(chan map[string]string);
    ts := 5
    p, err := configfile.ReadConfigFile("config.cfg"); 
    if err != nil {
	fmt.Printf("can't open file; err=%s\n",  err.String())
	os.Exit(1)
    }


    if tsc, err := p.GetString("GenSetting","OThread"); err == nil {
        fmt.Printf("Thread: %s\n",tsc)
        ts, err = strconv.Atoi(tsc)
    }

    jobsplit:=len(m)/ts;
    jobmod:=len(m)%ts;
    for i := 0; i < ts ; i++ {
        s := &URLParse{} //{ ed2kurldb: NewURLStore("ed2kurl.gmap") }
        s.id=i;
        if jobmod>0 {
            s.size=jobsplit+1
            jobmod--;
        }else{
            s.size=jobsplit
        }
        go s.urlparser(c,tf);
    }

    for _, url := range m {
        c <- url;
    }

    ed2kurldb:=NewURLStore("ed2kurl.gmap")
    lock:=0
    for i := 0; i < ts ; i++ {
        for k, v := range <-tf {
            bull:=""
            if err:=ed2kurldb.Get(&k,&bull);err!=nil {
                printamule(v,p)
                ed2kurldb.Put(&v, &k)
                lock=1
            }
        }
    }
    if lock==1 { ed2kurldb.save(); }
    //ed2kurldb.dirty <- true }
}

func printamule(el string,p *configfile.ConfigFile){
    ars,_:=p.GetString("amule","ARS")
    arp,_:=p.GetString("amule","ARP")
    arps,_:=p.GetString("amule","ARPS")

    fmt.Printf("amulecmd --host=%s -p %s -P %s -c \"add %s\"\n",ars,arp,arps,el)
}

func help(){
    fmt.Printf("%s\n","ed2kcrawler v0.0.2pre");
    fmt.Printf("%s\n","Usage: ed2kcrawler [Option]... [URL]...");
    fmt.Printf("%s\n","Commands:");
    flag.PrintDefaults();
}

func main() {
    var Loadpagelist = flag.Bool("l", false, "\tLoad the ed2klink page url list")
    flag.Parse();
    if *Loadpagelist {
        listfilename:=flag.Arg(0);
        m,err:=loadvv(listfilename);
        if err==nil {
            loadlist(m);
        } else {
            fmt.Printf("File %s not exist.\n",listfilename);
        }
    }else{
        if flag.NArg() >0 {
            m := make(map[int]string);
            for p,i := 0,0; i < flag.NArg(); i++ {
                if strings.HasPrefix(flag.Arg(i),"http://") {
                    m[p]=flag.Arg(i);
                    p++;
                }
            }
            loadlist(m);
        }else{
            help();
        }
    }
}

