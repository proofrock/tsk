<script>

  let { tasks, onComplete, onEdit, onReorder } = $props();

  let draggedTask = $state(null);
  let dragOverIndex = $state(null);
  let dropPosition = $state(null); // 'before' or 'after'
  let selectedTasks = $state(new Set());
  let touchStartY = $state(0);
  let touchStartIndex = $state(null);
  let listElement = $state(null);

  // Mouse drag handlers
  function handleDragStart(task) {
    draggedTask = task;
  }

  function handleDragOver(e, index) {
    e.preventDefault();
    const rect = e.currentTarget.getBoundingClientRect();
    const midpoint = rect.top + rect.height / 2;

    dragOverIndex = index;
    dropPosition = e.clientY < midpoint ? 'before' : 'after';
  }

  function handleDragEnd() {
    if (draggedTask && dragOverIndex !== null && dropPosition !== null) {
      const newTasks = [...tasks];
      const draggedIndex = newTasks.findIndex(t => t.id === draggedTask.id);

      // Calculate target index based on drop position
      let targetIndex = dragOverIndex;
      if (dropPosition === 'after') {
        targetIndex = dragOverIndex + 1;
      }

      // Adjust if dragging from before to after the target
      if (draggedIndex < targetIndex) {
        targetIndex--;
      }

      if (draggedIndex !== targetIndex) {
        newTasks.splice(draggedIndex, 1);
        newTasks.splice(targetIndex, 0, draggedTask);

        const taskIds = newTasks.map(t => t.id);
        onReorder(taskIds);
      }
    }

    draggedTask = null;
    dragOverIndex = null;
    dropPosition = null;
  }

  function handleDragLeave() {
    dragOverIndex = null;
    dropPosition = null;
  }

  // Touch drag handlers
  function handleTouchStart(e, task, index) {
    if (e.target.closest('.form-check') || e.target.closest('.task-content')) {
      return;
    }
    touchStartY = e.touches[0].clientY;
    touchStartIndex = index;
    draggedTask = task;
  }

  function handleTouchMove(e, index) {
    if (draggedTask && touchStartIndex !== null) {
      e.preventDefault();
      const touch = e.touches[0];
      const element = document.elementFromPoint(touch.clientX, touch.clientY);
      const taskElement = element?.closest('.task-card');

      if (taskElement) {
        const newIndex = Array.from(taskElement.parentElement.children).indexOf(taskElement);
        if (newIndex !== -1) {
          const rect = taskElement.getBoundingClientRect();
          const midpoint = rect.top + rect.height / 2;

          dragOverIndex = newIndex;
          dropPosition = touch.clientY < midpoint ? 'before' : 'after';
        }
      }
    }
  }

  function handleTouchEnd() {
    if (draggedTask && dragOverIndex !== null && dropPosition !== null && touchStartIndex !== null) {
      const newTasks = [...tasks];
      const draggedIndex = touchStartIndex;

      // Calculate target index based on drop position
      let targetIndex = dragOverIndex;
      if (dropPosition === 'after') {
        targetIndex = dragOverIndex + 1;
      }

      // Adjust if dragging from before to after the target
      if (draggedIndex < targetIndex) {
        targetIndex--;
      }

      if (draggedIndex !== targetIndex) {
        newTasks.splice(draggedIndex, 1);
        newTasks.splice(targetIndex, 0, draggedTask);

        const taskIds = newTasks.map(t => t.id);
        onReorder(taskIds);
      }
    }

    draggedTask = null;
    dragOverIndex = null;
    dropPosition = null;
    touchStartIndex = null;
  }

  function toggleTaskSelection(taskId) {
    const newSelection = new Set(selectedTasks);
    if (newSelection.has(taskId)) {
      newSelection.delete(taskId);
    } else {
      newSelection.add(taskId);
    }
    selectedTasks = newSelection;
  }

  function deleteSelectedTasks() {
    selectedTasks.forEach(taskId => {
      onComplete(taskId);
    });
    selectedTasks = new Set();
  }

  function handleEditClick(e, task) {
    if (e.target.closest('.form-check') || e.target.closest('.drag-handle')) {
      return;
    }
    onEdit(task);
  }

  $effect(() => {
    if (!listElement || tasks.length === 0) return;

    const cleanup = [];
    const taskCards = listElement.querySelectorAll('.task-card');

    taskCards.forEach((card, index) => {
      const dragHandle = card.querySelector('.drag-handle');
      if (!dragHandle) return;

      const taskId = parseInt(card.dataset.taskId);
      const task = tasks.find(t => t.id === taskId);
      if (!task) return;

      const touchStartHandler = (e) => handleTouchStart(e, task, index);
      const touchMoveHandler = (e) => handleTouchMove(e, index);
      const touchEndHandler = handleTouchEnd;

      dragHandle.addEventListener('touchstart', touchStartHandler, { passive: false });
      dragHandle.addEventListener('touchmove', touchMoveHandler, { passive: false });
      dragHandle.addEventListener('touchend', touchEndHandler, { passive: false });

      cleanup.push(() => {
        dragHandle.removeEventListener('touchstart', touchStartHandler);
        dragHandle.removeEventListener('touchmove', touchMoveHandler);
        dragHandle.removeEventListener('touchend', touchEndHandler);
      });
    });

    return () => {
      cleanup.forEach(fn => fn());
    };
  });
