<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { PUBLIC_API_URL } from '$env/static/public';

	let { data }: PageProps = $props();

	let tierlistId: string = $state('');
	let tierlistName: string = $state('');

	type TierImage = { id: Number; src: string };
	type Tier = { id: Number; name: string; entries: TierImage[] };
	type SourceType = 'uploaded' | number;

	let uploadedImages: TierImage[] = $state([]);
	let tiers: Tier[] = $state([]);

	let draggedImage: TierImage | null = null;
	let draggedFrom: SourceType | null = null;

	onMount(async () => {
		tierlistId = data.tierlistId;

		let response = await fetch(`${PUBLIC_API_URL}/tierlist/${tierlistId}`);
		let tierlist = await response.json();
		tierlistName = tierlist['Name'];
		uploadedImages = tierlist['UnassignedEntries'].map((x) => ({
			id: x['id'],
			src: `${PUBLIC_API_URL}/images/${x['file_key']}`
		}));
		tiers = tierlist['Tiers'].map((x) => ({
			id: x['id'],
			name: x['name'],
			entries: x['entries'].map((y) => ({
				id: y['id'],
				src: `${PUBLIC_API_URL}/images/${y['file_key']}`
			}))
		}));
	});

	async function handleUpload(event: Event) {
		const input = event.target as HTMLInputElement;
		const files = Array.from(input.files ?? []);
		let newImages = [];
		for (const file of files) {
			const formData = new FormData();
			formData.append('image', file);

			try {
				const response = await fetch(`${PUBLIC_API_URL}/tierlist/${tierlistId}/upload`, {
					method: 'POST',
					body: formData
				});
				const data = await response.json();
				newImages.push({ id: data['id'], src: `${PUBLIC_API_URL}/images/${data['filename']}` });
			} catch (err) {
				console.error(err);
			}
		}
		uploadedImages = [...uploadedImages, ...newImages];
	}

	async function handleDragStart(image: TierImage, from: SourceType) {
		draggedImage = image;
		draggedFrom = from;
	}

	async function handleDrop(target: SourceType) {
		if (!draggedImage || draggedFrom === null) return;

		if (draggedFrom === 'uploaded') {
			uploadedImages = uploadedImages.filter((img) => img.id !== draggedImage?.id);
		} else {
			tiers[draggedFrom].entries = tiers[draggedFrom].entries.filter(
				(img) => img.id !== draggedImage?.id
			);
		}

		if (target === 'uploaded') {
			uploadedImages = [...uploadedImages, draggedImage];
		} else {
			tiers[target].entries = [...tiers[target].entries, draggedImage];
		}

		const response = await fetch(`${PUBLIC_API_URL}/tierlist/${tierlistId}/move`, {
			method: 'POST',
			body: JSON.stringify({
				TierID: target === 'uploaded' ? null : tiers[target].id,
				ID: draggedImage?.id
			})
		});

		draggedImage = null;
		draggedFrom = null;
	}

	function allowDrop(event: DragEvent) {
		event.preventDefault();
	}

	function getTierColor(index: number): string {
		const hue = (index * 35) % 360;
		return `hsl(${hue}, 70%, 65%)`;
	}
</script>

<h1 class="my-6 select-none text-center text-3xl font-bold">Tierlist Creator</h1>

<!-- Tier rows -->
<div class="mx-4 flex flex-col divide-y-2 border-y-2">
	{#each tiers as tier, i}
		<div class="grid h-24 grid-cols-[90px_1fr]">
			<div
				class="flex select-none items-center justify-center text-xl font-bold text-gray-800"
				style="background-color: {getTierColor(i)}"
			>
				{tier.name}
			</div>
			<div
				role="list"
				aria-label="Tier {tier.name} drop zone"
				class="flex items-center gap-2 overflow-x-auto border-l border-black bg-neutral-900 p-2"
				ondragover={allowDrop}
				ondrop={() => handleDrop(i)}
			>
				{#each tier.entries as image}
					<img
						src={image.src}
						alt="tier item"
						class="h-16 w-16 cursor-pointer rounded object-cover"
						draggable="true"
						ondragstart={() => handleDragStart(image, i)}
					/>
				{:else}
					<p class="text-gray-500 italic select-none">Drop items here...</p>
				{/each}
			</div>
		</div>
	{/each}
</div>

<div
	role="list"
	aria-label="Uploaded items drop zone"
	class="mx-4 mb-6 mt-8 min-h-[120px] rounded-lg border-2 border-dashed border-gray-400 bg-gray-50 p-4"
	ondragover={allowDrop}
	ondrop={() => handleDrop('uploaded')}
>
	<h2 class="mb-2 select-none text-lg font-semibold">Uploaded Items</h2>
	<div class="flex flex-wrap gap-4">
		{#each uploadedImages as image}
			<img
				src={image.src}
				alt="uploaded item"
				class="h-16 w-16 cursor-pointer rounded object-cover"
				draggable="true"
				ondragstart={() => handleDragStart(image, 'uploaded')}
			/>
		{:else}
			<p class="text-gray-500 italic select-none">No images uploaded yet.</p>
		{/each}
	</div>
</div>

<!-- Upload button -->
<div class="mt-8 text-center">
	<label
		class="inline-block cursor-pointer select-none rounded bg-sky-500 px-4 py-2 font-bold text-white hover:bg-sky-700"
	>
		Upload Image
		<input type="file" accept="image/*" multiple class="hidden" onchange={handleUpload} />
	</label>
</div>
