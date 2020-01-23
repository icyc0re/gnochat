<script>
    import { ConnectionHandler } from './connector.js';
</script>

<div class="statusbar">
    <div class="statusmessage"></div>
    <div class="username">Connected: {ConnectionHandler.getUsername()}</div>
</div>

<style>
.statusmessage {
    float: left;
}
.username {
    float: right;
    color: green;
}
</style>