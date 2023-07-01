# Description
This project solves the task defined in Online preselection cloud challenge of Nuwe Zurich cloud hackathon.


`terraform` folder contains terraform infrastructure required to solve the task. It contains lambda, s3 bucket, dynamodb table, s3 notifications etc.


`lambda` folder contains lambda written in Golang. It expects to be triggered by s3 event, parses created/updated file into `Client` struct, than marshals it to dynamodb item and puts the item into `data_table`.