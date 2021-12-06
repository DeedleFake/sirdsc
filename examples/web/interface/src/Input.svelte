<style>
	.input {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;

		margin: 4px 4px;
	}
</style>

<script>
	const defaults = {
		text: '',
		number: 0,
		checkbox: false,
	}
	
	export let type = 'text'
	export let label
	
	export let defaultValue = defaults[type]
	export let value = defaultValue
</script>

<div class='input'>
	<span class='label'>{label}</span>
	{#if type === 'text'}
		<input type="text" {...$$restProps} bind:value />
	{:else if type === 'number'}
		<input type="number" {...$$restProps} bind:value />
	{:else if type === 'checkbox'}
		<input type="checkbox" {...$$restProps} bind:value />
	{/if}
</div>
