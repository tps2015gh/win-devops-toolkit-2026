<?php
class User extends CI_Controller {
    public function index() {
        $this->load->model('user_model');
        $data['users'] = $this->user_model->get_all();
        $this->load->view('user_list', $data);
    }
}
