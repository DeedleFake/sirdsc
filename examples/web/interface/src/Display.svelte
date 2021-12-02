<style>
	.display {
		flex: 0 1 100%;
		margin: 8px;
	}
</style>

<script>
	export let params

	$: query = Object.entries(params).map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`)
	$: src = `http://${location.host}/generate?${query.join('&')}`
</script>

{#if params['src']}
	<img class='display' alt='Display' {src} />
{/if}
