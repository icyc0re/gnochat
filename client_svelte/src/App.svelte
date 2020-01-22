<script>
	import JoinForm from './JoinForm.svelte';

	export let name = 'SomeApp';
	export let version = "0.0.0";
</script>

<main>
	<h2>{name} v{version}</h2>

	<JoinForm/>
</main>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>