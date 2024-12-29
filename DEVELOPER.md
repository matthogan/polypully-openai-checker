# Developer

Notes for development.

## Google Cloud environment

Google Cloud SDK

```shell
curl https://sdk.cloud.google.com | bash
exec -l $SHELL
gcloud init
```

Access

Key with read access to the secrets manager.

App uses the `GOOGLE_APPLICATION_CREDENTIALS_JSON` env var.

```shell
export GOOGLE_APPLICATION_CREDENTIALS_JSON=$(base64 -w 0 ~/sasecrdr1key.json)
```

## openai

API key for OpenAI stored in the `OPENAI_API_KEY` env var.

```shell
$ export OPENAI_API_KEY=...
$ gcloud secrets create OPENAI_API_KEY --replication-policy="automatic"
Created secret [OPENAI_API_KEY].
$echo -n "$OPENAI_API_KEY" | gcloud secrets versions add OPENAI_API_KEY --data-file=-
Created version [1] of the secret [OPENAI_API_KEY].
```

## Dependencies

```shell
go get cloud.google.com/go/secretmanager/apiv1
```

## Create assistant
Requires an api key with assistant write permission.

```shell
➜  polypully-openai-checker export OPENAI_API_KEY=...
➜  polypully-openai-checker ./assistant-create.sh                                                         
{
  "id": "asst_HdUKhBFNhhd5ErkZ0DnuKkuN",
  "object": "assistant",
  "created_at": 1721949495,
  "name": "Software Information Assistant",
  "description": null,
  "model": "gpt-4o",
  "instructions": "You are an assistant that provides information about software in a specific JSON format. Whenever you provide information, make sure it follows this structure: {\"name\":\"Example Software\",\"description\":\"A brief description of what the software does.\",\"benefits\":\"Description of benefits of using the software.\",\"categories\":[\"Category1\",\"Category2\"],\"alternatives\":[{\"name\":\"Alternative Software 1\",\"description\":\"Description of the alternative software.\",\"url\":\"https://example.com/alternative1\"}],\"classification\":{\"type\":\"Application\",\"subcategory\":\"Productivity\"},\"usage\":{\"environment\":[\"List of typical usage environments such as corporate, home, or educational.\"],\"instructions\":\"Instructions on how to use the software.\",\"age\":\"Minimum recommended age for usage.\",\"platforms\":[\"List of platforms the software runs on.\"],\"languages\":[\"List of languages the software supports.\"],\"license\":\"License information for the software.\",\"updates\":\"Information about how often the software is updated.\",\"installation\":\"Instructions on how to install the software.\",\"uninstallation\":\"Instructions on how to uninstall the software.\",\"features\":[\"List of features of the software.\"],\"limitations\":[\"List of limitations of the software.\"]},\"integration\":[{\"name\":\"Name of the software the software integrates with.\",\"description\":\"Description of the integration.\"}],\"requirements\":{\"minimum\":\"Minimum system requirements for the software.\",\"recommended\":\"Recommended system requirements for the software.\"},\"safety\":{\"corporate\":{\"safe\":true,\"description\":\"Explanation why it is safe for corporate use.\"},\"home\":{\"safe\":true,\"description\":\"Explanation why it is safe for home use.\"},\"school\":{\"safe\":true,\"description\":\"Explanation why it is safe for school use.\"}},\"complexity\":{\"score\":5,\"description\":\"Description of the complexity score.\"}}",
  "tools": [],
  "top_p": 1.0,
  "temperature": 1.0,
  "tool_resources": {},
  "metadata": {},
  "response_format": "auto"
}%    
```

## Create thread
Requires an api key with thread write permission.

```shell
➜  polypully-openai-checker export OPENAI_API_KEY=...
➜  polypully-openai-checker ./assistant-thread.sh                                                         
{
  "id": "thread_AhDW5k0Ood6SgyeyUwJnOjM0",
  "object": "thread",
  "created_at": 1721949711,
  "metadata": {},
  "tool_resources": {}
}%
```

## Add message to thread

