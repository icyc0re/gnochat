<script>
	import JoinForm from './JoinForm.svelte';
	import ChatView from './ChatView.svelte';

	import { connected } from './connector.js';

	export let name = 'SomeApp';
	export let version = "0.0.0";
</script>

<main>
	<h2>{name} v{version}</h2>

	{#if $connected}
	<ChatView/>
	{:else}
	<JoinForm/>
	{/if}
	<!-- <div>Connected: {$connected}</div>
	<JoinForm/>
	<ChatView/> -->
</main>

<style>
	main {
		display: flex;
		flex-direction: column;
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
		height: 100%;
		box-sizing: border-box;
		overflow: hidden;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>