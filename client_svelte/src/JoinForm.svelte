<script>
    let username = '';

    function joinChat() {
        console.log(username)
        // TODO send message to the server
    }
</script>

<form on:submit|preventDefault={joinChat}>
    <input type="value" placeholder="username" bind:value={username}>
    <button>JOIN!</button>
</form>
