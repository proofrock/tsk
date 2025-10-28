<!--
  Copyright (C) 2025 Germano Rizzo <oss /AT/ germanorizzo /DOT/ it>
  Licensed under the EUPL v. 1.2
  See LICENSE file for full license text
-->

<script>
  import { onMount } from 'svelte';
  import TaskList from './lib/TaskList.svelte';
  import AddTaskModal from './lib/AddTaskModal.svelte';
  import Toast from './lib/Toast.svelte';

  let categories = $state([]);
  let tasks = $state([]);
  let selectedCategory = $state(null);
  let showAddModal = $state(false);
  let editingTask = $state(null);
  let toastMessage = $state('');
  let showToast = $state(false);
  let version = $state('');
  let taskCounts = $state({});
  let taskListComponent = $state(null);

  const API_BASE = '/api';

  onMount(async () => {
    await loadCategories();
    await loadVersion();
  });

  async function loadVersion() {
    try {
      const res = await fetch(`${API_BASE}/version`);
      const data = await res.json();
      version = data.version;
    } catch (err) {
      console.error('Failed to load version:', err);
    }
  }

  async function loadCategories() {
    const res = await fetch(`${API_BASE}/categories`);
    categories = await res.json();
    await loadTaskCounts();
    if (categories.length > 0 && !selectedCategory) {
      // Select default category or fall back to first
      selectedCategory = categories.find(c => c.is_default) || categories[0];
      await loadTasks();
    }
  }

  async function loadTasks() {
    if (!selectedCategory) return;
    const res = await fetch(`${API_BASE}/tasks?category_id=${selectedCategory.id}`);
    const allTasks = await res.json();

    // Organize tasks: parents first with their subtasks following immediately after
    const organized = [];
    const parents = allTasks.filter(t => !t.parent_id);
    const children = allTasks.filter(t => t.parent_id);

    parents.forEach(parent => {
      organized.push(parent);
      const subtasks = children.filter(c => c.parent_id === parent.id);
      organized.push(...subtasks);
    });

    tasks = organized;
  }

  async function loadTaskCounts() {
    const counts = {};
    for (const category of categories) {
      const res = await fetch(`${API_BASE}/tasks?category_id=${category.id}`);
      const categoryTasks = await res.json();
      counts[category.id] = categoryTasks.length;
    }
    taskCounts = counts;
  }

  async function handleCategoryChange(e) {
    const categoryId = parseInt(e.target.value);
    selectedCategory = categories.find(c => c.id === categoryId);
    await loadTasks();
  }

  async function handleAddTask(task) {
    const res = await fetch(`${API_BASE}/tasks`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(task),
    });
    if (res.ok) {
      await loadTasks();
      await loadTaskCounts();
      showAddModal = false;
      showToastMessage('Task created!');
    }
  }

  async function handleEditTask(task) {
    const res = await fetch(`${API_BASE}/tasks/${task.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        title: task.title,
        description: task.description,
        category_id: task.category_id,
        parent_id: task.parent_id,
      }),
    });
    if (res.ok) {
      await loadTasks();
      await loadTaskCounts();
      showAddModal = false;
      editingTask = null;
    }
  }

  async function handleCompleteTask(taskId) {
    const res = await fetch(`${API_BASE}/tasks/${taskId}/complete`, {
      method: 'POST',
    });
    if (res.ok) {
      await loadTasks();
      await loadTaskCounts();
      showToastMessage('Task deleted!');
    }
  }

  function showToastMessage(message) {
    toastMessage = message;
    showToast = true;
    setTimeout(() => {
      showToast = false;
    }, 3000);
  }

  async function handleReorderTasks(reorderedTasks) {
    const tasksData = reorderedTasks.map(t => ({
      id: t.id,
      parent_id: t.parent_id || null
    }));

    const res = await fetch(`${API_BASE}/tasks/reorder`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ tasks: tasksData }),
    });
    if (res.ok) {
      await loadTasks();
    }
  }

  function openEditModal(task) {
    editingTask = task;
    showAddModal = true;
  }

  function openAddModal() {
    editingTask = null;
    showAddModal = true;
  }

  function closeModal() {
    showAddModal = false;
    editingTask = null;
  }
</script>

<div class="container-fluid py-3 py-md-4">
  <div class="row">
    <div class="col-12 col-lg-10 col-xl-8 mx-auto">
      <header class="mb-3 mb-md-4">
        <div class="d-flex align-items-baseline gap-2 mb-3">
          <h1 class="gradient-text fw-bold mb-0">tsk</h1>
          {#if version}
            <small class="text-secondary">{version}</small>
          {/if}
        </div>

        <div class="row g-2 align-items-center">
          <div class="col-12 col-sm-6 col-md-4">
            <select
              class="form-select bg-dark text-light border-secondary category-select"
              value={selectedCategory?.id || ''}
              onchange={handleCategoryChange}
            >
              {#each categories as category}
                <option value={category.id}>
                  {category.name} ({taskCounts[category.id] || 0})
                </option>
              {/each}
            </select>
          </div>

          <div class="col-auto ms-md-auto">
            <div class="btn-group" role="group">
              <button
                class="btn btn-sm btn-outline-secondary"
                onclick={() => taskListComponent?.expandAll()}
                title="Expand all"
                aria-label="Expand all tasks"
              >
                <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
                  <path d="M1 8a.5.5 0 0 1 .5-.5h13a.5.5 0 0 1 0 1h-13A.5.5 0 0 1 1 8zM7.646.146a.5.5 0 0 1 .708 0l2 2a.5.5 0 0 1-.708.708L8.5 1.707V5.5a.5.5 0 0 1-1 0V1.707L6.354 2.854a.5.5 0 1 1-.708-.708l2-2zM8 10a.5.5 0 0 1 .5.5v3.793l1.146-1.147a.5.5 0 0 1 .708.708l-2 2a.5.5 0 0 1-.708 0l-2-2a.5.5 0 0 1 .708-.708L7.5 14.293V10.5A.5.5 0 0 1 8 10z"/>
                </svg>
              </button>
              <button
                class="btn btn-sm btn-outline-secondary"
                onclick={() => taskListComponent?.collapseAll()}
                title="Collapse all"
                aria-label="Collapse all tasks"
              >
                <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
                  <path d="M1 8a.5.5 0 0 1 .5-.5h13a.5.5 0 0 1 0 1h-13A.5.5 0 0 1 1 8zm7-8a.5.5 0 0 1 .5.5v3.793l1.146-1.147a.5.5 0 0 1 .708.708l-2 2a.5.5 0 0 1-.708 0l-2-2a.5.5 0 1 1 .708-.708L7.5 4.293V.5A.5.5 0 0 1 8 0zm-.5 11.707l-1.146 1.147a.5.5 0 0 1-.708-.708l2-2a.5.5 0 0 1 .708 0l2 2a.5.5 0 0 1-.708.708L8.5 11.707V15.5a.5.5 0 0 1-1 0v-3.793z"/>
                </svg>
              </button>
            </div>
          </div>

          <div class="col-12 col-sm-auto">
            <button
              class="btn btn-warning w-100"
              onclick={openAddModal}
              style="background-color: #f97316; border-color: #f97316;"
            >
              <i class="bi bi-plus-lg"></i> Add Task
            </button>
          </div>
        </div>
      </header>

      <main>
        {#if selectedCategory}
          <TaskList
            bind:this={taskListComponent}
            {tasks}
            onComplete={handleCompleteTask}
            onEdit={openEditModal}
            onReorder={handleReorderTasks}
          />
        {/if}
      </main>
    </div>
  </div>
</div>

{#if showAddModal}
  <AddTaskModal
    categories={categories}
    defaultCategory={selectedCategory}
    tasks={tasks}
    task={editingTask}
    onSave={editingTask ? handleEditTask : handleAddTask}
    onClose={closeModal}
  />
{/if}

<Toast message={toastMessage} show={showToast} />

<style>
  .category-select {
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
  }
</style>