```shell
➜  polypully-openai-checker ./assistant-thread-message.sh
{
  "id": "msg_qxELb3hObdu7RSftxlGDO5F4",
  "object": "thread.message",
  "created_at": 1721950351,
  "assistant_id": null,
  "thread_id": "thread_AhDW5k0Ood6SgyeyUwJnOjM0",
  "run_id": null,
  "role": "user",
  "content": [
    {
      "type": "text",
      "text": {
        "value": "Sensible Soccer",
        "annotations": []
      }
    }
  ],
  "attachments": [],
  "metadata": {}
}
```

## Run the thread

```shell
➜  polypully-openai-checker ./assistant-thread-run.sh    
{
  "id": "run_vtw6jurRZskeMoNQtu8oUORx",
  "object": "thread.run",
  "created_at": 1721950856,
  "assistant_id": "asst_HdUKhBFNhhd5ErkZ0DnuKkuN",
  "thread_id": "thread_AhDW5k0Ood6SgyeyUwJnOjM0",
  "status": "queued",
  "started_at": null,
  "expires_at": 1721951456,
  "cancelled_at": null,
  "failed_at": null,
  "completed_at": null,
  "required_action": null,
  "last_error": null,
  "model": "gpt-4o",
  "instructions": "You are an assistant that provides information about software in a specific JSON format. Whenever you provide information, make sure it follows this structure: {\"name\":\"Example Software\",\"description\":\"A brief description of what the software does.\",\"benefits\":\"Description of benefits of using the software.\",\"categories\":[\"Category1\",\"Category2\"],\"alternatives\":[{\"name\":\"Alternative Software 1\",\"description\":\"Description of the alternative software.\",\"url\":\"https://example.com/alternative1\"}],\"classification\":{\"type\":\"Application\",\"subcategory\":\"Productivity\"},\"usage\":{\"environment\":[\"List of typical usage environments such as corporate, home, or educational.\"],\"instructions\":\"Instructions on how to use the software.\",\"age\":\"Minimum recommended age for usage.\",\"platforms\":[\"List of platforms the software runs on.\"],\"languages\":[\"List of languages the software supports.\"],\"license\":\"License information for the software.\",\"updates\":\"Information about how often the software is updated.\",\"installation\":\"Instructions on how to install the software.\",\"uninstallation\":\"Instructions on how to uninstall the software.\",\"features\":[\"List of features of the software.\"],\"limitations\":[\"List of limitations of the software.\"]},\"integration\":[{\"name\":\"Name of the software the software integrates with.\",\"description\":\"Description of the integration.\"}],\"requirements\":{\"minimum\":\"Minimum system requirements for the software.\",\"recommended\":\"Recommended system requirements for the software.\"},\"safety\":{\"corporate\":{\"safe\":true,\"description\":\"Explanation why it is safe for corporate use.\"},\"home\":{\"safe\":true,\"description\":\"Explanation why it is safe for home use.\"},\"school\":{\"safe\":true,\"description\":\"Explanation why it is safe for school use.\"}},\"complexity\":{\"score\":5,\"description\":\"Description of the complexity score.\"}}",
  "tools": [],
  "tool_resources": {},
  "metadata": {},
  "temperature": 1.0,
  "top_p": 1.0,
  "max_completion_tokens": null,
  "max_prompt_tokens": null,
  "truncation_strategy": {
    "type": "auto",
    "last_messages": null
  },
  "incomplete_details": null,
  "usage": null,
  "response_format": "auto",
  "tool_choice": "auto",
  "parallel_tool_calls": true
}
```

## List services

```shell
grpc_cli ls localhost:50051 -l
```

## Call with grpcurl

```shell
grpcurl -plaintext -d '{"name":"Notepad++"}' localhost:50051 software.SoftwareInfoService/GetSoftwareInfo
```

With mtls:
    
```shell
grpcurl -cert .secret/client.crt -key .secret/client.key -cacert .secret/ca.crt -d '{"name":"Notepad++"}' localhost:50051 software.SoftwareInfoService/GetSoftwareInfo
```
