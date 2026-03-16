<?php
class User_model extends CI_Model {
    public function get_all() {
        return $this->db->get('users')->result();
    }
    public function get_user($id) {
        return $this->db->get_where('users', ['id' => $id])->row();
    }
}
