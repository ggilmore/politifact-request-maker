#Politifact-Request-Maker

This is a simple command line utility that can automatically
make batch requests (only
[statements](http://static.politifact.com/api/doc.html#statements) from the [Politifact API](static.politifact.com/api/doc.html) and:

1. Write the responses to a `json` file

2. Forward the responses to a Firebase service 

##Usage 



###One-Time Request -> write response to `.json` file
   
   Right now, this can only handle requesting "Statements" from the
   API. There are two ways of going about doing so. 
   
__General Schema__:

```
>[METHOD (\"person\" or \"subject\")] [NAME_RESOURCE] [MAX_ITEMS] [OUTPUT_DIR]
```

`METHOD` - You can either request statmenets about a particular
`person`, or all statments that fall under a given `subject`

`NAME_RESOURCE` - The name of either the `subject` or `person` that
you're requesting statements about. The program will handle
"normalizing" the name or subject that you're asking about). For
example, asking about "berniesanders", "bernie Sanders, or "BerNiE
saNders" will all map to the same person. 

`MAX_ITEMS` - (integer) - The maximum number of statements that you
want returned from the Politifact API, if available

`OUTPUT_DIR` - The location of the file that you want the statements
saved to. The program will create a statement that looks like
`[NORMALIZED_NAME_RESOURCE]-statements.json` in the `OUTPUT_DIR` directory.

   
   Ex:
   ```
   > subject Guns 1000 ~/
   ```
   This will create a file called `~/guns-statements.json`.
   
   ```
   > person "bernie Sanders" 1000 ~/
   ```
   
   This will create a file called `~/bernie-s-statements.json`

   
###Act as a server, grabbing the lastest statements and forwarding the
responses to Firebase. 

__Command Schema__: 

```scheme
>serve [CONFIG_FILE_LOCATION]
```

`CONFIG_FILE_LOCATION` - The location of the configuration file for
the server mode of the program. The schema for the configuration file
is described below.

Ex: 

```
serve ~/politifact_config.json
```

This will start a server that fowards the latest politifact statements
to a Firebase instance, as described in `politifact_config.json`.

##Config File


The configuration file is a `json` file that looks like this: 

```json
{
	"politifact": {
		"request_rate_seconds": [N],
		"request_size":[N]
	},
	"firebase": {
		"max_concurrent_requests": [N],
		"root": [FIREBASE_URL],
		"people_child_name": [PEOPLE_ROUTE],
		"subjects_child_name" : [SUBJECT_ROUTE]"
	}
}
```


###`politifact` -> all the politifact related info

`request_rate_seconds`: (integer) ->  How often you want the the
program to query politifact (in seconds)

`request_size`: (integer) How many responses you want the program to request
from the politifact API at a given time 

###`firebase` -> all the firebase related info
`max_concurrent_requests`: (integer) The maxmium number of concurrent
connections that you want the server to have with firebase at a given
time 

`root`: The root URL for your Firebase instance

`people_child_name`: Firebase child name for the place where the program should
put the responses (grouped by people, and then by subject)

`subject_child_name`: Firebase child name for the place where the program should
put the responses (only grouped by subject)


###Full Example

```json
{
    "politifact": {
	"request_rate_seconds": 20,
        "request_size": 50

    },
    "firebase": {
	"max_concurrent_requests": 200,
        "root": "https://garbage.firebaseio.com/",
	"people_child_name": "people",
	"subjects_child_name" : "subjects"
    }
}
```

This tells the server that when it communicates with the Politifact
API, it should ask for the latest 50 statements every 20 seconds. 

It also tells the server that when it communicates with Firebase:

1. It can maintain at most 200 connections at a time.

2.The Firebase root URL is `"https://garbage.firebaseio.com/"`

3. The place where it should put the statements as sorted by people is 
`"https://garbage.firebaseio.com/people"`

4. The place where it should put the statements as sorted by subject is 
`"https://garbage.firebaseio.com/subjects"`

