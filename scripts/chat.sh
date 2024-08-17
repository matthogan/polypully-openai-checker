
package="Discord"

layout=$(cat <<EOF
{ 
    "name": "Example Software", 
    "description": "A brief description of what the software does.", 
    "benefits": "Description of benefits of using the software.",
    "categories": ["Category1", "Category2"], 
    "alternatives": [ 
        { 
            "name": "Alternative Software 1", 
            "description": "Description of the alternative software.", 
            "url": "https://example.com/alternative1"
        }
    ],
    "classification": { 
        "type": "Application", 
        "subcategory": "Productivity" 
    }, 
    "usage": { 
        "environment": ["List of typical usage environments such as corporate, home, or educational."],
        "instructions": "Instructions on how to use the software.",
        "age": "Minimum recommended age for usage.",
        "platforms": ["List of platforms the software runs on."],
        "languages": ["List of languages the software supports."],
        "license": "License information for the software.",
        "updates": "Information about how often the software is updated.",
        "installation": "Instructions on how to install the software.",
        "uninstallation": "Instructions on how to uninstall the software.",
        "features": ["List of features of the software."],
        "limitations": ["List of limitations of the software."]
    }, 
    "integration": [{ 
        "name": "Name of the software the software integrates with.", 
        "description": "Description of the integration." 
    }],
    "requirements": { 
        "minimum": "Minimum system requirements for the software.",
        "recommended": "Recommended system requirements for the software."
    },
    "safety": { 
        "corporate": { 
            "safe": true, 
            "description": "Explanation why it is safe for corporate use." 
        }, 
        "home": { 
            "safe": true, 
            "description": "Explanation why it is safe for home use." 
        }, 
        "school": { 
            "safe": true, 
            "description": "Explanation why it is safe for school use." 
        } 
    }, 
    "complexity": { 
        "score": 5, 
        "description": "Description of the complexity score." 
    } 
}
EOF
)

escaped=$(echo "$layout" | jq -c . | sed 's/"/\\"/g')
instructions="You are an assistant that provides information about software in a specific JSON format. "`
        `"Whenever you provide information, make sure it follows this structure: $escaped"

# echo "{
#     \"instructions\": \"$instructions\",
#   }"

body=$(cat <<EOF
{
  "model": "gpt-3.5-turbo",
  "messages": [
    {
      "role": "system",
      "content": [
        {
          "text": "$instructions",
          "type": "text"
        }
      ]
    },
    {
      "role": "user",
      "content": [
        {
          "type": "text",
          "text": "$package"
        }
      ]
    }
  ],
  "temperature": 0.5,
  "max_tokens": 1024,
  "top_p": 1,
  "frequency_penalty": 0,
  "presence_penalty": 0
}
EOF
)

# echo $body

curl https://api.openai.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -d "$body"
