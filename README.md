### An HTTP server to find, insert and delete values in BST  
Need a json file with initial values

#### Run
Put a `.json` file containing int array  
Build & run `go build -o ./vsru . && ./vsru [filename]` or default filename is `array.json`

#### Endpoints
GET `/api/v1/search?val=[int]` find a value in BST  
Response is 200 and message `OK`  

POST `/api/v1/insert` with body `{"val": [int]}` insert new value int BST  
Response is 201 and message `OK`  

DELETE `/api/v1/delete?val=[int]` delete a value from BST  
Response is 204 without any messages

#### Errors
When file is wrong (can't marshal json from the file) an error with message will be shown, program will not start
When trying to insert existing value, the response is an error with message  
When trying to remove non-existing value, the response is an error with message
