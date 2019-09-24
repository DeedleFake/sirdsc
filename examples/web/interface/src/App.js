// @format

import React, { useState, useMemo, useReducer, useEffect } from 'react'

import './App.css'

const formInputs = [
	{
		label: 'Depth Map',
		type: 'text',
		name: 'src',
	},
	{
		label: 'Pattern',
		type: 'text',
		name: 'pat',
	},
	{
		label: 'Seed',
		type: 'number',
		name: 'seed',

		defaultValue: 0,
	},
	{
		label: 'Part Size',
		type: 'number',
		name: 'partsize',

		defaultValue: 100,
		min: 0,
		max: 500,
	},
	{
		label: 'Max Depth',
		type: 'number',
		name: 'depth',

		defaultValue: 40,
		min: 0,
		max: 50,
	},
	{
		label: 'Symmetric Random Generation',
		type: 'checkbox',
		name: 'sym',
	},
	{
		label: 'Inverse',
		type: 'checkbox',
		name: 'inverse',
	},
	{
		label: 'Flat',
		type: 'checkbox',
		name: 'flat',
	},
]

const defaultValues = {
	text: '',
	number: 0,
	checkbox: false,
}

const getValue = {
	text: (ev) => ev.target.value,
	number: (ev) => (!isNaN(ev.target.value) ? parseFloat(ev.target.value) : ''),
	checkbox: (ev) => ev.target.checked,
}

const Form = ({ onSubmit }) => {
	const [values, setValues] = useReducer(
		(values, action) => ({ ...values, ...action }),
		formInputs,
		(inputs) =>
			Object.fromEntries(
				inputs.map(({ name, type, defaultValue }) => [
					name,
					defaultValue || defaultValues[type],
				]),
			),
	)

	useEffect(() => {
		onSubmit(values)
	}, [values, onSubmit])

	return (
		<div className="Form">
			<div className="row">
				{formInputs.map(({ label, defaultValue, ...attr }) => (
					<div key={attr.name} className="section">
						<span className="label">{label}</span>
						<input
							{...attr}
							value={values[attr.name]}
							onChange={(ev) =>
								setValues({ [attr.name]: getValue[attr.type](ev) })
							}
						/>
					</div>
				))}
			</div>
		</div>
	)
}

const Display = ({ params }) => {
	const query = useMemo(() => {
		let q = Object.entries(params).map(
			([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`,
		)
		return `http://localhost:8080/generate?${q.join('&')}`
	}, [params])

	return (
		<img
			className={params.src != null ? 'Display' : ''}
			alt="Display"
			src={query}
		/>
	)
}

const App = () => {
	const [params, setParams] = useState({})

	return (
		<div className="App">
			<Form onSubmit={(params) => setParams(params)} />
			<Display params={params} />
		</div>
	)
}

export default App
