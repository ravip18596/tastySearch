# tastySearch
Golang based REST API to search gourmet food reviews data

## How to run it.
1. Download and extarct foods.txt from `https://drive.google.com/file/d/0B8_VSW2-5XmpSTNlZXV4cVdLRUE/view`
2. run `export FILE="<Absolute path of the file>"` - path of the downloaded foods.txt
3. run `nohup ./tastySearchLinux > log.out 2>&1 &` for Linux users
4. run `tail -f log.out` to see logs
5. open any browser and hit - `localhost:8081/search/words?queries=cat,food,They,results`

query words are comma seperated
