package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// SearchTest defines a single test case for the CI3 searcher
type SearchTest struct {
	Query           string
	ExpectedMatches []string // Substrings expected in the search output
	MinMatches      int      // Minimum number of expected matches in the output
}

func main() {
	fmt.Println("--- STARTING CI3 SEARCHER COMPREHENSIVE TEST ---")

	// Ensure mock project exists (creating it if not, or overwriting)
	fmt.Println("Setting up mock CI3 project...")
	setupMockProject()

	// Define test cases
	tests := []SearchTest{
		{
			Query:           "user model",
			ExpectedMatches: []string{"model", "User_model.php", "controller", "User.php"},
			MinMatches:      2,
		},
		{
			Query:           "admin dashboard view",
			ExpectedMatches: []string{"view", "admin_dashboard.php", "controller", "Admin.php"},
			MinMatches:      2,
		},
		{
			Query:           "database password config",
			ExpectedMatches: []string{"config", "database.php", "db", "ci_app_db"},
			MinMatches:      1, // Just checking config file
		},
		{
			Query:           "validate form javascript",
			ExpectedMatches: []string{"js", "common.js", "validateForm"},
			MinMatches:      1,
		},
		{
			Query:           "css font background",
			ExpectedMatches: []string{"css", "style.css", "font-family", "background-color"},
			MinMatches:      1,
		},
		{
			Query:           "controller user data", // Should hit user controller/model
			ExpectedMatches: []string{"controller", "User.php", "model", "User_model.php"},
			MinMatches:      2,
		},
		{
			Query:           "panel administration", // Keywords from admin_dashboard.php
			ExpectedMatches: []string{"view", "admin_dashboard.php", "administration panel"},
			MinMatches:      1,
		},
	}

	allTestsPassed := true
	for i, test := range tests {
		fmt.Printf("\n--- Running Test %d: Query: '%s' ---\n", i+1, test.Query)

		// Simulate search query followed by "q" to quit
		input := test.Query + "\nq\n"
		cmd := exec.Command("go", "run", "ci3_searcher.go", "tests/mock_ci3")
		cmd.Stdin = strings.NewReader(input)

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr

		// Start command in background to get real-time progress update
		// go func() {
		// 	cmd.Run()
		// }()
		// For simpler test execution, run synchronously.
		err := cmd.Run()
		if err != nil {
			fmt.Printf("[ERROR] Command for query '%s' failed: %v\n", test.Query, err)
			allTestsPassed = false
			continue
		}

		output := out.String()
		// fmt.Println(output) // Uncomment to see full output for debugging

		// Verify expected matches
		foundCount := 0
		for _, expected := range test.ExpectedMatches {
			if strings.Contains(output, expected) {
				foundCount++
				fmt.Printf("[PASS] Found expected match: '%s'\n", expected)
			} else {
				fmt.Printf("[FAIL] Could not find expected match: '%s'\n", expected)
			}
		}

		if foundCount >= test.MinMatches {
			fmt.Printf("[TEST PASS] Query '%s' yielded at least %d expected matches.\n", test.Query, test.MinMatches)
		} else {
			fmt.Printf("[TEST FAIL] Query '%s' only yielded %d/%d minimum expected matches.\n", test.Query, foundCount, test.MinMatches)
			fmt.Println("--- Full Output for Failed Test ---")
			fmt.Println(output)
			allTestsPassed = false
		}
		time.Sleep(500 * time.Millisecond) // Give a small pause between tests
	}

	if allTestsPassed {
		fmt.Println("\n--- ALL COMPREHENSIVE TESTS PASSED SUCCESSFULLY ---")
	} else {
		fmt.Println("\n--- SOME COMPREHENSIVE TESTS FAILED ---")
		os.Exit(1)
	}
}

