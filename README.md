Stepup steps:
   1. Clone the repo
      git clone https://github.com/Ssenhub/transaction-crud-svc-go-postgres

   2.    
      

LLM experience:

  I have used Chatgpt as primary resource and Gemini as a secondary backup. I have satisfied all the requirements with these two LLMs. 
  
What was good with LLMs
   1. Able to define transaction schema
   2. Understand context really well. Did not need to mention go and postgres with every question
   3. Came up with basic API structures
   4. Good with pointing out which librbaries are regularly maintained and which are not
   5. Good with design doc and test plan
   6. Good with specific questions regarding particular errors or scenarios. If ChatGpt's answer cannot fix it, Gemini's answer did.
   7. Able to write basic docker and compose.ymal

What did not work well
   1. Did not provided indexes at first. But did it after asking
   2. Folder hirarchy was minimal
   3. Did not use libraries (chi and gorm) at first. Did it after asking.
   4. Did not use environment variables on docker and compose file. Did it after prompting.
