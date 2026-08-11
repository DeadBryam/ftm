<script lang="ts">
  import QRCode from "qrcode";
  import { cn } from "$lib/utils/cn";

  let {
    value,
    size = 148,
    dataUrl = $bindable(""),
    class: className = "",
  }: {
    value: string;
    size?: number;
    dataUrl?: string;
    class?: string;
  } = $props();

  $effect(() => {
    if (!value) {
      dataUrl = "";
      return;
    }

    let cancelled = false;
    QRCode.toDataURL(value, {
      width: size * 2,
      margin: 1,
      color: { dark: "#000000", light: "#ffffff" },
    })
      .then((url) => {
        if (!cancelled) dataUrl = url;
      })
      .catch(() => {
        if (!cancelled) dataUrl = "";
      });

    return () => {
      cancelled = true;
    };
  });
</script>

{#if dataUrl}
  <img
    src={dataUrl}
    alt={value}
    width={size}
    height={size}
    class={cn("rounded-control bg-white p-1", className)}
  />
{/if}
