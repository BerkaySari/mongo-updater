package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func updateMongoDBFromExcel(connectionString string, excelFilePath string, clearCollection bool, keyFields []string, databaseName string, collectionName string) error {
	// MongoDB bağlantı ayarları
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	// Excel dosyasını aç
	xlFile, err := xlsx.OpenFile(excelFilePath)
	if err != nil {
		return err
	}

	// Excel verilerini işle
	sheet := xlFile.Sheets[0]

	// İlk satırı alan adları olarak kullan
	headerRow := sheet.Rows[0]
	fieldNames := make([]string, len(headerRow.Cells))
	for i, cell := range headerRow.Cells {
		fieldNames[i] = cell.String()
	}

	collection := client.Database(databaseName).Collection(collectionName)

	// bool değişkenine göre MongoDB koleksiyonunu güncelle
	if clearCollection {
		// Eski verileri sil
		_, err := collection.DeleteMany(context.Background(), bson.D{})
		if err != nil {
			return err
		}
	}

	// Excel verilerini işle
	for rowIndex, row := range sheet.Rows {
		if rowIndex == 0 {
			// İlk satır alan adları, geç
			continue
		}

		// Yeni bir belge oluştur
		document := make(map[string]interface{})
		for i, cell := range row.Cells {
			fieldName := fieldNames[i]
			setField(document, fieldName, cell)
		}

		if clearCollection {
			_, err := collection.InsertOne(context.Background(), document)
			if err != nil {
				return err
			}
		} else {

			// MongoDB koleksiyonunu güncelle
			filter := bson.M{}
			for _, keyField := range keyFields {
				filter[keyField] = document[keyField].(string)
			}

			update := bson.M{"$set": document}
			_, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
			if err != nil {
				return err
			}
		}
	}

	fmt.Println("Veriler başarıyla MongoDB koleksiyonuna yüklendi.")
	return nil
}

func setField(data map[string]interface{}, field string, value *xlsx.Cell) {
	fieldAndType := strings.Split(field, "#")

	keys := strings.Split(fieldAndType[0], ".")
	obj := data

	for i := 0; i < len(keys)-1; i++ {
		key := keys[i]
		if _, exists := obj[key]; !exists {
			obj[key] = make(map[string]interface{})
		}
		obj = obj[key].(map[string]interface{})
	}

	if len(fieldAndType) > 1 {
		if strings.Contains(fieldAndType[1], "int") {
			obj[keys[len(keys)-1]], _ = value.Int()
		}
		if strings.Contains(fieldAndType[1], "bool") {
			obj[keys[len(keys)-1]] = false
			if strings.ToLower(value.String()) == "true" {
				obj[keys[len(keys)-1]] = true
			}
		}
		if strings.Contains(fieldAndType[1], "float") {
			obj[keys[len(keys)-1]], _ = value.Float()
		}
		if strings.Contains(fieldAndType[1], "int64") {
			obj[keys[len(keys)-1]], _ = value.Int64()
		}
		if strings.Contains(fieldAndType[1], "time") {
			obj[keys[len(keys)-1]], _ = value.GetTime(false)
		}
		if strings.Contains(fieldAndType[1], "time1904") {
			obj[keys[len(keys)-1]], _ = value.GetTime(true)
		}
		if strings.Contains(fieldAndType[1], "date") {
			layout := "06.11.2023"
			obj[keys[len(keys)-1]], _ = time.Parse(layout, value.String())
		}
	} else {
		obj[keys[len(keys)-1]] = value.String()
	}
}
