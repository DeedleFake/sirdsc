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
	
	const handle = (ev) => {
		switch (type) {
			case 'text':
				value = ev.target.value
				break
				
			case 'number':
				value = !isNaN(ev.target.value) ? parseFloat(ev.target.value) : 0
				break
				
			case 'checkbox':
				value = ev.target.checked
				break
				
			default:
				throw new Error(`invalid input type: ${type}`)
		}
	}
</script>

<div class='input'>
	<span class='label'>{label}</span>
	<input {...$$restProps} {type} {value} on:input={handle} />
</div>
