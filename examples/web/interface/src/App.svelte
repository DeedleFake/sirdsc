<script>
	import Input from './Input.svelte'
	import Display from './Display.svelte'
	
	let params = new Map()
</script>

<Input bind:value={params['src']} label='Depth Map' type='text' />
<Input bind:value={params['pat']} label='Pattern' type='text' />
<Input bind:value={params['seed']} label='Seed' type='number' defaultValue={0} />
<Input bind:value={params['partsize']} label='Part Size' type='number' defaultValue={100} min={0} max={500} />
<Input bind:value={params['depth']} label='Max Depth' type='number' defaultValue={40} min={0} max={50} />
<Input bind:value={params['sym']} label='Symmetric Random Generation' type='checkbox' />
<Input bind:value={params['inverse']} label='Inverse' type='checkbox' />
<Input bind:value={params['flat']} label='Flat' type='checkbox' />

<Display {params} />
