gridfs-stream-54817848
======================

This folder contains the code to [my answer][myanswer] to the "question"
[Store Uploaded File in MongoDB GridFS Using mgodb-go-driver without Saving to
Memory][question]

Usage
-----

```none
$ git clone https://github.com/mwmahlberg/stackoverflow-answers.git mwmahlberg-so-answers
Cloning into 'mwmahlberg-so-answers'...
$ cd mwmahlberg-so-answers/gridfs-stream-54817848
$ go build -o streamupload
$ set +o history 
$ export MONGO_USERNAME=$YOUR_MONGO_USERNAME
$ export MONGO_PASSWORD=$YOUR_MONGO_PASSWORD
$ export MONGO_HOST=$YOUR_MONGO_HOST
$ export MONGO_APPNAME=$YOUR_MONGO_APPNAME
$ set -o history
$ ./streamupload
2025/04/13 17:23:03 INFO New file uploaded name=main.go objectID=67fbd6d7859c44691365977d
2025/04/13 17:23:03 INFO File downloaded name=main_67fbd6d7859c44691365977d.go.dl objectID=67fbd6d7859c44691365977d
2025/04/13 17:23:03 WARN File deleted from gridfs name=main.go objectID=67fbd6d7859c44691365977d
```

> Note that your output may vary.

[myanswer]: https://stackoverflow.com/a/79571700/1296707
[question]: https://stackoverflow.com/questions/54817848/store-uploaded-file-in-mongodb-gridfs-using-mgodb-go-driver-without-saving-to-me/79571700#79571700
