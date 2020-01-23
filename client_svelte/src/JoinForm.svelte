<script>
    import { ConnectionHandler } from './connector.js';

    let usernameInput;
    let username = '';

    function joinChat() {
        if (usernameInput) {
            usernameInput.disabled = true;
        }

        ConnectionHandler.connect(function () {
            ConnectionHandler.handshake(username);
        });
    }
</script>

<form on:submit|preventDefault={joinChat}>
    <input type="value" placeholder="username" bind:value={username} bind:this={usernameInput}>
    <button>JOIN!</button>
</form>
