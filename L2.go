package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"unicode"
)

const DUOMENU_KIEKIS int = 28;
var threads int =0;

type Duom struct {
	Duomenys []Data `json:"users"`
}

type Data struct {
	Vardas string `json:"Vardas"`
	Amzius int  `json:"Amzius"`
	Atlyginimas float64 `json:"Atlyginimas"`
	hash string `json:"hash"`
}

func main() {

	jsonFile, err := os.Open("users.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully Opened .json file")
	}

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var Duomenys Duom

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &Duomenys)

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	DataChannel := make(chan Data)
	ResultChannel := make(chan Data)
	WorkingChannel := make(chan Data)
	MainChannel := make(chan Data)

	//Duomenys := Read("IFF72_PetraskaJ_L1_dat_2.json") // Nuskaitymas

	for i:=0; i<7; i++ { go take_exec_write(DataChannel,WorkingChannel) 	}

	go Data_Thread(MainChannel,DataChannel)

	go Rez_Thread(WorkingChannel,ResultChannel)

	for i:=0; i<len(Duomenys.Duomenys); i++{
		fmt.Sprintf("%v", Duomenys.Duomenys[i].Vardas)
		MainChannel <- Duomenys.Duomenys[i]
	}


	close(MainChannel)

	rezFile, err := os.Create("IFF-7-2_PetraskaJ-Lab2-Rez-2.txt")

	if err != nil {
		fmt.Println(err)
	}

	i := 0

	rezFile.WriteString("Rezultatai:\n")
	rezFile.WriteString("Numeris |  Vardas            |  Atlyginimas    |   Amzius      | Hash \n")

	for x:= range ResultChannel {
		i++
		rezFile.WriteString(fmt.Sprintf("%-6d  |%20s|%-17f|%-8.2d       |%15s\n", i, x.Vardas, x.Atlyginimas, x.Amzius, x.hash))
	}

}

func take_exec_write(in chan Data, out chan Data) {// in from Data_Thread ; out to Result_Thread

	fmt.Sprint("Started")
	for x := range in{ // pasiima(laukia) is Duomenu iki kol bus uzdarytas kanalas tada baigia cikla
		h := hash(x.Vardas + strconv.FormatFloat(x.Atlyginimas, 'E', -1, 64) + strconv.Itoa(x.Amzius) ) // skaiciuoja hash koda
		x.hash = h
		fmt.Println(h[0:5]) // patikrinimui pirmi penki simboliai
		r:= []rune(h)
		if unicode.IsDigit(r[0]) { // jeigu jis prasideda skaiciumi tuomet jis atitinka kriteriju
			out <- x // iveda i kanala rezultatu gijai
		}
	}
	close_this(out)
}

func close_this(ch chan Data){
	threads++
	if threads == 7 {
		close(ch)
	}
}

func Data_Thread(in chan Data,out chan Data) {// in from Main_Thread ; out to Work Thread

	var duomenys [DUOMENU_KIEKIS/2]Data
	count := 0
	for x := range in{
		if count < DUOMENU_KIEKIS/2 { // jei yra vietos pasiima is Main'o nuskaitytus duomenis
			duomenys[count] = x
			count++
		}
		if count > 0 {
			out <- duomenys[count]
			count--
		}
	}
	close(out)
}

func Rez_Thread(in chan Data, out chan Data) {// in from Result thread ; out to Main

	var rezultatai [DUOMENU_KIEKIS]Data
	count := 0
	for x := range in{
		rezultatai[count] = x
		count++
	}
	for i := 0; i <= count-1; i++{
		out <- rezultatai[i]
	}
	close(out)
}







func hash(s string) string {
	data := []byte("hello")
	return fmt.Sprintf("%x", md5.Sum(data))
}

func Read(path string) []Data {
	// we initialize our au array
	var Duomenys []Data
	// Open our jsonFile
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened " + path)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'Cars' which we defined above

	json.Unmarshal(byteValue, &Duomenys)
	fmt.Printf("stuff : %+v", Duomenys)
	return Duomenys
}