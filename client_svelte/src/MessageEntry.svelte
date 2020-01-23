<script>
    export let message;
</script>

<div>
    <span class="msguser {message.own ? 'own' : ''}">[{message.user}]</span>
    <span>{message.text}</span>
</div>

<style>
div {
    text-align: left;
}
.msguser {
    padding-right: .4em;
}
.msguser:not(.own) {
    color: green;
}
.msguser.own {
    color: lightblue;
}
</style>