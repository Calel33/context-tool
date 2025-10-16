# Prime Context for the AI Coding Assistant (catch it up to speed on the project when starting a new conversation)

Start with reading the AGENTS.md file in (root) to get an understanding of the project.
Read the README.md file to get an understanding of the project.
Read key files in the directories:

how to use code-index
User emphasizes using '**/.ext' for recursive search in Code Index MCP, not '.ext' which only matches the current directory.
🔍 Advanced Search Examples
Code Pattern Search
Search for all function calls matching "get.*Data" using regex
Finds: getData(), getUserData(), getFormData(), etc.

Fuzzy Function Search
Find authentication-related functions with fuzzy search for 'authUser'
Matches: authenticateUser, authUserToken, userAuthCheck, etc.

Language-Specific Search
Search for "API_ENDPOINT" only in Python files
Uses: search_code_advanced with file_pattern: "*.py"

Auto-refresh Configuration
Configure automatic index updates when files change
Uses: configure_file_watcher to enable/disable monitoring and set debounce timing

Project Maintenance
I added new components, please refresh the project index
Uses: refresh_index to update the searchable cache

### Knowledge Base
If our knowledge base dont have all the info or libraries you need you  can use Deepwiki and or octocode


Explain back to me:

Project structure
Project purpose and goals
Key files and their purposes
Any important dependencies
Any important configuration files
quick short summary
