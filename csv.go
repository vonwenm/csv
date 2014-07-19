package csv

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "reflect"
    "regexp"
    "strconv"
    "strings"
)

type Csv struct {
    mp       map[string]string //fileName --> title
    hasStart bool
    writer   *bufio.Writer
    needBom  bool

    separator string
    newline   string
}

func NewCsv(w io.Writer) *Csv {
    var csv = &Csv{}
    csv.writer = bufio.NewWriter(w)
    csv.needBom = true
    csv.separator = ","
    csv.newline = "\n"
    csv.hasStart = false

    return csv
}

func (this *Csv) GetWriter() *bufio.Writer {
    return this.writer
}

func (this *Csv) writeLine(str ...string) {
    var ss = strings.Join(str, this.separator)
    for _, s := range ss {
        this.writer.WriteRune(s)
    }

    this.writer.Write([]byte(this.newline))
}

func (this *Csv) SetBoom(flag bool) {
    this.needBom = flag
}

func (this *Csv) init() {
    this.mp = make(map[string]string)
    //this.writer.C
}

func (this *Csv) buildMap(entity interface{}) error {
    var v = reflect.ValueOf(entity)
    var t = reflect.TypeOf(entity)

    switch t.Kind() {
    case reflect.Slice:
        if v.Len() < 1 {
            return errors.New("[CSV:error] Slice cannot be empty")
        }
        return this.buildMap(v.Index(0))

    case reflect.Ptr:
        return this.buildMap(v.Elem().Interface())

    case reflect.Struct:

        for i := 0; i < t.NumField(); i++ {
            if !v.Field(i).CanInterface() {
                continue
            }
            var name = t.Field(i).Name
            var titleName = name
            var tag = string(t.Field(i).Tag)
            var regex = regexp.MustCompile(`csv:"([^\s]+)"`)
            var res = regex.FindSubmatch([]byte(tag))
            if len(res) > 1 {
                titleName = string(res[1])
            }

            if titleName != "-" && couldCsv(v.Field(i).Interface()) {
                this.mp[name] = titleName
            }

        }
        return nil

    default:
        return errors.New(fmt.Sprintf("[CSV:error] cannot suppor this struct -> %#v", t.Kind()))
    }
}

func (this *Csv) Parse(entity interface{}) (err error) {
    this.init()
    err = this.buildMap(entity)
    if err != nil {
        return err
    }

    err = this.parseCsv(entity)
    this.writer.Flush()
    return err
}

func (this *Csv) parseCsv(entity interface{}) (err error) {
    var v = reflect.ValueOf(entity)
    var t = reflect.TypeOf(entity)

    switch t.Kind() {
    case reflect.Slice:
        for i := 0; i < v.Len(); i++ {
            if err := this.Parse(v.Index(i).Interface()); err != nil {
                return err
            }
        }
        return nil
    case reflect.Ptr:
        return this.parseCsv(v.Elem().Interface())

    case reflect.Struct:
        if !this.hasStart {
            this.writeTitle(v.Interface())
        }
        var mp = this.mp
        var strs = make([]string, 0, t.NumField())
        for i := 0; i < t.NumField(); i++ {
            var name = t.Field(i).Name
            _, ok := mp[name]
            if ok {
                strs = append(strs, bean2Str(v.FieldByName(name).Interface()))
            }
        }
        this.writeLine(strs...)
        return nil

    }
    return errors.New(fmt.Sprintf("[CSV:error] cannot suppor this struct -> %#v", t.Kind()))
}

func (this *Csv) writeTitle(entity interface{}) {
    if this.hasStart {
        return
    }
    if this.needBom {
        this.writer.WriteString("\xEF\xBB\xBF")
    }
    var t = reflect.TypeOf(entity)

    var strs = make([]string, 0, t.NumField())
    for i := 0; i < t.NumField(); i++ {
        var name = t.Field(i).Name
        title, ok := this.mp[name]
        if ok {
            strs = append(strs, title)
        }
    }

    this.hasStart = true
    this.writeLine(strs...)
}

func couldCsv(bean interface{}) bool {
    switch reflect.TypeOf(bean).Kind() {
    case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Uint, reflect.Int64, reflect.String, reflect.Bool:
        return true
    }
    return false
}

func bean2Str(bean interface{}) string {
    var v = reflect.ValueOf(bean)

    switch reflect.TypeOf(bean).Kind() {
    case reflect.Float32, reflect.Float64:
        return strconv.FormatFloat(v.Float(), 'f', -1, 64)
    case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Uint, reflect.Int64:
        return strconv.FormatInt(v.Int(), 10)
    case reflect.String:
        return v.String()
    case reflect.Bool:
        if v.Bool() {
            return "true"
        }
        return "false"
    default:
        return "UNKNOW"
    }
}
