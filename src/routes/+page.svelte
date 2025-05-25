<script lang="ts">
	type TierImage = { id: string; src: string };
	type SourceType = 'uploaded' | number;

	const tiers: string[] = ['S', 'A', 'B', 'C', 'D'];

	let uploadedImages: TierImage[] = [];
	let tierItems: TierImage[][] = tiers.map(() => []);

	let draggedImage: TierImage | null = null;
	let draggedFrom: SourceType | null = null;

	function handleUpload(event: Event) {
		const input = event.target as HTMLInputElement;
		const files = Array.from(input.files ?? []);
		const newImages = files.map((file) => ({
			id: crypto.randomUUID(),
			src: URL.createObjectURL(file)
		}));
		uploadedImages = [...uploadedImages, ...newImages];
	}

	function handleDragStart(image: TierImage, from: SourceType) {
		draggedImage = image;
		draggedFrom = from;
	}

	function handleDrop(target: SourceType) {
		if (!draggedImage || draggedFrom === null) return;

		if (draggedFrom === 'uploaded') {
			uploadedImages = uploadedImages.filter((img) => img.id !== draggedImage?.id);
		} else {
			tierItems[draggedFrom] = tierItems[draggedFrom].filter((img) => img.id !== draggedImage?.id);
		}

		if (target === 'uploaded') {
			uploadedImages = [...uploadedImages, draggedImage];
		} else {
			tierItems[target] = [...tierItems[target], draggedImage];
		}

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

<h1 class="my-6 text-center text-3xl font-bold">Tierlist Creator</h1>

<!-- Tier rows -->
<div class="mx-4 flex flex-col divide-y-2 border-y-2">
	{#each tiers as label, i}
		<div class="grid h-24 grid-cols-[90px_1fr]">
			<div
				class="flex items-center justify-center text-xl font-bold text-gray-800"
				style="background-color: {getTierColor(i)}"
			>
				{label}
			</div>
			<div
				role="list"
				aria-label="Tier {label} drop zone"
				class="flex items-center gap-2 overflow-x-auto border-l border-black bg-neutral-900 p-2"
				on:dragover={allowDrop}
				on:drop={() => handleDrop(i)}
			>
				{#if tierItems[i].length === 0}
					<p class="text-gray-500 italic">Drop items here...</p>
				{:else}
					{#each tierItems[i] as image}
						<img
							src={image.src}
							alt="tier item"
							class="h-16 w-16 rounded object-cover"
							draggable="true"
							on:dragstart={() => handleDragStart(image, i)}
						/>
					{/each}
				{/if}
			</div>
		</div>
	{/each}
</div>

<div
	role="list"
	aria-label="Uploaded items drop zone"
	class="mx-4 mt-8 mb-6 min-h-[120px] rounded-lg border-2 border-dashed border-gray-400 bg-gray-50 p-4"
	on:dragover={allowDrop}
	on:drop={() => handleDrop('uploaded')}
>
	<h2 class="mb-2 text-lg font-semibold">Uploaded Items</h2>
	<div class="flex flex-wrap gap-4">
		{#if uploadedImages.length === 0}
			<p class="text-gray-500 italic">No images uploaded yet.</p>
		{:else}
			{#each uploadedImages as image}
				<img
					src={image.src}
					alt="uploaded item"
					class="h-16 w-16 rounded object-cover"
					draggable="true"
					on:dragstart={() => handleDragStart(image, 'uploaded')}
				/>
			{/each}
		{/if}
	</div>
</div>

<!-- Upload button -->
<div class="mt-8 text-center">
	<label
		class="inline-block cursor-pointer rounded bg-sky-500 px-4 py-2 font-bold text-white hover:bg-sky-700"
	>
		Upload Image
		<input type="file" accept="image/*" multiple class="hidden" on:change={handleUpload} />
	</label>
</div>
