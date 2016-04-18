#Politifact-Request-Maker

This is a simple command line utility that can automatically
make requests from the politifact API and:

1. Write the responses to a `json` file

2. Forward the responses to a Firebase service 

##Usage 

Two major modes: 

1. Make a one time request to the politifact API. Write the response
   to a json file.
   
   Right now, this can only handle requesting "Statements" from the
   API. There are two ways of going about doing so. 
   
   A. Requesting statements that all fall under a particular
   subject. 
   
   Ex:  `> subject Guns 1000 ~/`
   
   The first part of the command lets the program know that you want
   to search by subject. The second part is the actual subject that
   you want to search for. (note - the program can already deal with
   santizing the names of subejects/people into a regular format). The
   third part is the maxmium number of requests that you want the API
   to return. The fourth part is the location of the file that you
   want to responses to be saved to. In this case, the program will
   create `guns-statements.json` in your home directory.
   
   B. Requesting statements that were made by a particular person. 
   
   Ex:  `> person "bernie Sanders" 1000 ~/`
   
   The command is analogous to the one above. 
   
2. Act as a server, grabbing the lastest statements and forwarding the
   responses to Firebase. 

	Ex: `serve ~/politifact_config.json`
	
	The first part of the command tells to program to act as a server,
    and the second part is the configuration file for the server. 
	

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

`request_rate_seconds`: How often you want the the program to query
politifact (in seconds)

`request_size`: How many responses you want the program to request
from the politifact API at a given time 

###`firebase` -> all the firebase related info
`max_concurrent_requests`: the maxmium number of concurrent
connections that you want the server to have with firebase at a given
time 

`root`: the url for your firebase database

`people_child_name`: child name for the place where the program should
put the responses (grouped by people, and then by subject)

`subject_child_name`: child name for the place where the program should
put the responses (only grouped by subject),


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