// setupMockProject ensures the mock CI3 project structure and files are in place
func setupMockProject() {
	// Clean up existing mock project to ensure a fresh state
	os.RemoveAll("tests/mock_ci3")

	// Create directories
	dirs := []string{
		"tests/mock_ci3/application/controllers",
		"tests/mock_ci3/application/models",
		"tests/mock_ci3/application/views",
		"tests/mock_ci3/application/config",
		"tests/mock_ci3/assets/css",
		"tests/mock_ci3/assets/js",
		"tests/mock_ci3/system", // Should be skipped by searcher
		"tests/mock_ci3/vendor", // Should be skipped by searcher
	}
	for _, dir := range dirs {
		os.MkdirAll(dir, 0755)
	}

	// Create files
	files := map[string]string{
		"tests/mock_ci3/application/controllers/User.php": `<?php
class User extends CI_Controller {
    public function index() {
        $this->load->model('user_model');
        $data['users'] = $this->user_model->get_all();
        $this->load->view('user_list', $data);
    }
    public function profile($id) {
        $data['user'] = $this->user_model->get_user($id);
        $this->load->view('user_profile', $data);
    }
}
`,
		"tests/mock_ci3/application/models/User_model.php": `<?php
class User_model extends CI_Model {
    public function get_all() {
        return $this->db->get('users')->result();
    }
    public function get_user($id) {
        return $this->db->get_where('users', ['id' => $id])->row();
    }
}
`,
		"tests/mock_ci3/application/views/user_list.php": `<h1>User List</h1>
<ul>
    <?php foreach ($users as $user): ?>
        <li><?php echo $user->username; ?></li>
    <?php endforeach; ?>
</ul>
`,
		"tests/mock_ci3/application/views/user_profile.php": `<h1>User Profile</h1>
<p>Name: <?php echo $user->name; ?></p>
<p>Email: <?php echo $user->email; ?></p>
`,
		"tests/mock_ci3/application/config/database.php": `<?php
defined('BASEPATH') OR exit('No direct script access allowed');
$active_group = 'default';
$query_builder = TRUE;
$db['default'] = array(
    'dsn'	=> '',
    'hostname' => 'localhost',
    'username' => 'root',
    'password' => 'secret',
    'database' => 'ci_app_db',
    'dbdriver' => 'mysqli',
    'dbprefix' => '',
    'pconnect' => FALSE,
    'db_debug' => (ENVIRONMENT !== 'production'),
    'cache_on' => FALSE,
    'cachedir' => '',
    'char_set' => 'utf8',
    'dbcollat' => 'utf8_general_ci',
    'swap_pre' => '',
    'encrypt' => FALSE,
    'compress' => FALSE,
    'stricton' => FALSE,
    'failover' => array(),
    'save_queries' => TRUE
);
`,
		"tests/mock_ci3/application/controllers/Admin.php": `<?php
class Admin extends CI_Controller {
    public function dashboard() {
        $this->load->view('admin_dashboard');
    }
    public function users() {
        // Manage user data for admin
    }
}
`,
		"tests/mock_ci3/application/views/admin_dashboard.php": `<!DOCTYPE html>
<html>
<head>
    <title>Admin Dashboard</title>
    <link rel="stylesheet" href="../../assets/css/style.css">
</head>
<body>
    <h1>Welcome to Admin Dashboard</h1>
    <p>This is the administration panel.</p>
</body>
</html>
`,
		"tests/mock_ci3/assets/css/style.css": `body {
    font-family: Arial, sans-serif;
    background-color: #f0f0f0;
}
.container {
    width: 80%;
    margin: 0 auto;
    padding: 20px;
    border: 1px solid #ccc;
}
`,
		"tests/mock_ci3/assets/js/common.js": `function validateForm() {
    // Client-side validation logic
    console.log("Form validated!");
    return true;
}

$(document).ready(function() {
    // jQuery related scripts
    $('#myButton').click(function() {
        alert('Button clicked!');
    });
});
`,
		"tests/mock_ci3/system/core/CodeIgniter.php": `<?php // CI System file - should be skipped
	echo "This is the core CI file";
`,
		"tests/mock_ci3/vendor/autoload.php": `<?php // Composer vendor file - should be skipped
	require_once 'some_library/some_file.php';
`,
	}
	for path, content := range files {
		os.WriteFile(path, []byte(content), 0644)
	}
	fmt.Println("Mock CI3 project setup complete.")
}
