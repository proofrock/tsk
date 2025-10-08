<script>

  let { tasks, onComplete, onEdit, onReorder } = $props();

  // Export collapse functions for parent component
  export { collapseAll, expandAll };

  // Helper to check if a task has subtasks
  function hasSubtasks(taskId) {
    return tasks.some(t => t.parent_id === taskId);
  }

  let draggedTask = $state(null);
  let dragOverIndex = $state(null);
  let dropPosition = $state(null); // 'before', 'after', 'child', or 'sibling-child'
  let dropBeforeOrAfter = $state('after'); // For sibling-child positioning
  let selectedTasks = $state(new Set());
  let touchStartY = $state(0);
  let touchStartIndex = $state(null);
  let listElement = $state(null);
  let collapsedParents = $state(new Set());

  function toggleCollapse(taskId, e) {
    e.stopPropagation();
    const newCollapsed = new Set(collapsedParents);
    if (newCollapsed.has(taskId)) {
      newCollapsed.delete(taskId);
    } else {
      newCollapsed.add(taskId);
    }
    collapsedParents = newCollapsed;
  }

  function collapseAll() {
    const parents = tasks.filter(t => hasSubtasks(t.id));
    collapsedParents = new Set(parents.map(t => t.id));
  }

  function expandAll() {
    collapsedParents = new Set();
  }

  function isCollapsed(taskId) {
    return collapsedParents.has(taskId);
  }

  function shouldShowTask(task, index) {
    if (!task.parent_id) return true;
    return !isCollapsed(task.parent_id);
  }

  // Mouse drag handlers
  function handleDragStart(task) {
    draggedTask = task;
  }

  function handleDragOver(e, index, task) {
    e.preventDefault();
    const rect = e.currentTarget.getBoundingClientRect();
    const midpoint = rect.top + rect.height / 2;
    const sixthWidth = rect.width / 6;
    const offsetX = e.clientX - rect.left;

    dragOverIndex = index;

    // Store before/after for later use
    dropBeforeOrAfter = e.clientY < midpoint ? 'before' : 'after';

    // Determine if we should make it a child or a sibling
    // If target is a parent task and dragging to the right (> 1/6 width), make it a child
    // But only if the dragged task doesn't have subtasks
    if (offsetX > sixthWidth && !task.parent_id && draggedTask && !draggedTask.parent_id && !hasSubtasks(draggedTask.id)) {
      dropPosition = 'child';
    }
    // If target is a subtask and dragging to the right (> 1/6 width), also make it a child of the same parent
    else if (offsetX > sixthWidth && task.parent_id && draggedTask && !hasSubtasks(draggedTask.id)) {
      dropPosition = 'sibling-child';
    }
    else {
      dropPosition = dropBeforeOrAfter;
    }
  }

  function handleDragEnd() {
    if (draggedTask && dragOverIndex !== null && dropPosition !== null) {
      const newTasks = [...tasks];
      const draggedIndex = newTasks.findIndex(t => t.id === draggedTask.id);
      const targetTask = newTasks[dragOverIndex];

      if (dropPosition === 'child') {
        // Make dragged task a child of target
        draggedTask.parent_id = targetTask.id;
        newTasks[draggedIndex] = {...draggedTask};

        // Reorder: remove from current position and place after target
        newTasks.splice(draggedIndex, 1);
        const newTargetIndex = newTasks.findIndex(t => t.id === targetTask.id);
        newTasks.splice(newTargetIndex + 1, 0, draggedTask);
      } else if (dropPosition === 'sibling-child') {
        // Make dragged task a sibling of target (same parent)
        draggedTask.parent_id = targetTask.parent_id;
        newTasks[draggedIndex] = {...draggedTask};

        // Calculate target index based on stored before/after
        let targetIndex = dragOverIndex;
        if (dropBeforeOrAfter === 'after') {
          targetIndex = dragOverIndex + 1;
        }

        // Adjust if dragging from before to after the target
        if (draggedIndex < targetIndex) {
          targetIndex--;
        }

        newTasks.splice(draggedIndex, 1);
        newTasks.splice(targetIndex, 0, draggedTask);
      } else {
        // If dropping in normal position and target is not a subtask, clear parent_id
        if (!targetTask.parent_id && draggedTask.parent_id) {
          draggedTask.parent_id = null;
          newTasks[draggedIndex] = {...draggedTask};
        }

        // Calculate target index based on drop position
        let targetIndex = dragOverIndex;
        if (dropPosition === 'after') {
          targetIndex = dragOverIndex + 1;
        }

        // Adjust if dragging from before to after the target
        if (draggedIndex < targetIndex) {
          targetIndex--;
        }

        if (draggedIndex !== targetIndex || !targetTask.parent_id) {
          newTasks.splice(draggedIndex, 1);
          newTasks.splice(targetIndex, 0, draggedTask);
        }
      }

      onReorder(newTasks);
    }

    draggedTask = null;
    dragOverIndex = null;
    dropPosition = null;
    dropBeforeOrAfter = 'after';
  }

  function handleDragLeave() {
    dragOverIndex = null;
    dropPosition = null;
  }

  // Touch drag handlers
  function handleTouchStart(e, task, index) {
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
        const taskId = parseInt(taskElement.dataset.taskId);
        const task = tasks.find(t => t.id === taskId);
        const newIndex = Array.from(taskElement.parentElement.children).indexOf(taskElement.parentElement);

        if (newIndex !== -1 && task) {
          const rect = taskElement.getBoundingClientRect();
          const midpoint = rect.top + rect.height / 2;
          const quarterWidth = rect.width * 0.25;
          const offsetX = touch.clientX - rect.left;

          dragOverIndex = tasks.findIndex(t => t.id === taskId);

          const sixthWidth = rect.width / 6;

          // If dragging more than 1/6 width to the right and target has no parent, make it a child
          // But only if the dragged task doesn't have subtasks
          if (offsetX > sixthWidth && !task.parent_id && draggedTask && !draggedTask.parent_id && !hasSubtasks(draggedTask.id)) {
            dropPosition = 'child';
          }
          // If target is a subtask and dragging to the right (> 1/6 width), make it a sibling of the same parent
          else if (offsetX > sixthWidth && task.parent_id && draggedTask && !hasSubtasks(draggedTask.id)) {
            dropPosition = 'sibling-child';
          }
          else {
            dropPosition = touch.clientY < midpoint ? 'before' : 'after';
          }
        }
      }
    }
  }

  function handleTouchEnd() {
    if (draggedTask && dragOverIndex !== null && dropPosition !== null && touchStartIndex !== null) {
      const newTasks = [...tasks];
      const draggedIndex = touchStartIndex;
      const targetTask = newTasks[dragOverIndex];

      if (dropPosition === 'child') {
        // Make dragged task a child of target
        draggedTask.parent_id = targetTask.id;
        newTasks[draggedIndex] = {...draggedTask};

        // Reorder: remove from current position and place after target
        newTasks.splice(draggedIndex, 1);
        const newTargetIndex = newTasks.findIndex(t => t.id === targetTask.id);
        newTasks.splice(newTargetIndex + 1, 0, draggedTask);
      } else if (dropPosition === 'sibling-child') {
        // Make dragged task a sibling of target (same parent)
        draggedTask.parent_id = targetTask.parent_id;
        newTasks[draggedIndex] = {...draggedTask};

        // Calculate target index based on position (before/after)
        let targetIndex = dragOverIndex;
        // For touch, just use after for now (can enhance later)
        targetIndex = dragOverIndex + 1;

        // Adjust if dragging from before to after the target
        if (draggedIndex < targetIndex) {
          targetIndex--;
        }

        if (draggedIndex !== targetIndex) {
          newTasks.splice(draggedIndex, 1);
          newTasks.splice(targetIndex, 0, draggedTask);
        }
      } else {
        // If dropping in normal position and target is not a subtask, clear parent_id
        if (!targetTask.parent_id && draggedTask.parent_id) {
          draggedTask.parent_id = null;
          newTasks[draggedIndex] = {...draggedTask};
        }

        // Calculate target index based on drop position
        let targetIndex = dragOverIndex;
        if (dropPosition === 'after') {
          targetIndex = dragOverIndex + 1;
        }

        // Adjust if dragging from before to after the target
        if (draggedIndex < targetIndex) {
          targetIndex--;
        }

        if (draggedIndex !== targetIndex || !targetTask.parent_id) {
          newTasks.splice(draggedIndex, 1);
          newTasks.splice(targetIndex, 0, draggedTask);
        }
      }

      onReorder(newTasks);
    }

    draggedTask = null;
    dragOverIndex = null;
    dropPosition = null;
    dropBeforeOrAfter = 'after';
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

  function handleRowClick(e, task) {
    // Don't do anything if clicking on checkbox or edit button
    if (e.target.closest('.form-check') || e.target.closest('.edit-btn') || e.target.closest('.collapse-btn')) {
      return;
    }

    // Toggle collapse if parent task
    if (hasSubtasks(task.id)) {
      toggleCollapse(task.id, e);
    }
    // Non-parent tasks: do nothing (edit is handled by edit button)
  }

  $effect(() => {
    if (!listElement || tasks.length === 0) return;

    const cleanup = [];
    const taskCards = listElement.querySelectorAll('.task-card');

    taskCards.forEach((card, index) => {
      const taskId = parseInt(card.dataset.taskId);
      const task = tasks.find(t => t.id === taskId);
      if (!task) return;

      const touchStartHandler = (e) => {
        // Don't start drag if touching checkbox, edit button, or collapse button
        if (e.target.closest('.form-check') || e.target.closest('.edit-btn') || e.target.closest('.collapse-btn')) {
          return;
        }
        handleTouchStart(e, task, index);
      };
      const touchMoveHandler = (e) => handleTouchMove(e, index);
      const touchEndHandler = handleTouchEnd;

      card.addEventListener('touchstart', touchStartHandler, { passive: false });
      card.addEventListener('touchmove', touchMoveHandler, { passive: false });
      card.addEventListener('touchend', touchEndHandler, { passive: false });

      cleanup.push(() => {
        card.removeEventListener('touchstart', touchStartHandler);
        card.removeEventListener('touchmove', touchMoveHandler);
        card.removeEventListener('touchend', touchEndHandler);
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
      {#if shouldShowTask(task, index)}
      <div class="task-wrapper"
           class:show-drop-line-before={dragOverIndex === index && dropPosition === 'before'}
           class:show-drop-line-after={dragOverIndex === index && dropPosition === 'after'}
           class:show-drop-line-child={dragOverIndex === index && dropPosition === 'child'}
           class:show-drop-line-sibling={dragOverIndex === index && dropPosition === 'sibling-child'}>
        <div
          class="list-group-item list-group-item-action bg-dark text-light border-secondary p-3 mb-2 draggable-item task-card"
          class:opacity-50={draggedTask?.id === task.id}
          class:border-danger={selectedTasks.has(task.id)}
          class:subtask={task.parent_id}
          data-task-id={task.id}
          draggable="true"
          ondragstart={() => handleDragStart(task)}
          ondragover={(e) => handleDragOver(e, index, task)}
          ondragend={handleDragEnd}
          ondragleave={handleDragLeave}
          ondrop={(e) => e.preventDefault()}
          onclick={(e) => handleRowClick(e, task)}
          onkeydown={(e) => (e.key === 'Enter' || e.key === ' ') && handleRowClick(e, task)}
          style="cursor: {hasSubtasks(task.id) ? 'pointer' : 'grab'};"
          role="button"
          tabindex="0"
        >
        <div class="d-flex align-items-center gap-2 gap-md-3">
          {#if hasSubtasks(task.id)}
            <button
              class="btn btn-link text-secondary p-0 border-0 collapse-btn"
              style="flex-shrink: 0; width: 20px; height: 20px;"
              onclick={(e) => toggleCollapse(task.id, e)}
              aria-label={isCollapsed(task.id) ? 'Expand subtasks' : 'Collapse subtasks'}
            >
              {#if isCollapsed(task.id)}
                <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
                  <path d="M6 4l4 4-4 4V4z"/>
                </svg>
              {:else}
                <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
                  <path d="M4 6l4 4 4-4H4z"/>
                </svg>
              {/if}
            </button>
          {:else if !task.parent_id}
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
          {:else}
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
          {/if}

          <div class="flex-grow-1 min-w-0 task-content">
            <h6 class="mb-1 text-truncate">{task.title}</h6>
            {#if task.description}
              <p class="mb-0 text-secondary small text-truncate-2">{task.description}</p>
            {/if}
          </div>

          <button
            class="btn btn-link text-secondary p-0 border-0 edit-btn"
            style="flex-shrink: 0;"
            onclick={(e) => { e.stopPropagation(); onEdit(task); }}
            aria-label="Edit task"
            title="Edit task"
          >
            <svg width="18" height="18" viewBox="0 0 16 16" fill="currentColor">
              <path d="M12.146.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1 0 .708l-10 10a.5.5 0 0 1-.168.11l-5 2a.5.5 0 0 1-.65-.65l2-5a.5.5 0 0 1 .11-.168l10-10zM11.207 2.5 13.5 4.793 14.793 3.5 12.5 1.207 11.207 2.5zm1.586 3L10.5 3.207 4 9.707V10h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.293l6.5-6.5zm-9.761 5.175-.106.106-1.528 3.821 3.821-1.528.106-.106A.5.5 0 0 1 5 12.5V12h-.5a.5.5 0 0 1-.5-.5V11h-.5a.5.5 0 0 1-.468-.325z"/>
            </svg>
          </button>
        </div>
        </div>
      </div>
      {/if}
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

  .task-wrapper.show-drop-line-child::after {
    content: '';
    position: absolute;
    bottom: -4px;
    left: 16.67%;
    right: 0;
    height: 3px;
    background-color: #f97316;
    border-radius: 2px;
    z-index: 10;
    box-shadow: 0 0 8px rgba(249, 115, 22, 0.6);
  }

  .task-wrapper.show-drop-line-sibling::after {
    content: '';
    position: absolute;
    bottom: -4px;
    left: 16.67%;
    right: 0;
    height: 3px;
    background-color: #22c55e;
    border-radius: 2px;
    z-index: 10;
    box-shadow: 0 0 8px rgba(34, 197, 94, 0.6);
  }

  .task-card {
    transition: all 0.25s ease;
  }

  .task-wrapper:has(.subtask) {
    padding-left: 3rem;
  }

  .task-card:hover {
    background-color: #1c1c1f !important;
    transform: translateY(-1px);
    box-shadow: 0 2px 8px rgba(249, 115, 22, 0.1);
  }

  .task-card.border-danger {
    border-width: 2px !important;
  }

  .task-content {
    pointer-events: none;
  }

  .collapse-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    transition: opacity 0.2s ease;
  }

  .collapse-btn:hover {
    opacity: 0.7;
  }

  .edit-btn {
    transition: opacity 0.2s ease;
  }

  .edit-btn:hover {
    opacity: 0.7;
  }
</style>
