thread="thread_skwB2DYXV0lV3DUNNmQFbHvf"

curl https://api.openai.com/v1/threads/${thread}/messages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "OpenAI-Beta: assistants=v2"
