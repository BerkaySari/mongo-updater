//go run main.go -excelFilePath=/path/to/excel.xlsx -keyFields=field1,field2 -databaseName=mydb -collectionName=mycollection -clearCollection

package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	// Komut satırı argümanlarını tanımla
	connectionString := flag.String("connectionString", "", "Bağlantı adresini belirtin")
	excelFilePath := flag.String("excelFilePath", "", "Excel dosyasının yolunu belirtin")
	clearCollection := flag.Bool("clearCollection", false, "Koleksiyonu temizle (true/false)")
	keyFields := flag.String("keyFields", "", "Anahtar alanları (virgülle ayrılmış)")
	databaseName := flag.String("databaseName", "", "Veritabanı adı")
	collectionName := flag.String("collectionName", "", "Koleksiyon adı")

	// Argümanları parse et
	flag.Parse()

	// Her bir argümanı kullanarak MongoDB güncellemesini çağır
	if *connectionString == "" || *excelFilePath == "" || *databaseName == "" || *collectionName == "" {
		fmt.Println("Lütfen zorunlu argümanları (excelFilePath, databaseName, collectionName) belirtin.")
		return
	}

	keyFieldsSlice := strings.Split(*keyFields, ",")

	err := updateMongoDBFromExcel(*connectionString, *excelFilePath, *clearCollection, keyFieldsSlice, *databaseName, *collectionName)
	if err != nil {
		fmt.Printf("Hata: %v\n", err)
		return
	}

	fmt.Println("Veriler başarıyla MongoDB koleksiyonuna yüklendi.")
}
