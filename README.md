# mongo-updater

Oluşturulan excel dosyasındaki bilgileri mongodb içerisine ekleyen mini go projesi.


### Modlar ###
1) Collection içeriğini temizleyip tüm excel dosyasını import eder(clearCollection parametresi eklenirse).
2) Collection içeriğine dokunmadan excel dosyasında var olan verileri upsert eder.


### Çalıştırmak için ###
Proje gerekli değişkenleri args olarak almaktadır. 

**Terminalden** çalıştırmak için dosya dizininde aşağıdaki komutları çalıştırabilirsiniz,
```
go mod tidy
```
```
go run main.go updater.go -connectionString connStr -excelFilePath /excel/file/path -keyFields key1,key2 -databaseName dbname -collectionName collection -clearCollection
```

**Goland** gibi bir ide kullanıyorsanız projeyi açtıktan sonra build oluştururken program arguments kısmına gerekli değişkenleri yazabilirsiniz,
```
-connectionString connStr -excelFilePath /excel/file/path -keyFields key1,key2 -databaseName dbname -collectionName collection -clearCollection
```


### Parametreler ve örnek kullanım ###
**connectionString**: mongodb connectionstring. User eklenmiş halini kullanabilirsiniz.\
**excelFilePath**: excel konum.\
**keyFields**: **clearCollection** parametresi eklenmemişse update yaparken bu değişkendeki keylere bakıyor. Excelin ilk satırındaki field isimleriyle uyuşmalıdır. Virgül ile ayrılabilir.\
**databaseName**: mongodb içerisindeki db adı.\
**collectionName**: db içerisindeki collection adı.\
**clearCollection**: bu parametre eklenirse collection içeriği temizlenip tüm dosya import edilir. Parametre eklenmezse **keyFields** alanına göre update yapılır.


**Örneğin** localde projeye örnek olarak eklenen insert.xlsx dosyasını testdb/mongoUpdater collectionuna sıfırdan insert etmek için
```
go run main.go updater.go -connectionString mongodb://localhost:27017 -excelFilePath insert.xlsx -keyFields a -databaseName testdb -collectionName mongoUpdater -clearCollection
```


**Örneğin** localde projeye örnek olarak eklenen insert.xlsx dosyasındaki dökümanları testdb/mongoUpdater collectionunda upsert etmek için(a fieldini key olarak kullanıyoruz).
```
go run main.go updater.go -connectionString mongodb://localhost:27017 -excelFilePath insert.xlsx -keyFields a -databaseName testdb -collectionName mongoUpdater
```



### Excel dosya içeriği ###
Excel dosyasının ilk satırına collection field bilgileri yazılmalıdır. Altındaki satırlara insert edilecek bilgiler yazılmalıdır.

**Object insert etmek için,** field bilgilerini yazarken . ile ayırabilirsiniz. Örneğin, "b.x" ve "b.y" alanları b: {"x": "testx", "y":"testy"} olarak algılanacaktır.\
**Farklı tipleri insert etmek için,** field bilgilerinin sonuna #int, #bool şeklinde belirtebilirsiniz. # bilgisi eklenmezse alan string olarak algılanacaktır.


**Örnek excel;**\
![image](https://github.com/BerkaySari/mongo-updater/assets/20197226/fdf3a905-aa67-431d-8b67-938396293df4)


**mongodb çıktısı;**\
![image](https://github.com/BerkaySari/mongo-updater/assets/20197226/c4271aaa-2aa5-404b-910f-5f485cc5c0ca)
