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

	let draggedImage: TierImage | null = $state(null);
	let draggedFrom: SourceType | null = $state(null);
	let draggedIndex: number | null = $state(null);
	let activeDropTarget: { tier: SourceType; index: number } | null = $state(null);

	onMount(async () => {
		tierlistId = data.tierlistId;
		try {
			let response = await fetch(`${PUBLIC_API_URL}/tierlist/${tierlistId}`);
			if (response.ok) {
				let tierlist = await response.json();
				tierlistName = tierlist['Name'];
				uploadedImages = tierlist['UnassignedEntries'].map((x: any) => ({
					id: x['id'],
					src: `${PUBLIC_API_URL}/images/${x['file_key']}`
				}));
				tiers = tierlist['Tiers'].map((x: any) => ({
					id: x['id'],
					name: x['name'],
					entries: x['entries'].map((y: any) => ({
						id: y['id'],
						src: `${PUBLIC_API_URL}/images/${y['file_key']}`
					}))
				}));
			}
		} catch (error) {
			console.error("Failed to load tierlist:", error);
		}
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

	function handleDragStart(image: TierImage, from: SourceType, index: number) {
		draggedImage = image;
		draggedFrom = from;
		draggedIndex = index;
	}

	function handleDragEnd() {
		draggedImage = null;
		draggedFrom = null;
		draggedIndex = null;
		activeDropTarget = null;
	}

	function handleEnter(tier: SourceType, index: number) {
		activeDropTarget = { tier, index };
	}

	function handleLeave(tier: SourceType, index: number) {
		if (activeDropTarget?.tier === tier && activeDropTarget?.index === index) {
			activeDropTarget = null;
		}
	}

	async function handleDrop(target: SourceType, insertIndex: number) {
		if (!draggedImage || draggedFrom === null) return;

		if (draggedFrom === 'uploaded') {
			uploadedImages = uploadedImages.filter((img) => img.id !== draggedImage?.id);
		} else {
			tiers[draggedFrom].entries = tiers[draggedFrom].entries.filter(
				(img) => img.id !== draggedImage?.id
			);
		}

		if (target === 'uploaded') {
			const items = [...uploadedImages];
			items.splice(insertIndex, 0, draggedImage);
			uploadedImages = items;
		} else {
			const items = [...tiers[target].entries];
			items.splice(insertIndex, 0, draggedImage);
			tiers[target].entries = items;
		}

		await fetch(`${PUBLIC_API_URL}/tierlist/${tierlistId}/move`, {
			method: 'POST',
			body: JSON.stringify({
				TierID: target === 'uploaded' ? null : tiers[target].id,
				ID: draggedImage?.id
			})
		});

		handleDragEnd();
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
				class="flex items-center justify-start gap-2 overflow-x-auto border-l border-black bg-neutral-900 px-2"
				ondragover={allowDrop}
			>
				{#if tier.entries.length === 0}
					<div
						class="full-drop-target"
						ondragenter={() => handleEnter(i, 0)}
						ondragleave={() => handleLeave(i, 0)}
						ondragover={allowDrop}
						ondrop={() => handleDrop(i, 0)}
						role="listitem"
					>
						{#if activeDropTarget?.tier === i && activeDropTarget?.index === 0 && draggedImage}
							<img
								src={draggedImage.src}
								alt="ghost preview"
								class="h-16 w-16 rounded object-cover opacity-40 pointer-events-none"
							/>
						{:else}
							<p class="text-gray-500 italic select-none pointer-events-none">Drop items here...</p>
						{/if}
					</div>
				{:else}
					{#each tier.entries as image, index (image.id.toString())}
						<div
							class="drop-target"
							class:is-dragging={draggedImage !== null}
							class:drop-hover={activeDropTarget?.tier === i &&
								activeDropTarget?.index === index &&
								(draggedFrom !== i || (draggedIndex !== index && draggedIndex !== index - 1))}
							ondragenter={() => handleEnter(i, index)}
							ondragleave={() => handleLeave(i, index)}
							ondragover={allowDrop}
							ondrop={() => handleDrop(i, index)}
							role="listitem"
						>
							{#if activeDropTarget?.tier === i && activeDropTarget?.index === index && draggedImage && (draggedFrom !== i || (draggedIndex !== index && draggedIndex !== index - 1))}
								<img
									src={draggedImage.src}
									alt="ghost preview"
									class="h-16 w-16 rounded object-cover opacity-40 pointer-events-none"
								/>
							{/if}
						</div>

						<img
							src={image.src}
							alt="tier item"
							class="h-16 w-16 cursor-pointer rounded object-cover"
							draggable="true"
							ondragstart={() => handleDragStart(image, i, index)}
							ondragend={handleDragEnd}
						/>
					{/each}

					<div
						class="drop-target drop-expand"
						class:is-dragging={draggedImage !== null}
						class:drop-hover={activeDropTarget?.tier === i &&
							activeDropTarget?.index === tier.entries.length &&
							(draggedFrom !== i ||
								(draggedIndex !== tier.entries.length && draggedIndex !== tier.entries.length - 1))}
						ondragenter={() => handleEnter(i, tier.entries.length)}
						ondragleave={() => handleLeave(i, tier.entries.length)}
						ondragover={allowDrop}
						ondrop={() => handleDrop(i, tier.entries.length)}
						role="listitem"
					>
						{#if activeDropTarget?.tier === i && activeDropTarget?.index === tier.entries.length && draggedImage && (draggedFrom !== i || (draggedIndex !== tier.entries.length && draggedIndex !== tier.entries.length - 1))}
							<img
								src={draggedImage.src}
								alt="ghost preview"
								class="h-16 w-16 rounded object-cover opacity-40 pointer-events-none"
							/>
						{/if}
					</div>
				{/if}
			</div>
		</div>
	{/each}
</div>

<div
	role="list"
	aria-label="Uploaded items drop zone"
	class="mx-4 mb-6 mt-8 min-h-[120px] rounded-lg border-2 border-dashed border-gray-400 bg-gray-50 p-4"
	ondragover={allowDrop}
	ondrop={() => handleDrop('uploaded', uploadedImages.length)}
>
	<h2 class="mb-2 select-none text-lg font-semibold">Uploaded Items</h2>
	<div class="flex flex-wrap gap-4">
		{#each uploadedImages as image, index (image.id.toString())}
			<div
				class="drop-target"
				class:is-dragging={draggedImage !== null}
				class:drop-hover={activeDropTarget?.tier === 'uploaded' &&
					activeDropTarget?.index === index &&
					(draggedFrom !== 'uploaded' || (draggedIndex !== index && draggedIndex !== index - 1))}
				ondragenter={() => handleEnter('uploaded', index)}
				ondragleave={() => handleLeave('uploaded', index)}
				ondragover={allowDrop}
				ondrop={() => handleDrop('uploaded', index)}
				role="listitem"
			>
				{#if activeDropTarget?.tier === 'uploaded' && activeDropTarget?.index === index && draggedImage && (draggedFrom !== 'uploaded' || (draggedIndex !== index && draggedIndex !== index - 1))}
					<img
						src={draggedImage.src}
						alt="ghost preview"
						class="h-16 w-16 rounded object-cover opacity-40 pointer-events-none"
					/>
				{/if}
			</div>

			<img
				src={image.src}
				alt="uploaded item"
				class="h-16 w-16 cursor-pointer rounded object-cover"
				draggable="true"
				ondragstart={() => handleDragStart(image, 'uploaded', index)}
				ondragend={handleDragEnd}
			/>
		{/each}

		<div
			class="drop-target"
			class:is-dragging={draggedImage !== null}
			class:drop-hover={activeDropTarget?.tier === 'uploaded' &&
				activeDropTarget?.index === uploadedImages.length &&
				(draggedFrom !== 'uploaded' ||
					(draggedIndex !== uploadedImages.length && draggedIndex !== uploadedImages.length - 1))}
			ondragenter={() => handleEnter('uploaded', uploadedImages.length)}
			ondragleave={() => handleLeave('uploaded', uploadedImages.length)}
			ondragover={allowDrop}
			ondrop={() => handleDrop('uploaded', uploadedImages.length)}
			role="listitem"
		>
			{#if activeDropTarget?.tier === 'uploaded' && activeDropTarget?.index === uploadedImages.length && draggedImage && (draggedFrom !== 'uploaded' || (draggedIndex !== uploadedImages.length && draggedIndex !== uploadedImages.length - 1))}
				<img
					src={draggedImage.src}
					alt="ghost preview"
					class="h-16 w-16 rounded object-cover opacity-40 pointer-events-none"
				/>
			{/if}
		</div>
	</div>
</div>

<div class="mt-8 text-center">
	<label
		class="inline-block cursor-pointer select-none rounded bg-sky-500 px-4 py-2 font-bold text-white hover:bg-sky-700"
	>
		Upload Image
		<input type="file" accept="image/*" multiple class="hidden" onchange={handleUpload} />
	</label>
</div>

<style>
	.drop-target {
		position: relative;
		width: 4px;
		min-width: 4px;
		height: 64px;
		border: 2px dashed transparent;
		background-color: transparent;
		transition:
			width 0.2s ease,
			background-color 0.2s ease,
			border-color 0.2s ease;
		flex-shrink: 0;

		display: flex;
		align-items: center;
		justify-content: center;
	}

	.drop-target.is-dragging::after {
		content: '';
		position: absolute;
		top: 0;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);
		width: 60px;
		height: 100%;
		z-index: 20;
	}

	.drop-target.is-dragging.drop-hover::after {
		width: 120px;
	}

	.drop-hover {
		width: 64px;
		background-color: rgba(74, 222, 128, 0.1);
		border-color: #4ade80;
	}

	.full-drop-target {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: flex-start;
		border: 2px dashed #ccc;
		transition: border-color 0.2s;
		padding-left: 4px;
	}

	.full-drop-target:hover {
		border-color: #4ade80;
	}

	.drop-expand {
		flex-grow: 1;
		justify-content: flex-start;
		padding-left: 4px;
	}
</style>