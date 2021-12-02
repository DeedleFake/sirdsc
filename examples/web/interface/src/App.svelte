<style>
	.app {
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.form {
		display: flex;
		flex-direction: column;

		color: white;
		background-color: blue;
		margin: 8px;
		padding: 8px;
		padding-top: 100px;
		transform: translate(0px, -100px);
		border-radius: 8px;
		box-shadow: #555555 4px 4px 4px;
	}

	.form .row {
		display: flex;
		flex-direction: row;
		justify-content: flex-end;

		margin: 0px 4px;
	}
</style>

<script>
	import Input from './Input.svelte'
	import Display from './Display.svelte'
	
	let params = new Map()
</script>

<div class='app'>
	<div class='form'>
		<div class='row'>
			<Input bind:value={params['src']} label='Depth Map' type='text' />
			<Input bind:value={params['pat']} label='Pattern' type='text' />
			<Input bind:value={params['seed']} label='Seed' type='number' defaultValue={0} />
			<Input bind:value={params['partsize']} label='Part Size' type='number' defaultValue={100} min={0} max={500} />
			<Input bind:value={params['depth']} label='Max Depth' type='number' defaultValue={40} min={0} max={50} />
			<Input bind:value={params['sym']} label='Symmetric Random Generation' type='checkbox' />
			<Input bind:value={params['inverse']} label='Inverse' type='checkbox' />
			<Input bind:value={params['flat']} label='Flat' type='checkbox' />
		</div>
	</div>

	<Display {params} />
</div>
