package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

type Address struct {
	City string `json:"city"`
	State string `json:"state"`
	Country string `json:"country"`
	Pincode int `json:"pincode"`
}

type User struct {
	Name string
	Age json.Number
	Contact string
	Company string
	Address Address
}

const Version = "1.0.1"

type (Logger interface {
    Fatal(string, ...interface{})
    Error(string, ...interface{})
    Warn(string, ...interface{})
    Info(string, ...interface{})
    Debug(string, ...interface{})
    Trace(string, ...interface{})
}

 Driver struct {
    mutex   sync.Mutex
    mutexes map[string]*sync.Mutex
    dir     string
    log     Logger
	}
)

type Options struct {
	Logger 
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)
	opts := Options{}
	if options != nil {
	 	opts = *options
	}
	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}
	d := Driver{
	mutex: sync.Mutex{},
	mutexes: map[string]*sync.Mutex{},
	dir: dir,
	log: opts.Logger,
	}

	if _ , err :=  os.Stat(dir); err == nil {
		opts.Logger.Debug("directory already exists \n",dir)
		return &d, nil
	}

	opts.Logger.Debug("creating directory \n",dir)
	return &d, os.MkdirAll(dir, 0755)
	
}

func (d *Driver) Write(collection , resource string, data interface{}) error {
	if collection == "" {
		return fmt.Errorf("collection is required")
	}

	if resource == "" {
		return fmt.Errorf("resource is required")
	}
	
	mutex := d.getorCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	path := filepath.Join(d.dir,collection)
	fnlpath := filepath.Join(path,resource + ".json")

	tmpPath := fnlpath + ".tmp"

	if err := os.Mkdir(path, 0755); err != nil {
		return fmt.Errorf("err creating file %s",fnlpath)
	}	

	b,err := json.MarshalIndent(data,"","\t")
	if err != nil {
		return err
	}		

	b = append(b,byte('\n'))

	if err := os.WriteFile(tmpPath,b,0644); err != nil {
		return err
	}

	if err := os.Rename(tmpPath,fnlpath); err != nil {
		return err
	}

	return nil

}

func (d *Driver) Read(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("collection is required")
	}
	
	if resource == "" {
		return fmt.Errorf("resource is required")
	}

	path := filepath.Join(d.dir,collection,resource)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("resource not found")
	}

	 b , err := os.ReadFile(path + ".json")
	 if err != nil {
		return err
	 }

	 return json.Unmarshal(b,&v)
		
}

func(d *Driver) ReadAll(collection string) ([]string, error) {

	if collection == "" {
		return nil, fmt.Errorf("collection is required")
	}

	path := filepath.Join(d.dir,collection)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("resource not found")
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var data []string
	for _, f := range files {
		b , err := os.ReadFile(filepath.Join(path,f.Name()))
		if err != nil {
			return nil, err
		}
		data = append(data,string(b))
	}
	return data, nil
}

func(d *Driver) Delete(collection, resource string) error {
	if collection == "" {
		return fmt.Errorf("collection is required")
	}
	
	if resource == "" {
		return fmt.Errorf("resource is required")
	}

	path := filepath.Join(d.dir,collection,resource)

	mutex := d.getorCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir,path)

	switch fi,err := stat(dir); {
		case fi==nil, err != nil:
			return fmt.Errorf("resource not found")
		case fi.Mode().IsDir():
			return os.RemoveAll(dir)
		case fi.Mode().IsRegular():
			return os.RemoveAll(dir+".json")
	}
	return nil
}


func(d *Driver) getorCreateMutex(collection string) *sync.Mutex {

	d.mutex.Lock()
	defer d.mutex.Unlock()

	m , ok := d.mutexes[collection]
	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}
	return m
}

func stat (path string) (fi os.FileInfo,err error) {
	if fi,err = os.Stat(path); os.IsNotExist(err) {
		fi ,err = os.Stat(path+ ".json")
	}
	return 
}	


//main
func main() {
	dir := "./"
	db ,err := New(dir,nil)
	if err != nil {
		fmt.Println("err",err)
	}

	employees := []User{
		{"SpaceX",json.Number("42"),"SpaceX","SpaceX",Address{City:"Earth",State:"USA",Country:"USA",Pincode:1234}},
		{"Tesla",json.Number("42"),"Tesla","Tesla",Address{City:"Earth",State:"USA",Country:"USA",Pincode:1234}},
		{"Google",json.Number("42"),"Google","Google",Address{City:"Earth",State:"USA",Country:"USA",Pincode:1234}},
		{"Apple",json.Number("42"),"Apple","Apple",Address{City:"Earth",State:"USA",Country:"USA",Pincode:1234}},
	}

	for _,employee := range employees {
		err := db.Write("employee",employee.Name,User{
			Name: employee.Name,
			Age: employee.Age,
			Contact: employee.Contact,
			Company: employee.Company,
			Address: employee.Address,
		})
		if err != nil {
			fmt.Println("err writing",err)
		}
	}
	
	data , err := db.ReadAll("employee")
	if err != nil {
		fmt.Println("err",err)
	}	

	fmt.Println(data)

	allemployes := []User{}

	for _, f := range data {
		employee := User{}
		err := json.Unmarshal([]byte(f),&employee)
		if err != nil {
			fmt.Println("err",err)
		}
		allemployes = append(allemployes,employee)
	}
	fmt.Println(allemployes)
	
}