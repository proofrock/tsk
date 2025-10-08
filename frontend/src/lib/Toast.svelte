<!--
  Copyright (C) 2025 Germano Rizzo <oss /AT/ germanorizzo /DOT/ it>
  Licensed under the EUPL v. 1.2
  See LICENSE file for full license text
-->

<script>
  let { message, show, type = 'success' } = $props();
</script>

{#if show}
  <div class="toast-container position-fixed bottom-0 end-0 p-3" style="z-index: 1100;">
    <div class="toast show bg-{type} text-white border-0" role="alert">
      <div class="toast-body d-flex align-items-center gap-2">
        {#if type === 'success'}
          <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
          </svg>
        {/if}
        <span>{message}</span>
      </div>
    </div>
  </div>
{/if}

<style>
  .bg-success {
    background-color: #f97316 !important;
  }

  .toast {
    animation: slideIn 0.3s ease-out;
  }

  @keyframes slideIn {
    from {
      transform: translateX(100%);
      opacity: 0;
    }
    to {
      transform: translateX(0);
      opacity: 1;
    }
  }
</style>
