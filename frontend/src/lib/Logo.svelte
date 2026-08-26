<script lang="ts">
  let { class: className = "w-8 h-8" }: { class?: string } = $props();
</script>

<img src="/logo.png" alt="NutriPlan logo" class="{className} object-contain" />
