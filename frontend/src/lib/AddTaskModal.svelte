<script>
  import { onMount } from 'svelte';

  let { categories, defaultCategory, tasks, task, onSave, onClose } = $props();

  let title = $state(task?.title || '');
  let description = $state(task?.description || '');
  let categoryId = $state(task?.category_id || defaultCategory?.id || 1);
  let parentId = $state(task?.parent_id || null);
  let titleInput;

  // If task is a subtask, it can't change category
  let isSubtask = $derived(task ? !!task.parent_id : false);

  // Update category when parent changes (for new tasks with parent)
  $effect(() => {
    if (parentId) {
      const parent = tasks.find(t => t.id === parentId);
      if (parent) {
        categoryId = parent.category_id;
      }
    }
  });

  // Helper to check if a task has subtasks
  function hasSubtasks(taskId) {
    return tasks.some(t => t.parent_id === taskId);
  }

  // Check if current task being edited has subtasks (can't become a subtask itself)
  let taskHasSubtasks = $derived(task ? hasSubtasks(task.id) : false);

  // Get available parent tasks (top-level tasks only, excluding the task being edited)
  let availableParents = $derived(
    tasks.filter(t =>
      !t.parent_id && (!task || t.id !== task.id) && t.category_id === categoryId
    )
  );

  onMount(() => {
    if (titleInput) {
      titleInput.focus();
      titleInput.select();
    }
  });

  function handleSubmit(e) {
    e.preventDefault();
    if (!title.trim()) return;

    if (task) {
      onSave({
        id: task.id,
        title: title.trim(),
        description: description.trim(),
        category_id: categoryId,
        parent_id: parentId || null,
      });
    } else {
      onSave({
        title: title.trim(),
        description: description.trim(),
        category_id: categoryId,
        parent_id: parentId || null,
      });
    }
  }

  function handleBackdropClick(e) {
    if (e.target === e.currentTarget) {
      onClose();
    }
  }
</script>

<div
  class="modal d-block"
  tabindex="-1"
  role="dialog"
  aria-modal="true"
  onclick={handleBackdropClick}
  onkeydown={(e) => e.key === 'Escape' && onClose()}
>
  <div class="modal-dialog modal-dialog-centered modal-dialog-scrollable">
    <div class="modal-content bg-dark text-light border-secondary">
      <div class="modal-header border-secondary">
        <h5 class="modal-title">{task ? 'Edit Task' : 'New Task'}</h5>
        <button
          type="button"
          class="btn-close btn-close-white"
          aria-label="Close"
          onclick={onClose}
        ></button>
      </div>

      <form onsubmit={handleSubmit}>
        <div class="modal-body">
          <div class="mb-3">
            <label for="title" class="form-label">Title</label>
            <input
              type="text"
              class="form-control bg-dark text-light border-secondary"
              id="title"
              bind:value={title}
              bind:this={titleInput}
              placeholder="Enter task title"
              required
            />
          </div>

          <div class="mb-3">
            <label for="description" class="form-label">Description</label>
            <textarea
              class="form-control bg-dark text-light border-secondary"
              id="description"
              bind:value={description}
              placeholder="Enter task description (optional)"
              rows="4"
            ></textarea>
          </div>

          <div class="mb-3">
            <label for="category" class="form-label">Category</label>
            <select
              id="category"
              class="form-select bg-dark text-light border-secondary"
              bind:value={categoryId}
              disabled={isSubtask || !!parentId}
            >
              {#each categories as category}
                <option value={category.id}>{category.name}</option>
              {/each}
            </select>
            {#if isSubtask || parentId}
              <small class="text-muted d-block mt-1">Subtasks inherit their parent's category.</small>
            {/if}
          </div>

          <div class="mb-3">
            <label for="parent" class="form-label">Parent Task (Optional)</label>
            <select
              id="parent"
              class="form-select bg-dark text-light border-secondary"
              bind:value={parentId}
              disabled={taskHasSubtasks}
            >
              <option value={null}>(None)</option>
              {#each availableParents as parent}
                <option value={parent.id}>{parent.title}</option>
              {/each}
            </select>
            {#if taskHasSubtasks}
              <small class="text-muted d-block mt-1">This task has subtasks and cannot be made a subtask itself.</small>
            {/if}
          </div>
        </div>

        <div class="modal-footer border-secondary">
          <button type="button" class="btn btn-secondary" onclick={onClose}>
            Cancel
          </button>
          <button type="submit" class="btn btn-warning" style="background-color: #f97316; border-color: #f97316;">
            {task ? 'Save Changes' : 'Create Task'}
          </button>
        </div>
      </form>
    </div>
  </div>
</div>

<div class="modal-backdrop show"></div>

<style>
  .form-control:focus,
  .form-select:focus {
    border-color: #f97316;
    box-shadow: 0 0 0 0.25rem rgba(249, 115, 22, 0.25);
  }

  .modal-dialog {
    max-width: 500px;
  }

  @media (max-width: 575.98px) {
    .modal-dialog {
      margin: 0.5rem;
    }
  }
</style>
