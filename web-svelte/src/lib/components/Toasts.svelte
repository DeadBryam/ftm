<script>
  import { useToast } from "$lib/stores/toast.svelte";

  const toastStore = useToast();
</script>

<div class="toasts-container">
  {#each toastStore.toasts as toast (toast.id)}
    <div class="toast toast-{toast.type}" role="alert">
      <span class="toast-message">{toast.message}</span>
      <button
        class="toast-close"
        onclick={() => toastStore.remove(toast.id)}
        aria-label="Close notification"
      >
        &times;
      </button>
    </div>
  {/each}
</div>

<style>
  .toasts-container {
    position: fixed;
    top: 20px;
    right: 20px;
    z-index: 1000;
    display: flex;
    flex-direction: column;
    gap: 10px;
    pointer-events: none;
  }

  .toast {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 14px 18px;
    border-radius: 10px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.15);
    animation: slideIn 0.3s cubic-bezier(0.16, 1, 0.3, 1);
    pointer-events: auto;
    min-width: 280px;
    max-width: 400px;
  }

  @keyframes slideIn {
    from {
      opacity: 0;
      transform: translateX(100px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  .toast-success {
    background: linear-gradient(135deg, var(--btn-start-bg) 0%, var(--btn-start-hover-bg) 100%);
    border-color: var(--btn-start-bg);
    color: var(--badge-text);
  }

  .toast-error {
    background: linear-gradient(135deg, var(--btn-stop-bg) 0%, var(--btn-stop-hover-bg) 100%);
    border-color: var(--btn-stop-bg);
    color: var(--badge-text);
  }

  .toast-info {
    background: linear-gradient(135deg, var(--btn-primary-bg) 0%, var(--btn-primary-hover-bg) 100%);
    border-color: var(--btn-primary-bg);
    color: var(--badge-text);
  }

  .toast-message {
    flex: 1;
    font-size: 14px;
    font-weight: 500;
  }

  .toast-close {
    background: color-mix(in srgb, var(--badge-text) 20%, transparent);
    border: none;
    color: var(--badge-text);
    width: 28px;
    height: 28px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 18px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background 0.15s;
  }

  .toast-close:hover {
    background: color-mix(in srgb, var(--badge-text) 30%, transparent);
  }

  @media (max-width: 640px) {
    .toasts-container {
      top: 10px;
      right: 10px;
      left: 10px;
    }

    .toast {
      min-width: auto;
      max-width: none;
      width: 100%;
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .toast {
      animation: none;
    }
  }
</style>
