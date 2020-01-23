<script>
    import { onMount } from 'svelte';
    import StatusBar from './StatusBar.svelte';
    import MessageEntry from './MessageEntry.svelte';
    import { ConnectionHandler, messages } from './connector.js';

    let msgInput;
    let currentMessage = '';

    function sendMessage(event) {
        ConnectionHandler.sendMessage(currentMessage);

        currentMessage = '';
    }

    onMount(() => {
        if (msgInput) {
            msgInput.focus();
        }
    });
</script>

<div class="mainView">
    <StatusBar/>

    <!-- Area where the chat appears -->
    <div class="messagesView">
        {#each $messages as msg}
        <MessageEntry message={msg}/>
        {/each}
    </div>

    <form on:submit|preventDefault={sendMessage}>
        <input type="text" placeholder="message..." bind:value={currentMessage} bind:this={msgInput}>
        <button>SEND</button>
    </form>
</div>

<style>
.mainView {
    display: flex;
    flex-direction: column;
    flex-grow: 1;
    width: 100%;
}

.messagesView {
    padding: .2em .4em;
    box-sizing: border-box;
    height: 100%;
    background-color: whitesmoke;
    overflow-x: hidden;
    overflow-y: auto;
}

form {
    display: flex;
    margin: .4em 0;
}

form > input[type="text"] {
    display: flex;
    align-content: left;
    margin-right: .4em;
    flex-grow: 1;
    /* width: 100%; */
    /* display: inline-block; */
    /* width: 100%; */
    /* width: auto; */
}

form > button {
    display: flex;
    align-content: right;
    /* display: inline-block; */
}
</style>