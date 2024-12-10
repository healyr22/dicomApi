A simple REST API Server that allows uploading DICOM medical files, reading tag data, and converting to PNG, written in GoLang.

To run the server:
```
go build
./dicom
```

Apis:
`/upload`
`/extract-header`
`/convert-to-png`

Usage:
```
curl -X POST -F "file=@sample.dcm" http://localhost:8080/upload
```
```
curl -X GET "http://localhost:8080/extract-header?fileName=sample.dcm&tag=PatientName"
```
```
curl -X GET "http://localhost:8080/convert-to-png?fileName=sample.dcm" > output.png
```

System Design Diagram:
![alt text](https://github.com/healyr22/dicomApi/blob/main/system_design.png?raw=true)
