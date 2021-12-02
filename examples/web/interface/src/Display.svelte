<script>
	export let params

	$: query = Object.entries(params).map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`)
	$: src = `http://${location.host}/generate?${query.join('&')}`
</script>

{#if params['src']}
	<img alt='Display' {src} />
{/if}
