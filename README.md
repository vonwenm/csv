csv
=====

golang structs to csv file

## Install
go get github/sundy-li/csv

##  Tag to title mapping

  You could map the struct to csv indicating the title by using this
  
  
    type People struct {
        Name   string `csv:"姓名"`
        Age    int
        Sex    string
        Weight float32 `csv:"体重"`
        Marry  bool    `csv:"婚姻状况"`
    }

* ** If you have not indicate the csv tag , it will just use the field name **

## Write to file
    
    func main() {
    
        var peoples = make([]*People, 0, 10)
        for i := 0; i < 10; i++ {
            var people = &People{
                Name:   "someone",
                Age:    22,
                Sex:    "美女",
                Weight: 90.6,
                Marry:  false,
            }
            peoples = append(peoples, people)
        }
    
        file, _ := os.OpenFile("aa.csv", os.O_CREATE|os.O_RDWR, 0644)
        
        //if only file implements the io.Writer interface 
        var parser = csv.NewCsv(file)
        err := parser.Parse(peoples)
        if err != nil {
            println(err.Error())
        }
    }

run this code and it will build an csv file

## Write to the ResponseWriter
    
    func main() {
        http.HandleFunc("/", saveCsv)
        http.ListenAndServe(":8080", nil)
    }

    func saveCsv(w http.ResponseWriter, req *http.Request) {
        var randFile = time.Now().String() + ".csv"
        w.Header().Set("Content-Type", "application/octet-stream;charset=utf-8")
        w.Header().Set("Content-Disposition", "attachment;filename="+randFile)
        var peoples = make([]*People, 0, 10)
        for i := 0; i < 10; i++ {
            var people = &People{
                Name:   "someone",
                Age:    22,
                Sex:    "美女",
                Weight: 90.6,
                Marry:  false,
            }
            peoples = append(peoples, people)
        }
         //if only ResponseWriter implements the io.Writer interface 
        var parser = csv.NewCsv(w)
        err := parser.Parse(peoples)
        if err != nil {
            println(err.Error())
        }
        return
    }


* ** Run this code,while you visit `127.0.0.1:8080/` you will download the csv file **


