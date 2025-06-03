**Stepup steps**
   1. Clone the repo <br>
      `git clone https://github.com/Ssenhub/transaction-crud-svc-go-postgres` <br> 

   2. Test composed docker containter <br>
        a. `cd transaction-crud-svc-go-postgres` <br>
        b. `docker compose up` <br>
        c. Run automated tests <br>
         &nbsp;&nbsp;&nbsp;  i) Open a new terminal <br>
         &nbsp;&nbsp;&nbsp;  ii) `cd transaction-crud-svc-go-postgres/tests` <br>
         &nbsp;&nbsp;&nbsp; iii) `go test -v` <br>
        d. Run manual tests <br>
         &nbsp;&nbsp;&nbsp;  i) Use a command line tool (`curl` on bash or powershell) or GUI apps such as PostMan or Insomnia. <br>

         Powershell exacmple:
      
          > $loginbody = '{"username": "jondoe", "password": "passwd"}' #Make sure user name and password mathces with 'USER_NAME' and 'PASSWORD' variable in .env
          > $r = curl -Uri "http://localhost:3000/login" -Method POST -Body $loginbody  #Make sure the port matches with .env 'PORT' variable
          > $token = (convertfrom-json $r.Content).token
          > $headers = @{"Authorization" = "Bearer "+$token}
          > $metadata = @{"channel" = "mobile_app1"; "location" = "Seattle, WA"}
          > $tx = @{                                      
                  "AccountId"    = "act1";
                  "Type"         = 1;
                  "Amout"        = 123.45;
                  "Currency"     = "USD";
                  "Description"  = "test desc 1";
                  "Status"       = 2;
                  "MerchantId"   = "Merch_1";
                  "MerchantName" = "MerchName_1";
                  "Metadata"     = convertto-json $md
               }
          > $body = ConvertTo-Json $tx
          > $r = Invoke-WebRequest -Uri "http://localhost:3000/transactions" -Method POST -Header $headers -Body $body #Create
          > $r = Invoke-WebRequest -Uri "http://localhost:3000/transactions" -Method GET -Header $headers #Get list
          > $result = convertfrom-json  $r.Content
          > $result
          > $r = Invoke-WebRequest -Uri "http://localhost:3000/transactions?limit=1&page=0" -Method GET -Header $headers #Get pages
          > $result = convertfrom-json  $r.Content
          > $result
          > $r = Invoke-WebRequest -Uri "http://localhost:3000/transactions/505" -Method GET -Header $headers #Get tx by id  
          > $result = convertfrom-json  $r.Content
          > $result
          > $tx = @{                                      
                  "AccountId"    = "act1";
                  "Type"         = 1;
                  "Amout"        = 123.45;
                  "Currency"     = "USD";
                  "Description"  = "test desc updated";
                  "Status"       = 2;
                  "MerchantId"   = "Merch_1";
                  "MerchantName" = "MerchName_1";
                  "Metadata"     = convertto-json $md
               }
          > $body = ConvertTo-Json $tx
          > $r = Invoke-WebRequest -Uri "http://localhost:3000/transactions/505" -Method PUT -Header $headers -Body $body #Update tx
          > $result = convertfrom-json $r.Content
          > $result
          > $r = Invoke-WebRequest -Uri "http://localhost:3000/transactions/505" -Method DELETE -Header $headers #Delete tx by id
   
   3. Test docker containter with external DB <br>
      a. Update following vartiables in `.env` file to appropriate values of local PostGreSQL installation <br>
         &nbsp;&nbsp;&nbsp;`POSTGRES_PORT=5432` <br>
         &nbsp;&nbsp;&nbsp;`POSTGRES_USER=postgres` <br>
         &nbsp;&nbsp;&nbsp;`POSTGRES_PASSWORD=passwd` <br>
         &nbsp;&nbsp;&nbsp;`POSTGRES_DB=postgres` <br>
      b. `cd transaction-crud-svc-go-postgres` <br>
      c. `docker build  -t server-image .` <br>
      d. `docker run -d -p 3000:3000 server-image` #Make sure the port matches with .env 'PORT' variable. <br>
      e. Run automated tests <br>
         &nbsp;&nbsp;&nbsp;    i) Open a new terminal <br>
         &nbsp;&nbsp;&nbsp;   ii) `cd transaction-crud-svc-go-postgres/tests` <br>
         &nbsp;&nbsp;&nbsp;  iii) `go test -v` <br>
      f. Run manual tests <br>
         &nbsp;&nbsp;&nbsp;   i) Follow same steps as step 2d. <br>
                       
   5. Test locally <br>
      a. Update following vartiables in `.env` file to appropriate values of local PostGreSQL installation <br>
         &nbsp;&nbsp;&nbsp;`POSTGRES_PORT=5432` <br>
         &nbsp;&nbsp;&nbsp;`POSTGRES_USER=postgres` <br>
         &nbsp;&nbsp;&nbsp;`POSTGRES_PASSWORD=passwd` <br>
         &nbsp;&nbsp;&nbsp;`POSTGRES_DB=postgres` <br>
      b. `cd transaction-crud-svc-go-postgres` <br>
      c. `go run .` <br>
      d. Run automated tests <br>
         &nbsp;&nbsp;&nbsp;   i) Open a new terminal <br>
         &nbsp;&nbsp;&nbsp;  ii) `cd transaction-crud-svc-go-postgres/tests` <br>
         &nbsp;&nbsp;&nbsp; iii) `go test -v` <br>
      e. Run manual tests <br>
         &nbsp;&nbsp;&nbsp; i) Follow same steps as step 2d. <br>
         
          
**LLM experience**

  I have used Chatgpt as primary resource and Gemini as a secondary backup. I have satisfied all the requirements with these two LLMs. 
  
***What was good with LLMs***
   1. Able to define transaction schema
   2. Understand context really well. Did not need to mention go and postgres with every question
   3. Came up with basic API structures
   4. Good with pointing out which librbaries are regularly maintained and which are not
   5. Good with design doc and test plan
   6. Good with specific questions regarding particular errors or scenarios. If ChatGpt's answer cannot fix it, Gemini's answer did.
   7. Able to write basic docker and compose.ymal

***What did not work well***
   1. Did not provided indexes at first. But did it after asking
   2. Folder hirarchy was minimal
   3. Did not use libraries (chi and gorm) at first. Did it after asking.
   4. Did not use environment variables on docker and compose file. Did it after prompting.

Overall, it was smooth and positive experience using LLM.
