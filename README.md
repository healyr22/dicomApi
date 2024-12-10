A simple REST API Server that allows uploading DICOM medical files, reading tag data, and converting to PNG, written in GoLang.

Usage:
```
go build
./dicom
```

Apis:
`/upload`
`/extract-header`
`/convert-to-png`

System Design Diagram:
![alt text](https://github.com/healyr22/dicomApi/blob/main/system_design.png?raw=true)
