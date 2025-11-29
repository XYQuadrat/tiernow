<script lang="ts">
	type TierImage = { id: string; src: string };
	type SourceType = 'uploaded' | number;

	const tiers: string[] = ['S', 'A', 'B', 'C', 'D'];

	let uploadedImages: TierImage[] = [];
	let tierItems: TierImage[][] = tiers.map(() => []);

	let draggedImage: TierImage | null = null;
	let draggedFrom: SourceType | null = null;
	let activeDropTarget: { tier: SourceType; index: number } | null = null;

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

	function handleDragEnd() {
		draggedImage = null;
		draggedFrom = null;
		activeDropTarget = null;
	}

	function handleDrop(target: SourceType, index: number) {
		if (!draggedImage || draggedFrom === null) return;

		if (draggedFrom === 'uploaded') {
			uploadedImages = uploadedImages.filter((img) => img.id !== draggedImage?.id);
		} else {
			tierItems[draggedFrom] = tierItems[draggedFrom].filter((img) => img.id !== draggedImage?.id);
		}

		if (target === 'uploaded') {
			const items = [...uploadedImages];
			items.splice(index, 0, draggedImage);
			uploadedImages = items;
		} else {
			const items = [...tierItems[target]];
			items.splice(index, 0, draggedImage);
			tierItems[target] = items;
		}

		draggedImage = null;
		draggedFrom = null;
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

	function allowDrop(event: DragEvent) {
		event.preventDefault();
	}

	function getTierColor(index: number): string {
		const hue = (index * 35) % 360;
		return `hsl(${hue}, 70%, 65%)`;
	}
</script>

<h1 class="my-6 text-center text-3xl font-bold select-none">Tierlist Creator</h1>

<!-- Tier rows -->
<div class="mx-4 flex flex-col divide-y-2 border-y-2">
	{#each tiers as label, tierLevelIndex}
		<div class="grid h-24 grid-cols-[90px_1fr]">
			<div
				class="flex items-center justify-center text-xl font-bold text-gray-800 select-none"
				style="background-color: {getTierColor(tierLevelIndex)}"
			>
				{label}
			</div>
			<div
				role="list"
				aria-label="Tier {label} drop zone"
				class="flex items-center justify-start gap-2 overflow-x-auto border-l border-black bg-neutral-900 px-2"
				on:dragover={allowDrop}
			>
				{#if tierItems[tierLevelIndex].length === 0}
					<div
						class="full-drop-target"
						on:dragenter={() => handleEnter(tierLevelIndex, 0)}
						on:dragleave={() => handleLeave(tierLevelIndex, 0)}
						on:dragover={allowDrop}
						on:drop={() => handleDrop(tierLevelIndex, 0)}
					>
						{#if activeDropTarget?.tier === tierLevelIndex && activeDropTarget?.index === 0 && draggedImage}
							<img
								src={draggedImage.src}
								alt="ghost preview"
								class="h-16 w-16 rounded object-cover opacity-40 pointer-events-none"
							/>
						{:else}
							<span class="pointer-events-none text-gray-400 italic select-none">Drop items here...</span>
						{/if}
					</div>
				{:else}
					{#each tierItems[tierLevelIndex] as image, index (image.id)}
						{#if draggedImage?.id !== image.id}
							<div
								class="drop-target"
								class:drop-hover={activeDropTarget?.tier === tierLevelIndex && activeDropTarget?.index === index}
								on:dragenter={() => handleEnter(tierLevelIndex, index)}
								on:dragleave={() => handleLeave(tierLevelIndex, index)}
								on:dragover={allowDrop}
								on:drop={() => handleDrop(tierLevelIndex, index)}
							>
								{#if activeDropTarget?.tier === tierLevelIndex && activeDropTarget?.index === index && draggedImage}
									<img
										src={draggedImage.src}
										alt="ghost preview"
										class="h-16 w-16 rounded object-cover opacity-40 pointer-events-none"
									/>
								{/if}
							</div>
						{/if}

						<img
							src={image.src}
							alt="tier item"
							class="h-16 w-16 cursor-pointer rounded object-cover"
							draggable="true"
							on:dragstart={() => handleDragStart(image, tierLevelIndex)}
							on:dragend={handleDragEnd}
						/>
					{/each}

					<div
						class="drop-target drop-expand"
						class:drop-hover={activeDropTarget?.tier === tierLevelIndex && activeDropTarget?.index === tierItems[tierLevelIndex].length}
						on:dragenter={() => handleEnter(tierLevelIndex, tierItems[tierLevelIndex].length)}
						on:dragleave={() => handleLeave(tierLevelIndex, tierItems[tierLevelIndex].length)}
						on:dragover={allowDrop}
						on:drop={() => handleDrop(tierLevelIndex, tierItems[tierLevelIndex].length)}
					>
						{#if activeDropTarget?.tier === tierLevelIndex && activeDropTarget?.index === tierItems[tierLevelIndex].length && draggedImage}
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

<!-- Uploaded items -->
<div
	role="list"
	aria-label="Uploaded items drop zone"
	class="mx-4 mt-8 mb-6 min-h-[120px] rounded-lg border-2 border-dashed border-gray-400 bg-gray-50 p-4"
	on:dragover={allowDrop}
	on:drop={() => handleDrop('uploaded', uploadedImages.length)}
>
	<h2 class="mb-2 text-lg font-semibold select-none">Uploaded Items</h2>
	<div class="flex flex-wrap gap-4">
		{#each uploadedImages as image, index (image.id)}
			{#if draggedImage?.id !== image.id}
				<div
					class="drop-target"
					class:drop-hover={activeDropTarget?.tier === 'uploaded' && activeDropTarget?.index === index}
					on:dragenter={() => handleEnter('uploaded', index)}
					on:dragleave={() => handleLeave('uploaded', index)}
					on:dragover={allowDrop}
					on:drop={() => handleDrop('uploaded', index)}
				>
					{#if activeDropTarget?.tier === 'uploaded' && activeDropTarget?.index === index && draggedImage}
						<img
							src={draggedImage.src}
							alt="ghost preview"
							class="h-16 w-16 rounded object-cover opacity-40 pointer-events-none"
						/>
					{/if}
				</div>
			{/if}

			<img
				src={image.src}
				alt="uploaded item"
				class="h-16 w-16 cursor-pointer rounded object-cover"
				draggable="true"
				on:dragstart={() => handleDragStart(image, 'uploaded')}
				on:dragend={handleDragEnd}
			/>
		{/each}

		<div
			class="drop-target"
			class:drop-hover={activeDropTarget?.tier === 'uploaded' && activeDropTarget?.index === uploadedImages.length}
			on:dragenter={() => handleEnter('uploaded', uploadedImages.length)}
			on:dragleave={() => handleLeave('uploaded', uploadedImages.length)}
			on:dragover={allowDrop}
			on:drop={() => handleDrop('uploaded', uploadedImages.length)}
		>
			{#if activeDropTarget?.tier === 'uploaded' && activeDropTarget?.index === uploadedImages.length && draggedImage}
				<img
					src={draggedImage.src}
					alt="ghost preview"
					class="h-16 w-16 rounded object-cover opacity-40 pointer-events-none"
				/>
			{/if}
		</div>
	</div>
</div>

<!-- Upload button -->
<div class="mt-8 text-center">
	<label
		class="inline-block cursor-pointer rounded bg-sky-500 px-4 py-2 font-bold text-white select-none hover:bg-sky-700"
	>
		Upload Image
		<input type="file" accept="image/*" multiple class="hidden" on:change={handleUpload} />
	</label>
</div>

<style>
	.drop-target {
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
