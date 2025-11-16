## core

[![Deploy to Koyeb](https://www.koyeb.com/static/images/deploy/button.svg)](https://app.koyeb.com/deploy?name=sillycore&type=git&repository=VerbTeam%2Fcore&branch=main&builder=dockerfile&instance_type=free&regions=fra&instances_min=0&autoscaling_sleep_idle_delay=3600&env%5BGEMINI_API_KEY%5D=YOUR_GEMINI_API_KEY_HERE&env%5BREDIS_PASSWORDS%5D=YOUR_REDIS_PASSWORD_HERE&env%5BREDIS_PUBLIC_ENDPOINT%5D=YOUR_REDIS_ENDPOINT_HERE&env%5BREDIS_USERNAME%5D=YOUR_REDIS_USERNAME_HERE&env%5BSUPABASE_URL%5D=YOUR_SUPABASE_URL_HERE)

this is the main service powering the moderation system.

the application have **two-type classification pipeline**:

### 1. **machine-learning classification (sybauML)**

the first stage uses **sybauML**, a custom model fine-tuned from **facebookai/roberta-base** and trained on the **FYM dataset** from hugging face. this model performs fast, reliable text classification specifically optimized for roblox-style bios, detecting unsafe, inappropriate, or disallowed content.

### 2. **gemini ai classification**

the text is sent to **gemini ai** for secondary classification. gemini adds contextual reasoning, edge-case detection, and higher-level semantic checks to reduce false positives and handle cases that require more nuanced understanding.



the service is hosted on **koyeb**, and you can deploy it directly using the button above.
if you deploy it yourself, make sure to replace all environment variables with your own credentials.

public endpoint:

```
https://sillycore.koyeb.app/
```
