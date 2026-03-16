<h1>User List</h1>
<ul>
    <?php foreach ($users as $user): ?>
        <li><?php echo $user->username; ?></li>
    <?php endforeach; ?>
</ul>