</script>

{#if selectedTasks.size > 0}
  <div class="mb-3">
    <button class="btn btn-danger" onclick={deleteSelectedTasks}>
      Delete {selectedTasks.size} task{selectedTasks.size > 1 ? 's' : ''}
    </button>
  </div>
{/if}

{#if tasks.length === 0}
  <div class="alert alert-secondary text-center" role="alert">
    No tasks yet. Click "Add Task" to get started.
  </div>
{:else}
  <div class="list-group" bind:this={listElement}>
    {#each tasks as task, index (task.id)}
      <div class="task-wrapper" class:show-drop-line-before={dragOverIndex === index && dropPosition === 'before'} class:show-drop-line-after={dragOverIndex === index && dropPosition === 'after'}>
        <div
          class="list-group-item list-group-item-action bg-dark text-light border-secondary p-3 mb-2 draggable-item task-card"
          class:opacity-50={draggedTask?.id === task.id}
          class:border-danger={selectedTasks.has(task.id)}
          data-task-id={task.id}
          draggable="true"
          ondragstart={() => handleDragStart(task)}
          ondragover={(e) => handleDragOver(e, index)}
          ondragend={handleDragEnd}
          ondragleave={handleDragLeave}
          ondrop={(e) => e.preventDefault()}
          onclick={(e) => handleEditClick(e, task)}
          role="button"
          tabindex="0"
          onkeydown={(e) => (e.key === 'Enter' || e.key === ' ') && !e.target.closest('.form-check') && onEdit(task)}
        >
        <div class="d-flex align-items-center gap-2 gap-md-3">
          <div class="drag-handle text-secondary" style="cursor: move; flex-shrink: 0; touch-action: none;">
            <svg width="20" height="20" viewBox="0 0 16 16" fill="currentColor">
              <circle cx="6" cy="4" r="1.5"/>
              <circle cx="10" cy="4" r="1.5"/>
              <circle cx="6" cy="8" r="1.5"/>
              <circle cx="10" cy="8" r="1.5"/>
              <circle cx="6" cy="12" r="1.5"/>
              <circle cx="10" cy="12" r="1.5"/>
            </svg>
          </div>

          <div class="form-check" style="flex-shrink: 0;">
            <input
              class="form-check-input"
              type="checkbox"
              checked={selectedTasks.has(task.id)}
              onchange={() => toggleTaskSelection(task.id)}
              id="task-{task.id}"
            />
            <label class="form-check-label visually-hidden" for="task-{task.id}">
              Select task: {task.title}
            </label>
          </div>

          <div class="flex-grow-1 min-w-0 task-content">
            <h6 class="mb-1 text-truncate">{task.title}</h6>
            {#if task.description}
              <p class="mb-0 text-secondary small text-truncate-2">{task.description}</p>
            {/if}
          </div>
        </div>
        </div>
      </div>
    {/each}
  </div>
{/if}

<style>
  .text-truncate-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.4;
  }

  .min-w-0 {
    min-width: 0;
  }

  .form-check-input:checked {
    background-color: #f97316;
    border-color: #f97316;
  }

  .task-wrapper {
    position: relative;
  }

  .task-wrapper.show-drop-line-before::before {
    content: '';
    position: absolute;
    top: -4px;
    left: 0;
    right: 0;
    height: 3px;
    background-color: #f97316;
    border-radius: 2px;
    z-index: 10;
    box-shadow: 0 0 8px rgba(249, 115, 22, 0.6);
  }

  .task-wrapper.show-drop-line-after::after {
    content: '';
    position: absolute;
    bottom: -4px;
    left: 0;
    right: 0;
    height: 3px;
    background-color: #f97316;
    border-radius: 2px;
    z-index: 10;
    box-shadow: 0 0 8px rgba(249, 115, 22, 0.6);
  }

  .task-card {
    cursor: pointer;
    transition: all 0.25s ease;
  }

  .task-card:hover {
    background-color: #1c1c1f !important;
    transform: translateY(-1px);
    box-shadow: 0 2px 8px rgba(249, 115, 22, 0.1);
  }

  .task-card:hover .drag-handle {
    opacity: 1;
  }

  .task-card.border-danger {
    border-width: 2px !important;
  }

  .drag-handle {
    opacity: 0.3;
    transition: opacity 0.2s ease;
  }

  .task-content {
    cursor: pointer;
  }

  @media (hover: none) {
    .drag-handle {
      opacity: 1;
    }
  }
</style>
