<script>
	export let params

	$: query = Object.entries(params).map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`)
	$: src = `http://localhost:8080/generate?${query.join('&')}`
</script>

<img class='display' alt='Display' {src} />
