# Folder

## GET

Retrieve folder content:

* account
* breadcrumbs
* file
* entries

## POST

Create new content below folder (new entries):
* upload new file
* create empty file (params: name, type=F)
* create new folder
( copying: sollte über post geregelt werden. eventuell source param, der pfad zur Quellressource enthält)

### Params: 
* file: uploaded file (multipart)
* name: name of uploaded file (overwriting original name) or to be created file/folder
* type: 'F' for to be created file and 'D' for respective folder

## PUT


# File

## GET

Retrieve file content:

* account
* breadcrumbs
* file
* content