thread="thread_skwB2DYXV0lV3DUNNmQFbHvf"

curl https://api.openai.com/v1/threads/${thread}/runs \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "Content-Type: application/json" \
  -H "OpenAI-Beta: assistants=v2" \
  -d '{
    "assistant_id": "asst_HdUKhBFNhhd5ErkZ0DnuKkuN",
    "additional_instructions": null,
    "tool_choice": null
  }'
