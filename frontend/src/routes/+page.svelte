<script lang="ts">
	import { onMount } from 'svelte';

	let uuid: string;

	onMount(async () => {
		const response = await fetch('/api/tierlist', {
			method: 'POST',
			body: JSON.stringify({ name: 'New Tierlist', uuid: '' })
		});

		if (!response.ok) {
			console.error('Failed to create tierlist');
			return;
		}

		const data = await response.json();
		uuid = data.Uuid;
		if (!uuid) {
			console.error('Failed to get new tierlist UUID');
			return;
		}

		window.location.href = `/${uuid}`;
	});
</script>

<p>{uuid ? `Created tierlist ${uuid}` : 'Creating tierlist...'}</p>
