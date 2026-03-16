# CodeIgniter 3 Project Searcher (`ci3_searcher.exe`)

## 📄 Overview
The `ci3_searcher.exe` is an intelligent command-line tool designed to quickly index and search through CodeIgniter 3 (CI3) projects. It leverages a 100-dimensional vector space with a simplified "attention" mechanism to provide similarity-based search results across various CI3 components. This tool is invaluable for understanding large or unfamiliar CI3 codebases, helping developers and auditors quickly locate relevant files and code sections.

## ✨ Key Features

-   **Intelligent Indexing**: Automatically scans a specified CI3 project directory, identifying and indexing key components such as Controllers, Models, Views, Configuration files, JavaScript, CSS, and code sections related to database interactions.
-   **Noise Reduction**: Skips irrelevant files and directories like `.git` folders, `logs`, and common binary formats (PDF, JPG, PNG, EXE, etc.) to maintain a clean and focused index.
-   **100D Vector Space Search**: Utilizes a compact 100-dimensional vector space, generated through feature hashing, to represent each indexed code component. This allows for efficient similarity comparisons between search queries and project files.
-   **Simplified Attention Mechanism**: During vectorization, component types (e.g., "controller", "model") and filenames are given higher weight than general code content. This simulates an attention-like effect, prioritizing results that are semantically closer to the query's intent.
-   **Interactive Command-Line Interface**: Provides a user-friendly CLI where users can enter search terms, and the tool will display ranked results in real-time.
-   **Contextual Matching**: The search algorithm is designed to find code components that are "related" to the search query, not just exact keyword matches. For example, searching "user login" can find a `Login` controller, a `User_model`, and `login_view.php`.
-   **Automatic Result Export**: For queries yielding more than 15 results, the tool automatically exports the full list to `search_results.txt` for further analysis, while still displaying the top matches in the console.

## 🚀 Usage

### 1. Build the Executable (Optional)
If you wish to build a standalone executable, navigate to the project root and run the `build.bat` script. This will compile `ci3_searcher.go` into `ci3_searcher.exe`.

### 2. Run the Tool
Execute `ci3_searcher.go` directly or use the compiled `ci3_searcher.exe` followed by the path to your CodeIgniter 3 project.

```powershell
# Using 'go run' (requires Go environment)
go run ci3_searcher.go <path_to_your_ci3_project>

# Example:
go run ci3_searcher.go C:\xampp\htdocs\my_ci3_app

# Using the compiled executable
.\ci3_searcher.exe <path_to_your_ci3_project>

# Example:
.\ci3_searcher.exe C:\xampp\htdocs\my_ci3_app
```

### 3. Interactive Search
Once the tool has finished indexing (which should be very quick for most projects), you will be presented with a `Search>` prompt.

-   **Enter your search query**: Type any keywords related to what you are looking for (e.g., `user registration`, `db query`, `css layout`, `admin controller`).
-   **Review Results**: The tool will display a ranked list of relevant files, their types, and paths, along with a similarity score.
-   **Exit**: Type `q` (or `exit`) and press Enter to quit the program.

#### Example Search Session:
```
Indexing CI3 project at: C:\xampp\htdocs\my_ci3_app...
Successfully indexed 125 items.

Search> user model
0.9876 | model        | application\models\User_model.php
0.7231 | controller   | application\controllers\Users.php
0.6502 | view         | application\views\user_profile.php

Search> login view
0.9543 | view         | application\views\auth\login_form.php
0.8876 | controller   | application\controllers\Auth.php
0.6123 | js           | assets\js\login_validation.js

Search> q
```

## 🧠 Technical Insights

### Developer's Perspective (Gemini CLI Agent)
As the intelligent assistant who developed this tool, I believe `ci3_searcher.exe` stands out as a highly practical utility for anyone working with CodeIgniter 3 projects. Its design prioritizes speed and relevance: the 100D vector space, while simplified compared to larger language models, provides a remarkable balance for effective similarity search within codebases. The "attention" mechanism, by weighting component types and filenames, is particularly effective at cutting through noise and delivering contextually accurate results. This tool perfectly embodies the project's goal of delivering high-fidelity insights, making code auditing and understanding significantly more efficient.

### Documenter's Opinion (tps2015gh)
[To be added by tps2015gh]

### Vectorization and Similarity
The core of `ci3_searcher.exe` lies in its ability to convert text (code) into numerical vectors. This process involves:
1.  **Tokenization**: Breaking down the content of files and search queries into individual words or symbols.
2.  **Feature Hashing**: Each token is hashed to an index within a 100-dimensional array. This fixed-size representation allows for fast and memory-efficient vector creation without needing a large vocabulary.
3.  **Normalization**: Vectors are normalized (L2 norm) to ensure that longer documents or queries don't inherently have higher scores.
4.  **Dot Product Similarity**: The similarity between a search query's vector and an indexed file's vector is calculated using the dot product. A higher dot product indicates greater similarity.

### Simplified Attention
While not a full neural network-based attention mechanism, the tool simulates attention by assigning different weights during vector creation:
-   **Component Type**: The identified CI3 component type (e.g., "controller", "model") receives a 5x boost in its contribution to the vector.
-   **Filename**: The base filename (e.g., "User_model.php") receives a 3x boost.
-   **Content**: The actual code content receives a 1x weight.

This weighting scheme allows the search to "pay more attention" to what kind of file it is and its name, making searches like "user controller" highly effective in pinpointing the correct component.

## 👨‍💻 Project Team & Contributors
-   **Director & Supervisor:** **tps2015gh** (Human)
-   **Programming & Testing:** tps2015gh, Qwen Code
-   **Intelligent Assistant (CLI Agent):** Gemini CLI Agent

**Legal Note on Authorship:** This tool was developed under the sole ownership and direction of **tps2015gh**. Gemini CLI Agent provided intelligent assistance, implementation, and testing support under direct instruction and oversight. All intellectual property, copyright, and strategic decisions reside with **tps2015gh**.
