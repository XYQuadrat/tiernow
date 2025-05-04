<script lang="ts">
  let uploadedImages: string[] = [];

  function handleUpload(event: Event) {
    const input = event.target as HTMLInputElement;
    const files = Array.from(input.files ?? []);
    uploadedImages = files.map(file => URL.createObjectURL(file));
  }

  const tiers: string[] = ['S', 'A', 'B', 'C', 'D'];

  function getTierColor(index: number): string {
    const hue = (index * 35) % 360;
    return `hsl(${hue}, 70%, 65%)`;
  }
</script>

<h1 class="text-3xl font-bold text-center my-6">Tierlist Creator</h1>

<!-- Tier rows -->
<div class="flex flex-col mx-4 divide-y-2 border-y-2">
  {#each tiers as label, i}
    <div class="grid grid-cols-[90px_1fr] h-24">
      <div
        class="flex items-center justify-center text-xl font-bold text-gray-800"
        style="background-color: {getTierColor(i)}"
      >
        {label}
      </div>
      <div class="bg-neutral-900 border-l border-black p-2 flex items-center gap-2 overflow-x-auto">
        <!-- Dropped items will go here -->
        <p class="text-gray-500 italic">Drop items here...</p>
      </div>
    </div>
  {/each}
</div>

<!-- Uploaded items box -->
<div class="mt-8 mx-4 border-2 border-dashed border-gray-400 rounded-lg p-4 mb-6 min-h-[120px] bg-gray-50">
  <h2 class="text-lg font-semibold mb-2">Uploaded Items</h2>
  <div class="flex flex-wrap gap-4">
    {#if uploadedImages.length === 0}
      <p class="text-gray-500 italic">No images uploaded yet.</p>
    {:else}
      {#each uploadedImages as image}
        <img src="{image}" alt="uploaded item" class="w-16 h-16 object-cover rounded" draggable="true" />
      {/each}
    {/if}
  </div>
</div>

<!-- Upload button -->
<div class="mt-8 text-center">
  <label class="inline-block bg-sky-500 hover:bg-sky-700 text-white font-bold py-2 px-4 rounded cursor-pointer">
    Upload Image
    <input type="file" accept="image/*" multiple class="hidden" on:change={handleUpload} />
  </label>
</div>
