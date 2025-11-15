## core

[![Deploy to Koyeb](https://www.koyeb.com/static/images/deploy/button.svg)](https://app.koyeb.com/deploy?name=sillycore&type=git&repository=VerbTeam%2Fcore&branch=main&builder=dockerfile&instance_type=free&regions=fra&instances_min=0&autoscaling_sleep_idle_delay=3600&env%5BGEMINI_API_KEY%5D=YOUR_GEMINI_API_KEY_HERE&env%5BREDIS_PASSWORDS%5D=YOUR_REDIS_PASSWORD_HERE&env%5BREDIS_PUBLIC_ENDPOINT%5D=YOUR_REDIS_ENDPOINT_HERE&env%5BREDIS_USERNAME%5D=YOUR_REDIS_USERNAME_HERE&env%5BSUPABASE_URL%5D=YOUR_SUPABASE_URL_HERE)

this is the main service powering the moderation system.

the application uses a custom machine-learning model named **sybauML**, trained on the **FYM dataset** from hugging face. the model is fine-tuned from **facebookai/roberta-base**, allowing it to perform accurate text classification and detect inappropriate or unsafe roblox bios efficiently.

after sybauML performs its initial classification, the system uses **gemini ai** for secondary analysis and contextual validation. this layered approach improves accuracy, reduces false positives, and ensures consistent moderation results.

the service is hosted on **koyeb** for reliability and easy deployment.
if you choose to deploy it on koyeb, make sure to replace all environment variables with your own credentials.

the public endpoint for this service is:

```
https://sillycore.koyeb.app/
```
