<?php
class Admin extends CI_Controller {
    public function dashboard() {
        $this->load->view('admin_dashboard');
    }
    public function users() {
        // Manage user data for admin
    }
}
