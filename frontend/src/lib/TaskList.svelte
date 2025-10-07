<script>
  let { tasks, onComplete, onEdit, onReorder } = $props();

  let draggedTask = $state(null);
  let dragOverIndex = $state(null);
  let completingTaskId = $state(null);

  function handleDragStart(task) {
    draggedTask = task;
  }

  function handleDragOver(e, index) {
    e.preventDefault();
    dragOverIndex = index;
  }

  function handleDragEnd() {
    if (draggedTask && dragOverIndex !== null) {
      const newTasks = [...tasks];
      const draggedIndex = newTasks.findIndex(t => t.id === draggedTask.id);

      if (draggedIndex !== dragOverIndex) {
        newTasks.splice(draggedIndex, 1);
        newTasks.splice(dragOverIndex, 0, draggedTask);

        const taskIds = newTasks.map(t => t.id);
        onReorder(taskIds);
      }
    }

    draggedTask = null;
    dragOverIndex = null;
  }

  function handleDragLeave() {
    dragOverIndex = null;
  }

  function handleComplete(taskId) {
    completingTaskId = taskId;
    setTimeout(() => {
      onComplete(taskId);
      completingTaskId = null;
    }, 300);
  }

  function handleEditClick(e, task) {
    // Don't edit if clicking on checkbox or drag handle
    if (e.target.closest('.form-check') || e.target.closest('.drag-handle')) {
      return;
    }
    onEdit(task);
  }
</script>

{#if tasks.length === 0}
  <div class="alert alert-secondary text-center" role="alert">
    No tasks yet. Click "Add Task" to get started.
  </div>
{:else}
  <div class="list-group">
    {#each tasks as task, index (task.id)}
      <div
        class="list-group-item list-group-item-action bg-dark text-light border-secondary p-3 mb-2 draggable-item task-card"
        class:opacity-50={draggedTask?.id === task.id}
        class:border-warning={dragOverIndex === index}
        class:completing={completingTaskId === task.id}
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
          <div class="drag-handle text-secondary" style="cursor: move; flex-shrink: 0;">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
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
              checked={task.completed}
              onchange={() => handleComplete(task.id)}
              id="task-{task.id}"
            />
            <label class="form-check-label visually-hidden" for="task-{task.id}">
              Complete task: {task.title}
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

  .task-card.border-warning {
    border-width: 2px !important;
    transform: translateY(-2px);
  }

  .task-card.completing {
    animation: fadeOut 0.3s ease-out forwards;
  }

  @keyframes fadeOut {
    from {
      opacity: 1;
      transform: scale(1);
    }
    to {
      opacity: 0;
      transform: scale(0.95);
    }
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
