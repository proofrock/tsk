<script>
  let { tasks, onComplete, onEdit, onReorder } = $props();

  let draggedTask = $state(null);
  let dragOverIndex = $state(null);

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
</script>

{#if tasks.length === 0}
  <div class="alert alert-secondary text-center" role="alert">
    No tasks yet. Click "Add Task" to get started.
  </div>
{:else}
  <div class="list-group">
    {#each tasks as task, index (task.id)}
      <div
        class="list-group-item list-group-item-action bg-dark text-light border-secondary p-3 mb-2 draggable-item"
        class:opacity-50={draggedTask?.id === task.id}
        class:border-warning={dragOverIndex === index}
        draggable="true"
        ondragstart={() => handleDragStart(task)}
        ondragover={(e) => handleDragOver(e, index)}
        ondragend={handleDragEnd}
        ondragleave={handleDragLeave}
        ondrop={(e) => e.preventDefault()}
        role="listitem"
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
              onchange={() => onComplete(task.id)}
              id="task-{task.id}"
            />
          </div>

          <div class="flex-grow-1 min-w-0" role="button" tabindex="0" onclick={() => onEdit(task)} onkeydown={(e) => (e.key === 'Enter' || e.key === ' ') && onEdit(task)}>
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

  .draggable-item {
    cursor: move;
    transition: all 0.2s;
  }

  .draggable-item:hover {
    background-color: #1c1c1f !important;
  }

  .draggable-item.border-warning {
    border-width: 2px !important;
  }
</style>
