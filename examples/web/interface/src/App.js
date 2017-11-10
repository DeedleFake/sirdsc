import React, { Component } from 'react';
import './App.css';

class Form extends Component {
	inputs = [
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
	];

	onSubmit(ev) {
		ev.preventDefault()

		this.props.onSubmit(this.inputs.reduce((p, v) => {
			const attrOf = (name, attr = 'value') => document.querySelector(`.Form input[name=${name}]`)[attr];

			switch (v.type) {
				case 'text':
				case 'number':
					p[v.name] = attrOf(v.name);
					break;

				case 'checkbox':
					p[v.name] = attrOf(v.name, 'checked');
					break;

				default:
					console.error(`Unexpected input type for ${v.name}: ${v.type}`);
					break;
			}

			return p;
		}, {}));

		return false;
	}

	render() {
		return (
			<form className='Form' onSubmit={this.onSubmit.bind(this)}>
				{this.inputs.map((v, i) => {
					let {label, ...attr} = v;
					return <div key={attr.name}>{label}: <input {...attr} /></div>;
				})}

				<input type='submit' value='Display' />
			</form>
		);
	}
}

function Display({params, ...props}) {
	let q = [];
	for (let f in params) {
		q.push(encodeURIComponent(f) + '=' + encodeURIComponent(params[f]));
	}

	// TODO: Display something to indicate loading.
	return (
		<img className={params.src !== undefined ? 'Display' : ''} alt='' src={'/generate?' + q.join('&')} />
	);
}

class App extends Component {
	state = {
		params: {},
	};

	onSubmit(params) {
		this.setState({
			params: params,
		});
	}

  render() {
    return (
			<div className='App'>
				<Form onSubmit={this.onSubmit.bind(this)} />
				<Display params={this.state.params} />
			</div>
    );
  }
}

export default App;
