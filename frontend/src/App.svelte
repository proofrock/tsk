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
      selectedCategory = categories[0];
      await loadTasks();
    }
  }

  async function loadTasks() {
    if (!selectedCategory) return;
    const res = await fetch(`${API_BASE}/tasks?category_id=${selectedCategory.id}`);
    tasks = await res.json();
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

  async function handleReorderTasks(taskIds) {
    const res = await fetch(`${API_BASE}/tasks/reorder`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ task_ids: taskIds }),
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
        <h1 class="gradient-text fw-bold mb-3">tsk</h1>

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

          <div class="col-12 col-sm-6 col-md-auto ms-md-auto">
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
            {tasks}
            onComplete={handleCompleteTask}
            onEdit={openEditModal}
            onReorder={handleReorderTasks}
          />
        {/if}
      </main>

      {#if version}
        <footer class="mt-4 text-center">
          <small class="text-secondary">v{version}</small>
        </footer>
      {/if}
    </div>
  </div>
</div>

{#if showAddModal}
  <AddTaskModal
    categories={categories}
    defaultCategory={selectedCategory}
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
