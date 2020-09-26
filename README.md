## sending slack messages using oauth2 

app credentials :

https://api.slack.com/apps/A01BMJ27YDA

1 - get the app secrets (see previous snapshot)
2 - auth request




## Send messages using slack app webhook

slack enable webhook in selected channel :

![adding a webhook step](images/webhook.png)
then should be able to run :

curl -X POST -H 'Content-type: application/json' --data '{"text":"Hello, World!"}' https://hooks.slack.com/services/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



## run the container